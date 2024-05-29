package vault

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/hashicorp/vault/api"
	log "github.com/rs/zerolog/log"
	"github.com/soerenschneider/vault-ssh-cli/internal/signature"
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

func (c *SignatureClient) SignHostKey(req signature.SignHostKeyRequest) (string, error) {
	log.Info().Msg("Trying to authenticate against vault")
	token, err := c.auth.Authenticate()
	if err != nil {
		return "", fmt.Errorf("could not authenticate: %v", err)
	}

	c.client.SetToken(token)
	log.Info().Msgf("Signing public key using role '%s'", c.role)

	data := convertHostKeyRequest(req)
	path := fmt.Sprintf("%s/sign/%s", c.sshMountPath, c.role)
	secret, err := c.client.Logical().Write(path, data)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s", secret.Data["signed_key"]), nil
}

func convertHostKeyRequest(req signature.SignHostKeyRequest) map[string]any {
	data := map[string]interface{}{
		"public_key": req.PublicKey,
		"cert_type":  "host",
	}

	if req.Ttl > 0 {
		data["ttl"] = req.Ttl
	}

	if len(req.Principals) > 0 {
		data["valid_principals"] = strings.Join(req.Principals, ",")
	}

	return data
}

func (c *SignatureClient) SignUserKey(req signature.SignUserKeyRequest) (string, error) {
	log.Info().Msg("Trying to authenticate against vault")
	token, err := c.auth.Authenticate()
	if err != nil {
		return "", fmt.Errorf("could not authenticate: %v", err)
	}

	c.client.SetToken(token)
	log.Info().Msgf("Signing public key using role '%s'", c.role)

	data := convertUserKeyRequest(req)
	path := fmt.Sprintf("%s/sign/%s", c.sshMountPath, c.role)
	secret, err := c.client.Logical().Write(path, data)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s", secret.Data["signed_key"]), nil
}

func convertUserKeyRequest(req signature.SignUserKeyRequest) map[string]any {
	data := map[string]interface{}{
		"public_key": req.PublicKey,
		"cert_type":  "user",
	}

	if req.Ttl > 0 {
		data["ttl"] = req.Ttl
	}

	if len(req.Principals) > 0 {
		data["valid_principals"] = strings.Join(req.Principals, ",")
	}

	return data
}
