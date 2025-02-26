---
page_title: "keycloak_organization_identity_provider Resource"
---

# keycloak_organization_identity_provider Resource

Allows to link an identity provider to an organization.

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
    name = "example.com"
  }

  domain       {
    name = "anotherexample.com"
  }
}

resource "keycloak_oidc_identity_provider" "oidc" {
    realm_id     = keycloak_realm.example.id
    provider_id       = "oidc"
    alias             = "myidp"
    authorization_url = "https://example.com/auth"
    token_url         = "https://example.com/token"
    client_id         = "example_id"
    client_secret     = "example_token"
    default_scopes    = "openid random"
}

resource "keycloak_organization_identity_provider" "oidc" {
    realm_id = data.keycloak_realm.realm.id
    organization_id = keycloak_organization.engineering.id
    identity_provider_alias = keycloak_oidc_identity_provider.oidc.alias
    domain = "example.com"
    redirect_email_domain_matches = true
}

```

## Argument Reference

- `realm_id` - (Required) The realm this identity provider and organization exist in.
- `organization_id` - (Required) The unique ID of the organization.
- `identity_provider_alias` - (Required) The alias of the identity provider to link to the organization.
- `domain` - (Optional) the domain associated to the identity provider.
- `redirect_email_domain_matches` - (Optional) When true, redirects users to the identity provider if the user's email matches the domain.

## Import

Organization identity providers can be imported using the format `{{realm_id}}/{{organization_id}}/{{identity_provider_alias}}`, where `organization` is the unique ID that Keycloak assigns to the organization upon creation. This value can be found in the URI when editing this organization in the GUI, and is typically a GUID.

Example:

```bash
$ terraform import keycloak_organization_identity_provider.oidc my-realm/934a4a4e-28bd-4703-a0fa-332df153aabd/myidp
```
