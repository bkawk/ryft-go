package ryft

import (
	"context"
	"encoding/json"
	"net/http"
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
	Currency  string          `json:"currency,omitempty"`
	Available json.RawMessage `json:"available,omitempty"`
}

type BalanceTransactionList struct {
	Items []BalanceTransaction `json:"items"`
}

type BalanceTransaction struct {
	ID       string          `json:"id,omitempty"`
	Amount   int             `json:"amount,omitempty"`
	Currency string          `json:"currency,omitempty"`
	Type     string          `json:"type,omitempty"`
	Status   string          `json:"status,omitempty"`
	Origin   json.RawMessage `json:"origin,omitempty"`
}

type BalanceListParams struct {
	Currency string
}

type BalanceTransactionListParams struct {
	ListParams
	PayoutID string
}

func (s *BalancesService) List(ctx context.Context, params BalanceListParams, opts ...RequestOption) (*BalanceList, error) {
	query := url.Values{}
	if params.Currency != "" {
		query.Set("currency", params.Currency)
	}

	req, err := s.client.newRequestWithQuery(ctx, http.MethodGet, "balances", query, nil)
	if err != nil {
		return nil, err
	}
	s.client.applyRequestOptions(req, opts...)

	var balances BalanceList
	if err := s.client.doJSON(req, &balances); err != nil {
		return nil, err
	}

	return &balances, nil
}

func (s *BalanceTransactionsService) List(
	ctx context.Context,
	params BalanceTransactionListParams,
	opts ...RequestOption,
) (*BalanceTransactionList, error) {
	query := buildListQuery(params.ListParams)
	if params.PayoutID != "" {
		query.Set("payoutId", params.PayoutID)
	}

	req, err := s.client.newRequestWithQuery(ctx, http.MethodGet, "balance-transactions", query, nil)
	if err != nil {
		return nil, err
	}
	s.client.applyRequestOptions(req, opts...)

	var balanceTransactions BalanceTransactionList
	if err := s.client.doJSON(req, &balanceTransactions); err != nil {
		return nil, err
	}

	return &balanceTransactions, nil
}
