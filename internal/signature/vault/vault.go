package vault

import (
	"errors"
	"fmt"
	"io"

	"github.com/hashicorp/vault/api"
	log "github.com/rs/zerolog/log"
)

type AuthMethod interface {
	Authenticate() (string, error)
}

type SignatureClient struct {
	client  *api.Client
	auth    AuthMethod
	backend string
	pathSsh string
}

func NewVaultSigner(client *api.Client, auth AuthMethod, backendName string) (*SignatureClient, error) {
	if client == nil {
		return nil, errors.New("nil client passed")
	}

	if auth == nil {
		return nil, errors.New("nil auth passed")
	}

	return &SignatureClient{
		client: client,
		auth:   auth,
		// TODO: make  configurable
		pathSsh: "ssh",
		backend: backendName,
	}, nil
}

func (c *SignatureClient) ReadCaCert() (string, error) {
	path := fmt.Sprintf("%s/public_key", c.pathSsh)
	resp, err := c.client.Logical().ReadRaw(path)
	if err != nil {
		return "", fmt.Errorf("reading cert failed: %v", err)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("could not read body from response: %v", err)
	}

	return string(data), nil
}

func (c *SignatureClient) signPublicKey(backend string, publicKeyData string) (string, error) {
	path := fmt.Sprintf("%s/sign/%s", c.pathSsh, backend)
	data := map[string]interface{}{
		"public_key": publicKeyData,
		"cert_type":  "host",
	}
	secret, err := c.client.Logical().Write(path, data)
	if err != nil {
		return "", fmt.Errorf("could not sign ssh public key: %v", err)
	}

	signedData := fmt.Sprintf("%s", secret.Data["signed_key"])
	return signedData, nil
}

func (c *SignatureClient) SignPublicKey(publicKeyData string) (string, error) {
	log.Info().Msg("Trying to authenticate against vault")
	token, err := c.auth.Authenticate()
	if err != nil {
		return "", fmt.Errorf("could not authenticate: %v", err)
	}

	c.client.SetToken(token)
	log.Info().Msgf("Signing public key using backend %s", c.backend)
	return c.signPublicKey(c.backend, publicKeyData)
}
