package ryft

import "context"

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
) (*ApplePayWebDomain, error) {
	return s.RegisterDomainForAccount(ctx, request, "")
}

func (s *ApplePayService) RegisterDomainForAccount(
	ctx context.Context,
	request RegisterApplePayWebDomainRequest,
	accountID string,
) (*ApplePayWebDomain, error) {
	req, err := s.client.newRequest(ctx, "POST", "apple-pay/web-domains", request)
	if err != nil {
		return nil, err
	}
	if accountID != "" {
		req.Header.Set("Account", accountID)
	}

	var domain ApplePayWebDomain
	if err := s.client.doJSON(req, &domain); err != nil {
		return nil, err
	}

	return &domain, nil
}

func (s *ApplePayService) ListDomains(
	ctx context.Context,
	ascending bool,
	limit int,
	startsAfter string,
) (*ApplePayWebDomainList, error) {
	return s.ListDomainsForAccount(ctx, ascending, limit, startsAfter, "")
}

func (s *ApplePayService) ListDomainsForAccount(
	ctx context.Context,
	ascending bool,
	limit int,
	startsAfter string,
	accountID string,
) (*ApplePayWebDomainList, error) {
	query := buildListQuery(ascending, limit, startsAfter)
	req, err := s.client.newRequestWithQuery(ctx, "GET", "apple-pay/web-domains", query, nil)
	if err != nil {
		return nil, err
	}
	if accountID != "" {
		req.Header.Set("Account", accountID)
	}

	var domains ApplePayWebDomainList
	if err := s.client.doJSON(req, &domains); err != nil {
		return nil, err
	}

	return &domains, nil
}

func (s *ApplePayService) GetDomain(ctx context.Context, domainID string) (*ApplePayWebDomain, error) {
	return s.GetDomainForAccount(ctx, domainID, "")
}

func (s *ApplePayService) GetDomainForAccount(
	ctx context.Context,
	domainID string,
	accountID string,
) (*ApplePayWebDomain, error) {
	req, err := s.client.newRequest(ctx, "GET", "apple-pay/web-domains/"+domainID, nil)
	if err != nil {
		return nil, err
	}
	if accountID != "" {
		req.Header.Set("Account", accountID)
	}

	var domain ApplePayWebDomain
	if err := s.client.doJSON(req, &domain); err != nil {
		return nil, err
	}

	return &domain, nil
}

func (s *ApplePayService) DeleteDomain(ctx context.Context, domainID string) (*DeletedResource, error) {
	return s.DeleteDomainForAccount(ctx, domainID, "")
}

func (s *ApplePayService) DeleteDomainForAccount(
	ctx context.Context,
	domainID string,
	accountID string,
) (*DeletedResource, error) {
	req, err := s.client.newRequest(ctx, "DELETE", "apple-pay/web-domains/"+domainID, nil)
	if err != nil {
		return nil, err
	}
	if accountID != "" {
		req.Header.Set("Account", accountID)
	}

	var deleted DeletedResource
	if err := s.client.doJSON(req, &deleted); err != nil {
		return nil, err
	}

	return &deleted, nil
}

func (s *ApplePayService) CreateSession(
	ctx context.Context,
	request CreateApplePayWebSessionRequest,
) (*ApplePayWebSession, error) {
	return s.CreateSessionForAccount(ctx, request, "")
}

func (s *ApplePayService) CreateSessionForAccount(
	ctx context.Context,
	request CreateApplePayWebSessionRequest,
	accountID string,
) (*ApplePayWebSession, error) {
	req, err := s.client.newRequest(ctx, "POST", "apple-pay/sessions", request)
	if err != nil {
		return nil, err
	}
	if accountID != "" {
		req.Header.Set("Account", accountID)
	}

	var session ApplePayWebSession
	if err := s.client.doJSON(req, &session); err != nil {
		return nil, err
	}

	return &session, nil
}
