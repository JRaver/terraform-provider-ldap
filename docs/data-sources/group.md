---
subcategory: ""
page_title: "ldap_group Data Source - terraform-provider-ldap"
description: |-
  ldap_group is a data source for reading an LDAP group.
---

# ldap_group (Data Source)

`ldap_group` is a data source for reading an LDAP group.

## Example Usage

```hcl
data "ldap_group" "group" {
  ou          = "OU=MyOU,DC=domain,DC=tld"
  name        = "MyGroup"
}
```

## Argument Reference

* `ou` - (Required) OU where the LDAP group will be searched.
* `name` - (Required) LDAP group name.
* `scope` - (Optional) LDAP search scope: 0 = BaseObject, 1 = SingleLevel, 2 = WholeSubtree. Defaults to `0`.

## Attribute Reference

* `id` - The full DN of the LDAP group.
* `description` - Description attribute of the LDAP group.
* `group_type` - The groupType attribute value of the LDAP group.
* `members` - List of full DNs of the group members.
* `members_names` - Display names (CN values) of the group members.
* `managed_by` - The DN of the object that manages this group (managedBy LDAP attribute).
* `display_name` - The displayName attribute of the group.