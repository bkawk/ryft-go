package ryft

import (
	"context"
)

type PersonsService struct {
	client *Client
}

type Person struct {
	ID        string            `json:"id"`
	FirstName string            `json:"firstName,omitempty"`
	LastName  string            `json:"lastName,omitempty"`
	Email     string            `json:"email,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

type PersonList struct {
	Items []Person `json:"items"`
}

type CreatePersonRequest struct {
	FirstName     string            `json:"firstName"`
	LastName      string            `json:"lastName"`
	Email         string            `json:"email"`
	DateOfBirth   string            `json:"dateOfBirth"`
	Gender        string            `json:"gender"`
	Nationalities []string          `json:"nationalities"`
	Address       Address           `json:"address"`
	PhoneNumber   string            `json:"phoneNumber"`
	BusinessRoles []string          `json:"businessRoles"`
	Documents     []map[string]any  `json:"documents"`
	Metadata      map[string]string `json:"metadata,omitempty"`
}

type UpdatePersonRequest struct {
	FirstName      string            `json:"firstName,omitempty"`
	MiddleNames    string            `json:"middleNames,omitempty"`
	LastName       string            `json:"lastName,omitempty"`
	Email          string            `json:"email,omitempty"`
	DateOfBirth    string            `json:"dateOfBirth,omitempty"`
	CountryOfBirth string            `json:"countryOfBirth,omitempty"`
	Gender         string            `json:"gender,omitempty"`
	Nationalities  []string          `json:"nationalities,omitempty"`
	Address        *Address          `json:"address,omitempty"`
	PhoneNumber    string            `json:"phoneNumber,omitempty"`
	BusinessRoles  []string          `json:"businessRoles,omitempty"`
	Documents      []map[string]any  `json:"documents,omitempty"`
	Metadata       map[string]string `json:"metadata,omitempty"`
}

func (s *PersonsService) Create(ctx context.Context, accountID string, request CreatePersonRequest) (*Person, error) {
	req, err := s.client.newRequest(ctx, "POST", "accounts/"+accountID+"/persons", request)
	if err != nil {
		return nil, err
	}

	var person Person
	if err := s.client.doJSON(req, &person); err != nil {
		return nil, err
	}

	return &person, nil
}

func (s *PersonsService) Get(ctx context.Context, accountID string, personID string) (*Person, error) {
	req, err := s.client.newRequest(ctx, "GET", "accounts/"+accountID+"/persons/"+personID, nil)
	if err != nil {
		return nil, err
	}

	var person Person
	if err := s.client.doJSON(req, &person); err != nil {
		return nil, err
	}

	return &person, nil
}

func (s *PersonsService) List(
	ctx context.Context,
	accountID string,
	ascending bool,
	limit int,
	startsAfter string,
) (*PersonList, error) {
	query := buildListQuery(ascending, limit, startsAfter)

	req, err := s.client.newRequestWithQuery(ctx, "GET", "accounts/"+accountID+"/persons", query, nil)
	if err != nil {
		return nil, err
	}

	var people PersonList
	if err := s.client.doJSON(req, &people); err != nil {
		return nil, err
	}

	return &people, nil
}

func (s *PersonsService) Update(ctx context.Context, accountID string, personID string, request UpdatePersonRequest) (*Person, error) {
	req, err := s.client.newRequest(ctx, "PATCH", "accounts/"+accountID+"/persons/"+personID, request)
	if err != nil {
		return nil, err
	}

	var person Person
	if err := s.client.doJSON(req, &person); err != nil {
		return nil, err
	}

	return &person, nil
}

func (s *PersonsService) Delete(ctx context.Context, accountID string, personID string) (*DeletedResource, error) {
	req, err := s.client.newRequest(ctx, "DELETE", "accounts/"+accountID+"/persons/"+personID, nil)
	if err != nil {
		return nil, err
	}

	var deleted DeletedResource
	if err := s.client.doJSON(req, &deleted); err != nil {
		return nil, err
	}

	return &deleted, nil
}
