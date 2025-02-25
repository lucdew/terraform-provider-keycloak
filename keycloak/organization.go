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
	Alias       string               `json:"alias,omitempty"`
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

// GetOrganizationsPath returns the URL for the organization API endpoint
func (keycloakClient *KeycloakClient) GetOrganizationsPath(realmId string) string {
	return fmt.Sprintf("/realms/%s/organizations", realmId)
}

// GetOrganizations gets all organizations in a realm
// This function can be used to list all organizations or filter by search criteria
func (keycloakClient *KeycloakClient) GetOrganizations(ctx context.Context, realmId string, params ...map[string]string) ([]*Organization, error) {
	var organizations []*Organization
	queryParams := make(map[string]string)

	// Apply optional filter parameters
	if len(params) > 0 && params[0] != nil {
		for k, v := range params[0] {
			queryParams[k] = v
		}
	}

	err := keycloakClient.get(ctx, keycloakClient.GetOrganizationsPath(realmId), &organizations, queryParams)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch organizations: %w", err)
	}

	// Set RealmId for each organization
	for _, organization := range organizations {
		organization.RealmId = realmId
	}

	return organizations, nil
}

// GetOrganizationsPaginated gets organizations with pagination support
func (keycloakClient *KeycloakClient) GetOrganizationsPaginated(ctx context.Context, realmId string, first, max int, search string) ([]*Organization, error) {
	queryParams := map[string]string{
		"first": fmt.Sprintf("%d", first),
		"max":   fmt.Sprintf("%d", max),
	}

	if search != "" {
		queryParams["search"] = search
	}

	return keycloakClient.GetOrganizations(ctx, realmId, queryParams)
}

// GetOrganizationByName gets an organization by name
func (keycloakClient *KeycloakClient) GetOrganizationByName(ctx context.Context, realmId, name string) (*Organization, error) {
	var organizations []*Organization

	params := map[string]string{
		"search": name,
		"exact":  "true",
	}

	err := keycloakClient.get(ctx, keycloakClient.GetOrganizationsPath(realmId), &organizations, params)
	if err != nil {
		return nil, err
	}

	// Find the exact match by name
	for _, organization := range organizations {
		if organization.Name == name {
			organization.RealmId = realmId
			return organization, nil
		}
	}

	return nil, nil
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
