package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/keycloak/terraform-provider-keycloak/keycloak"
	"github.com/keycloak/terraform-provider-keycloak/keycloak/types"
)

func resourceKeycloakOrganizationIdentityProvider() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKeycloakOrganizationIdentityProviderCreate,
		ReadContext:   resourceKeycloakOrganizationIdentityProviderRead,
		UpdateContext: resourceKeycloakOrganizationIdentityProviderUpdate,
		DeleteContext: resourceKeycloakOrganizationIdentityProviderDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceKeycloakOrganizationIdentityProviderImport,
		},
		Schema: map[string]*schema.Schema{
			"realm_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"organization_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"identity_provider_alias": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"domain": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"redirect_email_domain_matches": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceKeycloakOrganizationIdentityProviderCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId, orgId, idpAlias := getOrganizationIdentityProviderFromData(data)

	idp, err := keycloakClient.GetIdentityProvider(ctx, realmId, idpAlias)
	if err != nil {
		return diag.FromErr(err)
	}

	idp.Config.OrgDomain = data.Get("domain").(string)
	idp.Config.OrgRedirectEmailMatches = types.KeycloakBoolQuoted(data.Get("redirect_email_domain_matches").(bool))
	keycloakClient.UpdateIdentityProvider(ctx, idp)

	err = keycloakClient.LinkIdentityProviderToOrganization(ctx, realmId, orgId, idpAlias)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(fmt.Sprintf("%s/%s/%s", realmId, orgId, idpAlias))

	return resourceKeycloakOrganizationIdentityProviderRead(ctx, data, meta)
}

func resourceKeycloakOrganizationIdentityProviderRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId, orgId, idpAlias := getOrganizationIdentityProviderFromData(data)

	err := keycloakClient.CheckIdentityProviderLinkToOrganization(ctx, realmId, orgId, idpAlias)
	if err != nil {
		return handleNotFoundError(ctx, err, data)
	}

	idp, err := keycloakClient.GetIdentityProvider(ctx, realmId, idpAlias)
	if err != nil {
		return diag.FromErr(err)
	}
	data.Set("domain", idp.Config.OrgDomain)
	data.Set("redirect_email_domain_matches", idp.Config.OrgRedirectEmailMatches)

	return nil
}

func resourceKeycloakOrganizationIdentityProviderUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId, orgId, idpAlias := getOrganizationIdentityProviderFromData(data)

	err := keycloakClient.CheckIdentityProviderLinkToOrganization(ctx, realmId, orgId, idpAlias)
	if err != nil {
		return handleNotFoundError(ctx, err, data)
	}

	idp, err := keycloakClient.GetIdentityProvider(ctx, realmId, idpAlias)
	if err != nil {
		return diag.FromErr(err)
	}
	idp.Config.OrgDomain = data.Get("domain").(string)
	idp.Config.OrgRedirectEmailMatches = types.KeycloakBoolQuoted(data.Get("redirect_email_domain_matches").(bool))
	keycloakClient.UpdateIdentityProvider(ctx, idp)

	return nil
}

func resourceKeycloakOrganizationIdentityProviderDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId, orgId, idpAlias := getOrganizationIdentityProviderFromData(data)

	err := keycloakClient.UnlinkIdentityProviderToOrganization(ctx, realmId, orgId, idpAlias)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceKeycloakOrganizationIdentityProviderImport(ctx context.Context, data *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(data.Id(), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid import. Supported format: {{realm}}/{{organizationId}}/{{identityProviderAlias}}")
	}

	data.Set("realm_id", parts[0])
	data.Set("organization_id", parts[1])
	data.Set("identity_provider_alias", parts[2])

	return []*schema.ResourceData{data}, nil
}

func getOrganizationIdentityProviderFromData(data *schema.ResourceData) (realmId, orgId, idpAlias string) {
	realmId = data.Get("realm_id").(string)
	orgId = data.Get("organization_id").(string)
	idpAlias = data.Get("identity_provider_alias").(string)

	return
}
