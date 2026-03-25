package ryft

import (
	"context"
	"net/url"
)

type PaymentSessionsService struct {
	client *Client
}

type PaymentSession struct {
	ID            string            `json:"id"`
	Amount        int               `json:"amount,omitempty"`
	Currency      string            `json:"currency,omitempty"`
	CustomerEmail string            `json:"customerEmail,omitempty"`
	CaptureFlow   string            `json:"captureFlow,omitempty"`
	ReturnURL     string            `json:"returnUrl,omitempty"`
	Metadata      map[string]string `json:"metadata,omitempty"`
}

type PaymentSessionTransaction struct {
	ID               string `json:"id,omitempty"`
	PaymentSessionID string `json:"paymentSessionId,omitempty"`
	Amount           int    `json:"amount,omitempty"`
	Currency         string `json:"currency,omitempty"`
	Type             string `json:"type,omitempty"`
	Status           string `json:"status,omitempty"`
	Reason           string `json:"reason,omitempty"`
	RefundedAmount   int    `json:"refundedAmount,omitempty"`
}

type PaymentSessionTransactionList struct {
	Items []PaymentSessionTransaction `json:"items"`
}

type CreatePaymentSessionRequest struct {
	Amount        int               `json:"amount"`
	Currency      string            `json:"currency"`
	CustomerEmail string            `json:"customerEmail,omitempty"`
	PaymentType   string            `json:"paymentType,omitempty"`
	EntryMode     string            `json:"entryMode,omitempty"`
	CaptureFlow   string            `json:"captureFlow,omitempty"`
	ReturnURL     string            `json:"returnUrl,omitempty"`
	Metadata      map[string]string `json:"metadata,omitempty"`
}

type UpdatePaymentSessionRequest struct {
	Amount        *int              `json:"amount,omitempty"`
	CustomerEmail string            `json:"customerEmail,omitempty"`
	CaptureFlow   string            `json:"captureFlow,omitempty"`
	Metadata      map[string]string `json:"metadata,omitempty"`
}

type RefundPaymentSessionRequest struct {
	Amount            int    `json:"amount,omitempty"`
	Reason            string `json:"reason,omitempty"`
	RefundPlatformFee bool   `json:"refundPlatformFee,omitempty"`
}

func (s *PaymentSessionsService) Create(ctx context.Context, request CreatePaymentSessionRequest) (*PaymentSession, error) {
	req, err := s.client.newRequest(ctx, "POST", "payment-sessions", request)
	if err != nil {
		return nil, err
	}

	var paymentSession PaymentSession
	if err := s.client.doJSON(req, &paymentSession); err != nil {
		return nil, err
	}

	return &paymentSession, nil
}

func (s *PaymentSessionsService) Get(ctx context.Context, paymentSessionID string) (*PaymentSession, error) {
	req, err := s.client.newRequest(ctx, "GET", "payment-sessions/"+paymentSessionID, nil)
	if err != nil {
		return nil, err
	}

	var paymentSession PaymentSession
	if err := s.client.doJSON(req, &paymentSession); err != nil {
		return nil, err
	}

	return &paymentSession, nil
}

func (s *PaymentSessionsService) Update(ctx context.Context, paymentSessionID string, request UpdatePaymentSessionRequest) (*PaymentSession, error) {
	req, err := s.client.newRequest(ctx, "PATCH", "payment-sessions/"+paymentSessionID, request)
	if err != nil {
		return nil, err
	}

	var paymentSession PaymentSession
	if err := s.client.doJSON(req, &paymentSession); err != nil {
		return nil, err
	}

	return &paymentSession, nil
}

func (s *PaymentSessionsService) Refund(
	ctx context.Context,
	paymentSessionID string,
	request RefundPaymentSessionRequest,
) (*PaymentSessionTransaction, error) {
	req, err := s.client.newRequest(ctx, "POST", "payment-sessions/"+paymentSessionID+"/refunds", request)
	if err != nil {
		return nil, err
	}

	var transaction PaymentSessionTransaction
	if err := s.client.doJSON(req, &transaction); err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (s *PaymentSessionsService) ListTransactions(
	ctx context.Context,
	paymentSessionID string,
	startTimestamp int,
	endTimestamp int,
	ascending bool,
	limit int,
) (*PaymentSessionTransactionList, error) {
	query := url.Values{}
	if startTimestamp > 0 {
		query.Set("startTimestamp", itoa(startTimestamp))
	}
	if endTimestamp > 0 {
		query.Set("endTimestamp", itoa(endTimestamp))
	}
	query.Set("ascending", boolString(ascending))
	if limit > 0 {
		query.Set("limit", itoa(limit))
	}

	req, err := s.client.newRequestWithQuery(ctx, "GET", "payment-sessions/"+paymentSessionID+"/transactions", query, nil)
	if err != nil {
		return nil, err
	}

	var transactions PaymentSessionTransactionList
	if err := s.client.doJSON(req, &transactions); err != nil {
		return nil, err
	}

	return &transactions, nil
}

func (s *PaymentSessionsService) GetTransaction(
	ctx context.Context,
	paymentSessionID string,
	transactionID string,
) (*PaymentSessionTransaction, error) {
	req, err := s.client.newRequest(ctx, "GET", "payment-sessions/"+paymentSessionID+"/transactions/"+transactionID, nil)
	if err != nil {
		return nil, err
	}

	var transaction PaymentSessionTransaction
	if err := s.client.doJSON(req, &transaction); err != nil {
		return nil, err
	}

	return &transaction, nil
}
