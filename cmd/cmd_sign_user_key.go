package main

import (
	"os"

	log "github.com/rs/zerolog/log"
	"github.com/soerenschneider/vault-ssh-cli/internal"
	"github.com/soerenschneider/vault-ssh-cli/internal/config"
	"github.com/soerenschneider/vault-ssh-cli/pkg/signature"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func getSignUserKeyCmd() *cobra.Command {
	var signCmd = &cobra.Command{
		Use:   "sign-user-key",
		Short: "Sign a SSH user public key",
		PreRun: func(cmd *cobra.Command, args []string) {
			log.Info().Msgf("Starting up version %s (%s)", internal.BuildVersion, internal.CommitHash)
		},
		Run: signUserKeyEntryPoint,
	}

	signCmd.Flags().StringP(config.FLAG_VAULT_SSH_ROLE, "r", "", "The Vault role to use")
	signCmd.Flags().Bool(config.FLAG_FORCE_NEW_SIGNATURE, config.FLAG_FORCE_NEW_SIGNATURE_DEFAULT, "Force signing a public key")
	signCmd.Flags().StringP(config.FLAG_PUBKEY_FILE, "p", "", "Public key file to sign")
	signCmd.Flags().StringP(config.FLAG_SIGNED_KEY_FILE, "s", "", "File to write signature to")
	signCmd.Flags().Float32(config.FLAG_RENEW_THRESHOLD_PERCENTAGE, config.FLAG_RENEW_THRESHOLD_PERCENTAGE_DEFAULT, "Sign key after passing lifetime threshold (in %)")
	signCmd.Flags().String(config.FLAG_METRICS_FILE, "", "File to write metrics to")
	signCmd.Flags().Int(config.FLAG_TTL, 0, "TTL for the signed certificate")
	signCmd.Flags().Int(config.FLAG_RETRIES, config.FLAG_RETRIES_DEFAULT, "The amount of retries to perform on non-permanent errors")
	signCmd.Flags().StringSlice(config.FLAG_PRINCIPALS, []string{}, "Principals")

	return signCmd
}

func signUserKeyEntryPoint(ccmd *cobra.Command, args []string) {
	viper.SetDefault(config.FLAG_RENEW_THRESHOLD_PERCENTAGE, config.FLAG_RENEW_THRESHOLD_PERCENTAGE_DEFAULT)
	viper.SetDefault(config.FLAG_RETRIES, config.FLAG_RETRIES_DEFAULT)

	conf, err := getConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("could not read config")
	}
	config.Print(conf)
	err = signUserKey(conf)
	if err != nil {
		log.Error().Msgf("signing key not successful, %v", err)
		internal.MetricSuccess.Set(0)
	} else {
		internal.MetricSuccess.Set(1)
	}

	internal.MetricRunTimestamp.SetToCurrentTime()
	if len(conf.MetricsFile) > 0 {
		if err := internal.WriteMetrics(conf.MetricsFile); err != nil {
			log.Warn().Err(err).Msg("could not write metrics")
		}
	}

	if err != nil {
		os.Exit(1)
	}
}

func signUserKey(conf *config.Config) error {
	if err := config.Validate(conf); err != nil {
		return err
	}

	app := buildApp(conf)
	keys := buildKeys(conf)

	request := signature.SignatureRequest{
		Ttl:        conf.Ttl,
		Principals: conf.Principals,
		Extensions: conf.Extensions,
		VaultRole:  conf.VaultSshRole,
	}

	outcome, err := app.signatureService.SignUserCert(request, keys.pub, keys.sign)
	writeLogs(outcome)
	updateCertMetrics(outcome)
	return err
}
