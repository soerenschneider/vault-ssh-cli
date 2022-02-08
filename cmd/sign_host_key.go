package main

import (
	"fmt"
	"github.com/soerenschneider/ssh-key-signer/internal"
	"github.com/soerenschneider/ssh-key-signer/internal/signature"
	"github.com/soerenschneider/ssh-key-signer/internal/signature/vault"
	"github.com/soerenschneider/ssh-key-signer/pkg/ssh"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/hashicorp/vault/api"
	log "github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func getSignHostKeyCmd() *cobra.Command {
	var signCmd = &cobra.Command{
		Use:   "sign-host-key",
		Short: "Sign a SSH host public key",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return initializeConfig(cmd)
		},
		Run: signHostKeyEntryPoint,
	}

	signCmd.PersistentFlags().BoolP("debug", "v", false, "Enable debug logging")

	signCmd.PersistentFlags().StringP(FLAG_VAULT_ADDRESS, "a", "", "Vault instance to connect to. If not specified, falls back to env var VAULT_ADDR.")
	viper.BindPFlag(FLAG_VAULT_ADDRESS, signCmd.PersistentFlags().Lookup(FLAG_VAULT_ADDRESS))

	signCmd.PersistentFlags().StringP(FLAG_VAULT_TOKEN, "t", "", "Vault token to use for authentication. Can not be used in conjunction with AppRole login data.")
	viper.BindPFlag(FLAG_VAULT_TOKEN, signCmd.PersistentFlags().Lookup(FLAG_VAULT_TOKEN))

	signCmd.PersistentFlags().StringP(FLAG_VAULT_ROLE_ID, "r", "", "Vault role_id to use for AppRole login. Can not be used in conjuction with Vault token flag.")
	viper.BindPFlag(FLAG_VAULT_ROLE_ID, signCmd.PersistentFlags().Lookup(FLAG_VAULT_ROLE_ID))

	signCmd.PersistentFlags().StringP(FLAG_VAULT_SECRET_ID, "s", "", "Vault secret_id to use for AppRole login. Can not be used in conjuction with Vault token flag.")
	viper.BindPFlag(FLAG_VAULT_SECRET_ID, signCmd.PersistentFlags().Lookup(FLAG_VAULT_SECRET_ID))

	signCmd.PersistentFlags().StringP(FLAG_VAULT_SECRET_ID_FILE, "", "", "Flat file to read Vault secret_id from. Can not be used in conjuction with Vault token flag.")
	viper.BindPFlag(FLAG_VAULT_SECRET_ID_FILE, signCmd.PersistentFlags().Lookup(FLAG_VAULT_SECRET_ID_FILE))

	signCmd.PersistentFlags().StringP(FLAG_VAULT_MOUNT_SSH, "", "ssh", "Path where the SSH secret engine is mounted.")
	viper.BindPFlag(FLAG_VAULT_MOUNT_SSH, signCmd.PersistentFlags().Lookup(FLAG_VAULT_MOUNT_SSH))

	signCmd.PersistentFlags().StringP(FLAG_VAULT_MOUNT_APPROLE, "", "approle", "Path where the AppRole auth method is mounted.")
	viper.BindPFlag(FLAG_VAULT_MOUNT_APPROLE, signCmd.PersistentFlags().Lookup(FLAG_VAULT_MOUNT_APPROLE))

	signCmd.PersistentFlags().StringP(FLAG_VAULT_SSH_BACKEND_ROLE, "", "host_key_sign", "The name of the SSH role backend.")
	viper.BindPFlag(FLAG_VAULT_SSH_BACKEND_ROLE, signCmd.PersistentFlags().Lookup(FLAG_VAULT_SSH_BACKEND_ROLE))

	signCmd.PersistentFlags().StringP(FLAG_CONFIG_FILE, "c", "", "File to load configuration from")
	viper.BindPFlag(FLAG_CONFIG_FILE, signCmd.PersistentFlags().Lookup(FLAG_CONFIG_FILE))

	signCmd.PersistentFlags().BoolP(FLAG_FORCE_NEW_SIGNATURE, "", FLAG_FORCE_NEW_SIGNATURE_DEFAULT, "Force new signature, regardless of its lifetime")
	viper.BindPFlag(FLAG_FORCE_NEW_SIGNATURE, signCmd.PersistentFlags().Lookup(FLAG_FORCE_NEW_SIGNATURE))

	signCmd.PersistentFlags().Float64P(FLAG_LIFETIME_THRESHOLD_PERCENTAGE, "", FLAG_LIFETIME_THRESHOLD_PERCENTAGE_DEFAULT, "Create new signature after lifetime x percent of lifetime has been reached")
	viper.BindPFlag(FLAG_LIFETIME_THRESHOLD_PERCENTAGE, signCmd.PersistentFlags().Lookup(FLAG_LIFETIME_THRESHOLD_PERCENTAGE))

	signCmd.PersistentFlags().StringP(FLAG_PUBKEY_FILE, "p", "", "SSH Public Host Key to sign")
	viper.BindPFlag(FLAG_PUBKEY_FILE, signCmd.PersistentFlags().Lookup(FLAG_PUBKEY_FILE))

	signCmd.PersistentFlags().StringP(FLAG_SIGNED_KEY_FILE, "o", "", "File to write the signed key to")
	viper.BindPFlag(FLAG_SIGNED_KEY_FILE, signCmd.PersistentFlags().Lookup(FLAG_SIGNED_KEY_FILE))

	signCmd.PersistentFlags().StringP(FLAG_METRICS_FILE, "", FLAG_METRICS_FILE_DEFAULT, "File to write metrics to")
	viper.BindPFlag(FLAG_METRICS_FILE, signCmd.PersistentFlags().Lookup(FLAG_METRICS_FILE))

	signCmd.MarkFlagRequired(FLAG_PUBKEY_FILE)
	signCmd.MarkFlagRequired(FLAG_SIGNED_KEY_FILE)

	return signCmd
}

func signHostKeyEntryPoint(ccmd *cobra.Command, args []string) {
	log.Info().Msgf("Starting up version %s (%s)", internal.BuildVersion, internal.CommitHash)
	configFile := viper.GetViper().GetString(FLAG_CONFIG_FILE)
	if len(configFile) > 0 {
		err := readConfig(configFile)
		if err != nil {
			log.Fatal().Msgf("Could not load desired config file: %s: %v", configFile, err)
		}
		log.Info().Msgf("Read config from file %s", viper.ConfigFileUsed())
	}

	config := NewConfigFromViper()
	config.PrintConfig()

	err := signHostKey(config)
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

func signHostKey(config Config) error {
	errors := config.Validate()
	if len(errors) > 0 {
		fmtErrors := make([]string, len(errors))
		for i, er := range errors {
			fmtErrors[i] = fmt.Sprintf("\"%s\"", er)
		}
		return fmt.Errorf("invalid config, %d errors: %s", len(errors), strings.Join(fmtErrors, ", "))
	}

	vaultClient, err := api.NewClient(getVaultConfig(&config))
	if err != nil {
		return fmt.Errorf("could not build vault client: %v", err)
	}

	authStrategy, err := buildAuthImpl(vaultClient, &config)
	if err != nil {
		return fmt.Errorf("could not build auth strategy: %v", err)
	}

	signingImpl, err := vault.NewVaultSigner(vaultClient, authStrategy, config.VaultSshBackend)
	if err != nil {
		return fmt.Errorf("could not build rotation client: %v", err)
	}

	strat, err := buildSignatureStrategy(&config)
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
