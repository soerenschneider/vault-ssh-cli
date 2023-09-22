package auth

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/rs/zerolog/log"
)

const (
	defaultEnvVar    = "VAULT_TOKEN"
	defaultTokenFile = ".vault-token"
)

type TokenImplicitAuth struct {
	envVar    string
	tokenFile string
}

func NewTokenImplicitAuth() *TokenImplicitAuth {
	return &TokenImplicitAuth{
		envVar:    defaultEnvVar,
		tokenFile: defaultTokenFile,
	}
}

func (t *TokenImplicitAuth) Authenticate() (string, error) {
	token := os.Getenv(t.envVar)
	if len(token) > 0 {
		log.Info().Msgf("Using vault token from env var %s", t.envVar)
		return token, nil
	}

	dirname, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("can't get user home dir: %v", err)
	}

	tokenPath := path.Join(dirname, t.tokenFile)
	if _, err := os.Stat(tokenPath); errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("file '%s' to read vault token from does not exist", t.tokenFile)
	}

	read, err := os.ReadFile(tokenPath)
	if err != nil {
		return "", fmt.Errorf("error reading file '%s': %v", defaultTokenFile, err)
	}

	log.Info().Msgf("Using vault token from file '%s'", t.tokenFile)
	return string(read), nil
}
