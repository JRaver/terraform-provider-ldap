---
subcategory: ""
page_title: "ldap_group_attribute Resource - terraform-provider-ldap"
description: |-
  ldap_group_attribute manages a single LDAP attribute on a group without affecting other attributes.
---

# ldap_group_attribute

`ldap_group_attribute` manages a single LDAP attribute on a group. It only controls its own attribute and does not affect any other attributes — including those managed by the `ldap_group` resource or external systems.

Use this resource to set arbitrary LDAP attributes (e.g. `telephoneNumber`, `department`, `extensionAttribute1`) that are not directly supported by the `ldap_group` resource.

> **Note:** The `ldap_group` resource manages only its own fixed set of attributes (`description`, `displayName`, `groupType`, `managedBy`). It will never overwrite attributes set via `ldap_group_attribute`.

## Example Usage

Single-value attribute on an existing group:

```hcl
resource "ldap_group_attribute" "example_ext" {
  group_dn         = "CN=my-app-admins,OU=Groups,DC=example,DC=com"
  attribute_name   = "extensionAttribute1"
  attribute_values = ["some-value"]
}
```

Multi-value attribute:

```hcl
resource "ldap_group_attribute" "example_phones" {
  group_dn         = "CN=my-app-admins,OU=Groups,DC=example,DC=com"
  attribute_name   = "telephoneNumber"
  attribute_values = ["111-222-3333", "444-555-6666"]
}
```

Used together with `ldap_group`:

```hcl
resource "ldap_group" "example" {
  ou           = "OU=Groups,DC=example,DC=com"
  name         = "my-app-admins"
  description  = "Managed by Terraform"
  display_name = "my-app-admins"
}

resource "ldap_group_attribute" "example_ext" {
  group_dn         = ldap_group.example.id
  attribute_name   = "extensionAttribute1"
  attribute_values = ["some-value"]
}

resource "ldap_group_attribute" "example_dept" {
  group_dn         = ldap_group.example.id
  attribute_name   = "department"
  attribute_values = ["Engineering"]
}
```

## Argument Reference

* `group_dn` - (Required) The full DN of the LDAP group to manage the attribute on. Changes force recreation.
* `attribute_name` - (Required) The LDAP attribute name to manage (e.g. `telephoneNumber`, `department`, `extensionAttribute1`). Changes force recreation.
* `attribute_values` - (Required) List of values for the attribute. At least one value is required. Can be updated in-place.

## Attribute Reference

* `id` - Unique identifier in the format `<group_dn>/<attribute_name>`.

## Import

A group attribute can be imported using the composite ID `<group_dn>/<attribute_name>`, e.g.

```
$ terraform import ldap_group_attribute.example_ext \
  "CN=my-app-admins,OU=Groups,DC=example,DC=com/extensionAttribute1"
```
