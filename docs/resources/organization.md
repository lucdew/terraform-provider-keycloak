---
page_title: "keycloak_organization Resource"
---

# keycloak_organization Resource

Allows for creating and managing Organizations within Keycloak.

It only configures the Organization own data but not the association to the identity providers or organization members.

## Example Usage

```hcl
resource "keycloak_realm" "realm" {
  realm   = "my-realm"
  enabled = true
}

resource "keycloak_organization" "engineering" {
  realm_id     = keycloak_realm.example.id
  name         = "engineering"
  alias        = "engineering"
  description  = "Organization for the engineering department"

  domain       {
	name = "engineering.example.com"
	verified = true
  }

  domain       {
	name = "engineering-lab.example.com"
  }

  attributes   = {
    department = "technical"
    location   = "headquarter"
  }
}



```

## Argument Reference

- `realm_id` - (Required) The realm this organization exists in.
- `name` - (Required) The name of the organization
- `alias` - (Optional) The alias of the organization. It cannot be updated and must not have spaces. If not set it is computed.
- `enabled` - (Optional) When `false`, members will not be able to access this organization. Defaults to `true`.
- `description` - (Optional) the organization description.
- `redirect_url` - (Optional) the organization redirect url.
- `attributes` - (Optional) A map representing attributes for the organization. In order to add multivalued attributes, use `##` to separate the values. Max length for each value is 255 chars

### Domains

Associated domains can be configured by using 1 or more `domain` block, which supports the following arguments:

- `name` - (Required) The domain name. Must be unique.
- `verified` - (Optional) When `true`, indicates that the domain has been verified. Defaults to `false`.

## Import

Organizations can be imported using the format `{{realm_id}}/{{organization_id}}`, where `organization` is the unique ID that Keycloak assigns to the organization upon creation. This value can be found in the URI when editing this organization in the GUI, and is typically a GUID.

Example:

```bash
$ terraform import keycloak_organization.engineering my-realm/934a4a4e-28bd-4703-a0fa-332df153aabd
```
