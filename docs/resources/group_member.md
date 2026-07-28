---
subcategory: ""
page_title: "ldap_group_member Resource - terraform-provider-ldap"
description: |-
  ldap_group_member manages a single member of an LDAP group without replacing other members.
---

# ldap_group_member

`ldap_group_member` manages a single member of an LDAP group. Unlike `ldap_group.members`, this resource only controls its own member entry and does not affect other members added outside of Terraform.

Use this resource when you need to add members to a group that is also managed by other systems (e.g. Azure AD sync, manual AD administration). Do not combine this resource with `ldap_group.members` on the same group, as `ldap_group.members` performs a full replacement of the member list.

## Example Usage

```hcl
resource "ldap_group_member" "example" {
  group_dn  = "CN=My-Group,OU=Groups,DC=domain,DC=tld"
  member_dn = "CN=John Doe,OU=Users,DC=domain,DC=tld"
}
```

Multiple members on the same group:

```hcl
resource "ldap_group_member" "app_user1" {
  group_dn  = "CN=MyApp-SSO,OU=Groups,DC=domain,DC=tld"
  member_dn = "CN=Alice,OU=Users,DC=domain,DC=tld"
}

resource "ldap_group_member" "app_user2" {
  group_dn  = "CN=MyApp-SSO,OU=Groups,DC=domain,DC=tld"
  member_dn = "CN=Bob,OU=Users,DC=domain,DC=tld"
}
```

## Argument Reference

* `group_dn` - (Required) The full DN of the LDAP group to add the member to. Changes force recreation.
* `member_dn` - (Required) The full DN of the LDAP object (user, group, or computer) to add as a member. Changes force recreation.

## Attribute Reference

* `id` - Unique identifier for this membership in the format `<group_dn>/<member_dn>`.

## Import

A group membership can be imported using the composite ID `<group_dn>/<member_dn>`, e.g.

```
$ terraform import ldap_group_member.example "CN=My-Group,OU=Groups,DC=domain,DC=tld/CN=John Doe,OU=Users,DC=domain,DC=tld"
```
