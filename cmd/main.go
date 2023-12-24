package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
	log "github.com/rs/zerolog/log"
	"github.com/soerenschneider/vault-ssh-cli/internal"
	config "github.com/soerenschneider/vault-ssh-cli/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

var cliDescription = fmt.Sprintf("vault-ssh-cli - %s", internal.BuildVersion)

const (
	cliName               = "vault-ssh-cli"
	envPrefix             = "VAULT_SSH_CLI"
	defaultConfigFilename = "config"
)

func main() {
	initLogging()
	root := &cobra.Command{
		Use:   cliName,
		Short: cliDescription,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var errs []error
			cmd.Flags().Visit(func(flag *pflag.Flag) {
				err := viper.BindPFlag(flag.Name, cmd.Flags().Lookup(flag.Name))
				if err != nil {
					errs = append(errs, err)
				}
			})
			if len(errs) > 0 {
				return fmt.Errorf("can't bind flags: %v", errs)
			}
			return nil
		},
	}

	root.PersistentFlags().StringP(config.FLAG_VAULT_ADDRESS, "a", "", "Vault instance to connect to. If not specified, falls back to env var VAULT_ADDR.")
	root.PersistentFlags().StringP(config.FLAG_VAULT_AUTH_TOKEN, "t", "", "Vault token to use for authentication. Can not be used in conjunction with AppRole login data.")
	root.PersistentFlags().BoolP(config.FLAG_VAULT_AUTH_IMPLICIT, "i", false, "Try to implicitly authenticate to vault using VAULT_TOKEN env var or ~/.vault-token file.")
	root.PersistentFlags().StringP(config.FLAG_VAULT_AUTH_APPROLE_ROLE_ID, "", "", "Vault role_id to use for AppRole login. Can not be used in conjuction with Vault token flag.")
	root.PersistentFlags().StringP(config.FLAG_VAULT_AUTH_APPROLE_SECRET_ID, "", "", "Vault secret_id to use for AppRole login. Can not be used in conjuction with Vault token flag.")
	root.PersistentFlags().StringP(config.FLAG_VAULT_AUTH_APPROLE_SECRET_ID_FILE, "", "", "Flat file to read Vault secret_id from. Can not be used in conjuction with Vault token flag.")
	root.PersistentFlags().StringP(config.FLAG_VAULT_AUTH_APPROLE_MOUNT, "", config.FLAG_VAULT_AUTH_APPROLE_MOUNT_DEFAULT, "Path where the AppRole auth method is mounted.")
	root.PersistentFlags().StringP(config.FLAG_VAULT_SSH_MOUNT, "m", config.FLAG_VAULT_SSH_MOUNT_DEFAULT, "Path where the PKI secret engine is mounted.")
	root.PersistentFlags().StringP(config.FLAG_CONFIG_FILE, "c", "", "File to read config from")
	root.PersistentFlags().BoolP(config.FLAG_DEBUG, "", false, "Set loglevel to debug")

	root.AddCommand(readCaCertCmd())
	root.AddCommand(getSignHostKeyCmd())
	root.AddCommand(getSignUserKeyCmd())
	root.AddCommand(versionCmd)

	if err := root.Execute(); err != nil {
		log.Fatal().Err(err).Msgf("could not run %s", cliName)
	}
}

func getConfig() (*config.Config, error) {
	viper.SetDefault(config.FLAG_VAULT_SSH_MOUNT, config.FLAG_VAULT_SSH_MOUNT_DEFAULT)
	viper.SetDefault(config.FLAG_VAULT_AUTH_APPROLE_MOUNT, config.FLAG_VAULT_AUTH_APPROLE_MOUNT_DEFAULT)

	viper.SetConfigName(defaultConfigFilename)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.config/vault-ssh-cli")
	viper.AddConfigPath("/etc/vault-ssh-cli/")

	if viper.IsSet(config.FLAG_CONFIG_FILE) {
		configFile := viper.GetString(config.FLAG_CONFIG_FILE)
		log.Info().Msgf("Trying to read config from '%s'", configFile)
		viper.SetConfigFile(configFile)
	}

	viper.SetEnvPrefix(envPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil && viper.IsSet(config.FLAG_CONFIG_FILE) {
		log.Fatal().Msgf("Can't read config: %v", err)
	}

	var config *config.Config
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatal().Msgf("unable to decode into struct, %v", err)
	}

	config.ExpandPaths()
	setupLogLevel(config.Debug)
	return config, nil
}

func setupLogLevel(debug bool) {
	level := zerolog.InfoLevel
	if debug {
		level = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(level)
}

func initLogging() {
	if term.IsTerminal(int(os.Stdout.Fd())) {
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: "15:04:05",
		})
	}
}
