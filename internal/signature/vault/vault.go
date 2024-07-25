package vault

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/cenkalti/backoff/v3"
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
		if err != nil {
			var respErr *api.ResponseError
			if errors.As(err, &respErr) && !shouldRetry(respErr.StatusCode) {
				return "", backoff.Permanent(err)
			}
			return "", err
		}

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
		var respErr *api.ResponseError
		if errors.As(err, &respErr) && !shouldRetry(respErr.StatusCode) {
			return "", backoff.Permanent(err)
		}
		return "", err
	}

	return fmt.Sprintf("%s", secret.Data["signed_key"]), nil
}

func convertHostKeyRequest(req signature.SignHostKeyRequest) map[string]any {
	data := map[string]interface{}{
		"public_key": req.PublicKey,
		"cert_type":  "host",
	}

	if len(req.Ttl) > 0 {
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
		var respErr *api.ResponseError
		if errors.As(err, &respErr) && !shouldRetry(respErr.StatusCode) {
			return "", backoff.Permanent(err)
		}
		return "", err
	}

	return fmt.Sprintf("%s", secret.Data["signed_key"]), nil
}

func convertUserKeyRequest(req signature.SignUserKeyRequest) map[string]any {
	data := map[string]interface{}{
		"public_key": req.PublicKey,
		"cert_type":  "user",
	}

	if len(req.Ttl) > 0 {
		data["ttl"] = req.Ttl
	}

	if len(req.Principals) > 0 {
		data["valid_principals"] = strings.Join(req.Principals, ",")
	}

	return data
}

func shouldRetry(statusCode int) bool {
	switch statusCode {
	case 400, // Bad Request
		401, // Unauthorized
		403, // Forbidden
		404, // Not Found
		405, // Method Not Allowed
		406, // Not Acceptable
		407, // Proxy Authentication Required
		409, // Conflict
		410, // Gone
		411, // Length Required
		412, // Precondition Failed
		413, // Payload Too Large
		414, // URI Too Long
		415, // Unsupported Media Type
		416, // Range Not Satisfiable
		417, // Expectation Failed
		418, // I'm a Teapot
		421, // Misdirected Request
		422, // Unprocessable Entity
		423, // Locked (WebDAV)
		424, // Failed Dependency (WebDAV)
		425, // Too Early
		426, // Upgrade Required
		428, // Precondition Required
		429, // Too Many Requests
		431, // Request Header Fields Too Large
		451: // Unavailable For Legal Reasons
		return false
	default:
		return true
	}
}
