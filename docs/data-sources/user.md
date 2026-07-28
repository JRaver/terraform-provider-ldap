---
subcategory: ""
page_title: "ldap_user Data Source - terraform-provider-ldap"
description: |-
  ldap_user is a data source for retrieving an LDAP user.
---

# ldap_user (Data Source)

`ldap_user` is a data source for retrieving an LDAP user.

## Example Usage

```hcl
data "ldap_user" "user" {
  ou   = "OU=MyOU,DC=domain,DC=tld"
  name = "MyUser"
}
```

## Argument Reference

* `ou` - (Required) OU where the LDAP user will be searched.
* `name` - (Optional) The name of the LDAP user. At least one of `name`, `sam_account_name`, or `user_principal_name` must be set.
* `sam_account_name` - (Optional) The sAMAccountName of the LDAP user. At least one of `name`, `sam_account_name`, or `user_principal_name` must be set.
* `user_principal_name` - (Optional) The userPrincipalName (UPN) of the LDAP user. At least one of `name`, `sam_account_name`, or `user_principal_name` must be set.

## Attribute Reference

* `id` - The full DN of the LDAP user.
* `description` - Description attribute of the LDAP user.
* `mail` - Mail attribute of the LDAP user.
