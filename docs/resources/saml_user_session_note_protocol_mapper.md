---
page_title: "keycloak_saml_user_session_note_protocol_mapper Resource"
---

# keycloak\_saml\_user\_session\_note\_protocol\_mapper Resource

Allows for creating and managing user session protocol mappers for SAML clients within Keycloak.

SAML user session note protocol mappers allow you to set in a SAML assertion attribute a value from a user session note.

Protocol mappers can be defined for a single client, or they can be defined for a client scope which can be shared between
multiple different clients.

## Example Usage

```hcl
resource "keycloak_realm" "realm" {
  realm   = "my-realm"
  enabled = true
}

resource "keycloak_saml_client" "saml_client" {
  realm_id  = keycloak_realm.realm.id
  client_id = "saml-client"
  name      = "saml-client"
}

resource "keycloak_saml_user_session_note_protocol_mapper" "saml_user_session_note_mapper" {
  realm_id  = keycloak_realm.realm.id
  client_id = keycloak_saml_client.saml_client.id
  name      = "idp-user-session-note-mapper"

  note_name                  = "idp"
  saml_attribute_name        = "idp"
  saml_attribute_name_format = "Unspecified"
}
```

## Argument Reference

- `realm_id` - (Required) The realm this protocol mapper exists within.
- `name` - (Required) The display name of this protocol mapper in the GUI.
- `note_name` - (Required) The name of the user session note to get the value from for the SAML attribute.
- `saml_attribute_name` - (Required) The name of the SAML attribute.
- `saml_attribute_name_format` - (Required) The SAML attribute Name Format. Can be one of `Unspecified`, `Basic`, or `URI Reference`.
- `client_id` - (Optional) The client this protocol mapper should be attached to. Conflicts with `client_scope_id`. One of `client_id` or `client_scope_id` must be specified.
- `client_scope_id` - (Optional) The client scope this protocol mapper should be attached to. Conflicts with `client_id`. One of `client_id` or `client_scope_id` must be specified.
- `friendly_name` - (Optional) An optional human-friendly name for this attribute.

## Import

Protocol mappers can be imported using one of the following formats:
- Client: `{{realm_id}}/client/{{client_keycloak_id}}/{{protocol_mapper_id}}`
- Client Scope: `{{realm_id}}/client-scope/{{client_scope_keycloak_id}}/{{protocol_mapper_id}}`

Example:

```bash
$ terraform import keycloak_saml_user_session_note_protocol_mapper.saml_user_session_note_mapper my-realm/client/a7202154-8793-4656-b655-1dd18c181e14/71602afa-f7d1-4788-8c49-ef8fd00af0f4
$ terraform import keycloak_saml_user_session_note_protocol_mapper.saml_user_session_note_mapper my-realm/client-scope/b799ea7e-73ee-4a73-990a-1eafebe8e20a/71602afa-f7d1-4788-8c49-ef8fd00af0f4
```
