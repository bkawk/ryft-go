package ryft

import "context"

type AccountLinksService struct {
	client *Client
}

type TemporaryAccountLink struct {
	URL              string `json:"url,omitempty"`
	CreatedTimestamp int64  `json:"createdTimestamp,omitempty"`
	ExpiresTimestamp int64  `json:"expiresTimestamp,omitempty"`
}

type CreateTemporaryAccountLinkRequest struct {
	AccountID   string `json:"accountId"`
	RedirectURL string `json:"redirectUrl"`
}

func (s *AccountLinksService) GenerateTemporaryAccountLink(
	ctx context.Context,
	request CreateTemporaryAccountLinkRequest,
) (*TemporaryAccountLink, error) {
	req, err := s.client.newRequest(ctx, "POST", "account-links", request)
	if err != nil {
		return nil, err
	}

	var link TemporaryAccountLink
	if err := s.client.doJSON(req, &link); err != nil {
		return nil, err
	}

	return &link, nil
}
