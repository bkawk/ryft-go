package ryft

import (
	"context"
)

type SubscriptionsService struct {
	client *Client
}

type Subscription struct {
	ID              string         `json:"id"`
	Status          string         `json:"status,omitempty"`
	Description     string         `json:"description,omitempty"`
	Customer        map[string]any `json:"customer,omitempty"`
	PaymentMethod   map[string]any `json:"paymentMethod,omitempty"`
	BillingDetail   map[string]any `json:"billingDetail,omitempty"`
	Metadata        map[string]any `json:"metadata,omitempty"`
	PaymentSessions map[string]any `json:"paymentSessions,omitempty"`
	CancelDetail    map[string]any `json:"cancelDetail,omitempty"`
}

type SubscriptionList struct {
	Items []Subscription `json:"items"`
}

type SubscriptionPaymentSessionList struct {
	Items []PaymentSession `json:"items"`
}

type SubscriptionCustomerReference struct {
	ID string `json:"id"`
}

type SubscriptionPaymentMethodReference struct {
	ID string `json:"id,omitempty"`
}

type SubscriptionInterval struct {
	Unit  string `json:"unit"`
	Count int    `json:"count"`
	Times int    `json:"times,omitempty"`
}

type SubscriptionPrice struct {
	Amount   int                  `json:"amount"`
	Currency string               `json:"currency"`
	Interval SubscriptionInterval `json:"interval"`
}

type CreateSubscriptionRequest struct {
	Customer              SubscriptionCustomerReference      `json:"customer"`
	PaymentMethod         SubscriptionPaymentMethodReference `json:"paymentMethod,omitempty"`
	Description           string                             `json:"description,omitempty"`
	Metadata              map[string]string                  `json:"metadata,omitempty"`
	Price                 SubscriptionPrice                  `json:"price"`
	BillingCycleTimestamp int                                `json:"billingCycleTimestamp,omitempty"`
	PaymentSettings       map[string]any                     `json:"paymentSettings,omitempty"`
	ShippingDetails       map[string]any                     `json:"shippingDetails,omitempty"`
}

type UpdateSubscriptionRequest struct {
	Price                 *SubscriptionPrice                  `json:"price,omitempty"`
	PaymentMethod         *SubscriptionPaymentMethodReference `json:"paymentMethod,omitempty"`
	Description           string                              `json:"description,omitempty"`
	BillingCycleTimestamp int                                 `json:"billingCycleTimestamp,omitempty"`
	Metadata              map[string]string                   `json:"metadata,omitempty"`
	PaymentSettings       map[string]any                      `json:"paymentSettings,omitempty"`
	ShippingDetails       map[string]any                      `json:"shippingDetails,omitempty"`
}

type PauseSubscriptionRequest struct {
	Reason          string `json:"reason,omitempty"`
	ResumeTimestamp int    `json:"resumeTimestamp,omitempty"`
	Unschedule      bool   `json:"unschedule,omitempty"`
}

func (s *SubscriptionsService) Create(ctx context.Context, request CreateSubscriptionRequest) (*Subscription, error) {
	req, err := s.client.newRequest(ctx, "POST", "subscriptions", request)
	if err != nil {
		return nil, err
	}

	var subscription Subscription
	if err := s.client.doJSON(req, &subscription); err != nil {
		return nil, err
	}

	return &subscription, nil
}

func (s *SubscriptionsService) Get(ctx context.Context, subscriptionID string) (*Subscription, error) {
	req, err := s.client.newRequest(ctx, "GET", "subscriptions/"+subscriptionID, nil)
	if err != nil {
		return nil, err
	}

	var subscription Subscription
	if err := s.client.doJSON(req, &subscription); err != nil {
		return nil, err
	}

	return &subscription, nil
}

func (s *SubscriptionsService) Update(
	ctx context.Context,
	subscriptionID string,
	request UpdateSubscriptionRequest,
) (*Subscription, error) {
	req, err := s.client.newRequest(ctx, "PATCH", "subscriptions/"+subscriptionID, request)
	if err != nil {
		return nil, err
	}

	var subscription Subscription
	if err := s.client.doJSON(req, &subscription); err != nil {
		return nil, err
	}

	return &subscription, nil
}

func (s *SubscriptionsService) List(
	ctx context.Context,
	startTimestamp int,
	endTimestamp int,
	ascending bool,
	limit int,
	startsAfter string,
) (*SubscriptionList, error) {
	query := buildListQuery(ascending, limit, startsAfter)
	if startTimestamp > 0 {
		query.Set("startTimestamp", itoa(startTimestamp))
	}
	if endTimestamp > 0 {
		query.Set("endTimestamp", itoa(endTimestamp))
	}

	req, err := s.client.newRequestWithQuery(ctx, "GET", "subscriptions", query, nil)
	if err != nil {
		return nil, err
	}

	var subscriptions SubscriptionList
	if err := s.client.doJSON(req, &subscriptions); err != nil {
		return nil, err
	}

	return &subscriptions, nil
}

func (s *SubscriptionsService) GetPaymentSessions(
	ctx context.Context,
	subscriptionID string,
	startTimestamp int,
	endTimestamp int,
	ascending bool,
	limit int,
	startsAfter string,
) (*SubscriptionPaymentSessionList, error) {
	query := buildListQuery(ascending, limit, startsAfter)
	if startTimestamp > 0 {
		query.Set("startTimestamp", itoa(startTimestamp))
	}
	if endTimestamp > 0 {
		query.Set("endTimestamp", itoa(endTimestamp))
	}

	req, err := s.client.newRequestWithQuery(
		ctx,
		"GET",
		"subscriptions/"+subscriptionID+"/payment-sessions",
		query,
		nil,
	)
	if err != nil {
		return nil, err
	}

	var paymentSessions SubscriptionPaymentSessionList
	if err := s.client.doJSON(req, &paymentSessions); err != nil {
		return nil, err
	}

	return &paymentSessions, nil
}

func (s *SubscriptionsService) Pause(
	ctx context.Context,
	subscriptionID string,
	request PauseSubscriptionRequest,
) (*Subscription, error) {
	req, err := s.client.newRequest(ctx, "PATCH", "subscriptions/"+subscriptionID+"/pause", request)
	if err != nil {
		return nil, err
	}

	var subscription Subscription
	if err := s.client.doJSON(req, &subscription); err != nil {
		return nil, err
	}

	return &subscription, nil
}

func (s *SubscriptionsService) Resume(ctx context.Context, subscriptionID string) (*Subscription, error) {
	req, err := s.client.newRequest(ctx, "PATCH", "subscriptions/"+subscriptionID+"/resume", nil)
	if err != nil {
		return nil, err
	}

	var subscription Subscription
	if err := s.client.doJSON(req, &subscription); err != nil {
		return nil, err
	}

	return &subscription, nil
}

func (s *SubscriptionsService) Cancel(ctx context.Context, subscriptionID string) (*DeletedResource, error) {
	req, err := s.client.newRequest(ctx, "DELETE", "subscriptions/"+subscriptionID+"/cancel", nil)
	if err != nil {
		return nil, err
	}

	var deleted DeletedResource
	if err := s.client.doJSON(req, &deleted); err != nil {
		return nil, err
	}

	return &deleted, nil
}
