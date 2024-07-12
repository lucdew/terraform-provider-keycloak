package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/lucdew/terraform-provider-keycloak/keycloak"
)

// Tests for attaching SAML mappers to SAML client scopes are omitted
// because the keycloak_saml_client_scope resource does not exist yet.

func TestAccKeycloakSamlHardcodedAttributeProtocolMapper_basicClient(t *testing.T) {
	t.Parallel()
	clientId := acctest.RandomWithPrefix("tf-acc")
	mapperName := acctest.RandomWithPrefix("tf-acc")

	resourceName := "keycloak_saml_hardcoded_attribute_protocol_mapper.saml_hardcoded_attribute_mapper"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccKeycloakSamlHardcodedAttributeProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakSamlHardcodedAttributeProtocolMapper_basic_client(clientId, mapperName),
				Check:  testKeycloakSamlHardcodedAttributeProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakSamlHardcodedAttributeProtocolMapper_import(t *testing.T) {
	t.Parallel()
	clientId := acctest.RandomWithPrefix("tf-acc")
	mapperName := acctest.RandomWithPrefix("tf-acc")

	clientResourceName := "keycloak_saml_hardcoded_attribute_protocol_mapper.saml_hardcoded_attribute_mapper"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccKeycloakSamlHardcodedAttributeProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakSamlHardcodedAttributeProtocolMapper_basic_client(clientId, mapperName),
				Check:  testKeycloakSamlHardcodedAttributeProtocolMapperExists(clientResourceName),
			},
			{
				ResourceName:      clientResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getGenericProtocolMapperIdForClient(clientResourceName),
			},
		},
	})
}

func TestAccKeycloakSamlHardcodedAttributeProtocolMapper_update(t *testing.T) {
	t.Parallel()
	clientId := acctest.RandomWithPrefix("tf-acc")
	mapperName := acctest.RandomWithPrefix("tf-acc")

	userAttribute := acctest.RandomWithPrefix("tf-acc")
	updatedUserAttribute := acctest.RandomWithPrefix("tf-acc")
	resourceName := "keycloak_saml_hardcoded_attribute_protocol_mapper.saml_hardcoded_attribute_mapper"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccKeycloakSamlHardcodedAttributeProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakSamlHardcodedAttributeProtocolMapper_attributeValue(clientId, mapperName, userAttribute),
				Check:  testKeycloakSamlHardcodedAttributeProtocolMapperExists(resourceName),
			},
			{
				Config: testKeycloakSamlHardcodedAttributeProtocolMapper_attributeValue(clientId, mapperName, updatedUserAttribute),
				Check:  testKeycloakSamlHardcodedAttributeProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakSamlHardcodedAttributeProtocolMapper_createAfterManualDestroy(t *testing.T) {
	t.Parallel()
	mapper := &keycloak.SamlHardcodedAttributeProtocolMapper{}

	clientId := acctest.RandomWithPrefix("tf-acc")
	mapperName := acctest.RandomWithPrefix("tf-acc")

	resourceName := "keycloak_saml_hardcoded_attribute_protocol_mapper.saml_hardcoded_attribute_mapper"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccKeycloakSamlHardcodedAttributeProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakSamlHardcodedAttributeProtocolMapper_basic_client(clientId, mapperName),
				Check:  testKeycloakSamlHardcodedAttributeProtocolMapperFetch(resourceName, mapper),
			},
			{
				PreConfig: func() {
					err := keycloakClient.DeleteSamlHardcodedAttributeProtocolMapper(testCtx, mapper.RealmId, mapper.ClientId, mapper.ClientScopeId, mapper.Id)
					if err != nil {
						t.Error(err)
					}
				},
				Config: testKeycloakSamlHardcodedAttributeProtocolMapper_basic_client(clientId, mapperName),
				Check:  testKeycloakSamlHardcodedAttributeProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakSamlHardcodedAttributeProtocolMapper_validateClaimValueType(t *testing.T) {
	t.Parallel()
	clientId := acctest.RandomWithPrefix("tf-acc")
	mapperName := acctest.RandomWithPrefix("tf-acc")
	invalidSamlNameFormat := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccKeycloakSamlHardcodedAttributeProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakSamlHardcodedAttributeProtocolMapper_samlAttributeNameFormat(clientId, mapperName, invalidSamlNameFormat),
				ExpectError: regexp.MustCompile("expected saml_attribute_name_format to be one of .+ got " + invalidSamlNameFormat),
			},
		},
	})
}

