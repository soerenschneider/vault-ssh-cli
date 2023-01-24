package main

import (
	"fmt"
	"github.com/soerenschneider/ssh-key-signer/internal"
	"github.com/soerenschneider/ssh-key-signer/internal/signature"
	"github.com/soerenschneider/ssh-key-signer/internal/signature/vault"
	"github.com/soerenschneider/ssh-key-signer/pkg/ssh"
	"github.com/spf13/viper"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/hashicorp/vault/api"
	log "github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func getSignHostKeyCmd() *cobra.Command {
	var signCmd = &cobra.Command{
		Use:   "sign-host-key",
		Short: "Sign a SSH host public key",
		Run:   signHostKeyEntryPoint,
	}

	signCmd.PersistentFlags().String(FLAG_VAULT_SSH_BACKEND_ROLE, "", "Write the ca certificate to this output file")
	signCmd.PersistentFlags().Bool(FLAG_FORCE_NEW_SIGNATURE, FLAG_FORCE_NEW_SIGNATURE_DEFAULT, "Force signing a public key")
	signCmd.PersistentFlags().String(FLAG_PUBKEY_FILE, "", "Force signing a public key")
	signCmd.PersistentFlags().Float32(FLAG_RENEW_THRESHOLD_PERCENTAGE, FLAG_RENEW_THRESHOLD_PERCENTAGE_DEFAULT, "Sign key after passing lifetime threshold (in %)")
	signCmd.PersistentFlags().String(FLAG_METRICS_FILE, FLAG_METRICS_FILE_DEFAULT, "File to write metrics to")

	viper.SetDefault(FLAG_RENEW_THRESHOLD_PERCENTAGE, FLAG_RENEW_THRESHOLD_PERCENTAGE_DEFAULT)
	viper.SetDefault(FLAG_METRICS_FILE, FLAG_METRICS_FILE_DEFAULT)

	return signCmd
}

func signHostKeyEntryPoint(ccmd *cobra.Command, args []string) {
	log.Info().Msgf("Starting up version %s (%s)", internal.BuildVersion, internal.CommitHash)
	config, err := config()
	if err != nil {
		log.Fatal().Err(err).Msg("could not read config")
	}
	config.Print()
	err = signHostKey(config)
	if err != nil {
		log.Error().Msgf("signing key not successful, %v", err)
		internal.MetricSuccess.Set(0)
	} else {
		internal.MetricSuccess.Set(1)
	}
	internal.MetricRunTimestamp.SetToCurrentTime()
	if len(config.MetricsFile) > 0 {
		internal.WriteMetrics(config.MetricsFile)
	}

	if err == nil {
		os.Exit(0)
	}
	os.Exit(1)
}

func signHostKey(config *Config) error {
	errors := config.Validate()
	if len(errors) > 0 {
		fmtErrors := make([]string, len(errors))
		for i, er := range errors {
			fmtErrors[i] = fmt.Sprintf("\"%s\"", er)
		}
		return fmt.Errorf("invalid config, %d errors: %s", len(errors), strings.Join(fmtErrors, ", "))
	}

	vaultClient, err := api.NewClient(getVaultConfig(config))
	if err != nil {
		return fmt.Errorf("could not build vault client: %v", err)
	}

	authStrategy, err := buildAuthImpl(vaultClient, config)
	if err != nil {
		return fmt.Errorf("could not build auth strategy: %v", err)
	}

	signingImpl, err := vault.NewVaultSigner(vaultClient, authStrategy, config.VaultSshBackend)
	if err != nil {
		return fmt.Errorf("could not build rotation client: %v", err)
	}

	strat, err := buildSignatureStrategy(config)
	if err != nil {
		return fmt.Errorf("could not build signature refresh stratey: %v", err)
	}

	pubKeyPod := &signature.FsPod{FilePath: config.PublicKeyFile}
	err = pubKeyPod.CanRead()
	if err != nil {
		return fmt.Errorf("can not read from public key %s: %v", config.PublicKeyFile, err)
	}

	var signedKeyPod signature.KeyPod = &signature.BufferPod{}
	if len(config.SignedKeyFile) > 0 {
		signedKeyPod = &signature.FsPod{FilePath: config.SignedKeyFile}
	}
	err = signedKeyPod.CanWrite()
	if err != nil {
		return fmt.Errorf("%s '%s' is not writable", FLAG_SIGNED_KEY_FILE, config.SignedKeyFile)
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
	vaultConfig.MaxRetries = 5
	vaultConfig.Address = conf.VaultAddress
	return vaultConfig
}

func buildAuthImpl(client *api.Client, conf *Config) (vault.AuthMethod, error) {
	token := conf.VaultToken
	if len(token) > 0 {
		return vault.NewTokenAuth(token)
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
