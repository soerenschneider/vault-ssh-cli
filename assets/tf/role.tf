resource "vault_ssh_secret_backend_role" "host_key_sign" {
  name                    = "host_key_sign"
  backend                 = vault_mount.ssh.path
  key_type                = "ca"
  allow_user_certificates = false
  allow_host_certificates = true
  ttl                     = var.ttl
  max_ttl                 = var.max_ttl
  allowed_domains         = join(",", var.allowed_domains)
  allow_subdomains        = coalesce(var.allow_subdomains, false)
  algorithm_signer        = var.algorithm
}
