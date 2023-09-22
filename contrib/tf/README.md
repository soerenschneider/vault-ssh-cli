## Requirements

No requirements.

## Providers

| Name | Version |
|------|---------|
| <a name="provider_vault"></a> [vault](#provider\_vault) | n/a |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [vault_mount.ssh](https://registry.terraform.io/providers/hashicorp/vault/latest/docs/resources/mount) | resource |
| [vault_policy.host_key_sign](https://registry.terraform.io/providers/hashicorp/vault/latest/docs/resources/policy) | resource |
| [vault_ssh_secret_backend_ca.ca](https://registry.terraform.io/providers/hashicorp/vault/latest/docs/resources/ssh_secret_backend_ca) | resource |
| [vault_ssh_secret_backend_role.host_key_sign](https://registry.terraform.io/providers/hashicorp/vault/latest/docs/resources/ssh_secret_backend_role) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_algorithm"></a> [algorithm](#input\_algorithm) | n/a | `string` | `"rsa-sha2-512"` | no |
| <a name="input_allow_subdomains"></a> [allow\_subdomains](#input\_allow\_subdomains) | n/a | `bool` | `false` | no |
| <a name="input_allowed_domains"></a> [allowed\_domains](#input\_allowed\_domains) | n/a | `list(string)` | n/a | yes |
| <a name="input_max_ttl"></a> [max\_ttl](#input\_max\_ttl) | n/a | `number` | `604800` | no |
| <a name="input_ssh_path"></a> [ssh\_path](#input\_ssh\_path) | n/a | `string` | `"ssh"` | no |
| <a name="input_ttl"></a> [ttl](#input\_ttl) | n/a | `number` | `172800` | no |

## Outputs

No outputs.
