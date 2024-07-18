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

func TestAccKeycloakSamlUserSessionNoteProtocolMapper_basicClient(t *testing.T) {
	t.Parallel()
	clientId := acctest.RandomWithPrefix("tf-acc")
	mapperName := acctest.RandomWithPrefix("tf-acc")

	resourceName := "keycloak_saml_user_session_note_protocol_mapper.saml_user_session_note_mapper"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccKeycloakSamlUserSessionNoteProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakSamlUserSessionNoteProtocolMapper_basic_client(clientId, mapperName),
				Check:  testKeycloakSamlUserSessionNoteProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakSamlUserSessionNoteProtocolMapper_import(t *testing.T) {
	t.Parallel()
	clientId := acctest.RandomWithPrefix("tf-acc")
	mapperName := acctest.RandomWithPrefix("tf-acc")

	clientResourceName := "keycloak_saml_user_session_note_protocol_mapper.saml_user_session_note_mapper"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccKeycloakSamlUserSessionNoteProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakSamlUserSessionNoteProtocolMapper_basic_client(clientId, mapperName),
				Check:  testKeycloakSamlUserSessionNoteProtocolMapperExists(clientResourceName),
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

func TestAccKeycloakSamlUserSessionNoteProtocolMapper_update(t *testing.T) {
	t.Parallel()
	clientId := acctest.RandomWithPrefix("tf-acc")
	mapperName := acctest.RandomWithPrefix("tf-acc")

	userAttribute := acctest.RandomWithPrefix("tf-acc")
	updatedUserAttribute := acctest.RandomWithPrefix("tf-acc")
	resourceName := "keycloak_saml_user_session_note_protocol_mapper.saml_user_session_note_mapper"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccKeycloakSamlUserSessionNoteProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakSamlUserSessionNoteProtocolMapper_userProperty(clientId, mapperName, userAttribute),
				Check:  testKeycloakSamlUserSessionNoteProtocolMapperExists(resourceName),
			},
			{
				Config: testKeycloakSamlUserSessionNoteProtocolMapper_userProperty(clientId, mapperName, updatedUserAttribute),
				Check:  testKeycloakSamlUserSessionNoteProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakSamlUserSessionNoteProtocolMapper_createAfterManualDestroy(t *testing.T) {
	t.Parallel()
	mapper := &keycloak.SamlUserSessionNoteProtocolMapper{}

	clientId := acctest.RandomWithPrefix("tf-acc")
	mapperName := acctest.RandomWithPrefix("tf-acc")

	resourceName := "keycloak_saml_user_session_note_protocol_mapper.saml_user_session_note_mapper"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccKeycloakSamlUserSessionNoteProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakSamlUserSessionNoteProtocolMapper_basic_client(clientId, mapperName),
				Check:  testKeycloakSamlUserSessionNoteProtocolMapperFetch(resourceName, mapper),
			},
			{
				PreConfig: func() {
					err := keycloakClient.DeleteSamlUserSessionNoteProtocolMapper(testCtx, mapper.RealmId, mapper.ClientId, mapper.ClientScopeId, mapper.Id)
					if err != nil {
						t.Error(err)
					}
				},
				Config: testKeycloakSamlUserSessionNoteProtocolMapper_basic_client(clientId, mapperName),
				Check:  testKeycloakSamlUserSessionNoteProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakSamlUserSessionNoteProtocolMapper_validateClaimValueType(t *testing.T) {
	t.Parallel()
	clientId := acctest.RandomWithPrefix("tf-acc")
	mapperName := acctest.RandomWithPrefix("tf-acc")
	invalidSamlNameFormat := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccKeycloakSamlUserSessionNoteProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakSamlUserSessionNoteProtocolMapper_samlAttributeNameFormat(clientId, mapperName, invalidSamlNameFormat),
				ExpectError: regexp.MustCompile("expected saml_attribute_name_format to be one of .+ got " + invalidSamlNameFormat),
			},
		},
	})
}

