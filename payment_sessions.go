package ryft

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
)

type PaymentSessionsService struct {
	client *Client
}

type PaymentSession struct {
	ID                 string                   `json:"id"`
	Amount             int                      `json:"amount,omitempty"`
	Currency           string                   `json:"currency,omitempty"`
	CustomerEmail      string                   `json:"customerEmail,omitempty"`
	CaptureFlow        string                   `json:"captureFlow,omitempty"`
	ReturnURL          string                   `json:"returnUrl,omitempty"`
	Metadata           map[string]string        `json:"metadata,omitempty"`
	PlatformFee        int                      `json:"platformFee,omitempty"`
	CustomerDetails    *PaymentSessionCustomer  `json:"customerDetails,omitempty"`
	PreviousPayment    *PaymentSessionReference `json:"previousPayment,omitempty"`
	RebillingDetail    *RebillingDetail         `json:"rebillingDetail,omitempty"`
	SplitPaymentDetail *SplitPaymentDetail      `json:"splitPaymentDetail,omitempty"`
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

type PaymentSessionTransactionListParams struct {
	ListParams
	StartTimestamp int
	EndTimestamp   int
}

type CreatePaymentSessionRequest struct {
	Amount          int                        `json:"amount"`
	Currency        string                     `json:"currency"`
	CustomerEmail   string                     `json:"customerEmail,omitempty"`
	PaymentType     string                     `json:"paymentType,omitempty"`
	EntryMode       string                     `json:"entryMode,omitempty"`
	CaptureFlow     string                     `json:"captureFlow,omitempty"`
	ReturnURL       string                     `json:"returnUrl,omitempty"`
	Metadata        map[string]string          `json:"metadata,omitempty"`
	PlatformFee     int                        `json:"platformFee,omitempty"`
	CustomerDetails *PaymentSessionCustomer    `json:"customerDetails,omitempty"`
	PreviousPayment *PaymentSessionReference   `json:"previousPayment,omitempty"`
	RebillingDetail *RebillingDetail           `json:"rebillingDetail,omitempty"`
	Splits          *CreateSplitPaymentRequest `json:"splits,omitempty"`
	AttemptPayment  json.RawMessage            `json:"attemptPayment,omitempty"`
}

type PaymentSessionCustomer struct {
	ID string `json:"id,omitempty"`
}

type PaymentSessionReference struct {
	ID string `json:"id,omitempty"`
}

type RebillingDetail struct {
	AmountVariance              string `json:"amountVariance,omitempty"`
	NumberOfDaysBetweenPayments int    `json:"numberOfDaysBetweenPayments,omitempty"`
	TotalNumberOfPayments       int    `json:"totalNumberOfPayments,omitempty"`
	CurrentPaymentNumber        int    `json:"currentPaymentNumber,omitempty"`
}

type SplitFee struct {
	Amount int `json:"amount,omitempty"`
}

type SplitItem struct {
	AccountID   string            `json:"accountId"`
	Amount      int               `json:"amount"`
	Description string            `json:"description,omitempty"`
	Fee         *SplitFee         `json:"fee,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

type CreateSplitPaymentRequest struct {
	Items []SplitItem `json:"items"`
}

type SplitPaymentDetail struct {
	Items []SplitItem `json:"items"`
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

func (s *PaymentSessionsService) Create(
	ctx context.Context,
	request CreatePaymentSessionRequest,
	opts ...RequestOption,
) (*PaymentSession, error) {
	req, err := s.client.newRequest(ctx, http.MethodPost, "payment-sessions", request)
	if err != nil {
		return nil, err
	}
	s.client.applyRequestOptions(req, opts...)

	var paymentSession PaymentSession
	if err := s.client.doJSON(req, &paymentSession); err != nil {
		return nil, err
	}

	return &paymentSession, nil
}

func (s *PaymentSessionsService) Get(
	ctx context.Context,
	paymentSessionID string,
	opts ...RequestOption,
) (*PaymentSession, error) {
	req, err := s.client.newRequest(ctx, http.MethodGet, "payment-sessions/"+paymentSessionID, nil)
	if err != nil {
		return nil, err
	}
	s.client.applyRequestOptions(req, opts...)

	var paymentSession PaymentSession
	if err := s.client.doJSON(req, &paymentSession); err != nil {
		return nil, err
	}

	return &paymentSession, nil
}

func (s *PaymentSessionsService) Update(
	ctx context.Context,
	paymentSessionID string,
	request UpdatePaymentSessionRequest,
	opts ...RequestOption,
) (*PaymentSession, error) {
	return s.UpdateWithOptions(ctx, paymentSessionID, request, opts...)
}

func (s *PaymentSessionsService) UpdateWithOptions(
	ctx context.Context,
	paymentSessionID string,
	request UpdatePaymentSessionRequest,
	opts ...RequestOption,
) (*PaymentSession, error) {
	req, err := s.client.newRequest(ctx, http.MethodPatch, "payment-sessions/"+paymentSessionID, request)
	if err != nil {
		return nil, err
	}
	s.client.applyRequestOptions(req, opts...)

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
	opts ...RequestOption,
) (*PaymentSessionTransaction, error) {
	return s.RefundWithOptions(ctx, paymentSessionID, request, opts...)
}

func (s *PaymentSessionsService) RefundWithOptions(
	ctx context.Context,
	paymentSessionID string,
	request RefundPaymentSessionRequest,
	opts ...RequestOption,
) (*PaymentSessionTransaction, error) {
	req, err := s.client.newRequest(ctx, http.MethodPost, "payment-sessions/"+paymentSessionID+"/refunds", request)
	if err != nil {
		return nil, err
	}
	s.client.applyRequestOptions(req, opts...)

	var transaction PaymentSessionTransaction
	if err := s.client.doJSON(req, &transaction); err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (s *PaymentSessionsService) ListTransactions(
	ctx context.Context,
	paymentSessionID string,
	params PaymentSessionTransactionListParams,
	opts ...RequestOption,
) (*PaymentSessionTransactionList, error) {
	query := buildListQuery(params.ListParams)
	if params.StartTimestamp > 0 {
		query.Set("startTimestamp", strconv.Itoa(params.StartTimestamp))
	}
	if params.EndTimestamp > 0 {
		query.Set("endTimestamp", strconv.Itoa(params.EndTimestamp))
	}

	req, err := s.client.newRequestWithQuery(ctx, http.MethodGet, "payment-sessions/"+paymentSessionID+"/transactions", query, nil)
	if err != nil {
		return nil, err
	}
	s.client.applyRequestOptions(req, opts...)

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
	opts ...RequestOption,
) (*PaymentSessionTransaction, error) {
	req, err := s.client.newRequest(ctx, http.MethodGet, "payment-sessions/"+paymentSessionID+"/transactions/"+transactionID, nil)
	if err != nil {
		return nil, err
	}
	s.client.applyRequestOptions(req, opts...)

	var transaction PaymentSessionTransaction
	if err := s.client.doJSON(req, &transaction); err != nil {
		return nil, err
	}

	return &transaction, nil
}
