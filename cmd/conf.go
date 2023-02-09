package main

import (
	"errors"
	"fmt"
	log "github.com/rs/zerolog/log"
	"os"
	"reflect"
	"strings"
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
	VaultSshRole      string `mapstructure:"vault-ssh-role"`

	ForceNewSignature                      bool    `mapstructure:"force-signature"`
	CertificateLifetimeThresholdPercentage float32 `mapstructure:"renew-threshold-percent"`

	PublicKeyFile string `mapstructure:"pub-key-file"`
	SignedKeyFile string `mapstructure:"signed-key-file"`

	CaFile string `mapstructure:"ca-file"`

	MetricsFile string `mapstructure:"metrics-file"`
}

func (c *Config) PostValidation() {
	if len(c.SignedKeyFile) == 0 && len(c.PublicKeyFile) > 0 {
		auto := strings.Replace(c.PublicKeyFile, ".pub", "", 1)
		auto = getExpandedFile(fmt.Sprintf("%s-cert.pub", auto))
		log.Info().Msgf("Automatically derived value for '%s' (%s) from supplied '%s' (%s)", FLAG_SIGNED_KEY_FILE, auto, FLAG_PUBKEY_FILE, c.PublicKeyFile)
		c.SignedKeyFile = auto
	}
}

func (c *Config) ValidateCommon() []error {
	errs := make([]error, 0)

	if len(c.VaultAddress) == 0 {
		errs = append(errs, fmt.Errorf("empty '%s' provided", FLAG_VAULT_ADDRESS))
	}

	if len(c.VaultMountSsh) == 0 {
		errs = append(errs, fmt.Errorf("empty '%s' provided", FLAG_VAULT_SSH_MOUNT))
	}

	return errs
}

// ValidateSignCommand is called for commands that actually sign a public key.
func (c *Config) ValidateSignCommand() []error {
	errs := make([]error, 0)

	errs = append(errs, c.ValidateCommon()...)

	if len(c.PublicKeyFile) == 0 {
		errs = append(errs, fmt.Errorf("empty '%s' provided", FLAG_PUBKEY_FILE))
	} else {
		_, err := os.Stat(c.PublicKeyFile)
		if err != nil {
			errs = append(errs, fmt.Errorf("couldn't access pub-key-file at '%s'", c.PublicKeyFile))
		}
	}

	if len(c.VaultSshRole) == 0 {
		errs = append(errs, fmt.Errorf("empty '%s' provided", FLAG_VAULT_SSH_ROLE))
	}

	emptyVaultToken := len(c.VaultToken) == 0
	emptyRoleId := len(c.VaultRoleId) == 0
	emptySecretId := len(c.VaultSecretId) == 0 && len(c.VaultSecretIdFile) == 0
	emptyAppRoleAuth := emptySecretId || emptyRoleId

	numAuthMethodsProvided := 0
	if !emptyVaultToken {
		numAuthMethodsProvided += 1
	}
	if !emptyAppRoleAuth {
		numAuthMethodsProvided += 1
	}
	if c.VaultAuthImplicit {
		numAuthMethodsProvided += 1
	}

	if numAuthMethodsProvided == 0 {
		errs = append(errs, errors.New("no vault auth info provided. supply either token, AppRole or k8s auth info"))
	} else if numAuthMethodsProvided > 1 {
		errs = append(errs, fmt.Errorf("must provide only a single vault auth method, %d were provided", numAuthMethodsProvided))
	}

	if len(c.VaultSecretId) > 0 && len(c.VaultSecretIdFile) > 0 {
		errs = append(errs, fmt.Errorf("both '%s' and '%s' auth info provided, don't know what to pick", FLAG_VAULT_AUTH_APPROLE_SECRET_ID, FLAG_VAULT_AUTH_APPROLE_SECRET_ID_FILE))
	}

	if !emptyAppRoleAuth && len(c.VaultMountApprole) == 0 {
		errs = append(errs, fmt.Errorf("empty '%s' provided", FLAG_VAULT_AUTH_APPROLE_MOUNT_DEFAULT))
	}

	if c.CertificateLifetimeThresholdPercentage < 5 || c.CertificateLifetimeThresholdPercentage > 90 {
		errs = append(errs, fmt.Errorf("'%s' must be [5, 90]", FLAG_RENEW_THRESHOLD_PERCENTAGE))
	}

	return errs
}

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
