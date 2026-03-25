package ryft

import (
	"context"
	"net/http"
)

type PayoutMethodsService struct {
	client *Client
}

type PayoutMethod struct {
	ID          string `json:"id"`
	Currency    string `json:"currency,omitempty"`
	CountryCode string `json:"countryCode,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
}

type PayoutMethodList struct {
	Items []PayoutMethod `json:"items"`
}

type BankAccountDetails struct {
	AccountNumberType string `json:"accountNumberType"`
	AccountNumber     string `json:"accountNumber"`
	BankIDType        string `json:"bankIdType"`
	BankID            string `json:"bankId"`
}

type CreatePayoutMethodRequest struct {
	Type        string             `json:"type"`
	DisplayName string             `json:"displayName"`
	Currency    string             `json:"currency"`
	Country     string             `json:"country"`
	BankAccount BankAccountDetails `json:"bankAccount"`
}

type UpdatePayoutMethodRequest struct {
	DisplayName string             `json:"displayName,omitempty"`
	BankAccount BankAccountDetails `json:"bankAccount"`
}

type PayoutMethodListParams struct {
	ListParams
}

func (s *PayoutMethodsService) Create(ctx context.Context, accountID string, request CreatePayoutMethodRequest) (*PayoutMethod, error) {
	req, err := s.client.newRequest(ctx, http.MethodPost, "accounts/"+accountID+"/payout-methods", request)
	if err != nil {
		return nil, err
	}

	var payoutMethod PayoutMethod
	if err := s.client.doJSON(req, &payoutMethod); err != nil {
		return nil, err
	}

	return &payoutMethod, nil
}

func (s *PayoutMethodsService) Get(ctx context.Context, accountID string, payoutMethodID string) (*PayoutMethod, error) {
	req, err := s.client.newRequest(ctx, http.MethodGet, "accounts/"+accountID+"/payout-methods/"+payoutMethodID, nil)
	if err != nil {
		return nil, err
	}

	var payoutMethod PayoutMethod
	if err := s.client.doJSON(req, &payoutMethod); err != nil {
		return nil, err
	}

	return &payoutMethod, nil
}

func (s *PayoutMethodsService) List(
	ctx context.Context,
	accountID string,
	params PayoutMethodListParams,
) (*PayoutMethodList, error) {
	query := buildListQuery(params.ListParams)

	req, err := s.client.newRequestWithQuery(ctx, http.MethodGet, "accounts/"+accountID+"/payout-methods", query, nil)
	if err != nil {
		return nil, err
	}

	var payoutMethods PayoutMethodList
	if err := s.client.doJSON(req, &payoutMethods); err != nil {
		return nil, err
	}

	return &payoutMethods, nil
}

func (s *PayoutMethodsService) Update(
	ctx context.Context,
	accountID string,
	payoutMethodID string,
	request UpdatePayoutMethodRequest,
) (*PayoutMethod, error) {
	req, err := s.client.newRequest(ctx, http.MethodPatch, "accounts/"+accountID+"/payout-methods/"+payoutMethodID, request)
	if err != nil {
		return nil, err
	}

	var payoutMethod PayoutMethod
	if err := s.client.doJSON(req, &payoutMethod); err != nil {
		return nil, err
	}

	return &payoutMethod, nil
}

func (s *PayoutMethodsService) Delete(ctx context.Context, accountID string, payoutMethodID string) (*DeletedResource, error) {
	req, err := s.client.newRequest(ctx, http.MethodDelete, "accounts/"+accountID+"/payout-methods/"+payoutMethodID, nil)
	if err != nil {
		return nil, err
	}

	var deleted DeletedResource
	if err := s.client.doJSON(req, &deleted); err != nil {
		return nil, err
	}

	return &deleted, nil
}
