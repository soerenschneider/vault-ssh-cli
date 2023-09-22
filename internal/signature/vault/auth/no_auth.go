package auth

import "errors"

type NoAuth struct {
}

func (t *NoAuth) Authenticate() (string, error) {
	return "", errors.New("auth not supported")
}
