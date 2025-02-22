package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/keycloak/terraform-provider-keycloak/keycloak"
)

func dataSourceKeycloakGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceKeycloakGroupsRead,
		Schema: map[string]*schema.Schema{
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"full_hierarchy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"realm_id": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"subgroup_count": {
							Type:     schema.TypeInt,
							Computed: true,
							Optional: true,
						},
						"path": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"parent_id": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"attributes": {
							Type:     schema.TypeMap,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func flattenGroups(groups []*keycloak.Group) []map[string]interface{} {
	groupsMap := make([]map[string]interface{}, 0)

	for _, group := range groups {
		element := make(map[string]interface{})
		element["id"] = group.Id
		element["name"] = group.Name
		element["realm_id"] = group.RealmId
		element["path"] = group.Path
		element["subgroup_count"] = group.SubGroupCount
		if group.ParentId != "" {
			element["parent_id"] = group.ParentId
		}

		attributes := map[string]string{}
		for k, v := range group.Attributes {
			attributes[k] = strings.Join(v, MULTIVALUE_ATTRIBUTE_SEPARATOR)
		}
		element["attributes"] = attributes

		groupsMap = append(groupsMap, element)
	}

	return groupsMap
}

func setGroupsData(data *schema.ResourceData, groups []*keycloak.Group) error {
	data.SetId(data.Get("realm_id").(string))

	err := data.Set("groups", flattenGroups(groups))
	if err != nil {
		return fmt.Errorf("could not set 'groups' with values '%+v'\n%+v", groups, err)
	}

	return nil
}

func dataSourceKeycloakGroupsRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	fullHierarchy := data.Get("full_hierarchy").(bool)

	groups, err := keycloakClient.GetFlattenedGroupsHierarchy(ctx, realmId, fullHierarchy)
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.FromErr(setGroupsData(data, groups))
}
