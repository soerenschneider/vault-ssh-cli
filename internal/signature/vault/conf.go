package vault

import (
	"github.com/hashicorp/vault/api"
	"github.com/soerenschneider/vault-ssh-cli/internal/config"
)

func DeriveVaultConfig(conf *config.Config) *api.Config {
	vaultConfig := api.DefaultConfig()
	vaultConfig.MaxRetries = 5
	vaultConfig.Address = conf.VaultAddress
	return vaultConfig
}
