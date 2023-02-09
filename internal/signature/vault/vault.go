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
	role    string
	pathSsh string
}

func NewVaultSigner(client *api.Client, auth AuthMethod, vaultMount, backendName string) (*SignatureClient, error) {
	if client == nil {
		return nil, errors.New("nil client passed")
	}

	if auth == nil {
		return nil, errors.New("nil auth passed")
	}

	if len(vaultMount) == 0 {
		return nil, errors.New("empty vault mount passed")
	}

	return &SignatureClient{
		client:  client,
		auth:    auth,
		pathSsh: vaultMount,
		role:    backendName,
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

func (c *SignatureClient) signPublicKey(role string, publicKeyData, certType string) (string, error) {
	path := fmt.Sprintf("%s/sign/%s", c.pathSsh, role)
	data := map[string]interface{}{
		"public_key": publicKeyData,
		"cert_type":  certType,
	}
	secret, err := c.client.Logical().Write(path, data)
	if err != nil {
		return "", fmt.Errorf("could not sign ssh public key: %v", err)
	}

	signedData := fmt.Sprintf("%s", secret.Data["signed_key"])
	return signedData, nil
}

func (c *SignatureClient) SignHostKey(publicKeyData string) (string, error) {
	log.Info().Msg("Trying to authenticate against vault")
	token, err := c.auth.Authenticate()
	if err != nil {
		return "", fmt.Errorf("could not authenticate: %v", err)
	}

	c.client.SetToken(token)
	log.Info().Msgf("Signing public key using role '%s'", c.role)
	return c.signPublicKey(c.role, publicKeyData, "host")
}

func (c *SignatureClient) SignUserKey(publicKeyData string) (string, error) {
	log.Info().Msg("Trying to authenticate against vault")
	token, err := c.auth.Authenticate()
	if err != nil {
		return "", fmt.Errorf("could not authenticate: %v", err)
	}

	c.client.SetToken(token)
	log.Info().Msgf("Signing public key using role '%s'", c.role)
	return c.signPublicKey(c.role, publicKeyData, "user")
}
