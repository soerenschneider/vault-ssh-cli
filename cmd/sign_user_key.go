package main

import (
	"fmt"
	"github.com/soerenschneider/ssh-key-signer/internal"
	"github.com/soerenschneider/ssh-key-signer/internal/signature"
	"github.com/soerenschneider/ssh-key-signer/internal/signature/vault"
	"github.com/spf13/viper"
	"os"
	"strings"

	"github.com/hashicorp/vault/api"
	log "github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func getSignUserKeyCmd() *cobra.Command {
	var signCmd = &cobra.Command{
		Use:   "sign-user-key",
		Short: "Sign a SSH user public key",
		Run:   signUserKeyEntryPoint,
	}

	signCmd.PersistentFlags().String(FLAG_VAULT_SSH_ROLE, "", "Write the ca certificate to this output file")
	signCmd.PersistentFlags().Bool(FLAG_FORCE_NEW_SIGNATURE, FLAG_FORCE_NEW_SIGNATURE_DEFAULT, "Force signing a public key")
	signCmd.PersistentFlags().StringP(FLAG_PUBKEY_FILE, "p", "", "Public key file to sign")
	signCmd.PersistentFlags().StringP(FLAG_SIGNED_KEY_FILE, "s", "", "File to write signature to")
	signCmd.PersistentFlags().Float32(FLAG_RENEW_THRESHOLD_PERCENTAGE, FLAG_RENEW_THRESHOLD_PERCENTAGE_DEFAULT, "Sign key after passing lifetime threshold (in %)")
	signCmd.PersistentFlags().String(FLAG_METRICS_FILE, FLAG_METRICS_FILE_DEFAULT, "File to write metrics to")

	viper.SetDefault(FLAG_RENEW_THRESHOLD_PERCENTAGE, FLAG_RENEW_THRESHOLD_PERCENTAGE_DEFAULT)
	viper.SetDefault(FLAG_METRICS_FILE, FLAG_METRICS_FILE_DEFAULT)

	return signCmd
}

func signUserKeyEntryPoint(ccmd *cobra.Command, args []string) {
	log.Info().Msgf("Starting up version %s (%s)", internal.BuildVersion, internal.CommitHash)
	config, err := config()
	if err != nil {
		log.Fatal().Err(err).Msg("could not read config")
	}
	config.Print()
	err = signUserKey(config)
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

func signUserKey(config *Config) error {
	errors := config.ValidateSignCommand()
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

	signingImpl, err := vault.NewVaultSigner(vaultClient, authStrategy, config.VaultMountSsh, config.VaultSshRole)
	if err != nil {
		return fmt.Errorf("could not build rotation client: %v", err)
	}

	strat, err := buildSignatureStrategy(config)
	if err != nil {
		return fmt.Errorf("could not build signature refresh stratey: %v", err)
	}

	pubKeyPod, err := buildPublicKeyPod(config)
	if err != nil {
		return fmt.Errorf("can't build sink to read public key from: %v", err)
	}
	err = pubKeyPod.CanRead()
	if err != nil {
		return fmt.Errorf("can not read from public key %s: %v", config.PublicKeyFile, err)
	}

	signedKeyPod, err := buildSignedKeyPod(config)
	if err != nil {
		return fmt.Errorf("can't build sink to write signature to: %v", err)
	}
	err = signedKeyPod.CanWrite()
	if err != nil {
		return fmt.Errorf("%s '%s' is not writable", FLAG_SIGNED_KEY_FILE, config.SignedKeyFile)
	}

	issuer, err := signature.NewIssuer(signingImpl, strat)
	if err != nil {
		return fmt.Errorf("could not build issuer: %v", err)
	}

	err = issuer.SignClientCert(pubKeyPod, signedKeyPod)
	if err != nil {
		return fmt.Errorf("could not sign public key: %v", err)
	}

	return nil
}
