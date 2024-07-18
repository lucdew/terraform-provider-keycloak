package provider

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/lucdew/terraform-provider-keycloak/keycloak"
)

func TestAccKeycloakRealmKeystoreRsaEncGenerated_basic(t *testing.T) {
	t.Parallel()

	rsaName := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckRealmKeystoreRsaEncGeneratedDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakRealmKeystoreRsaEncGenerated_basic(rsaName),
				Check:  testAccCheckRealmKeystoreRsaEncGeneratedExists("keycloak_realm_keystore_rsa_enc_generated.realm_rsa_enc"),
			},
			{
				ResourceName:      "keycloak_realm_keystore_rsa_enc_generated.realm_rsa_enc",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getRealmKeystoreGenericImportId("keycloak_realm_keystore_rsa_enc_generated.realm_rsa_enc"),
			},
		},
	})
}

func TestAccKeycloakRealmKeystoreRsaEncGenerated_createAfterManualDestroy(t *testing.T) {
	t.Parallel()

	rsa := &keycloak.RealmKeystoreRsaEncGenerated{}

	fullNameKeystoreName := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckRealmKeystoreRsaEncGeneratedDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakRealmKeystoreRsaEncGenerated_basic(fullNameKeystoreName),
				Check:  testAccCheckRealmKeystoreRsaEncGeneratedFetch("keycloak_realm_keystore_rsa_enc_generated.realm_rsa_enc", rsa),
			},
			{
				PreConfig: func() {
					err := keycloakClient.DeleteRealmKeystoreRsaEncGenerated(testCtx, rsa.RealmId, rsa.Id)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testKeycloakRealmKeystoreRsaEncGenerated_basic(fullNameKeystoreName),
				Check:  testAccCheckRealmKeystoreRsaEncGeneratedFetch("keycloak_realm_keystore_rsa_enc_generated.realm_rsa_enc", rsa),
			},
		},
	})
}

func TestAccKeycloakRealmKeystoreRsaEncGenerated_keySizeValidation(t *testing.T) {
	t.Parallel()

	rsaName := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckRealmKeystoreRsaEncGeneratedDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakRealmKeystoreRsaEncGenerated_basicWithAttrValidation(rsaName, "key_size",
					strconv.Itoa(acctest.RandIntRange(0, 1000)*2+1)),
				ExpectError: regexp.MustCompile("expected key_size to be one of .+ got .+"),
			},
			{
				Config: testKeycloakRealmKeystoreRsaEncGenerated_basicWithAttrValidation(rsaName, "key_size", "2048"),
				Check:  testAccCheckRealmKeystoreRsaEncGeneratedExists("keycloak_realm_keystore_rsa_enc_generated.realm_rsa_enc"),
			},
		},
	})
}

func TestAccKeycloakRealmKeystoreRsaEncGenerated_algorithmValidation(t *testing.T) {
	t.Parallel()

	algorithm := randomStringInSlice(keycloakRealmKeystoreRsaEncGeneratedAlgorithm)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckRealmKeystoreRsaEncGeneratedDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakRealmKeystoreRsaEncGenerated_basicWithAttrValidation(algorithm, "algorithm",
					acctest.RandString(10)),
				ExpectError: regexp.MustCompile("expected algorithm to be one of .+ got .+"),
			},
			{
				Config: testKeycloakRealmKeystoreRsaEncGenerated_basicWithAttrValidation(algorithm, "algorithm", algorithm),
				Check:  testAccCheckRealmKeystoreRsaEncGeneratedExists("keycloak_realm_keystore_rsa_enc_generated.realm_rsa_enc"),
			},
		},
	})
}

