package main

import (
	"fmt"
	"github.com/soerenschneider/ssh-key-signer/internal"
	"os"
	"strings"

	"github.com/rs/zerolog"
	log "github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	envPrefix             = "SSH_KEY_SIGNER"
	defaultConfigFilename = "config"

	FLAG_VAULT_ADDRESS          = "vault-address"
	FLAG_VAULT_TOKEN            = "vault-token"
	FLAG_VAULT_ROLE_ID          = "vault-role-id"
	FLAG_VAULT_SECRET_ID        = "vault-secret-id"
	FLAG_VAULT_SECRET_ID_FILE   = "vault-secret-id-file"
	FLAG_VAULT_MOUNT_SSH        = "vault-mount-ssh"
	FLAG_VAULT_MOUNT_APPROLE    = "vault-mount-approle"
	FLAG_VAULT_SSH_BACKEND_ROLE = "vault-ssh-role-name"

	FLAG_FORCE_NEW_SIGNATURE                   = "force-new-signature"
	FLAG_FORCE_NEW_SIGNATURE_DEFAULT           = false
	FLAG_LIFETIME_THRESHOLD_PERCENTAGE         = "lifetime-threshold-percent"
	FLAG_LIFETIME_THRESHOLD_PERCENTAGE_DEFAULT = 33.

	FLAG_PUBKEY_FILE     = "pub-key-file"
	FLAG_SIGNED_KEY_FILE = "signed-key-file"
	FLAG_CONFIG_FILE     = "config-file"

	FLAG_METRICS_FILE         = "metrics-file"
	FLAG_METRICS_FILE_DEFAULT = "/var/lib/node_exporter/ssh_key_sign.prom"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Info().Msgf("Starting up version %s (%s)", internal.BuildVersion, internal.CommitHash)

	root := &cobra.Command{
		Use: "boing",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return initializeConfig(cmd)
		},
		Run: func(ccmd *cobra.Command, args []string) {
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

			err := signPublicKey(config)
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
		},
	}

	root.PersistentFlags().BoolP("debug", "v", false, "Enable debug logging")

	root.PersistentFlags().StringP(FLAG_VAULT_ADDRESS, "a", "", "Vault instance to connect to. If not specified, falls back to env var VAULT_ADDR.")
	viper.BindPFlag(FLAG_VAULT_ADDRESS, root.PersistentFlags().Lookup(FLAG_VAULT_ADDRESS))

	root.PersistentFlags().StringP(FLAG_VAULT_TOKEN, "t", "", "Vault token to use for authentication. Can not be used in conjunction with AppRole login data.")
	viper.BindPFlag(FLAG_VAULT_TOKEN, root.PersistentFlags().Lookup(FLAG_VAULT_TOKEN))

	root.PersistentFlags().StringP(FLAG_VAULT_ROLE_ID, "r", "", "Vault role_id to use for AppRole login. Can not be used in conjuction with Vault token flag.")
	viper.BindPFlag(FLAG_VAULT_ROLE_ID, root.PersistentFlags().Lookup(FLAG_VAULT_ROLE_ID))

	root.PersistentFlags().StringP(FLAG_VAULT_SECRET_ID, "s", "", "Vault secret_id to use for AppRole login. Can not be used in conjuction with Vault token flag.")
	viper.BindPFlag(FLAG_VAULT_SECRET_ID, root.PersistentFlags().Lookup(FLAG_VAULT_SECRET_ID))

	root.PersistentFlags().StringP(FLAG_VAULT_SECRET_ID_FILE, "", "", "Flat file to read Vault secret_id from. Can not be used in conjuction with Vault token flag.")
	viper.BindPFlag(FLAG_VAULT_SECRET_ID_FILE, root.PersistentFlags().Lookup(FLAG_VAULT_SECRET_ID_FILE))

	root.PersistentFlags().StringP(FLAG_VAULT_MOUNT_SSH, "", "ssh", "Path where the SSH secret engine is mounted.")
	viper.BindPFlag(FLAG_VAULT_MOUNT_SSH, root.PersistentFlags().Lookup(FLAG_VAULT_MOUNT_SSH))

	root.PersistentFlags().StringP(FLAG_VAULT_MOUNT_APPROLE, "", "approle", "Path where the AppRole auth method is mounted.")
	viper.BindPFlag(FLAG_VAULT_MOUNT_APPROLE, root.PersistentFlags().Lookup(FLAG_VAULT_MOUNT_APPROLE))

	root.PersistentFlags().StringP(FLAG_VAULT_SSH_BACKEND_ROLE, "", "host_key_sign", "The name of the SSH role backend.")
	viper.BindPFlag(FLAG_VAULT_SSH_BACKEND_ROLE, root.PersistentFlags().Lookup(FLAG_VAULT_SSH_BACKEND_ROLE))

	root.PersistentFlags().StringP(FLAG_CONFIG_FILE, "c", "", "File to load configuration from")
	viper.BindPFlag(FLAG_CONFIG_FILE, root.PersistentFlags().Lookup(FLAG_CONFIG_FILE))

	root.PersistentFlags().BoolP(FLAG_FORCE_NEW_SIGNATURE, "", FLAG_FORCE_NEW_SIGNATURE_DEFAULT, "Force new signature, regardless of its lifetime")
	viper.BindPFlag(FLAG_FORCE_NEW_SIGNATURE, root.PersistentFlags().Lookup(FLAG_FORCE_NEW_SIGNATURE))

	root.PersistentFlags().Float64P(FLAG_LIFETIME_THRESHOLD_PERCENTAGE, "", FLAG_LIFETIME_THRESHOLD_PERCENTAGE_DEFAULT, "Create new signature after lifetime x percent of lifetime has been reached")
	viper.BindPFlag(FLAG_LIFETIME_THRESHOLD_PERCENTAGE, root.PersistentFlags().Lookup(FLAG_LIFETIME_THRESHOLD_PERCENTAGE))

	root.PersistentFlags().StringP(FLAG_PUBKEY_FILE, "p", "", "SSH Public Host Key to sign")
	viper.BindPFlag(FLAG_PUBKEY_FILE, root.PersistentFlags().Lookup(FLAG_PUBKEY_FILE))

	root.PersistentFlags().StringP(FLAG_SIGNED_KEY_FILE, "o", "", "File to write the signed key to")
	viper.BindPFlag(FLAG_SIGNED_KEY_FILE, root.PersistentFlags().Lookup(FLAG_SIGNED_KEY_FILE))

	root.PersistentFlags().StringP(FLAG_METRICS_FILE, "", FLAG_METRICS_FILE_DEFAULT, "File to write metrics to")
	viper.BindPFlag(FLAG_METRICS_FILE, root.PersistentFlags().Lookup(FLAG_METRICS_FILE))

	root.MarkFlagRequired(FLAG_PUBKEY_FILE)
	root.MarkFlagRequired(FLAG_SIGNED_KEY_FILE)

	if err := root.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initializeConfig(cmd *cobra.Command) error {
	v := viper.GetViper()

	v.SetConfigName(defaultConfigFilename)

	v.AddConfigPath(".")
	v.AddConfigPath("$HOME/.config/ssh-key-signer")
	v.AddConfigPath("/etc/ssh-key-signer/")
	v.AddConfigPath("/etc/")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	v.SetEnvPrefix(envPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	return nil
}

func readConfig(filepath string) error {
	viper.SetConfigFile(filepath)
	return viper.ReadInConfig()
}