func TestAccKeycloakSamlUserSessionNoteProtocolMapper_updateClientIdForceNew(t *testing.T) {
	t.Parallel()
	clientId := acctest.RandomWithPrefix("tf-acc")
	updatedClientId := acctest.RandomWithPrefix("tf-acc")
	mapperName := acctest.RandomWithPrefix("tf-acc")

	userAttribute := acctest.RandomWithPrefix("tf-acc")
	resourceName := "keycloak_saml_user_session_note_protocol_mapper.saml_user_session_note_mapper"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccKeycloakSamlUserSessionNoteProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakSamlUserSessionNoteProtocolMapper_userProperty(clientId, mapperName, userAttribute),
				Check:  testKeycloakSamlUserSessionNoteProtocolMapperExists(resourceName),
			},
			{
				Config: testKeycloakSamlUserSessionNoteProtocolMapper_userProperty(updatedClientId, mapperName, userAttribute),
				Check:  testKeycloakSamlUserSessionNoteProtocolMapperExists(resourceName),
			},
		},
	})
}

func testAccKeycloakSamlUserSessionNoteProtocolMapperDestroy() resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for resourceName, rs := range state.RootModule().Resources {
			if rs.Type != "keycloak_saml_user_session_note_protocol_mapper" {
				continue
			}

			mapper, _ := getSamlUserSessionNoteMapperUsingState(state, resourceName)

			if mapper != nil {
				return fmt.Errorf("saml user property protocol mapper with id %s still exists", rs.Primary.ID)
			}
		}

		return nil
	}
}

func testKeycloakSamlUserSessionNoteProtocolMapperExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		_, err := getSamlUserSessionNoteMapperUsingState(state, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testKeycloakSamlUserSessionNoteProtocolMapperFetch(resourceName string, mapper *keycloak.SamlUserSessionNoteProtocolMapper) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		fetchedMapper, err := getSamlUserSessionNoteMapperUsingState(state, resourceName)
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

func getSamlUserSessionNoteMapperUsingState(state *terraform.State, resourceName string) (*keycloak.SamlUserSessionNoteProtocolMapper, error) {
	rs, ok := state.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found in TF state: %s ", resourceName)
	}

	id := rs.Primary.ID
	realm := rs.Primary.Attributes["realm_id"]
	clientId := rs.Primary.Attributes["client_id"]
	clientScopeId := rs.Primary.Attributes["client_scope_id"]

	return keycloakClient.GetSamlUserSessionNoteProtocolMapper(testCtx, realm, clientId, clientScopeId, id)
}

func testKeycloakSamlUserSessionNoteProtocolMapper_basic_client(clientId, mapperName string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_saml_client" "saml_client" {
	realm_id  = data.keycloak_realm.realm.id
	client_id = "%s"
}

resource "keycloak_saml_user_session_note_protocol_mapper" "saml_user_session_note_mapper" {
	name                       = "%s"
	realm_id                   = data.keycloak_realm.realm.id
	client_id                  = keycloak_saml_client.saml_client.id

	note_name                  = "idp"
	saml_attribute_name        = "idp"
	saml_attribute_name_format = "Unspecified"
}`, testAccRealm.Realm, clientId, mapperName)
}

func testKeycloakSamlUserSessionNoteProtocolMapper_userProperty(clientId, mapperName, userProperty string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_saml_client" "saml_client" {
	realm_id  = data.keycloak_realm.realm.id
	client_id = "%s"
}

resource "keycloak_saml_user_session_note_protocol_mapper" "saml_user_session_note_mapper" {
	name                       = "%s"
	realm_id                   = data.keycloak_realm.realm.id
	client_id                  = keycloak_saml_client.saml_client.id

	note_name                  = "%s"
	saml_attribute_name        = "test"
	saml_attribute_name_format = "Unspecified"
}`, testAccRealm.Realm, clientId, mapperName, userProperty)
}

func testKeycloakSamlUserSessionNoteProtocolMapper_samlAttributeNameFormat(clientName, mapperName, samlAttributeNameFormat string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_saml_client" "saml_client" {
	realm_id  = data.keycloak_realm.realm.id
	client_id = "%s"
}

resource "keycloak_saml_user_session_note_protocol_mapper" "saml_user_session_note_mapper" {
	name                       = "%s"
	realm_id                   = data.keycloak_realm.realm.id
	client_id                  = keycloak_saml_client.saml_client.id

	note_name                  = "idp"
	saml_attribute_name        = "idp"
	saml_attribute_name_format = "%s"
}`, testAccRealm.Realm, clientName, mapperName, samlAttributeNameFormat)
}
