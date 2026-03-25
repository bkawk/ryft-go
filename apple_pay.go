package ryft

import (
	"context"
	"net/http"
)

type ApplePayService struct {
	client *Client
}

type ApplePayWebDomain struct {
	ID               string `json:"id,omitempty"`
	DomainName       string `json:"domainName,omitempty"`
	CreatedTimestamp int64  `json:"createdTimestamp,omitempty"`
}

type ApplePayWebDomainList struct {
	Items           []ApplePayWebDomain `json:"items"`
	PaginationToken string              `json:"paginationToken,omitempty"`
}

type ApplePayWebDomainListParams struct {
	ListParams
}

type RegisterApplePayWebDomainRequest struct {
	DomainName string `json:"domainName"`
}

type CreateApplePayWebSessionRequest struct {
	DisplayName string `json:"displayName"`
	DomainName  string `json:"domainName"`
}

type ApplePayWebSession struct {
	SessionObject string `json:"sessionObject,omitempty"`
}

func (s *ApplePayService) RegisterDomain(
	ctx context.Context,
	request RegisterApplePayWebDomainRequest,
	opts ...RequestOption,
) (*ApplePayWebDomain, error) {
	req, err := s.client.newRequest(ctx, http.MethodPost, "apple-pay/web-domains", request)
	if err != nil {
		return nil, err
	}
	s.client.applyRequestOptions(req, opts...)

	var domain ApplePayWebDomain
	if err := s.client.doJSON(req, &domain); err != nil {
		return nil, err
	}

	return &domain, nil
}

func (s *ApplePayService) ListDomains(
	ctx context.Context,
	params ApplePayWebDomainListParams,
	opts ...RequestOption,
) (*ApplePayWebDomainList, error) {
	query := buildListQuery(params.ListParams)
	req, err := s.client.newRequestWithQuery(ctx, http.MethodGet, "apple-pay/web-domains", query, nil)
	if err != nil {
		return nil, err
	}
	s.client.applyRequestOptions(req, opts...)

	var domains ApplePayWebDomainList
	if err := s.client.doJSON(req, &domains); err != nil {
		return nil, err
	}

	return &domains, nil
}

func (s *ApplePayService) GetDomain(
	ctx context.Context,
	domainID string,
	opts ...RequestOption,
) (*ApplePayWebDomain, error) {
	req, err := s.client.newRequest(ctx, http.MethodGet, "apple-pay/web-domains/"+domainID, nil)
	if err != nil {
		return nil, err
	}
	s.client.applyRequestOptions(req, opts...)

	var domain ApplePayWebDomain
	if err := s.client.doJSON(req, &domain); err != nil {
		return nil, err
	}

	return &domain, nil
}

func (s *ApplePayService) DeleteDomain(
	ctx context.Context,
	domainID string,
	opts ...RequestOption,
) (*DeletedResource, error) {
	req, err := s.client.newRequest(ctx, http.MethodDelete, "apple-pay/web-domains/"+domainID, nil)
	if err != nil {
		return nil, err
	}
	s.client.applyRequestOptions(req, opts...)

	var deleted DeletedResource
	if err := s.client.doJSON(req, &deleted); err != nil {
		return nil, err
	}

	return &deleted, nil
}

func (s *ApplePayService) CreateSession(
	ctx context.Context,
	request CreateApplePayWebSessionRequest,
	opts ...RequestOption,
) (*ApplePayWebSession, error) {
	req, err := s.client.newRequest(ctx, http.MethodPost, "apple-pay/sessions", request)
	if err != nil {
		return nil, err
	}
	s.client.applyRequestOptions(req, opts...)

	var session ApplePayWebSession
	if err := s.client.doJSON(req, &session); err != nil {
		return nil, err
	}

	return &session, nil
}
