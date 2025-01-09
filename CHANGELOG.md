## 4.8.4 (Jan 9, 2025)

BUG FIX:

- Fix the hideOnLogin attribute wrongly set on the identity provider for Keycloak < 26.


## 4.8.3 (Dec 21, 2024)

FEATURES:

- Update identity provider resource to support Keycloak 26 hideOnLogin attribute.
    Keep the hide_on_login_page resource configuration parameter.

## 4.7.0 (July 18, 2024)

FEATURES:

- Add the new resource keycloak_groups

## 4.6.0 (July 18, 2024)

FEATURES:

- Add the new resource keycloak_realm_keystore_rsa_enc_generated
- Add the new resource keycloak_saml_user_session_note_protocol_mapper

## 4.5.0 (July 12, 2024)

FEATURES:

- Add the new resource keycloak_saml_hardcoded_attribute_protocol_mapper

## 4.4.2 (June 12, 2024)

FEATURES:

- Rename parameters for mTLS support to be consistent with existing tls_insecure_skip_verify parameter

## 4.4.1 (June 9, 2024)

FEATURES:

- Add mTLS support
