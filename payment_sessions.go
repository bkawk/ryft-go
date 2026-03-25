package ryft

import (
	"context"
	"net/url"
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
	AttemptPayment  map[string]any             `json:"attemptPayment,omitempty"`
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

func (s *PaymentSessionsService) Create(ctx context.Context, request CreatePaymentSessionRequest) (*PaymentSession, error) {
	return s.CreateForAccount(ctx, request, "")
}

func (s *PaymentSessionsService) CreateForAccount(
	ctx context.Context,
	request CreatePaymentSessionRequest,
	accountID string,
) (*PaymentSession, error) {
	req, err := s.client.newRequest(ctx, "POST", "payment-sessions", request)
	if err != nil {
		return nil, err
	}
	if accountID != "" {
		req.Header.Set("Account", accountID)
	}

	var paymentSession PaymentSession
	if err := s.client.doJSON(req, &paymentSession); err != nil {
		return nil, err
	}

	return &paymentSession, nil
}

func (s *PaymentSessionsService) Get(ctx context.Context, paymentSessionID string) (*PaymentSession, error) {
	return s.GetForAccount(ctx, paymentSessionID, "")
}

func (s *PaymentSessionsService) GetForAccount(
	ctx context.Context,
	paymentSessionID string,
	accountID string,
) (*PaymentSession, error) {
	req, err := s.client.newRequest(ctx, "GET", "payment-sessions/"+paymentSessionID, nil)
	if err != nil {
		return nil, err
	}
	if accountID != "" {
		req.Header.Set("Account", accountID)
	}

	var paymentSession PaymentSession
	if err := s.client.doJSON(req, &paymentSession); err != nil {
		return nil, err
	}

	return &paymentSession, nil
}

func (s *PaymentSessionsService) Update(ctx context.Context, paymentSessionID string, request UpdatePaymentSessionRequest) (*PaymentSession, error) {
	return s.UpdateForAccount(ctx, paymentSessionID, request, "")
}

func (s *PaymentSessionsService) UpdateForAccount(
	ctx context.Context,
	paymentSessionID string,
	request UpdatePaymentSessionRequest,
	accountID string,
) (*PaymentSession, error) {
	req, err := s.client.newRequest(ctx, "PATCH", "payment-sessions/"+paymentSessionID, request)
	if err != nil {
		return nil, err
	}
	if accountID != "" {
		req.Header.Set("Account", accountID)
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
	return s.RefundForAccount(ctx, paymentSessionID, request, "")
}

func (s *PaymentSessionsService) RefundForAccount(
	ctx context.Context,
	paymentSessionID string,
	request RefundPaymentSessionRequest,
	accountID string,
) (*PaymentSessionTransaction, error) {
	req, err := s.client.newRequest(ctx, "POST", "payment-sessions/"+paymentSessionID+"/refunds", request)
	if err != nil {
		return nil, err
	}
	if accountID != "" {
		req.Header.Set("Account", accountID)
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
	return s.ListTransactionsForAccount(ctx, paymentSessionID, startTimestamp, endTimestamp, ascending, limit, "")
}

func (s *PaymentSessionsService) ListTransactionsForAccount(
	ctx context.Context,
	paymentSessionID string,
	startTimestamp int,
	endTimestamp int,
	ascending bool,
	limit int,
	accountID string,
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
	if accountID != "" {
		req.Header.Set("Account", accountID)
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
	return s.GetTransactionForAccount(ctx, paymentSessionID, transactionID, "")
}

func (s *PaymentSessionsService) GetTransactionForAccount(
	ctx context.Context,
	paymentSessionID string,
	transactionID string,
	accountID string,
) (*PaymentSessionTransaction, error) {
	req, err := s.client.newRequest(ctx, "GET", "payment-sessions/"+paymentSessionID+"/transactions/"+transactionID, nil)
	if err != nil {
		return nil, err
	}
	if accountID != "" {
		req.Header.Set("Account", accountID)
	}

	var transaction PaymentSessionTransaction
	if err := s.client.doJSON(req, &transaction); err != nil {
		return nil, err
	}

	return &transaction, nil
}
