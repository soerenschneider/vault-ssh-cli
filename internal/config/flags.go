package config

const (
	FLAG_VAULT_ADDRESS                     = "vault-address"
	FLAG_VAULT_AUTH_TOKEN                  = "vault-auth-token" // #nosec: G101
	FLAG_VAULT_AUTH_KUBERNETES_ROLE        = "vault-auth-kubernetes-role"
	FLAG_VAULT_AUTH_KUBERNETES_MOUNT       = "vault-auth-kubernetes-mount"
	FLAG_VAULT_AUTH_APPROLE_ROLE_ID        = "vault-auth-approle-role-id"
	FLAG_VAULT_AUTH_APPROLE_SECRET_ID      = "vault-auth-approle-secret-id"      // #nosec: G101
	FLAG_VAULT_AUTH_APPROLE_SECRET_ID_FILE = "vault-auth-approle-secret-id-file" // #nosec: G101
	FLAG_VAULT_AUTH_APPROLE_MOUNT          = "vault-auth-approle-mount"          // #nosec: G101
	FLAG_VAULT_AUTH_APPROLE_MOUNT_DEFAULT  = "approle"
	FLAG_VAULT_SSH_MOUNT                   = "vault-ssh-mount"
	FLAG_VAULT_SSH_MOUNT_DEFAULT           = "ssh"
	FLAG_VAULT_SSH_ROLE                    = "vault-ssh-role"

	FLAG_TTL        = "ttl"
	FLAG_PRINCIPALS = "principals"

	FLAG_FORCE_NEW_SIGNATURE                = "force-new-signature"
	FLAG_FORCE_NEW_SIGNATURE_DEFAULT        = false
	FLAG_RETRIES                            = "retries"
	FLAG_RETRIES_DEFAULT                    = 15
	FLAG_RENEW_THRESHOLD_PERCENTAGE         = "renew-threshold-percent"
	FLAG_RENEW_THRESHOLD_PERCENTAGE_DEFAULT = 33.

	FLAG_CA_FILE         = "ca-file"
	FLAG_PUBKEY_FILE     = "pub-key-file"
	FLAG_SIGNED_KEY_FILE = "signed-key-file"
	FLAG_CONFIG_FILE     = "config-file"

	FLAG_METRICS_FILE         = "metrics-file"
	FLAG_METRICS_FILE_DEFAULT = "/var/lib/node_exporter/vault-ssh-cli.prom"
	FLAG_DEBUG                = "debug"
)
