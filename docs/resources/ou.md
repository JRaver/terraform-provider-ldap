---
subcategory: ""
page_title: "ldap_ou Resource - terraform-provider-ldap"
description: |-
  ldap_ou is a resource for managing an LDAP OU.
---

# ldap_ou

`ldap_ou` is a resource for managing an LDAP OU.

## Example Usage

```hcl
resource "ldap_ou" "ou" {
  name        = "MyOU"
  ou          = "OU=MyCompany,DC=domain,DC=tld"
  description = "My OU description"
}
```

## Argument Reference

* `ou` - (Required) OU where the LDAP OU will be created. Changes force recreation.
* `name` - (Required) LDAP OU name. Changes force recreation.
* `description` - (Optional) Description attribute for the LDAP OU.
* `managed_by` - (Optional) The DN of the object that manages this OU (managedBy LDAP attribute).

## Attribute Reference

* `id` - The full DN of the LDAP OU.

## Import

LDAP OU can be imported using the full LDAP DN (id), e.g.

```
$ terraform import ldap_ou.example OU=Myou,OU=MyCompany,DC=domain,DC=tld
```