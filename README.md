# vault-ssh-cli
[![Go Report Card](https://goreportcard.com/badge/github.com/soerenschneider/vault-ssh-cli)](https://goreportcard.com/report/github.com/soerenschneider/vault-ssh-cli)
![release-workflow](https://github.com/soerenschneider/vault-ssh-cli/actions/workflows/release-container.yaml/badge.svg)
![golangci-lint-workflow](https://github.com/soerenschneider/vault-ssh-cli/actions/workflows/golangci-lint.yaml/badge.svg)
![test-workflow](https://github.com/soerenschneider/vault-ssh-cli/actions/workflows/test.yaml/badge.svg)

Automate signing SSH host- and client certificates for a more secure and scalable infrastructure. 

## Features

🏭 Sign SSH host public keys to verify that a client securely connects to an (previously unknown) box

👨‍💻 Sign SSH user public keys to allow a client connecting via certificate without the need to synchronize `authorized_keys` on all boxes

🔗 Read a CA from a given path from Vault

🛂 Authenticate against Vault using AppRole, (explicit) token or implicit__ auth

💻 Runs effortlessly both on your workstation's CLI via command line flags or automated via systemd and config files on your server

⏰ Automatically renews certificates based on its lifetime

🔭 Provides metrics to increase observability for robust automation

## Example
![asciicinema demo](demo.svg)

## Installation

### Pre-compiled Binaries

Pre-compiled binaries can be found at the releases section. They are signed using a cryptographic signature made by [signify](https://man.openbsd.org/signify.1) using the following public key: 
```
untrusted comment: signify public key
RWSFxNuvQMx07H1IC6sUxJvlsdtfDlY39EdoHMG/ZpivtOmp8sJ3DMEg
```

To verify the cryptographic signature, run
```bash
$ signify -V -p /path/to/downloaded/pubkey -m checksum.sha256
$ sha256sum -c checksum.sha256
```

### Building it from source

```sh
$ go install github.com/soerenschneider/vault-ssh-cli@latest
```

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

## Automating Key Signatures
`vault-ssh-cli` is suited to be scheduled continuously by an external actor such as systemd or (Kubernetes) cron jobs and only renew a certificate after its expiration period has passed a certain threshold.

## Metrics

### Collecting the metrics

✅ Dumping metrics to disk to be picked up by Prometheus node_exporter

❌ Pushing metrics to Prometheus Pushgateway

### Available metrics

All metrics are exposed using the prefix `ssh_key_signer`

| Name                        | Type    | Description                                            |
|-----------------------------|---------|--------------------------------------------------------|
| success_bool                | Gauge   | Whether the tool ran successful                        |
| cert_expiry_seconds         | Gauge   | The date after the cert is not valid anymore           |
| cert_lifetime_seconds_total | Gauge   | The total number of seconds this certificate is valid  |
| cert_lifetime_percent       | Gauge   | The passed lifetime of the certificate in percent      | 
| run_timestamp_seconds       | Gauge   | The date after the cert is not valid anymore           |


## Configuring 3rd party Systems

### Vault Configuration
Vault needs to be configured with a SSH secret engine, see [this TF module](https://github.com/soerenschneider/tf-vault/tree/main/secret_ssh). 

### Configuring OpenSSH Server
https://man.openbsd.org/sshd_config#HostCertificate

### Configuring OpenSSH Client
https://www.vaultproject.io/docs/secrets/ssh/signed-ssh-certificates#client-side-host-verification