func TestAccKeycloakRealmKeystoreRsaEncGenerated_updateRsaKeystoreGenerated(t *testing.T) {
	t.Parallel()

	enabled := randomBool()
	active := randomBool()

	groupKeystoreOne := &keycloak.RealmKeystoreRsaEncGenerated{
		Name:      acctest.RandString(10),
		RealmId:   testAccRealmUserFederation.Realm,
		Enabled:   enabled,
		Active:    active,
		Priority:  acctest.RandIntRange(0, 100),
		KeySize:   1024,
		Algorithm: randomStringInSlice(keycloakRealmKeystoreRsaEncGeneratedAlgorithm),
	}

	groupKeystoreTwo := &keycloak.RealmKeystoreRsaEncGenerated{
		Name:      acctest.RandString(10),
		RealmId:   testAccRealmUserFederation.Realm,
		Enabled:   enabled,
		Active:    active,
		Priority:  acctest.RandIntRange(0, 100),
		KeySize:   2048,
		Algorithm: randomStringInSlice(keycloakRealmKeystoreRsaEncGeneratedAlgorithm),
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckRealmKeystoreRsaEncGeneratedDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakRealmKeystoreRsaEncGenerated_basicFromInterface(groupKeystoreOne),
				Check:  testAccCheckRealmKeystoreRsaEncGeneratedExists("keycloak_realm_keystore_rsa_enc_generated.realm_rsa_enc"),
			},
			{
				Config: testKeycloakRealmKeystoreRsaEncGenerated_basicFromInterface(groupKeystoreTwo),
				Check:  testAccCheckRealmKeystoreRsaEncGeneratedExists("keycloak_realm_keystore_rsa_enc_generated.realm_rsa_enc"),
			},
		},
	})
}

func testAccCheckRealmKeystoreRsaEncGeneratedExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := getKeycloakRealmKeystoreRsaEncGeneratedFromState(s, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckRealmKeystoreRsaEncGeneratedFetch(resourceName string, keystore *keycloak.RealmKeystoreRsaEncGenerated) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		fetchedKeystore, err := getKeycloakRealmKeystoreRsaEncGeneratedFromState(s, resourceName)
		if err != nil {
			return err
		}

		keystore.Id = fetchedKeystore.Id
		keystore.RealmId = fetchedKeystore.RealmId

		return nil
	}
}

func testAccCheckRealmKeystoreRsaEncGeneratedDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_realm_keystore_rsa_enc_generated" {
				continue
			}

			id := rs.Primary.ID
			realm := rs.Primary.Attributes["realm_id"]

			ldapGroupKeystore, _ := keycloakClient.GetRealmKeystoreRsaEncGenerated(testCtx, realm, id)
			if ldapGroupKeystore != nil {
				return fmt.Errorf("rsa keystore with id %s still exists", id)
			}
		}

		return nil
	}
}

func getKeycloakRealmKeystoreRsaEncGeneratedFromState(s *terraform.State,
	resourceName string) (*keycloak.RealmKeystoreRsaEncGenerated,
	error,
) {
	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	id := rs.Primary.ID
	realm := rs.Primary.Attributes["realm_id"]

	realmKeystore, err := keycloakClient.GetRealmKeystoreRsaEncGenerated(testCtx, realm, id)
	if err != nil {
		return nil, fmt.Errorf("error getting rsa keystore with id %s: %s", id, err)
	}

	return realmKeystore, nil
}

func testKeycloakRealmKeystoreRsaEncGenerated_basic(rsaName string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_realm_keystore_rsa_enc_generated" "realm_rsa_enc" {
	name      = "%s"
	realm_id  = data.keycloak_realm.realm.id

    priority  = 100
    algorithm = "RSA-OAEP"
}
	`, testAccRealmUserFederation.Realm, rsaName)
}

func testKeycloakRealmKeystoreRsaEncGenerated_basicWithAttrValidation(rsaName, attr, val string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_realm_keystore_rsa_enc_generated" "realm_rsa_enc" {
	name      = "%s"
	realm_id  = data.keycloak_realm.realm.id

	%s        = "%s"
}
	`, testAccRealmUserFederation.Realm, rsaName, attr, val)
}

func testKeycloakRealmKeystoreRsaEncGenerated_basicFromInterface(keystore *keycloak.RealmKeystoreRsaEncGenerated) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_realm_keystore_rsa_enc_generated" "realm_rsa_enc" {
	name      = "%s"
	realm_id  = data.keycloak_realm.realm.id

    priority  = %s
    algorithm = "%s"
    key_size  = %s
}
	`, testAccRealmUserFederation.Realm, keystore.Name, strconv.Itoa(keystore.Priority), keystore.Algorithm,
		strconv.Itoa(keystore.KeySize))
}
