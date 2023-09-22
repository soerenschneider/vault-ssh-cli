variable "ssh_path" {
  type    = string
  default = "ssh"
}

variable "ttl" {
  type    = number
  default = 3600 * 24 * 14
}

variable "max_ttl" {
  type    = number
  default = 3600 * 24 * 30 * 2
}

variable "allowed_domains" {
  type = list(string)
  default = ["example.com"]
}

variable "allow_subdomains" {
  type    = bool
  default = false
}

variable "algorithm" {
  type    = string
  default = "rsa-sha2-512"
}
