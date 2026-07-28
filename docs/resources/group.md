---
subcategory: ""
page_title: "ldap_group Resource - terraform-provider-ldap"
description: |-
  ldap_group is a resource for managing an LDAP group.
---

# ldap_group

`ldap_group` is a resource for managing an LDAP group.

## Example Usage

```hcl
resource "ldap_group" "group" {
  ou          = "OU=MyOU,DC=domain,DC=tld"
  name        = "MyGroup"
  members     = ["CN=MyUser,OU=MyOU,DC=domain,DC=tld"]
  description = "My group description"
}
```

## Argument Reference

* `ou` - (Required) OU where the LDAP group will be created. Changes force recreation.
* `name` - (Required) LDAP group name (CN). Changes force recreation.
* `description` - (Optional) Description attribute for the LDAP group.
* `members` - (Optional, Computed) List of full DNs of LDAP objects that are members of this group. When set, **replaces all existing members**. To manage members without full replacement, use `ldap_group_member` resources instead and omit this field.
* `group_type` - (Optional, Computed) The groupType attribute value (e.g. `-2147483646` for a global security group in Active Directory). Changes force recreation.
* `managed_by` - (Optional) The DN of the object that manages this group (managedBy LDAP attribute).
* `display_name` - (Optional) The displayName attribute of the group.

## Attribute Reference

* `id` - The full DN of the LDAP group.
* `members_names` - Display names (CN values) of the current group members. Computed automatically.

## Import

LDAP group can be imported using the full LDAP DN (id), e.g.

```
$ terraform import ldap_group.example CN=MyGroup,OU=MyOU,DC=domain,DC=tld
```