package config

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	log "github.com/rs/zerolog/log"
	"github.com/soerenschneider/vault-ssh-cli/pkg"
)

var (
	sensitiveVars = map[string]struct{}{
		FLAG_VAULT_AUTH_APPROLE_ROLE_ID:   {},
		FLAG_VAULT_AUTH_APPROLE_SECRET_ID: {},
		FLAG_VAULT_AUTH_TOKEN:             {},
	}

	validate = validator.New()
)

type Config struct {
	ForceNewSignature                      bool    `mapstructure:"force-new-signature"`
	CertificateLifetimeThresholdPercentage float32 `mapstructure:"renew-threshold-percent" validate:"lte=80,gte=20"`

	CaFile        string `mapstructure:"ca-file" validate:"omitempty,filepath"`
	PublicKeyFile string `mapstructure:"pub-key-file" validate:"omitempty,file"`
	SignedKeyFile string `mapstructure:"signed-key-file" validate:"omitempty,filepath"`

	Extensions map[string]string `mapstructure:"extensions" validate:"omitempty"`
	Principals []string          `mapstructure:"principals" validate:"omitempty"`

	MetricsFile string `mapstructure:"metrics-file" validate:"omitempty,filepath"`
	Debug       bool   `mapstructure:"debug"`

	Ttl int `mapstructure:"ttl" validate:"gte=0"`

	VaultAddress      string `mapstructure:"vault-address"`
	VaultToken        string `mapstructure:"vault-auth-token"`
	VaultRoleId       string `mapstructure:"vault-auth-approle-role-id"`
	VaultSecretId     string `mapstructure:"vault-auth-approle-secret-id"`
	VaultSecretIdFile string `mapstructure:"vault-auth-approle-secret-id-file" validate:"omitempty,file"`
	VaultMountApprole string `mapstructure:"vault-auth-approle-mount"`
	VaultMountSsh     string `mapstructure:"vault-ssh-mount" validate:"required"`
	VaultSshRole      string `mapstructure:"vault-ssh-role" validate:"required"`
}

func (c *Config) ExpandPaths() {
	if len(c.PublicKeyFile) > 0 {
		c.PublicKeyFile = pkg.GetExpandedFile(c.PublicKeyFile)
	}

	if len(c.CaFile) > 0 {
		c.CaFile = pkg.GetExpandedFile(c.CaFile)
	}

	if len(c.SignedKeyFile) == 0 && len(c.PublicKeyFile) > 0 {
		auto := strings.Replace(c.PublicKeyFile, ".pub", "", 1)
		auto = pkg.GetExpandedFile(fmt.Sprintf("%s-cert.pub", auto))
		log.Info().Msgf("Automatically derived value for '%s' (%s) from supplied '%s' (%s)", FLAG_SIGNED_KEY_FILE, auto, FLAG_PUBKEY_FILE, c.PublicKeyFile)
		c.SignedKeyFile = auto
	}
}

func Validate(s any) error {
	return validate.Struct(s)
}

func Print(c any) {
	log.Debug().Msg("Active config values:")
	val := reflect.ValueOf(c).Elem()
	for i := 0; i < val.NumField(); i++ {
		if !val.Field(i).IsZero() {
			fieldName := val.Type().Field(i).Tag.Get("mapstructure")
			_, isSensitive := sensitiveVars[fieldName]
			if isSensitive {
				log.Debug().Msgf("%s=*** (redacted)", fieldName)
			} else {
				log.Debug().Msgf("%s=%v", fieldName, val.Field(i))
			}
		}
	}
}
