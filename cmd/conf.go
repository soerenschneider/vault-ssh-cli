package main

import (
	"fmt"
	"reflect"

	log "github.com/rs/zerolog/log"
)

var sensitiveVars = map[string]struct{}{
	FLAG_VAULT_AUTH_APPROLE_ROLE_ID:   {},
	FLAG_VAULT_AUTH_APPROLE_SECRET_ID: {},
	FLAG_VAULT_AUTH_TOKEN:             {},
}

type Config struct {
	VaultAddress      string `mapstructure:"vault-address"`
	VaultToken        string `mapstructure:"vault-auth-token"`
	VaultRoleId       string `mapstructure:"vault-auth-approle-role-id"`
	VaultSecretId     string `mapstructure:"vault-auth-approle-secret-id"`
	VaultSecretIdFile string `mapstructure:"vault-auth-approle-secret-id-file"`
	VaultMountApprole string `mapstructure:"vault-auth-approle-mount"`
	VaultAuthImplicit bool   `mapstructure:"vault-auth-implicit"`
	VaultMountSsh     string `mapstructure:"vault-ssh-mount"`
	VaultSshBackend   string `mapstructure:"vault-ssh-backend"`

	ForceNewSignature                      bool    `mapstructure:"force-signature"`
	CertificateLifetimeThresholdPercentage float32 `mapstructure:"renew-threshold-percent"`

	PublicKeyFile string `mapstructure:"public-key-file"`
	SignedKeyFile string `mapstructure:"signed-key-file"`

	CaFile string `mapstructure:"ca-file"`

	MetricsFile string `mapstructure:"metrics-file"`
}

func (c *Config) Validate() []error {
	errs := make([]error, 0)

	if len(c.VaultAddress) == 0 {
		errs = append(errs, fmt.Errorf("empty '%s' provided", FLAG_VAULT_ADDRESS))
	}

	if len(c.VaultMountSsh) == 0 {
		errs = append(errs, fmt.Errorf("empty '%s' provided", FLAG_VAULT_SSH_MOUNT))
	}

	return errs
}

/*
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


	if len(c.VaultSshBackend) == 0 {
		errs = append(errs, fmt.Errorf("empty '%s' provided", FLAG_VAULT_SSH_BACKEND_ROLE))
	}

	emptyVaultToken := len(c.VaultToken) == 0
	emptyRoleId := len(c.VaultRoleId) == 0
	emptySecretId := len(c.VaultSecretId) == 0 && len(c.VaultSecretIdFile) == 0
	emptyAppRoleAuth := emptySecretId || emptyRoleId
	if emptyAppRoleAuth && emptyVaultToken {
		errs = append(errs, fmt.Errorf("neither '%s' nor AppRole auth info provided", FLAG_VAULT_AUTH_TOKEN))
	}

	if !emptyAppRoleAuth && !emptyVaultToken {
		errs = append(errs, fmt.Errorf("both '%s' and AppRole auth info provided, don't know what to pick", FLAG_VAULT_AUTH_TOKEN))
	}

	if len(c.VaultSecretId) > 0 && len(c.VaultSecretIdFile) > 0 {
		errs = append(errs, fmt.Errorf("both '%s' and '%s' auth info provided, don't know what to pick", FLAG_VAULT_AUTH_APPROLE_SECRET_ID, FLAG_VAULT_AUTH_APPROLE_SECRET_ID_FILE))
	}


	if len(c.VaultMountApprole) == 0 {
		errs = append(errs, fmt.Errorf("empty '%s' provided", FLAG_VAULT_MOUNT_APPROLE))
	}

	if c.CertificateLifetimeThresholdPercentage < 5 || c.CertificateLifetimeThresholdPercentage > 90 {
		errs = append(errs, fmt.Errorf("'%s' must be [5, 90]", FLAG_RENEW_THRESHOLD_PERCENTAGE))
	}

	return errs
}

*/

func (c *Config) Print() {
	log.Info().Msg("---")
	log.Info().Msg("Active config values:")
	val := reflect.ValueOf(c).Elem()
	for i := 0; i < val.NumField(); i++ {
		if !val.Field(i).IsZero() {
			fieldName := val.Type().Field(i).Tag.Get("mapstructure")
			_, isSensitive := sensitiveVars[fieldName]
			if isSensitive {
				log.Info().Msgf("%s=*** (redacted)", fieldName)
			} else {
				log.Info().Msgf("%s=%v", fieldName, val.Field(i))
			}
		}
	}
	log.Info().Msg("---")
}
