## Configuration
Configuration is supported via CLI arguments, ENV variables and through yaml-encoded config files.

### Configuration Options

| Name                              | Description                                                                                                                    | Default         |
|-----------------------------------|--------------------------------------------------------------------------------------------------------------------------------|-----------------|
| vault-address                     | Vault instance to connect to. If not specified, falls back to env var VAULT_ADDR.                                              |                 |
| vault-auth-token                  | Vault token to use for authentication. Can not be used in conjunction with AppRole login data                                  |                 |
| vault-auth-approle-role-id        | Vault role_id to use for AppRole login. Can not be used in conjuction with Vault token flag.                                   |                 |
| vault-auth-approle-secret-id      | Vault secret_id to use for AppRole login. Can not be used in conjuction with Vault token flag.                                 |                 |
| vault-auth-approle-secret-id-file | Flat file to read Vault secret_id from. Can not be used in conjuction with Vault token flag.                                   |                 |
| vault-auth-approle-mount          | Path where the AppRole auth method is mounted.                                                                                 | "approle"       |
| vault-auth-implicit               | Flat file to read Vault secret_id from. Can not be used in conjuction with Vault token flag.                                   |                 |
| vault-ssh-mount                   | Path where the SSH secret engine is mounted.                                                                                   | "ssh"           |
| vault-ssh-role                    | The name of the SSH role to use.                                                                                               | "host_key_sign" |
| force-new-signature               | Force a new signature, no matter if a previously existing signature is still valid.                                            | "host_key_sign" |
| renew-threshold-percent           | Renew a certificate when its lifetime is less than x percent.                                                                  | "host_key_sign" |
| pub-key-file                      | The path to the file containing the public key to be signed.                                                                   | "host_key_sign" |
| signed-key-file                   | The path of the file to write the certificate to.                                                                              | "host_key_sign" |
| ca-file                           | The path of the file to write the CA certificate to.                                                                           | "host_key_sign" |
| metrics-file                      | The path of the file to dump prometheus metrics to. This is usually a file under node_exporter's textfile collector directory. | "host_key_sign" |
| debug                             | Print debug statements.                                                                                                        | "host_key_sign" |

## Example

### Sign a user key
```bash
vault-ssh-cli sign-user-key \
     -a https://my-vault:8200 \
     --vault-ssh-role=user \
     --vault-auth-implicit=true \
     --pub-key-file=/etc/ssh/ssh_host_ed25519_key.pub
```

### Configuration via ENV variables
All configuration variables must be prefixed with `SSH_KEY_SIGNER`.

### Configuration Files
By default, config files named `config.yaml` are sought for at locations `$HOME/.config/vault-ssh-cli` and `/etc/vault-ssh-cli`.
