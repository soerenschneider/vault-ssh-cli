package main

import (
	"fmt"
	"github.com/soerenschneider/ssh-key-signer/internal"
	"github.com/spf13/pflag"
	"os"
	"strings"

	"github.com/rs/zerolog"
	log "github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	name                  = "ssh-key-signer"
	envPrefix             = "SSH_KEY_SIGNER"
	defaultConfigFilename = "config"

	FLAG_VAULT_ADDRESS                     = "vault-address"
	FLAG_VAULT_AUTH_TOKEN                  = "vault-auth-token"
	FLAG_VAULT_AUTH_APPROLE_ROLE_ID        = "vault-auth-role-id"
	FLAG_VAULT_AUTH_APPROLE_SECRET_ID      = "vault-auth-secret-id"
	FLAG_VAULT_AUTH_APPROLE_SECRET_ID_FILE = "vault-auth-secret-id-file"
	FLAG_VAULT_AUTH_APPROLE_MOUNT          = "vault-auth-approle-mount"
	FLAG_VAULT_AUTH_APPROLE_MOUNT_DEFAULT  = "approle"
	FLAG_VAULT_AUTH_IMPLICIT               = "vault-auth-implicit"
	FLAG_VAULT_SSH_MOUNT                   = "vault-ssh-mount"
	FLAG_VAULT_SSH_MOUNT_DEFAULT           = "ssh"

	FLAG_VAULT_SSH_BACKEND_ROLE = "vault-ssh-role-name"

	FLAG_FORCE_NEW_SIGNATURE                   = "force-new-signature"
	FLAG_FORCE_NEW_SIGNATURE_DEFAULT           = false
	FLAG_LIFETIME_THRESHOLD_PERCENTAGE         = "lifetime-threshold-percent"
	FLAG_LIFETIME_THRESHOLD_PERCENTAGE_DEFAULT = 33.

	FLAG_CA_FILE         = "ca-file"
	FLAG_PUBKEY_FILE     = "pub-key-file"
	FLAG_SIGNED_KEY_FILE = "signed-key-file"
	FLAG_CONFIG_FILE     = "config-file"

	FLAG_METRICS_FILE         = "metrics-file"
	FLAG_METRICS_FILE_DEFAULT = "/var/lib/node_exporter/ssh_key_sign.prom"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	root := &cobra.Command{
		Use:   name,
		Short: fmt.Sprintf("Sign SSH keys - %s", internal.BuildVersion),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

			var errs []error
			cmd.Flags().Visit(func(flag *pflag.Flag) {

				err := viper.BindPFlag(flag.Name, cmd.Flags().Lookup(flag.Name))
				if err != nil {
					errs = append(errs, err)
				}
				log.Info().Msgf("%s=%v", flag.Name, flag.Value)

			})
			if len(errs) > 0 {
				return fmt.Errorf("can't bind flags: %v", errs)
			}
			return nil
		},
	}

	root.PersistentFlags().StringP(FLAG_VAULT_ADDRESS, "a", "", "Vault instance to connect to. If not specified, falls back to env var VAULT_ADDR.")
	root.PersistentFlags().StringP(FLAG_VAULT_AUTH_TOKEN, "t", "", "Vault token to use for authentication. Can not be used in conjunction with AppRole login data.")
	root.PersistentFlags().StringP(FLAG_VAULT_AUTH_IMPLICIT, "i", "", "Try to implicitly authenticate to vault using VAULT_TOKEN env var or ~/.vault-token file.")
	root.PersistentFlags().StringP(FLAG_VAULT_AUTH_APPROLE_ROLE_ID, "r", "", "Vault role_id to use for AppRole login. Can not be used in conjuction with Vault token flag.")
	root.PersistentFlags().StringP(FLAG_VAULT_AUTH_APPROLE_SECRET_ID, "", "", "Vault secret_id to use for AppRole login. Can not be used in conjuction with Vault token flag.")
	root.PersistentFlags().StringP(FLAG_VAULT_AUTH_APPROLE_SECRET_ID_FILE, "s", "", "Flat file to read Vault secret_id from. Can not be used in conjuction with Vault token flag.")
	root.PersistentFlags().StringP(FLAG_VAULT_AUTH_APPROLE_MOUNT, "", FLAG_VAULT_AUTH_APPROLE_MOUNT_DEFAULT, "Path where the AppRole auth method is mounted.")
	root.PersistentFlags().StringP(FLAG_VAULT_SSH_MOUNT, "", FLAG_VAULT_SSH_MOUNT_DEFAULT, "Path where the PKI secret engine is mounted.")
	root.PersistentFlags().StringP(FLAG_CONFIG_FILE, "", "", "File to read the config from")

	root.AddCommand(readCaCertCmd())
	root.AddCommand(getSignHostKeyCmd())
	root.AddCommand(versionCmd)

	if err := root.Execute(); err != nil {
		log.Fatal().Err(err).Msgf("could not run %s", name)
	}
}

func config() (*Config, error) {
	viper.SetDefault(FLAG_VAULT_SSH_MOUNT, FLAG_VAULT_SSH_MOUNT_DEFAULT)
	viper.SetDefault(FLAG_VAULT_AUTH_APPROLE_MOUNT, FLAG_VAULT_AUTH_APPROLE_MOUNT_DEFAULT)

	viper.SetConfigName(defaultConfigFilename)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.config/ssh-key-signer")
	viper.AddConfigPath("/etc/ssh-key-signer/")

	if viper.IsSet(FLAG_CONFIG_FILE) {
		configFile := viper.GetString(FLAG_CONFIG_FILE)
		log.Info().Msgf("Trying to read config from '%s'", configFile)
		viper.SetConfigFile(configFile)
	}

	viper.SetEnvPrefix(envPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil && viper.IsSet(FLAG_CONFIG_FILE) {
		log.Fatal().Msgf("Can't read config: %v", err)
	}

	var config *Config

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatal().Msgf("unable to decode into struct, %v", err)
	}

	return config, nil
}
