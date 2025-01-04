package auth

import (
	"context"
	"os"

	"github.com/hashicorp/vault/api"
)

type NoAuth struct {
}

func (t *NoAuth) Login(ctx context.Context, client *api.Client) (*api.Secret, error) {
	return &api.Secret{
		Auth: &api.SecretAuth{
			ClientToken: os.Getenv("VAULT_TOKEN"),
		},
	}, nil
}
