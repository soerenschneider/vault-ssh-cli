package auth

import (
	"context"

	"github.com/hashicorp/vault/api"
)

type NoAuth struct {
}

func (t *NoAuth) Login(ctx context.Context, client *api.Client) (*api.Secret, error) {
	return &api.Secret{}, nil
}
