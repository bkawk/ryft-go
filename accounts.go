package ryft

import "context"

type AccountsService struct {
	client *Client
}

type Account struct {
	ID         string            `json:"id"`
	EntityType string            `json:"entityType,omitempty"`
	Email      string            `json:"email,omitempty"`
	Metadata   map[string]string `json:"metadata,omitempty"`
	Business   map[string]any    `json:"business,omitempty"`
	Individual map[string]any    `json:"individual,omitempty"`
}

type Address struct {
	LineOne    string `json:"lineOne"`
	City       string `json:"city"`
	Country    string `json:"country"`
	PostalCode string `json:"postalCode"`
}

type BusinessAccountDetails struct {
	Name               string  `json:"name"`
	Type               string  `json:"type"`
	RegistrationNumber string  `json:"registrationNumber"`
	RegisteredAddress  Address `json:"registeredAddress"`
	ContactEmail       string  `json:"contactEmail"`
}

type IndividualAccountDetails struct {
	FirstName     string   `json:"firstName"`
	LastName      string   `json:"lastName"`
	Email         string   `json:"email"`
	DateOfBirth   string   `json:"dateOfBirth"`
	Gender        string   `json:"gender"`
	Nationalities []string `json:"nationalities"`
	Address       Address  `json:"address"`
}

type TermsOfService struct {
	Acceptance Acceptance `json:"acceptance"`
}

type Acceptance struct {
	IPAddress string `json:"ipAddress"`
}

type CreateAccountRequest struct {
	OnboardingFlow string                    `json:"onboardingFlow"`
	EntityType     string                    `json:"entityType"`
	Email          string                    `json:"email"`
	Metadata       map[string]string         `json:"metadata,omitempty"`
	Business       *BusinessAccountDetails   `json:"business,omitempty"`
	Individual     *IndividualAccountDetails `json:"individual,omitempty"`
	TermsOfService TermsOfService            `json:"termsOfService"`
}

type UpdateBusinessAccountDetails struct {
	Name               string   `json:"name,omitempty"`
	Type               string   `json:"type,omitempty"`
	RegistrationNumber string   `json:"registrationNumber,omitempty"`
	RegistrationDate   string   `json:"registrationDate,omitempty"`
	RegisteredAddress  *Address `json:"registeredAddress,omitempty"`
	ContactEmail       string   `json:"contactEmail,omitempty"`
	PhoneNumber        string   `json:"phoneNumber,omitempty"`
	TradingName        string   `json:"tradingName,omitempty"`
	TradingAddress     *Address `json:"tradingAddress,omitempty"`
	TradingCountries   []string `json:"tradingCountries,omitempty"`
	WebsiteURL         string   `json:"websiteUrl,omitempty"`
}

type UpdateIndividualAccountDetails struct {
	FirstName      string   `json:"firstName,omitempty"`
	MiddleNames    string   `json:"middleNames,omitempty"`
	LastName       string   `json:"lastName,omitempty"`
	Email          string   `json:"email,omitempty"`
	DateOfBirth    string   `json:"dateOfBirth,omitempty"`
	CountryOfBirth string   `json:"countryOfBirth,omitempty"`
	Gender         string   `json:"gender,omitempty"`
	Nationalities  []string `json:"nationalities,omitempty"`
	Address        *Address `json:"address,omitempty"`
	PhoneNumber    string   `json:"phoneNumber,omitempty"`
}

type AccountSettings struct {
	Payouts map[string]any `json:"payouts,omitempty"`
}

type UpdateAccountRequest struct {
	EntityType     string                          `json:"entityType,omitempty"`
	Business       *UpdateBusinessAccountDetails   `json:"business,omitempty"`
	Individual     *UpdateIndividualAccountDetails `json:"individual,omitempty"`
	Metadata       map[string]string               `json:"metadata,omitempty"`
	Settings       map[string]any                  `json:"settings,omitempty"`
	TermsOfService *TermsOfService                 `json:"termsOfService,omitempty"`
}

type CreateAccountAuthorizationRequest struct {
	Email       string `json:"email"`
	RedirectURL string `json:"redirectUrl"`
}

type AccountAuthorization struct {
	URL string `json:"url,omitempty"`
}

func (s *AccountsService) Create(ctx context.Context, request CreateAccountRequest) (*Account, error) {
	req, err := s.client.newRequest(ctx, "POST", "accounts", request)
	if err != nil {
		return nil, err
	}

	var account Account
	if err := s.client.doJSON(req, &account); err != nil {
		return nil, err
	}

	return &account, nil
}

func (s *AccountsService) Get(ctx context.Context, accountID string) (*Account, error) {
	req, err := s.client.newRequest(ctx, "GET", "accounts/"+accountID, nil)
	if err != nil {
		return nil, err
	}

	var account Account
	if err := s.client.doJSON(req, &account); err != nil {
		return nil, err
	}

	return &account, nil
}

func (s *AccountsService) Update(ctx context.Context, accountID string, request UpdateAccountRequest) (*Account, error) {
	req, err := s.client.newRequest(ctx, "PATCH", "accounts/"+accountID, request)
	if err != nil {
		return nil, err
	}

	var account Account
	if err := s.client.doJSON(req, &account); err != nil {
		return nil, err
	}

	return &account, nil
}

func (s *AccountsService) Verify(ctx context.Context, accountID string) (*Account, error) {
	req, err := s.client.newRequest(ctx, "POST", "accounts/"+accountID+"/verify", nil)
	if err != nil {
		return nil, err
	}

	var account Account
	if err := s.client.doJSON(req, &account); err != nil {
		return nil, err
	}

	return &account, nil
}

func (s *AccountsService) CreateAuthLink(
	ctx context.Context,
	request CreateAccountAuthorizationRequest,
) (*AccountAuthorization, error) {
	req, err := s.client.newRequest(ctx, "POST", "accounts/authorize", request)
	if err != nil {
		return nil, err
	}

	var authorization AccountAuthorization
	if err := s.client.doJSON(req, &authorization); err != nil {
		return nil, err
	}

	return &authorization, nil
}
