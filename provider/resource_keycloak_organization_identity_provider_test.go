package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/keycloak/terraform-provider-keycloak/keycloak"
)

func TestAccKeycloakOrganizationIdentityProvider(t *testing.T) {
	if ok, _ := keycloakClient.VersionIsLessThan(testCtx, keycloak.Version_26); ok {
		t.Skip()
	}
	t.Parallel()

	orgName := acctest.RandomWithPrefix("tf-acc")
	idpAlias := orgName

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckKeycloakOrganizationIdentityProviderDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOrganizationIdentityProvider(orgName, idpAlias),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOrganizationIdentityProvider("keycloak_organization_identity_provider.this", orgName+".example.com", true),
					resource.TestCheckResourceAttr("keycloak_organization_identity_provider.this", "realm_id", testAccRealm.Realm),
					resource.TestCheckResourceAttr("keycloak_organization_identity_provider.this", "identity_provider_alias", idpAlias),
				),
			},
		},
	})
}

func TestAccKeycloakOrganizationIdentityProvider_import(t *testing.T) {
	if ok, _ := keycloakClient.VersionIsLessThan(testCtx, keycloak.Version_26); ok {
		t.Skip()
	}
	t.Parallel()
	orgName := acctest.RandomWithPrefix("tf-acc")
	idpAlias := orgName

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckKeycloakOrganizationIdentityProviderDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOrganizationIdentityProvider(orgName, idpAlias),
				Check:  testAccCheckKeycloakOrganizationIdentityProvider("keycloak_organization_identity_provider.this", orgName+".example.com", true),
			},
			{
				ResourceName:      "keycloak_organization_identity_provider.this",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getOrganizationIdentityProviderImportId("keycloak_organization_identity_provider.this"),
			},
		},
	})
}

func TestAccKeycloakOrganizationIdentityProvider_delete(t *testing.T) {
	if ok, _ := keycloakClient.VersionIsLessThan(testCtx, keycloak.Version_26); ok {
		t.Skip()
	}
	t.Parallel()
	orgName := acctest.RandomWithPrefix("tf-acc")
	idpAlias := orgName

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckKeycloakOrganizationIdentityProviderDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOrganizationIdentityProvider(orgName, idpAlias),
				Check:  testAccCheckKeycloakOrganizationIdentityProvider("keycloak_organization_identity_provider.this", orgName+".example.com", true),
			},
			{
				Config: testKeycloakOrganizationIdentityProviderDependencies(orgName, idpAlias),
				Check:  testAccCheckKeycloakOrganizationIdentityProviderDestroyOnly(),
			},
		},
	})
}

func testKeycloakOrganizationIdentityProviderDependencies(orgName string, idpAlias string) string {

	return fmt.Sprintf(`

data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_organization" "organization" {
	realm_id = data.keycloak_realm.realm.id
	name     = "%s"

	domain {
		name = "%s.example.com"
	}
}

resource "keycloak_oidc_identity_provider" "oidc" {
	realm             = data.keycloak_realm.realm.id
	provider_id       = "oidc"
	alias             = "%s"
	authorization_url = "https://example.com/auth"
	token_url         = "https://example.com/token"
	client_id         = "example_id"
	client_secret     = "example_token"
	default_scopes    = "openid random"
}
	`, testAccRealm.Realm, orgName, orgName, idpAlias)

}

// Test Terraform configurations
func testKeycloakOrganizationIdentityProvider(orgName string, idpAlias string) string {
	return testKeycloakOrganizationIdentityProviderDependencies(orgName, idpAlias) + fmt.Sprintf(`

resource "keycloak_organization_identity_provider" "this" {
	realm_id = data.keycloak_realm.realm.id
	organization_id = keycloak_organization.organization.id
	identity_provider_alias = keycloak_oidc_identity_provider.oidc.alias
	domain = "%s.example.com"
    redirect_email_domain_matches = true
}
	`, orgName)
}

// Helper functions for tests
func testAccCheckKeycloakOrganizationIdentityProvider(resourceName string, expectedDomain string, expectedRedirectEmailMatch bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		domain, redirectEmailsMatches, err := testAccCheckKeycloakOrganizationIdentityProviderFetch(resourceName, s)
		if err != nil {
			return err
		}
		if domain != expectedDomain {
			return fmt.Errorf("Invalid domain, expected=%s, got=%s", expectedDomain, domain)
		}
		if redirectEmailsMatches != expectedRedirectEmailMatch {
			return fmt.Errorf("Invalid redirectsEmailMatches, expected=%t, got=%t", expectedRedirectEmailMatch, redirectEmailsMatches)
		}
		return nil
	}
}

func testAccCheckKeycloakOrganizationIdentityProviderFetch(resourceName string, s *terraform.State) (string, bool, error) {
	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return "", false, fmt.Errorf("resource not found: %s", resourceName)
	}

	realmId := rs.Primary.Attributes["realm_id"]
	orgId := rs.Primary.Attributes["organization_id"]
	idpAlias := rs.Primary.Attributes["identity_provider_alias"]

	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	err := keycloakClient.CheckIdentityProviderLinkToOrganization(context.TODO(), realmId, orgId, idpAlias)
	if err != nil {
		return "", false, fmt.Errorf("error getting organization identity provider: %s", err)
	}

	idp, err := keycloakClient.GetIdentityProvider(context.TODO(), realmId, idpAlias)
	if err != nil {
		return "", false, err
	}

	return idp.Config.OrgDomain, bool(idp.Config.OrgRedirectEmailMatches), nil
}

func testAccCheckKeycloakOrganizationIdentityProviderDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_organization_identity_provider" {
				continue
			}

			realmId := rs.Primary.Attributes["realm_id"]
			orgId := rs.Primary.Attributes["organization_id"]
			idpAlias := rs.Primary.Attributes["identity_provider_alias"]

			keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

			err := keycloakClient.CheckIdentityProviderLinkToOrganization(context.TODO(), realmId, orgId, idpAlias)
			if err == nil {
				return fmt.Errorf("organization identity provider link still exists")
			}

			_, err = keycloakClient.GetIdentityProvider(context.TODO(), realmId, idpAlias)
			if err == nil {
				return fmt.Errorf("organization identity provider still exists")
			}

		}

		return nil
	}
}

func testAccCheckKeycloakOrganizationIdentityProviderDestroyOnly() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_organization_identity_provider" {
				continue
			}

			realmId := rs.Primary.Attributes["realm_id"]
			orgId := rs.Primary.Attributes["organization_id"]
			idpAlias := rs.Primary.Attributes["identity_provider_alias"]

			keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

			err := keycloakClient.CheckIdentityProviderLinkToOrganization(context.TODO(), realmId, orgId, idpAlias)
			if err == nil {
				return fmt.Errorf("organization identity provider link still exists")
			}

			idp, err := keycloakClient.GetIdentityProvider(context.TODO(), realmId, idpAlias)
			if err != nil {
				return err
			}
			if idp.OrganizationId != "" {
				return fmt.Errorf("identity provider is still linked to organization")
			}
			if idp.Config.OrgDomain != "" {
				return fmt.Errorf("identity provider has still organization domain")
			}

		}

		return nil
	}
}

func getOrganizationIdentityProviderImportId(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}

		return rs.Primary.ID, nil
	}
}
