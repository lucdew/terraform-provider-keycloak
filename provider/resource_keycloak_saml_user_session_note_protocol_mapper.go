package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/lucdew/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakSamlUserSessionNoteProtocolMapper() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKeycloakSamlUserSessionNoteProtocolMapperCreate,
		ReadContext:   resourceKeycloakSamlUserSessionNoteProtocolMapperRead,
		UpdateContext: resourceKeycloakSamlUserSessionNoteProtocolMapperUpdate,
		DeleteContext: resourceKeycloakSamlUserSessionNoteProtocolMapperDelete,
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
			"note_name": {
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

func mapFromDataToSamlUserSessionNoteProtocolMapper(data *schema.ResourceData) *keycloak.SamlUserSessionNoteProtocolMapper {
	return &keycloak.SamlUserSessionNoteProtocolMapper{
		Id:            data.Id(),
		Name:          data.Get("name").(string),
		RealmId:       data.Get("realm_id").(string),
		ClientId:      data.Get("client_id").(string),
		ClientScopeId: data.Get("client_scope_id").(string),

		NoteName:                data.Get("note_name").(string),
		FriendlyName:            data.Get("friendly_name").(string),
		SamlAttributeName:       data.Get("saml_attribute_name").(string),
		SamlAttributeNameFormat: data.Get("saml_attribute_name_format").(string),
	}
}

func mapFromSamlUserSessionNoteProtocolMapperToData(mapper *keycloak.SamlUserSessionNoteProtocolMapper, data *schema.ResourceData) {
	data.SetId(mapper.Id)
	data.Set("name", mapper.Name)
	data.Set("realm_id", mapper.RealmId)

	if mapper.ClientId != "" {
		data.Set("client_id", mapper.ClientId)
	} else {
		data.Set("client_scope_id", mapper.ClientScopeId)
	}

	data.Set("note_name", mapper.NoteName)
	data.Set("friendly_name", mapper.FriendlyName)
	data.Set("saml_attribute_name", mapper.SamlAttributeName)
	data.Set("saml_attribute_name_format", mapper.SamlAttributeNameFormat)
}

func resourceKeycloakSamlUserSessionNoteProtocolMapperCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	samlUserPropertyMapper := mapFromDataToSamlUserSessionNoteProtocolMapper(data)

	err := keycloakClient.ValidateSamlUserSessionNoteProtocolMapper(ctx, samlUserPropertyMapper)
	if err != nil {
		return diag.FromErr(err)
	}

	err = keycloakClient.NewSamlUserSessionNoteProtocolMapper(ctx, samlUserPropertyMapper)
	if err != nil {
		return diag.FromErr(err)
	}

	mapFromSamlUserSessionNoteProtocolMapperToData(samlUserPropertyMapper, data)

	return resourceKeycloakSamlUserSessionNoteProtocolMapperRead(ctx, data, meta)
}

func resourceKeycloakSamlUserSessionNoteProtocolMapperRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)
	clientScopeId := data.Get("client_scope_id").(string)

	samlUserPropertyMapper, err := keycloakClient.GetSamlUserSessionNoteProtocolMapper(ctx, realmId, clientId, clientScopeId, data.Id())
	if err != nil {
		return handleNotFoundError(ctx, err, data)
	}

	mapFromSamlUserSessionNoteProtocolMapperToData(samlUserPropertyMapper, data)

	return nil
}

func resourceKeycloakSamlUserSessionNoteProtocolMapperUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	samlUserPropertyMapper := mapFromDataToSamlUserSessionNoteProtocolMapper(data)

	err := keycloakClient.ValidateSamlUserSessionNoteProtocolMapper(ctx, samlUserPropertyMapper)
	if err != nil {
		return diag.FromErr(err)
	}

	err = keycloakClient.UpdateSamlUserSessionNoteProtocolMapper(ctx, samlUserPropertyMapper)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceKeycloakSamlUserSessionNoteProtocolMapperRead(ctx, data, meta)
}

func resourceKeycloakSamlUserSessionNoteProtocolMapperDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)
	clientScopeId := data.Get("client_scope_id").(string)

	return diag.FromErr(keycloakClient.DeleteSamlUserSessionNoteProtocolMapper(ctx, realmId, clientId, clientScopeId, data.Id()))
}
