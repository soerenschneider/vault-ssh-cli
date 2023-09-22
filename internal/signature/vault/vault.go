package vault

import (
	"errors"
	"fmt"
	"io"

	"github.com/hashicorp/vault/api"
	log "github.com/rs/zerolog/log"
	"go.uber.org/multierr"
)

const defaultSshMountPath = "ssh"

type AuthMethod interface {
	Authenticate() (string, error)
}

type SignatureClient struct {
	client       *api.Client
	auth         AuthMethod
	role         string
	sshMountPath string
}

type VaultOpts func(client *SignatureClient) error

func NewVaultSigner(client *api.Client, auth AuthMethod, opts ...VaultOpts) (*SignatureClient, error) {
	if client == nil {
		return nil, errors.New("nil client passed")
	}

	if auth == nil {
		return nil, errors.New("nil auth method passed")
	}

	vault := &SignatureClient{
		client:       client,
		auth:         auth,
		sshMountPath: defaultSshMountPath,
	}

	var errs error
	for _, opt := range opts {
		if err := opt(vault); err != nil {
			errs = multierr.Append(errs, err)
		}
	}

	return vault, errs
}

func (c *SignatureClient) ReadCaCert() (string, error) {
	path := fmt.Sprintf("%s/public_key", c.sshMountPath)
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
	path := fmt.Sprintf("%s/sign/%s", c.sshMountPath, role)
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
