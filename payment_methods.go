package ryft

import "context"

type PaymentMethodsService struct {
	client *Client
}

type PaymentMethod struct {
	ID             string         `json:"id"`
	CustomerID     string         `json:"customerId,omitempty"`
	Type           string         `json:"type,omitempty"`
	BillingAddress map[string]any `json:"billingAddress,omitempty"`
}

type UpdatePaymentMethodRequest struct {
	BillingAddress map[string]string `json:"billingAddress,omitempty"`
}

type DeletedResource struct {
	ID string `json:"id"`
}

func (s *PaymentMethodsService) Get(ctx context.Context, paymentMethodID string) (*PaymentMethod, error) {
	req, err := s.client.newRequest(ctx, "GET", "payment-methods/"+paymentMethodID, nil)
	if err != nil {
		return nil, err
	}

	var paymentMethod PaymentMethod
	if err := s.client.doJSON(req, &paymentMethod); err != nil {
		return nil, err
	}

	return &paymentMethod, nil
}

func (s *PaymentMethodsService) Update(
	ctx context.Context,
	paymentMethodID string,
	request UpdatePaymentMethodRequest,
) (*PaymentMethod, error) {
	req, err := s.client.newRequest(ctx, "PATCH", "payment-methods/"+paymentMethodID, request)
	if err != nil {
		return nil, err
	}

	var paymentMethod PaymentMethod
	if err := s.client.doJSON(req, &paymentMethod); err != nil {
		return nil, err
	}

	return &paymentMethod, nil
}

func (s *PaymentMethodsService) Delete(ctx context.Context, paymentMethodID string) (*DeletedResource, error) {
	req, err := s.client.newRequest(ctx, "DELETE", "payment-methods/"+paymentMethodID, nil)
	if err != nil {
		return nil, err
	}

	var deleted DeletedResource
	if err := s.client.doJSON(req, &deleted); err != nil {
		return nil, err
	}

	return &deleted, nil
}
