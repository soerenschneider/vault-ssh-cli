# vault-ssh-cli
[![Go Report Card](https://goreportcard.com/badge/github.com/soerenschneider/vault-ssh-cli)](https://goreportcard.com/report/github.com/soerenschneider/vault-ssh-cli)
![release-workflow](https://github.com/soerenschneider/vault-ssh-cli/actions/workflows/release-container.yaml/badge.svg)
![golangci-lint-workflow](https://github.com/soerenschneider/vault-ssh-cli/actions/workflows/golangci-lint.yaml/badge.svg)
![test-workflow](https://github.com/soerenschneider/vault-ssh-cli/actions/workflows/test.yaml/badge.svg)

Automate signing SSH host- and client certificates for a more secure and scalable infrastructure. 

## Features

🏭 Sign SSH host public keys

👨‍💻 Sign SSH user public keys

🔗 Read CA from a given Vault ssh mount

🛂 Authenticate against Vault using AppRole, (explicit) token or implicit__ auth

💻 Both your workstation's CLI and your servers up in the cloud are 1st class citizens

⏰ Automatically renews certificates based on its lifetime

🔭 Provides metrics to increase observability for robust automation

## Why would I need this?

SSH client certificates make sense 
- to avoid the chore of synchronizing `authorized_keys` files across servers
- to avoid theft of public key pairs

SSH host certificates help prevent MitM attacks for clients that have not established trust yet for a server

Both client and host certificates allow for efficient scaling regarding the number of clients and servers.

vault-ssh-cli, leveraging its automation and observability capabilities, allows using SSH certificates while obeying security best practices such as short-lived certificates and timely re-generation.

## Example
![asciicinema demo](docs/asciicinema.svg)

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
