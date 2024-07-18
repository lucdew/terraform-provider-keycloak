package keycloak

import (
	"context"
	"fmt"
)

type SamlUserSessionNoteProtocolMapper struct {
	Id            string
	Name          string
	RealmId       string
	ClientId      string
	ClientScopeId string

	NoteName                string
	FriendlyName            string
	SamlAttributeName       string
	SamlAttributeNameFormat string
}

func (mapper *SamlUserSessionNoteProtocolMapper) convertToGenericProtocolMapper() *protocolMapper {
	return &protocolMapper{
		Id:             mapper.Id,
		Name:           mapper.Name,
		Protocol:       "saml",
		ProtocolMapper: "saml-user-session-note-mapper",
		Config: map[string]string{
			attributeNameField:       mapper.SamlAttributeName,
			attributeNameFormatField: mapper.SamlAttributeNameFormat,
			friendlyNameField:        mapper.FriendlyName,
			noteAttributeField:       mapper.NoteName,
		},
	}
}

func (protocolMapper *protocolMapper) convertToSamlUserSessionNoteProtocolMapper(realmId, clientId, clientScopeId string) *SamlUserSessionNoteProtocolMapper {
	return &SamlUserSessionNoteProtocolMapper{
		Id:            protocolMapper.Id,
		Name:          protocolMapper.Name,
		RealmId:       realmId,
		ClientId:      clientId,
		ClientScopeId: clientScopeId,

		NoteName:                protocolMapper.Config[noteAttributeField],
		FriendlyName:            protocolMapper.Config[friendlyNameField],
		SamlAttributeName:       protocolMapper.Config[attributeNameField],
		SamlAttributeNameFormat: protocolMapper.Config[attributeNameFormatField],
	}
}

func (keycloakClient *KeycloakClient) GetSamlUserSessionNoteProtocolMapper(ctx context.Context, realmId, clientId, clientScopeId, mapperId string) (*SamlUserSessionNoteProtocolMapper, error) {
	var protocolMapper *protocolMapper

	err := keycloakClient.get(ctx, individualProtocolMapperPath(realmId, clientId, clientScopeId, mapperId), &protocolMapper, nil)
	if err != nil {
		return nil, err
	}

	return protocolMapper.convertToSamlUserSessionNoteProtocolMapper(realmId, clientId, clientScopeId), nil
}

func (keycloakClient *KeycloakClient) DeleteSamlUserSessionNoteProtocolMapper(ctx context.Context, realmId, clientId, clientScopeId, mapperId string) error {
	return keycloakClient.delete(ctx, individualProtocolMapperPath(realmId, clientId, clientScopeId, mapperId), nil)
}

func (keycloakClient *KeycloakClient) NewSamlUserSessionNoteProtocolMapper(ctx context.Context, mapper *SamlUserSessionNoteProtocolMapper) error {
	path := protocolMapperPath(mapper.RealmId, mapper.ClientId, mapper.ClientScopeId)

	_, location, err := keycloakClient.post(ctx, path, mapper.convertToGenericProtocolMapper())
	if err != nil {
		return err
	}

	mapper.Id = getIdFromLocationHeader(location)

	return nil
}

func (keycloakClient *KeycloakClient) UpdateSamlUserSessionNoteProtocolMapper(ctx context.Context, mapper *SamlUserSessionNoteProtocolMapper) error {
	path := individualProtocolMapperPath(mapper.RealmId, mapper.ClientId, mapper.ClientScopeId, mapper.Id)

	return keycloakClient.put(ctx, path, mapper.convertToGenericProtocolMapper())
}

func (keycloakClient *KeycloakClient) ValidateSamlUserSessionNoteProtocolMapper(ctx context.Context, mapper *SamlUserSessionNoteProtocolMapper) error {
	if mapper.ClientId == "" && mapper.ClientScopeId == "" {
		return fmt.Errorf("validation error: one of ClientId or ClientScopeId must be set")
	}

	protocolMappers, err := keycloakClient.listGenericProtocolMappers(ctx, mapper.RealmId, mapper.ClientId, mapper.ClientScopeId)
	if err != nil {
		return err
	}

	for _, protocolMapper := range protocolMappers {
		if protocolMapper.Name == mapper.Name && protocolMapper.Id != mapper.Id {
			return fmt.Errorf("validation error: a protocol mapper with name %s already exists for this client", mapper.Name)
		}
	}

	return nil
}
