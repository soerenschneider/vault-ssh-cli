package signature

import "errors"

func WithSshMountPath(path string) VaultOpts {
	return func(v *SignatureClient) error {
		if len(path) == 0 {
			return errors.New("empty path provided")
		}

		v.sshMountPath = path
		return nil
	}
}

func WithVaultRole(role string) VaultOpts {
	return func(v *SignatureClient) error {
		if len(role) == 0 {
			return errors.New("empty role provided")
		}

		v.role = role
		return nil
	}
}
