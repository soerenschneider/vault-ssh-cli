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
