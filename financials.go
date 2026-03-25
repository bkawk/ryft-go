package ryft

import (
	"context"
	"net/url"
)

type BalancesService struct {
	client *Client
}

type BalanceTransactionsService struct {
	client *Client
}

type BalanceList struct {
	Items []Balance `json:"items"`
}

type Balance struct {
	Currency  string         `json:"currency,omitempty"`
	Available map[string]any `json:"available,omitempty"`
}

type BalanceTransactionList struct {
	Items []BalanceTransaction `json:"items"`
}

type BalanceTransaction struct {
	ID       string         `json:"id,omitempty"`
	Amount   int            `json:"amount,omitempty"`
	Currency string         `json:"currency,omitempty"`
	Type     string         `json:"type,omitempty"`
	Status   string         `json:"status,omitempty"`
	Origin   map[string]any `json:"origin,omitempty"`
}

func (s *BalancesService) List(ctx context.Context, currency string, accountID string) (*BalanceList, error) {
	query := url.Values{}
	query.Set("currency", currency)

	req, err := s.client.newRequestWithQuery(ctx, "GET", "balances", query, nil)
	if err != nil {
		return nil, err
	}
	if accountID != "" {
		req.Header.Set("Account", accountID)
	}

	var balances BalanceList
	if err := s.client.doJSON(req, &balances); err != nil {
		return nil, err
	}

	return &balances, nil
}

func (s *BalanceTransactionsService) List(
	ctx context.Context,
	limit int,
	startsAfter string,
	payoutID string,
	accountID string,
) (*BalanceTransactionList, error) {
	query := url.Values{}
	if limit > 0 {
		query.Set("limit", itoa(limit))
	}
	if startsAfter != "" {
		query.Set("startsAfter", startsAfter)
	}
	if payoutID != "" {
		query.Set("payoutId", payoutID)
	}

	req, err := s.client.newRequestWithQuery(ctx, "GET", "balance-transactions", query, nil)
	if err != nil {
		return nil, err
	}
	if accountID != "" {
		req.Header.Set("Account", accountID)
	}

	var balanceTransactions BalanceTransactionList
	if err := s.client.doJSON(req, &balanceTransactions); err != nil {
		return nil, err
	}

	return &balanceTransactions, nil
}
