package main

import (
	"fmt"
	"github.com/soerenschneider/ssh-key-signer/internal/signature"
	"github.com/soerenschneider/ssh-key-signer/internal/signature/vault"
	"github.com/soerenschneider/ssh-key-signer/pkg/ssh"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/hashicorp/vault/api"
)

func signPublicKey(config Config) error {
	errors := config.Validate()
	if len(errors) > 0 {
		fmtErrors := make([]string, len(errors))
		for i, er := range errors {
			fmtErrors[i] = fmt.Sprintf("\"%s\"", er)
		}
		return fmt.Errorf("invalid config, %d errors: %s", len(errors), strings.Join(fmtErrors, ", "))
	}

	pubKeyFile := config.PublicKeyFile
	err := testIsFile(pubKeyFile)
	if err != nil || len(pubKeyFile) == 0 {
		return fmt.Errorf("invalid pubkey file: %v", err)
	}

	vaultClient, err := api.NewClient(getVaultConfig(&config))
	if err != nil {
		return fmt.Errorf("could not build vault client: %v", err)
	}

	authStrategy, err := buildAuthImpl(vaultClient, &config)
	if err != nil {
		return fmt.Errorf("could not build auth strategy: %v", err)
	}

	backendName := config.VaultSshBackend
	signingImpl, err := vault.NewVaultSigner(vaultClient, authStrategy, backendName)
	if err != nil {
		return fmt.Errorf("could not build rotation client: %v", err)
	}

	strat, err := buildSignatureStrategy(&config)
	if err != nil {
		return fmt.Errorf("could not build signature refresh stratey: %v", err)
	}

	pubKeyPod := &signature.FsPod{FilePath: config.PublicKeyFile}
	var signedKeyPod signature.KeyPod = &signature.BufferPod{}
	if len(config.SignedKeyFile) > 0 {
		signedKeyPod = &signature.FsPod{FilePath: config.SignedKeyFile}
	}

	issuer, err := signature.NewIssuer(signingImpl, strat)
	if err != nil {
		return fmt.Errorf("could not build issuer: %v", err)
	}

	err = issuer.SignHostCert(pubKeyPod, signedKeyPod)
	if err != nil {
		return fmt.Errorf("could not sign public key: %v", err)
	}

	return nil
}

func getVaultConfig(conf *Config) *api.Config {
	vaultConfig := api.DefaultConfig()

	vaultAddress := conf.VaultAddress
	if len(vaultAddress) == 0 {
		vaultAddress = os.Getenv("VAULT_ADDRESS")
	}

	vaultConfig.Address = vaultAddress
	return vaultConfig
}

func buildAuthImpl(client *api.Client, conf *Config) (vault.AuthMethod, error) {
	token := conf.VaultToken
	if len(token) > 0 {
		return vault.NewTokenAuth(token)
	}

	approleData := make(map[string]string)
	approleData["role_id"] = conf.VaultRoleId
	approleData["secret_id"] = conf.VaultSecretId
	approleData["secret_id_file"] = conf.VaultSecretIdFile

	return vault.NewAppRoleAuth(client, approleData)
}

func buildSignatureStrategy(config *Config) (ssh.RefreshSignatureStrategy, error) {
	if config.ForceNewSignature {
		return ssh.NewSimpleStrategy(true), nil
	}

	return ssh.NewPercentageStrategy(config.CertificateLifetimeThresholdPercentage)
}

func testIsFile(file string) error {
	_, err := os.Stat(file)
	return err
}

func getExpandedFile(filename string) string {
	usr, _ := user.Current()
	dir := usr.HomeDir

	if strings.HasPrefix(filename, "~/") {
		return filepath.Join(dir, filename[2:])
	}

	if strings.HasPrefix(filename, "$HOME/") {
		return filepath.Join(dir, filename[6:])
	}

	return filename
}
