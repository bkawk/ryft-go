package ryft

import (
	"context"
	"net/http"
	"strconv"
)

type CustomersService struct {
	client *Client
}

type Customer struct {
	ID        string            `json:"id"`
	Email     string            `json:"email,omitempty"`
	FirstName string            `json:"firstName,omitempty"`
	LastName  string            `json:"lastName,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

type CustomerPaymentMethodList struct {
	Items []PaymentMethod `json:"items"`
}

type CustomerList struct {
	Items []Customer `json:"items"`
}

type CustomerListParams struct {
	ListParams
	Email          string
	StartTimestamp int
	EndTimestamp   int
}

type CreateCustomerRequest struct {
	Email     string            `json:"email"`
	FirstName string            `json:"firstName,omitempty"`
	LastName  string            `json:"lastName,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

type UpdateCustomerRequest struct {
	FirstName string            `json:"firstName,omitempty"`
	LastName  string            `json:"lastName,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

func (s *CustomersService) Create(ctx context.Context, request CreateCustomerRequest) (*Customer, error) {
	req, err := s.client.newRequest(ctx, http.MethodPost, "customers", request)
	if err != nil {
		return nil, err
	}

	var customer Customer
	if err := s.client.doJSON(req, &customer); err != nil {
		return nil, err
	}

	return &customer, nil
}

func (s *CustomersService) List(ctx context.Context, params CustomerListParams) (*CustomerList, error) {
	query := buildListQuery(params.ListParams)
	if params.Email != "" {
		query.Set("email", params.Email)
	}
	if params.StartTimestamp > 0 {
		query.Set("startTimestamp", strconv.Itoa(params.StartTimestamp))
	}
	if params.EndTimestamp > 0 {
		query.Set("endTimestamp", strconv.Itoa(params.EndTimestamp))
	}

	req, err := s.client.newRequestWithQuery(ctx, http.MethodGet, "customers", query, nil)
	if err != nil {
		return nil, err
	}

	var customers CustomerList
	if err := s.client.doJSON(req, &customers); err != nil {
		return nil, err
	}

	return &customers, nil
}

func (s *CustomersService) Get(ctx context.Context, customerID string) (*Customer, error) {
	req, err := s.client.newRequest(ctx, http.MethodGet, "customers/"+customerID, nil)
	if err != nil {
		return nil, err
	}

	var customer Customer
	if err := s.client.doJSON(req, &customer); err != nil {
		return nil, err
	}

	return &customer, nil
}

func (s *CustomersService) Update(ctx context.Context, customerID string, request UpdateCustomerRequest) (*Customer, error) {
	req, err := s.client.newRequest(ctx, http.MethodPatch, "customers/"+customerID, request)
	if err != nil {
		return nil, err
	}

	var customer Customer
	if err := s.client.doJSON(req, &customer); err != nil {
		return nil, err
	}

	return &customer, nil
}

func (s *CustomersService) Delete(ctx context.Context, customerID string) (*DeletedResource, error) {
	req, err := s.client.newRequest(ctx, http.MethodDelete, "customers/"+customerID, nil)
	if err != nil {
		return nil, err
	}

	var deleted DeletedResource
	if err := s.client.doJSON(req, &deleted); err != nil {
		return nil, err
	}

	return &deleted, nil
}

func (s *CustomersService) GetPaymentMethods(ctx context.Context, customerID string) (*CustomerPaymentMethodList, error) {
	req, err := s.client.newRequest(ctx, http.MethodGet, "customers/"+customerID+"/payment-methods", nil)
	if err != nil {
		return nil, err
	}

	var paymentMethods CustomerPaymentMethodList
	if err := s.client.doJSON(req, &paymentMethods); err != nil {
		return nil, err
	}

	return &paymentMethods, nil
}
