---
page_title: "keycloak_groups Data Source"
---

# keycloak_groups Data Source

This data source can be used to retrieve all groups.

Information for Keycloak versions < 23:
the datasource will only return the top level groups.

## Example Usage

```hcl
resource "keycloak_realm" "realm" {
    realm   = "my-realm"
    enabled = true
}


data "keycloak_group" "groups" {
    realm_id = keycloak_realm.realm.id
    full_hierarchy = true
}

output "group_paths" {
	value = data.keycloak_group.groups[*].path
}

```

## Argument Reference

- `realm_id` - (Required) The realm this group exists within.
- `full_hierarchy` - (Optional) Retrieve the whole groups hierachy.

## Attributes Reference

The datasource returns an array of keycloak group.
