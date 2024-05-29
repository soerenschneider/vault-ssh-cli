package main

import (
	"fmt"
	"os"

	"github.com/soerenschneider/vault-ssh-cli/internal"
	"github.com/soerenschneider/vault-ssh-cli/internal/config"
	"github.com/spf13/viper"

	log "github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func getSignHostKeyCmd() *cobra.Command {
	var signCmd = &cobra.Command{
		Use:   "sign-host-key",
		Short: "Sign a SSH host public key",
		PreRun: func(cmd *cobra.Command, args []string) {
			log.Info().Msgf("Starting up version %s (%s)", internal.BuildVersion, internal.CommitHash)
		},
		Run: signHostKeyEntryPoint,
	}

	signCmd.Flags().StringP(config.FLAG_VAULT_SSH_ROLE, "r", "", "The Vault role to use")
	signCmd.Flags().Bool(config.FLAG_FORCE_NEW_SIGNATURE, config.FLAG_FORCE_NEW_SIGNATURE_DEFAULT, "Force signing a public key")
	signCmd.Flags().StringP(config.FLAG_PUBKEY_FILE, "p", "", "Public key file to sign")
	signCmd.Flags().StringP(config.FLAG_SIGNED_KEY_FILE, "s", "", "File to write signature to")
	signCmd.Flags().Float32(config.FLAG_RENEW_THRESHOLD_PERCENTAGE, config.FLAG_RENEW_THRESHOLD_PERCENTAGE_DEFAULT, "Sign key after passing lifetime threshold (in %)")
	signCmd.Flags().String(config.FLAG_METRICS_FILE, config.FLAG_METRICS_FILE_DEFAULT, "File to write metrics to")
	signCmd.Flags().Int(config.FLAG_TTL, 0, "TTL for the signed certificate")
	signCmd.Flags().StringSlice(config.FLAG_PRINCIPALS, []string{}, "Principals")

	return signCmd
}

func signHostKeyEntryPoint(ccmd *cobra.Command, args []string) {
	viper.SetDefault(config.FLAG_RENEW_THRESHOLD_PERCENTAGE, config.FLAG_RENEW_THRESHOLD_PERCENTAGE_DEFAULT)
	viper.SetDefault(config.FLAG_METRICS_FILE, config.FLAG_METRICS_FILE_DEFAULT)

	conf, err := getConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("could not read config")
	}
	config.Print(conf)
	err = signHostKey(conf)
	if err != nil {
		log.Error().Msgf("signing key not successful, %v", err)
		internal.MetricSuccess.Set(0)
	} else {
		internal.MetricSuccess.Set(1)
	}
	internal.MetricRunTimestamp.SetToCurrentTime()
	if len(config.MetricsFile) > 0 {
		if err := internal.WriteMetrics(config.MetricsFile); err != nil {
			log.Warn().Err(err).Msg("could not write metrics")
		}
	}

	if err != nil {
		os.Exit(1)
	}
}

func signHostKey(config *config.Config) error {
	err := config.ValidateSignCommand()
	if err != nil {
		return fmt.Errorf("invalid conf: %w", err)
	}
	app := buildApp(config)
	keys := buildKeys(config)

	if err = app.issuer.SignHostCert(keys.pub, keys.sign); err != nil {
		return fmt.Errorf("could not sign public key: %v", err)
	}
	return nil
}
