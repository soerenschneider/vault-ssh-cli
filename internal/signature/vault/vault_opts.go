package vault

import "errors"

func SshMountPath(path string) VaultOpts {
	return func(v *SignatureClient) error {
		if len(path) == 0 {
			return errors.New("empty path provided")
		}

		v.sshMountPath = path
		return nil
	}
}

func VaultRole(role string) VaultOpts {
	return func(v *SignatureClient) error {
		if len(role) == 0 {
			return errors.New("empty role provided")
		}

		v.role = role
		return nil
	}
}
