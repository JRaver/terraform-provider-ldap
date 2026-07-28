# `terraform-provider-ldap`

A [Terraform](https://www.terraform.io) provider for managing Active Directory / LDAP groups and organizational units.

Fork of [Ouest-France/terraform-provider-ldap](https://github.com/Ouest-France/terraform-provider-ldap) with the following additions:
- `ldap_group_member` resource — manages a single group member without replacing the full member list

## Requirements

- Terraform >= 1.0
- Go >= 1.22 (for building from source)

## Installation

```hcl
terraform {
  required_providers {
    ldap = {
      source  = "JRaver/ldap"
      version = "~> 0.1"
    }
  }
}
```

## Quick Start

```hcl
provider "ldap" {
  host          = "ad.example.com"
  port          = 389
  bind_user     = "svc-terraform"
  bind_password = "password"
}

# Create a group (do not manage members here)
resource "ldap_group" "app" {
  ou          = "OU=Groups,DC=example,DC=com"
  name        = "MyApp-SSO"
  description = "MyApp SSO group"
}

# Add members individually — does not affect members added outside Terraform
resource "ldap_group_member" "alice" {
  group_dn  = ldap_group.app.id
  member_dn = "CN=Alice,OU=Users,DC=example,DC=com"
}

resource "ldap_group_member" "bob" {
  group_dn  = ldap_group.app.id
  member_dn = "CN=Bob,OU=Users,DC=example,DC=com"
}
```

## Resources

| Resource | Description |
|----------|-------------|
| `ldap_group` | Manages an LDAP group. When `members` is set, replaces all existing members on update. |
| `ldap_group_member` | Manages a single member of an LDAP group. Safe to use alongside other membership sources. |
| `ldap_ou` | Manages an LDAP Organizational Unit. |

## Data Sources

| Data Source | Description |
|-------------|-------------|
| `data.ldap_group` | Reads an existing LDAP group and its members. |
| `data.ldap_user` | Reads an existing LDAP user. |
| `data.ldap_ou` | Reads an existing LDAP OU. |

## Building from Source

```bash
git clone https://github.com/JRaver/terraform-provider-ldap
cd terraform-provider-ldap
make install   # builds and installs for current OS/arch
```

## Documentation

Full documentation: https://registry.terraform.io/providers/JRaver/ldap/latest/docs

## License

[MIT](LICENSE) — fork of [Ouest-France/terraform-provider-ldap](https://github.com/Ouest-France/terraform-provider-ldap)
