resource "vault_mount" "ssh" {
  type = "ssh"
  path = var.ssh_path
}
