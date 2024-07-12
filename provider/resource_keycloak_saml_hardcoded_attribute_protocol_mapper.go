package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/lucdew/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakSamlHardcodedAttributeProtocolMapper() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKeycloakSamlHardcodedAttributeProtocolMapperCreate,
		ReadContext:   resourceKeycloakSamlHardcodedAttributeProtocolMapperRead,
		UpdateContext: resourceKeycloakSamlHardcodedAttributeProtocolMapperUpdate,
		DeleteContext: resourceKeycloakSamlHardcodedAttributeProtocolMapperDelete,
		Importer: &schema.ResourceImporter{
			// import a mapper tied to a client:
			// {{realmId}}/client/{{clientId}}/{{protocolMapperId}}
			// or a client scope:
			// {{realmId}}/client-scope/{{clientScopeId}}/{{protocolMapperId}}
			StateContext: genericProtocolMapperImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"client_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"client_scope_id"},
			},
			"client_scope_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"client_id"},
			},
			"attribute_value": {
				Type:     schema.TypeString,
				Required: true,
			},
			"friendly_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"saml_attribute_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"saml_attribute_name_format": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(keycloakSamlUserAttributeProtocolMapperNameFormats, false),
			},
		},
	}
}

func mapFromDataToSamlHardcodedAttributeProtocolMapper(data *schema.ResourceData) *keycloak.SamlHardcodedAttributeProtocolMapper {
	return &keycloak.SamlHardcodedAttributeProtocolMapper{
		Id:            data.Id(),
		Name:          data.Get("name").(string),
		RealmId:       data.Get("realm_id").(string),
		ClientId:      data.Get("client_id").(string),
		ClientScopeId: data.Get("client_scope_id").(string),

		AttributeValue:          data.Get("attribute_value").(string),
		FriendlyName:            data.Get("friendly_name").(string),
		SamlAttributeName:       data.Get("saml_attribute_name").(string),
		SamlAttributeNameFormat: data.Get("saml_attribute_name_format").(string),
	}
}

func mapFromSamlHardcodedAttributeProtocolMapperToData(mapper *keycloak.SamlHardcodedAttributeProtocolMapper, data *schema.ResourceData) {
	data.SetId(mapper.Id)
	data.Set("name", mapper.Name)
	data.Set("realm_id", mapper.RealmId)

	if mapper.ClientId != "" {
		data.Set("client_id", mapper.ClientId)
	} else {
		data.Set("client_scope_id", mapper.ClientScopeId)
	}

	data.Set("attribute_value", mapper.AttributeValue)
	data.Set("friendly_name", mapper.FriendlyName)
	data.Set("saml_attribute_name", mapper.SamlAttributeName)
	data.Set("saml_attribute_name_format", mapper.SamlAttributeNameFormat)
}

func resourceKeycloakSamlHardcodedAttributeProtocolMapperCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	samlHardcodedAttributeMapper := mapFromDataToSamlHardcodedAttributeProtocolMapper(data)

	err := keycloakClient.ValidateSamlHardcodedAttributeProtocolMapper(ctx, samlHardcodedAttributeMapper)
	if err != nil {
		return diag.FromErr(err)
	}

	err = keycloakClient.NewSamlHardcodedAttributeProtocolMapper(ctx, samlHardcodedAttributeMapper)
	if err != nil {
		return diag.FromErr(err)
	}

	mapFromSamlHardcodedAttributeProtocolMapperToData(samlHardcodedAttributeMapper, data)

	return resourceKeycloakSamlHardcodedAttributeProtocolMapperRead(ctx, data, meta)
}

func resourceKeycloakSamlHardcodedAttributeProtocolMapperRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)
	clientScopeId := data.Get("client_scope_id").(string)

	samlHardcodedAttributeMapper, err := keycloakClient.GetSamlHardcodedAttributeProtocolMapper(ctx, realmId, clientId, clientScopeId, data.Id())
	if err != nil {
		return handleNotFoundError(ctx, err, data)
	}

	mapFromSamlHardcodedAttributeProtocolMapperToData(samlHardcodedAttributeMapper, data)

	return nil
}

func resourceKeycloakSamlHardcodedAttributeProtocolMapperUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	samlHardcodedAttributeMapper := mapFromDataToSamlHardcodedAttributeProtocolMapper(data)

	err := keycloakClient.ValidateSamlHardcodedAttributeProtocolMapper(ctx, samlHardcodedAttributeMapper)
	if err != nil {
		return diag.FromErr(err)
	}

	err = keycloakClient.UpdateSamlHardcodedAttributeProtocolMapper(ctx, samlHardcodedAttributeMapper)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceKeycloakSamlHardcodedAttributeProtocolMapperRead(ctx, data, meta)
}

func resourceKeycloakSamlHardcodedAttributeProtocolMapperDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)
	clientScopeId := data.Get("client_scope_id").(string)

	return diag.FromErr(keycloakClient.DeleteSamlHardcodedAttributeProtocolMapper(ctx, realmId, clientId, clientScopeId, data.Id()))
}
