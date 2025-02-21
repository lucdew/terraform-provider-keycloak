package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccKeycloakDataSourceGroups_full_hierarchy(t *testing.T) {
	t.Parallel()

	groupPrefix := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakRoleDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDataSourceKeycloakGroups_full_hierarchy(groupPrefix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakGroupExists("keycloak_group.group_0_a"),
					testAccCheckKeycloakGroupExists("keycloak_group.group_0_a_1_a_2_b"),
					testAccCheckDataKeycloakGroups(groupPrefix, "data.keycloak_groups.all_groups"),
				),
			},
		},
	})
}

func testAccCheckDataKeycloakGroups(groupPrefix string, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		id := rs.Primary.ID
		if id == "" {
			return fmt.Errorf("no id for groups")
		}

		rs_group, ok := s.RootModule().Resources["keycloak_group.group_0_a_1_a_2_a"]
		if ok {
			fmt.Printf("Found group_0_a_1_a_2_a in resources=%v", rs_group.Primary.Attributes)
		} else {
			fmt.Println("Not Found group_0_a_1_a_2_a in resources")
		}

		// debug failed tests in CI
		rs_realm, ok := s.RootModule().Resources["data.keycloak_realm.realm"]
		if !ok {
			return fmt.Errorf("resource not found: %s", "data.keycloak_realm.realm")
		}
		fmt.Printf("realmid=%s, groups=%v", rs_realm.Primary.Attributes["id"], rs.Primary.Attributes)

		if len(rs.Primary.Attributes["groups.#"]) == 0 {
			return fmt.Errorf("no groups exist")
		}

		name_group_0_a := rs.Primary.Attributes["groups.0.name"]
		if name_group_0_a != fmt.Sprintf("%s_0_a", groupPrefix) {
			return fmt.Errorf("group %s_0_a is missing", groupPrefix)
		}

		name_group_0_a_1_a_2_b := rs.Primary.Attributes["groups.4.name"]
		if name_group_0_a_1_a_2_b != fmt.Sprintf("%s_0_a_1_a_2_b", groupPrefix) {
			return fmt.Errorf("%s_0_a_1_a_2_b is missing", groupPrefix)
		}

		path_group_0_a_1_a_2_b := rs.Primary.Attributes["groups.4.path"]
		if path_group_0_a_1_a_2_b != fmt.Sprintf("/%s_0_a/%s_0_a_1_a/%s_0_a_1_a_2_b", groupPrefix, groupPrefix, groupPrefix) {
			return fmt.Errorf("%s_0_a_1_a_2_b path is invalid got %s", groupPrefix, path_group_0_a_1_a_2_b)
		}

		return nil
	}
}

func testDataSourceKeycloakGroups_full_hierarchy(group string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_group" "group_0_a" {
	name     = "%s_0_a"
	realm_id = data.keycloak_realm.realm.id
}

resource "keycloak_group" "group_0_b" {
	name     = "%s_0_b"
	realm_id = data.keycloak_realm.realm.id
}

resource "keycloak_group" "group_0_a_1_a" {
	name     	= "%s_0_a_1_a"
	parent_id = keycloak_group.group_0_a.id
	realm_id 	= data.keycloak_realm.realm.id
}

resource "keycloak_group" "group_0_a_1_a_2_a" {
	name     	= "%s_0_a_1_a_2_a"
	parent_id = keycloak_group.group_0_a_1_a.id
	realm_id 	= data.keycloak_realm.realm.id
}

resource "keycloak_group" "group_0_a_1_a_2_b" {
	name     	= "%s_0_a_1_a_2_b"
	parent_id = keycloak_group.group_0_a_1_a.id
	realm_id 	= data.keycloak_realm.realm.id
}

data "keycloak_groups" "all_groups" {
	realm_id = data.keycloak_realm.realm.id
	full_hierarchy= true

	depends_on = [keycloak_group.group_0_b, keycloak_group.group_0_a_1_a_2_a, keycloak_group.group_0_a_1_a_2_b]
}

	`, testAccRealmAllGroups.Realm, group, group, group, group, group)
}
