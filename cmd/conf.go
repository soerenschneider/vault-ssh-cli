package main

import (
	"fmt"
	"os"

	log "github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	VaultAddress      string
	VaultToken        string
	VaultRoleId       string
	VaultSecretId     string
	VaultSecretIdFile string
	VaultMountSsh     string
	VaultMountApprole string
	VaultSshBackend   string

	ForceNewSignature                      bool
	CertificateLifetimeThresholdPercentage float32

	PublicKeyFile string
	SignedKeyFile string

	MetricsFile string
}

func (c *Config) Validate() []error {
	errs := make([]error, 0)
	if len(c.PublicKeyFile) == 0 {
		errs = append(errs, fmt.Errorf("empty '%s' provided", FLAG_PUBKEY_FILE))
	} else {
		_, err := os.Stat(c.PublicKeyFile)
		if err != nil {
			errs = append(errs, fmt.Errorf("couldn't access pub-key-file at '%s'", c.PublicKeyFile))
		}
	}

	emptyVaultToken := len(c.VaultToken) == 0
	emptyRoleId := len(c.VaultRoleId) == 0
	emptySecretId := len(c.VaultSecretId) == 0 && len(c.VaultSecretIdFile) == 0
	emptyAppRoleAuth := emptySecretId || emptyRoleId
	if emptyAppRoleAuth && emptyVaultToken {
		errs = append(errs, fmt.Errorf("neither '%s' nor AppRole auth info provided", FLAG_VAULT_TOKEN))
	}

	if !emptyAppRoleAuth && !emptyVaultToken {
		errs = append(errs, fmt.Errorf("both '%s' and AppRole auth info provided, don't know what to pick", FLAG_VAULT_TOKEN))
	}

	if len(c.VaultSecretId) > 0 && len(c.VaultSecretIdFile) > 0 {
		errs = append(errs, fmt.Errorf("both '%s' and '%s' auth info provided, don't know what to pick", FLAG_VAULT_SECRET_ID, FLAG_VAULT_SECRET_ID_FILE))
	}

	if len(c.VaultAddress) == 0 {
		errs = append(errs, fmt.Errorf("empty '%s' provided", FLAG_VAULT_ADDRESS))
	}

	if len(c.VaultMountApprole) == 0 {
		errs = append(errs, fmt.Errorf("empty '%s' provided", FLAG_VAULT_MOUNT_APPROLE))
	}

	if len(c.VaultMountSsh) == 0 {
		errs = append(errs, fmt.Errorf("empty '%s' provided", FLAG_VAULT_MOUNT_SSH))
	}

	if len(c.VaultSshBackend) == 0 {
		errs = append(errs, fmt.Errorf("empty '%s' provided", FLAG_VAULT_SSH_BACKEND_ROLE))
	}

	if c.CertificateLifetimeThresholdPercentage < 5 || c.CertificateLifetimeThresholdPercentage > 90 {
		errs = append(errs, fmt.Errorf("'%s' must be [5, 90]", FLAG_LIFETIME_THRESHOLD_PERCENTAGE))
	}

	return errs
}

func (c *Config) PrintConfig() {
	log.Info().Msg("Printing config values")
	log.Info().Msgf("%s=%s", FLAG_VAULT_ADDRESS, c.VaultAddress)
	if len(c.VaultToken) > 0 {
		log.Info().Msgf("%s=%s", FLAG_VAULT_TOKEN, c.VaultToken)
	}
	if len(c.VaultRoleId) > 0 {
		log.Info().Msgf("%s=*** (sensitive output)", FLAG_VAULT_ROLE_ID)
	}
	if len(c.VaultSecretId) > 0 {
		log.Info().Msgf("%s=*** (sensitive output)", FLAG_VAULT_SECRET_ID)
	}
	if len(c.VaultSecretIdFile) > 0 {
		log.Info().Msgf("%s=%s", FLAG_VAULT_SECRET_ID_FILE, c.VaultSecretIdFile)
	}
	log.Info().Msgf("%s=%s", FLAG_VAULT_MOUNT_SSH, c.VaultMountSsh)
	log.Info().Msgf("%s=%s", FLAG_VAULT_MOUNT_APPROLE, c.VaultMountApprole)
	log.Info().Msgf("%s=%s", FLAG_VAULT_SSH_BACKEND_ROLE, c.VaultSshBackend)
	log.Info().Msgf("%s=%s", FLAG_PUBKEY_FILE, c.PublicKeyFile)
	log.Info().Msgf("%s=%s", FLAG_SIGNED_KEY_FILE, c.SignedKeyFile)
	log.Info().Msgf("%s=%s", FLAG_METRICS_FILE, c.MetricsFile)
	log.Info().Msgf("%s=%t", FLAG_FORCE_NEW_SIGNATURE, c.ForceNewSignature)
	log.Info().Msgf("%s=%f", FLAG_LIFETIME_THRESHOLD_PERCENTAGE, c.CertificateLifetimeThresholdPercentage)
}

func NewConfigFromViper() Config {
	conf := Config{}

	conf.VaultAddress = viperOrEnv(FLAG_VAULT_ADDRESS, "VAULT_ADDR")
	conf.VaultToken = viper.GetViper().GetString(FLAG_VAULT_TOKEN)
	conf.VaultRoleId = viper.GetViper().GetString(FLAG_VAULT_ROLE_ID)
	conf.VaultSecretId = viper.GetViper().GetString(FLAG_VAULT_SECRET_ID)
	conf.VaultSecretIdFile = getExpandedFile(viper.GetViper().GetString(FLAG_VAULT_SECRET_ID_FILE))
	conf.VaultMountApprole = viper.GetViper().GetString(FLAG_VAULT_MOUNT_APPROLE)
	conf.VaultMountSsh = viper.GetViper().GetString(FLAG_VAULT_MOUNT_SSH)
	conf.VaultSshBackend = viper.GetViper().GetString(FLAG_VAULT_SSH_BACKEND_ROLE)
	conf.PublicKeyFile = getExpandedFile(viper.GetViper().GetString(FLAG_PUBKEY_FILE))
	conf.SignedKeyFile = getExpandedFile(viper.GetViper().GetString(FLAG_SIGNED_KEY_FILE))
	conf.ForceNewSignature = viper.GetViper().GetBool(FLAG_FORCE_NEW_SIGNATURE)
	conf.CertificateLifetimeThresholdPercentage = float32(viper.GetViper().GetFloat64(FLAG_LIFETIME_THRESHOLD_PERCENTAGE))
	conf.MetricsFile = viper.GetViper().GetString(FLAG_METRICS_FILE)

	return conf
}

func viperOrEnv(viperKey, envKey string) string {
	val := viper.GetViper().GetString(viperKey)
	if len(val) == 0 {
		return os.Getenv(envKey)
	}
	return val
}
