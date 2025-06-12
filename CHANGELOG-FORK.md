## 5.2.0-1.4.1 (June 12, 2025)

- Fix the smtp server authentication broken on realm updates on Keycloak 26.2 ([#12](https://github.com/lucdew/terraform-provider-keycloak/pull/12))

## 5.2.0-1.4.0 (April 21, 2025)

- Merge the official keycloak provider v5.2.0 ([#11](https://github.com/lucdew/terraform-provider-keycloak/pull/11))

## 5.1.1-1.4.0 (March 27, 2025)

- Remove useless feature ([#10](https://github.com/lucdew/terraform-provider-keycloak/pull/10)). Client certificate authentication works by setting a
  dummy client_secret. Setting the client_assertion_type is not needed.

## 5.1.1-1.3.0 (March 27, 2025)

FEATURES:

- feat: Add tls_client_auth in the provider for authenticating with a client certificate and private key ([#10](https://github.com/lucdew/terraform-provider-keycloak/pull/10))

## 5.1.1-1.2.0 (February 26, 2025)

FEATURES:

- feat: Add keycloak_organization_identity_provider resource ([#9](https://github.com/lucdew/terraform-provider-keycloak/pull/9))

## 5.1.1-1.1.0 (February 25, 2025)

FEATURES:

- feat: Add partial organization support ([#8](https://github.com/lucdew/terraform-provider-keycloak/pull/8))

## 5.1.1-1.0.0 (February 20, 2025)

FEATURES:

- feat: Add mTLS support ([#3](https://github.com/lucdew/terraform-provider-keycloak/pull/3))
- feat: Add keycloak_groups datasource ([#4](https://github.com/lucdew/terraform-provider-keycloak/pull/4))
- feat: Add new resource_keyloak_realm_keystore_rsa_enc_generated ([#5](https://github.com/lucdew/terraform-provider-keycloak/pull/5))
