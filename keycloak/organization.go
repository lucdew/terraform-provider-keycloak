package keycloak

import (
	"context"
	"fmt"
	"strings"
)

type OrganizationDomain struct {
	Name     string `json:"name,omitempty"`
	Verified bool   `json:"verified,omitempty"`
}

type Organization struct {
	Id          string               `json:"id,omitempty"`
	RealmId     string               `json:"-"`
	Name        string               `json:"name"`
	RedirectUrl string               `json:"redirectUrl,omitempty"`
	Description string               `json:"description,omitempty"`
	Domains     []OrganizationDomain `json:"domains,omitempty"`
	Attributes  map[string][]string  `json:"attributes,omitempty"`
}

func (keycloakClient *KeycloakClient) CreateOrganization(ctx context.Context, organization *Organization) error {
	path := fmt.Sprintf("/realms/%s/organizations", organization.RealmId)

	_, location, err := keycloakClient.post(ctx, path, organization)
	if err != nil {
		return err
	}

	// Extract ID from location URL
	parts := strings.Split(location, "/")
	organization.Id = parts[len(parts)-1]

	return nil
}

func (keycloakClient *KeycloakClient) GetOrganization(ctx context.Context, realmId, id string) (*Organization, error) {
	organization := Organization{}
	organization.RealmId = realmId

	path := fmt.Sprintf("/realms/%s/organizations/%s", realmId, id)
	err := keycloakClient.get(ctx, path, &organization, nil)
	if err != nil {
		return nil, err
	}

	return &organization, nil
}

func (keycloakClient *KeycloakClient) UpdateOrganization(ctx context.Context, organization *Organization) error {
	path := fmt.Sprintf("/realms/%s/organizations/%s", organization.RealmId, organization.Id)
	return keycloakClient.put(ctx, path, organization)
}

func (keycloakClient *KeycloakClient) DeleteOrganization(ctx context.Context, realmId, id string) error {
	path := fmt.Sprintf("/realms/%s/organizations/%s", realmId, id)
	return keycloakClient.delete(ctx, path, nil)
}
