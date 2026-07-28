---
subcategory: ""
page_title: "ldap_ou Data Source - terraform-provider-ldap"
description: |-
  ldap_ou is a data source for reading an LDAP OU.
---

# ldap_ou (Data Source)

`ldap_ou` is a data source for reading an LDAP OU.

## Example Usage

```hcl
data "ldap_ou" "ou" {
  ou    = "OU=MyCompany,DC=domain,DC=tld"
  name  = "MyOU"
  scope = 2
}
```

## Argument Reference

* `ou` - (Required) OU where the LDAP OU will be searched.
* `name` - (Required) LDAP OU name.
* `scope` - (Optional) LDAP search scope: 0 = BaseObject, 1 = SingleLevel, 2 = WholeSubtree. Defaults to `0`.

## Attribute Reference

* `id` - The full DN of the LDAP OU.
* `description` - Description attribute of the LDAP OU.
* `managed_by` - The DN of the object that manages this OU (managedBy LDAP attribute).
