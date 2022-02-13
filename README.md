
# ssh-key-signer

Sign SSH (host) keys using Hashicorp Vault

## Features

✅ Sign SSH host public keys using Hashicorp Vault

✅ Follow best-practices and automate public key signing after passing a configurable certificate expiry thresholds

✅ Collect metrics for observability

❌ Sign SSH client public keys

## Running it

Configuration is supported via CLI arguments, ENV variables and through yaml-encoded config files. 

## Configuration

### Configuration Options

| Name                 | Description                                                                                   | Default   |
|----------------------|-----------------------------------------------------------------------------------------------|-----------|
| vault-address        | Vault instance to connect to. If not specified, falls back to env var VAULT_ADDR.             |           |
| vault-token          | Vault token to use for authentication. Can not be used in conjunction with AppRole login data |           |
| vault-role-id        | Vault role_id to use for AppRole login. Can not be used in conjuction with Vault token flag.  |           |
| vault-secret-id      | Vault secret_id to use for AppRole login. Can not be used in conjuction with Vault token flag.|           |
| vault-secret-id-file | Flat file to read Vault secret_id from. Can not be used in conjuction with Vault token flag.  |           |
| vault-mount-ssh      | Path where the SSH secret engine is mounted.                                                  | "ssh"     |
| vault-mount-approle  | Path where the AppRole auth method is mounted.                                                | "approle" |
| vault-ssh-role-name  | The name of the SSH role auth backend.                                                        | "host_key_sign" |

| Name                 | Description                                                                                   | Default   |
|----------------------|-----------------------------------------------------------------------------------------------|-----------|
| signed-key-file      | File to write the signed key to                                                               |           |
| pub-key-file         | SSH Public Host Key to sign                                                                   |           |
| metrics-file         | Dump metrics to given file to be picked up by prometheus node_exporter                        | /var/lib/node_exporter/ssh_key_sign.prom |
| lifetime-threshold-percent | If there's already a signed certificate at `signed-key-file`, only sign public key again if its lifetime period is less than the given threshold.                                                              | 33        |
| force-new-signature  | Sign public key regardless of it's validity period                                            | false     |
| config-file          | File to read configuration from

### Configuration via ENV variables
All configuration variables must be prefixed with `SSH_KEY_SIGNER`.

### Configuration Files
By default, config files named `config.yaml` are sought for at locations `$HOME/.config/ssh-key-signer` and `/etc/ssh-key-signer`.

## Automating Key Signatures
`ssh-key-signer` is suited to be scheduled continuously by an external actor such as systemd or (Kubernetes) cron jobs and only renew a certificate after its expiration period has passed a certain threshold.

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

## Building it from source

```sh
$ make build
```