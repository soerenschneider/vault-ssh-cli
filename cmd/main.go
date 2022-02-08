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

	root := &cobra.Command{Use: "ssh-key-signer", Short: fmt.Sprintf("Sign SSH keys - %s", internal.BuildVersion)}

	root.AddCommand(getSignHostKeyCmd())
	root.AddCommand(versionCmd)

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
