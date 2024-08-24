package vault

import (
	"github.com/hashicorp/vault/api"
	"github.com/soerenschneider/vault-ssh-cli/internal/config"
)

func FromConfig(conf *config.Config) *api.Config {
	vaultConfig := api.DefaultConfig()
	vaultConfig.MaxRetries = 5
	if len(conf.VaultAddress) > 0 {
		vaultConfig.Address = conf.VaultAddress
	}
	return vaultConfig
}
