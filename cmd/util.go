package main

import (
	"errors"
	"github.com/hashicorp/vault/api"
	"github.com/soerenschneider/ssh-key-signer/internal/signature"
	"github.com/soerenschneider/ssh-key-signer/internal/signature/vault"
	"github.com/soerenschneider/ssh-key-signer/pkg/ssh"
	"os/user"
	"path/filepath"
	"strings"
)

func getVaultConfig(conf *Config) *api.Config {
	vaultConfig := api.DefaultConfig()
	vaultConfig.MaxRetries = 5
	vaultConfig.Address = conf.VaultAddress
	return vaultConfig
}

func buildAuthImpl(client *api.Client, conf *Config) (vault.AuthMethod, error) {
	token := conf.VaultToken
	if len(token) > 0 {
		return vault.NewTokenAuth(token)
	}

	if conf.VaultAuthImplicit {
		return vault.NewTokenImplicitAuth(), nil
	}

	approleData := make(map[string]string)
	approleData[vault.KeyRoleId] = conf.VaultRoleId
	approleData[vault.KeySecretId] = conf.VaultSecretId
	approleData[vault.KeySecretIdFile] = conf.VaultSecretIdFile

	return vault.NewAppRoleAuth(client, approleData)
}

func buildSignatureStrategy(config *Config) (ssh.RefreshSignatureStrategy, error) {
	if config.ForceNewSignature {
		return ssh.NewSimpleStrategy(true), nil
	}

	return ssh.NewPercentageStrategy(config.CertificateLifetimeThresholdPercentage)
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

func buildPublicKeySink(config *Config) (signature.Sink, error) {
	if nil == config {
		return nil, errors.New("empty config supplied")
	}

	if len(config.PublicKeyFile) == 0 {
		return nil, errors.New("no public key file supplied")
	}

	expanded := getExpandedFile(config.PublicKeyFile)
	return &signature.FileSink{FilePath: expanded}, nil
}

func buildSignedKeySink(config *Config) (signature.Sink, error) {
	if nil == config {
		return nil, errors.New("empty config supplied")
	}

	if len(config.SignedKeyFile) == 0 {
		return nil, errors.New("no signed key file supplied")
	}

	expanded := getExpandedFile(config.SignedKeyFile)
	return &signature.FileSink{FilePath: expanded}, nil
}
