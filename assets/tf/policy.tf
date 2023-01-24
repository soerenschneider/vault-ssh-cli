resource "vault_policy" "host_key_sign" {
  name = "ssh_host_key_sign"

  policy = <<EOT
path "${var.ssh_path}/sign/${vault_ssh_secret_backend_role.host_key_sign.name}" {
  capabilities = ["update"]
}
EOT
}
