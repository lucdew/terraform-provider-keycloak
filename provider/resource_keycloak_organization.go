package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/keycloak/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakOrganization() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKeycloakOrganizationCreate,
		ReadContext:   resourceKeycloakOrganizationRead,
		UpdateContext: resourceKeycloakOrganizationUpdate,
		DeleteContext: resourceKeycloakOrganizationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceKeycloakOrganizationImport,
		},
		Schema: map[string]*schema.Schema{
			"realm_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"alias": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"redirect_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"domain": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"verified": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
				// Custom validation function to ensure domain names are unique
				//ValidateFunc: validateUniqueDomainNames,
			},
			"attributes": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceKeycloakOrganizationCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	organization := getOrganizationFromData(data)

	err := keycloakClient.CreateOrganization(ctx, organization)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(organization.Id)

	return resourceKeycloakOrganizationRead(ctx, data, meta)
}

func resourceKeycloakOrganizationRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	organization, err := keycloakClient.GetOrganization(ctx, realmId, id)
	if err != nil {
		return handleNotFoundError(ctx, err, data)
	}

	setOrganizationData(data, organization)

	return nil
}

func resourceKeycloakOrganizationUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	organization := getOrganizationFromData(data)

	err := keycloakClient.UpdateOrganization(ctx, organization)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceKeycloakOrganizationRead(ctx, data, meta)
}

func resourceKeycloakOrganizationDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	err := keycloakClient.DeleteOrganization(ctx, realmId, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceKeycloakOrganizationImport(ctx context.Context, data *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(data.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid import. Supported format: {{realm}}/{{organizationId}}")
	}

	data.Set("realm_id", parts[0])
	data.SetId(parts[1])

	return []*schema.ResourceData{data}, nil
}

func getOrganizationFromData(data *schema.ResourceData) *keycloak.Organization {

	return &keycloak.Organization{
		Id:          data.Id(),
		RealmId:     data.Get("realm_id").(string),
		Name:        data.Get("name").(string),
		Alias:       data.Get("alias").(string),
		Enabled:     data.Get("enabled").(bool),
		RedirectUrl: data.Get("redirect_url").(string),
		Description: data.Get("description").(string),
		Domains:     expandOrganizationDomains(data.Get("domain").(*schema.Set)),
		Attributes:  expandStringMap(data.Get("attributes").(map[string]interface{})),
	}

}

func setOrganizationData(data *schema.ResourceData, organization *keycloak.Organization) {
	data.Set("realm_id", organization.RealmId)
	data.Set("name", organization.Name)
	data.Set("alias", organization.Alias)
	data.Set("enabled", organization.Enabled)
	data.Set("redirect_url", organization.RedirectUrl)
	data.Set("description", organization.Description)
	data.Set("domain", flattenOrganizationDomains(organization.Domains))
	data.Set("attributes", flattenStringMap(organization.Attributes))
}

// validateUniqueDomainNames ensures that all domain names within the set are unique
func validateUniqueDomainNames(v interface{}, k string) (warnings []string, errors []error) {

	domains := v.(*schema.Set).List()

	// Create a map to track domain names
	domainNames := make(map[string]bool)

	for _, domain := range domains {
		domainMap := domain.(map[string]interface{})
		domainName := domainMap["name"].(string)

		if domainName == "" {
			continue // Skip empty domain names
		}

		// Check if this domain name is already in use
		if _, exists := domainNames[domainName]; exists {
			errors = append(errors, fmt.Errorf("duplicate domain name found: %s. Domain names must be unique", domainName))
		} else {
			domainNames[domainName] = true
		}
	}

	return warnings, errors
}

// Helper functions for domain handling
func expandOrganizationDomains(set *schema.Set) []keycloak.OrganizationDomain {
	domains := make([]keycloak.OrganizationDomain, 0, set.Len())

	for _, value := range set.List() {
		domainMap := value.(map[string]interface{})
		domain := keycloak.OrganizationDomain{
			Name:     domainMap["name"].(string),
			Verified: domainMap["verified"].(bool),
		}
		domains = append(domains, domain)
	}

	return domains
}

func flattenOrganizationDomains(domains []keycloak.OrganizationDomain) *schema.Set {
	set := schema.NewSet(schema.HashResource(resourceKeycloakOrganization().Schema["domain"].Elem.(*schema.Resource)), []interface{}{})

	for _, domain := range domains {
		domainMap := map[string]interface{}{
			"name":     domain.Name,
			"verified": domain.Verified,
		}
		set.Add(domainMap)
	}

	return set
}

// Helper for expanding a map to a string-to-string-slice map
func expandStringMap(m map[string]interface{}) map[string][]string {
	result := make(map[string][]string)
	for k, v := range m {
		result[k] = []string{v.(string)}
	}
	return result
}

// Helper for flattening a string-to-string-slice map to a regular map
func flattenStringMap(m map[string][]string) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range m {
		if len(v) > 0 {
			result[k] = v[0]
		}
	}
	return result
}

func importStringParse(importString string) []string {
	parts := strings.Split(importString, "/")
	return parts
}
