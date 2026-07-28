---
page_title: "Provider: LDAP"
description: |-
  The LDAP provider is used to interact with an Active Directory or LDAP server.
---

# LDAP Provider

The LDAP provider is used to interact with an Active Directory or LDAP server.

## Example Usage

```hcl
terraform {
  required_providers {
    ldap = {
      source = "Ouest-France/ldap"
    }
  }
}

provider "ldap" {
  host          = "ldap.mycompany.tld"
  port          = 389
  bind_user     = "ldap_user"
  bind_password = "ldap_password"
}
```

## Argument Reference

* `host` - (Required) LDAP host. Can also be set with the `LDAP_HOST` environment variable.
* `port` - (Required) LDAP port. Can also be set with the `LDAP_PORT` environment variable.
* `bind_user` - (Required) LDAP bind username. Can also be set with the `LDAP_USER` environment variable.
* `bind_password` - (Required) LDAP bind password. Can also be set with the `LDAP_PASSWORD` environment variable.
* `tls` - (Optional) Enable TLS encryption for LDAP (LDAPS). Defaults to `false`.
* `tls_ca_certificate` - (Optional) PEM-encoded TLS CA certificate to trust for the LDAPS connection.
* `tls_insecure` - (Optional) Skip server TLS certificate verification. Defaults to `false`.
