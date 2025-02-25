package provider

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/keycloak/terraform-provider-keycloak/keycloak"
)

// Test sweeper to clean up test organizations
func TestAccKeycloakOrganization_sweepers(t *testing.T) {
	if ok, _ := keycloakClient.VersionIsLessThan(testCtx, keycloak.Version_26); ok {
		t.Skip()
	}
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOrganization_basic(),
				Check:  testAccCheckKeycloakOrganizationExists("keycloak_organization.organization"),
			},
			{
				PreConfig: func() {
					err := sweepOrganizations()
					if err != nil {
						t.Error(err)
					}
				},
				Config: testKeycloakOrganization_basic(),
				Check: func(state *terraform.State) error {
					_, err := testAccCheckKeycloakOrganizationFetch("keycloak_organization.organization", state)
					if err != nil {
						return fmt.Errorf("organization was not created")
					}
					return nil
				},
			},
		},
	})
}

// Basic CRUD tests
func TestAccKeycloakOrganization_basic(t *testing.T) {
	if ok, _ := keycloakClient.VersionIsLessThan(testCtx, keycloak.Version_26); ok {
		t.Skip()
	}
	t.Parallel()

	orgName := acctest.RandomWithPrefix("tf-acc")
	orgUpdatedName := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakOrganizationDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOrganization_basic_with_name(orgName, orgName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOrganizationExists("keycloak_organization.organization"),
					resource.TestCheckResourceAttr("keycloak_organization.organization", "realm_id", testAccRealm.Realm),
					resource.TestCheckResourceAttr("keycloak_organization.organization", "name", orgName),
				),
			},
			{
				Config: testKeycloakOrganization_basic_with_name(orgUpdatedName, orgName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOrganizationExists("keycloak_organization.organization"),
					resource.TestCheckResourceAttr("keycloak_organization.organization", "name", orgUpdatedName),
				),
			},
		},
	})
}

// Test with all fields
func TestAccKeycloakOrganization_withAllFields(t *testing.T) {
	if ok, _ := keycloakClient.VersionIsLessThan(testCtx, keycloak.Version_26); ok {
		t.Skip()
	}
	t.Parallel()

	orgName := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckKeycloakOrganizationDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOrganization_withAllFields(orgName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOrganizationExists("keycloak_organization.organization"),
					resource.TestCheckResourceAttr("keycloak_organization.organization", "realm_id", testAccRealm.Realm),
					resource.TestCheckResourceAttr("keycloak_organization.organization", "name", orgName),
					resource.TestCheckResourceAttr("keycloak_organization.organization", "redirect_url", "https://localhost:8080"),
					resource.TestCheckResourceAttr("keycloak_organization.organization", "description", "Description for "+orgName),
					resource.TestCheckResourceAttr("keycloak_organization.organization", "alias", orgName),
					resource.TestCheckResourceAttr("keycloak_organization.organization", "domain.#", "2"),
					resource.TestCheckResourceAttr("keycloak_organization.organization", "domain.0.name", "example.com"),
					resource.TestCheckResourceAttr("keycloak_organization.organization", "domain.0.verified", "true"),
					resource.TestCheckResourceAttr("keycloak_organization.organization", "domain.1.name", "example.org"),
					resource.TestCheckResourceAttr("keycloak_organization.organization", "domain.1.verified", "false"),
					resource.TestCheckResourceAttr("keycloak_organization.organization", "attributes.testKey", "testValue"),
				),
			},
		},
	})
}

// Import test
func TestAccKeycloakOrganization_import(t *testing.T) {
	if ok, _ := keycloakClient.VersionIsLessThan(testCtx, keycloak.Version_26); ok {
		t.Skip()
	}
	orgName := acctest.RandomWithPrefix("terraform-org")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckKeycloakOrganizationDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOrganization_basic_with_name(orgName, orgName),
				Check:  testAccCheckKeycloakOrganizationExists("keycloak_organization.organization"),
			},
			{
				ResourceName:      "keycloak_organization.organization",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getOrganizationImportId("keycloak_organization.organization"),
			},
		},
	})
}

// Test Terraform configurations
func testKeycloakOrganization_basic() string {
	return fmt.Sprintf(`

data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_organization" "organization" {
	realm_id = data.keycloak_realm.realm.id
	name     = "terraform-organization"

	domain {
		name = "examplebasic.com"
	}
}
	`, testAccRealm.Realm)
}

func testKeycloakOrganization_basic_with_name(orgName string, alias string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}
resource "keycloak_organization" "organization" {
	realm_id = data.keycloak_realm.realm.id
	name     = "%s"
	alias    = "%s"

	domain {
		name = "examplebasicwithname.com"
	}
}
	`, testAccRealm.Realm, orgName, alias)
}

func testKeycloakOrganization_withAllFields(name string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}
resource "keycloak_organization" "organization" {
	realm_id = data.keycloak_realm.realm.id
	name         = "%s"
	alias        = "%s"
	redirect_url = "https://localhost:8080"
	description  = "Description for %s"

	domain {
		name = "example.com"
		verified = true
	}

	domain {
		name = "example.org"
	}

	attributes   = {
		"testKey" = "testValue"
	}
}
	`, testAccRealm.Realm, name, name, name)
}

// Helper functions for tests
func testAccCheckKeycloakOrganizationExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := testAccCheckKeycloakOrganizationFetch(resourceName, s)
		if err != nil {
			return err
		}
		return nil
	}
}

func testAccCheckKeycloakOrganizationFetch(resourceName string, s *terraform.State) (*keycloak.Organization, error) {
	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	id := rs.Primary.ID
	realmId := rs.Primary.Attributes["realm_id"]

	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	organization, err := keycloakClient.GetOrganization(context.TODO(), realmId, id)
	if err != nil {
		return nil, fmt.Errorf("error getting organization: %s", err)
	}

	return organization, nil
}

func testAccCheckKeycloakOrganizationDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_organization" {
				continue
			}

			id := rs.Primary.ID
			realmId := rs.Primary.Attributes["realm_id"]

			keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

			organization, _ := keycloakClient.GetOrganization(context.TODO(), realmId, id)
			if organization != nil {
				return fmt.Errorf("organization still exists")
			}
		}

		return nil
	}
}

func getOrganizationImportId(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}

		id := rs.Primary.ID
		realmId := rs.Primary.Attributes["realm_id"]

		return fmt.Sprintf("%s/%s", realmId, id), nil
	}
}

// Sweeper function to clean up test organizations
func sweepOrganizations() error {

	organizations, err := keycloakClient.GetOrganizations(context.TODO(), testAccRealm.Realm)
	if err != nil {
		return err
	}

	for _, organization := range organizations {
		if strings.HasPrefix(organization.Name, "tf-acc") {
			err = keycloakClient.DeleteOrganization(context.TODO(), organization.RealmId, organization.Id)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
