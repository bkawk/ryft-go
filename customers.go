package ryft

import "context"

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
	req, err := s.client.newRequest(ctx, "POST", "customers", request)
	if err != nil {
		return nil, err
	}

	var customer Customer
	if err := s.client.doJSON(req, &customer); err != nil {
		return nil, err
	}

	return &customer, nil
}

func (s *CustomersService) List(
	ctx context.Context,
	email string,
	startTimestamp int,
	endTimestamp int,
	ascending bool,
	limit int,
	startsAfter string,
) (*CustomerList, error) {
	query := buildListQuery(ascending, limit, startsAfter)
	if email != "" {
		query.Set("email", email)
	}
	if startTimestamp > 0 {
		query.Set("startTimestamp", itoa(startTimestamp))
	}
	if endTimestamp > 0 {
		query.Set("endTimestamp", itoa(endTimestamp))
	}

	req, err := s.client.newRequestWithQuery(ctx, "GET", "customers", query, nil)
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
	req, err := s.client.newRequest(ctx, "GET", "customers/"+customerID, nil)
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
	req, err := s.client.newRequest(ctx, "PATCH", "customers/"+customerID, request)
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
	req, err := s.client.newRequest(ctx, "DELETE", "customers/"+customerID, nil)
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
	req, err := s.client.newRequest(ctx, "GET", "customers/"+customerID+"/payment-methods", nil)
	if err != nil {
		return nil, err
	}

	var paymentMethods CustomerPaymentMethodList
	if err := s.client.doJSON(req, &paymentMethods); err != nil {
		return nil, err
	}

	return &paymentMethods, nil
}