func TestAccKeycloakSamlHardcodedAttributeProtocolMapper_updateClientIdForceNew(t *testing.T) {
	t.Parallel()
	clientId := acctest.RandomWithPrefix("tf-acc")
	updatedClientId := acctest.RandomWithPrefix("tf-acc")
	mapperName := acctest.RandomWithPrefix("tf-acc")

	attributeValue := acctest.RandomWithPrefix("tf-acc")
	resourceName := "keycloak_saml_hardcoded_attribute_protocol_mapper.saml_hardcoded_attribute_mapper"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccKeycloakSamlHardcodedAttributeProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakSamlHardcodedAttributeProtocolMapper_attributeValue(clientId, mapperName, attributeValue),
				Check:  testKeycloakSamlHardcodedAttributeProtocolMapperExists(resourceName),
			},
			{
				Config: testKeycloakSamlHardcodedAttributeProtocolMapper_attributeValue(updatedClientId, mapperName, attributeValue),
				Check:  testKeycloakSamlHardcodedAttributeProtocolMapperExists(resourceName),
			},
		},
	})
}

func testAccKeycloakSamlHardcodedAttributeProtocolMapperDestroy() resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for resourceName, rs := range state.RootModule().Resources {
			if rs.Type != "keycloak_saml_hardcoded_attribute_protocol_mapper" {
				continue
			}

			mapper, _ := getSamlHardcodedAttributeMapperUsingState(state, resourceName)

			if mapper != nil {
				return fmt.Errorf("saml user property protocol mapper with id %s still exists", rs.Primary.ID)
			}
		}

		return nil
	}
}

func testKeycloakSamlHardcodedAttributeProtocolMapperExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		_, err := getSamlHardcodedAttributeMapperUsingState(state, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testKeycloakSamlHardcodedAttributeProtocolMapperFetch(resourceName string, mapper *keycloak.SamlHardcodedAttributeProtocolMapper) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		fetchedMapper, err := getSamlHardcodedAttributeMapperUsingState(state, resourceName)
		if err != nil {
			return err
		}

		mapper.Id = fetchedMapper.Id
		mapper.ClientId = fetchedMapper.ClientId
		mapper.ClientScopeId = fetchedMapper.ClientScopeId
		mapper.RealmId = fetchedMapper.RealmId

		return nil
	}
}

func getSamlHardcodedAttributeMapperUsingState(state *terraform.State, resourceName string) (*keycloak.SamlHardcodedAttributeProtocolMapper, error) {
	rs, ok := state.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found in TF state: %s ", resourceName)
	}

	id := rs.Primary.ID
	realm := rs.Primary.Attributes["realm_id"]
	clientId := rs.Primary.Attributes["client_id"]
	clientScopeId := rs.Primary.Attributes["client_scope_id"]

	return keycloakClient.GetSamlHardcodedAttributeProtocolMapper(testCtx, realm, clientId, clientScopeId, id)
}

func testKeycloakSamlHardcodedAttributeProtocolMapper_basic_client(clientId, mapperName string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_saml_client" "saml_client" {
	realm_id  = data.keycloak_realm.realm.id
	client_id = "%s"
}

resource "keycloak_saml_hardcoded_attribute_protocol_mapper" "saml_hardcoded_attribute_mapper" {
	name                       = "%s"
	realm_id                   = data.keycloak_realm.realm.id
	client_id                  = keycloak_saml_client.saml_client.id

	attribute_value            = "a_static_value"
	saml_attribute_name        = "role"
	saml_attribute_name_format = "Unspecified"
}`, testAccRealm.Realm, clientId, mapperName)
}

func testKeycloakSamlHardcodedAttributeProtocolMapper_attributeValue(clientId, mapperName, attributeValue string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_saml_client" "saml_client" {
	realm_id  = data.keycloak_realm.realm.id
	client_id = "%s"
}

resource "keycloak_saml_hardcoded_attribute_protocol_mapper" "saml_hardcoded_attribute_mapper" {
	name                       = "%s"
	realm_id                   = data.keycloak_realm.realm.id
	client_id                  = keycloak_saml_client.saml_client.id

	attribute_value              = "%s"
	saml_attribute_name        = "test"
	saml_attribute_name_format = "Unspecified"
}`, testAccRealm.Realm, clientId, mapperName, attributeValue)
}

func testKeycloakSamlHardcodedAttributeProtocolMapper_samlAttributeNameFormat(clientName, mapperName, samlAttributeNameFormat string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_saml_client" "saml_client" {
	realm_id  = data.keycloak_realm.realm.id
	client_id = "%s"
}

resource "keycloak_saml_hardcoded_attribute_protocol_mapper" "saml_hardcoded_attribute_mapper" {
	name                       = "%s"
	realm_id                   = data.keycloak_realm.realm.id
	client_id                  = keycloak_saml_client.saml_client.id

	attribute_value            = "a_static_value"
	saml_attribute_name        = "role"
	saml_attribute_name_format = "%s"
}`, testAccRealm.Realm, clientName, mapperName, samlAttributeNameFormat)
}
