package keycloak

import (
	"context"
	"fmt"
)

type SamlHardcodedAttributeProtocolMapper struct {
	Id            string
	Name          string
	RealmId       string
	ClientId      string
	ClientScopeId string

	AttributeValue          string
	FriendlyName            string
	SamlAttributeName       string
	SamlAttributeNameFormat string
}

func (mapper *SamlHardcodedAttributeProtocolMapper) convertToGenericProtocolMapper() *protocolMapper {
	return &protocolMapper{
		Id:             mapper.Id,
		Name:           mapper.Name,
		Protocol:       "saml",
		ProtocolMapper: "saml-hardcode-attribute-mapper",
		Config: map[string]string{
			attributeNameField:       mapper.SamlAttributeName,
			attributeNameFormatField: mapper.SamlAttributeNameFormat,
			friendlyNameField:        mapper.FriendlyName,
			attributeValueField:      mapper.AttributeValue,
		},
	}
}

func (protocolMapper *protocolMapper) convertToSamlHardcodedAttributeProtocolMapper(realmId, clientId, clientScopeId string) *SamlHardcodedAttributeProtocolMapper {
	return &SamlHardcodedAttributeProtocolMapper{
		Id:            protocolMapper.Id,
		Name:          protocolMapper.Name,
		RealmId:       realmId,
		ClientId:      clientId,
		ClientScopeId: clientScopeId,

		AttributeValue:          protocolMapper.Config[attributeValueField],
		FriendlyName:            protocolMapper.Config[friendlyNameField],
		SamlAttributeName:       protocolMapper.Config[attributeNameField],
		SamlAttributeNameFormat: protocolMapper.Config[attributeNameFormatField],
	}
}

func (keycloakClient *KeycloakClient) GetSamlHardcodedAttributeProtocolMapper(ctx context.Context, realmId, clientId, clientScopeId, mapperId string) (*SamlHardcodedAttributeProtocolMapper, error) {
	var protocolMapper *protocolMapper

	err := keycloakClient.get(ctx, individualProtocolMapperPath(realmId, clientId, clientScopeId, mapperId), &protocolMapper, nil)
	if err != nil {
		return nil, err
	}

	return protocolMapper.convertToSamlHardcodedAttributeProtocolMapper(realmId, clientId, clientScopeId), nil
}

func (keycloakClient *KeycloakClient) DeleteSamlHardcodedAttributeProtocolMapper(ctx context.Context, realmId, clientId, clientScopeId, mapperId string) error {
	return keycloakClient.delete(ctx, individualProtocolMapperPath(realmId, clientId, clientScopeId, mapperId), nil)
}

func (keycloakClient *KeycloakClient) NewSamlHardcodedAttributeProtocolMapper(ctx context.Context, mapper *SamlHardcodedAttributeProtocolMapper) error {
	path := protocolMapperPath(mapper.RealmId, mapper.ClientId, mapper.ClientScopeId)

	_, location, err := keycloakClient.post(ctx, path, mapper.convertToGenericProtocolMapper())
	if err != nil {
		return err
	}

	mapper.Id = getIdFromLocationHeader(location)

	return nil
}

func (keycloakClient *KeycloakClient) UpdateSamlHardcodedAttributeProtocolMapper(ctx context.Context, mapper *SamlHardcodedAttributeProtocolMapper) error {
	path := individualProtocolMapperPath(mapper.RealmId, mapper.ClientId, mapper.ClientScopeId, mapper.Id)

	return keycloakClient.put(ctx, path, mapper.convertToGenericProtocolMapper())
}

func (keycloakClient *KeycloakClient) ValidateSamlHardcodedAttributeProtocolMapper(ctx context.Context, mapper *SamlHardcodedAttributeProtocolMapper) error {
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
