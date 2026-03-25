package ryft

import (
	"context"
	"net/url"
)

type InPersonProductsService struct {
	client *Client
}

type InPersonSkusService struct {
	client *Client
}

type InPersonOrdersService struct {
	client *Client
}

type InPersonLocationsService struct {
	client *Client
}

type InPersonTerminalsService struct {
	client *Client
}

type InPersonProductList struct {
	Items           []InPersonProduct `json:"items"`
	PaginationToken string            `json:"paginationToken,omitempty"`
}

type InPersonProduct struct {
	ID          string            `json:"id,omitempty"`
	Name        string            `json:"name,omitempty"`
	Description string            `json:"description,omitempty"`
	Status      string            `json:"status,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

type InPersonSkuList struct {
	Items           []InPersonSku `json:"items"`
	PaginationToken string        `json:"paginationToken,omitempty"`
}

type InPersonSku struct {
	ID        string            `json:"id,omitempty"`
	ProductID string            `json:"productId,omitempty"`
	Country   string            `json:"country,omitempty"`
	Currency  string            `json:"currency,omitempty"`
	Amount    int               `json:"amount,omitempty"`
	Status    string            `json:"status,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

type InPersonOrderList struct {
	Items           []InPersonOrder `json:"items"`
	PaginationToken string          `json:"paginationToken,omitempty"`
}

type InPersonOrder struct {
	ID       string           `json:"id,omitempty"`
	Status   string           `json:"status,omitempty"`
	Amount   int              `json:"amount,omitempty"`
	Currency string           `json:"currency,omitempty"`
	Items    []map[string]any `json:"items,omitempty"`
	Customer map[string]any   `json:"customer,omitempty"`
	Shipping map[string]any   `json:"shipping,omitempty"`
	Tracking map[string]any   `json:"tracking,omitempty"`
}

type InPersonLocationList struct {
	Items           []InPersonLocation `json:"items"`
	PaginationToken string             `json:"paginationToken,omitempty"`
}

type InPersonLocation struct {
	ID             string                   `json:"id,omitempty"`
	Name           string                   `json:"name,omitempty"`
	Address        *InPersonLocationAddress `json:"address,omitempty"`
	GeoCoordinates *GeoCoordinates          `json:"geoCoordinates,omitempty"`
	Metadata       map[string]string        `json:"metadata,omitempty"`
}

type InPersonLocationAddress struct {
	FirstLine  string `json:"firstLine,omitempty"`
	SecondLine string `json:"secondLine,omitempty"`
	City       string `json:"city,omitempty"`
	Region     string `json:"region,omitempty"`
	PostalCode string `json:"postalCode,omitempty"`
	Country    string `json:"country,omitempty"`
}

type GeoCoordinates struct {
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
}

type CreateInPersonLocationRequest struct {
	Name           string                  `json:"name"`
	Address        InPersonLocationAddress `json:"address"`
	GeoCoordinates *GeoCoordinates         `json:"geoCoordinates,omitempty"`
	Metadata       map[string]string       `json:"metadata,omitempty"`
}

type UpdateInPersonLocationRequest struct {
	Name     string            `json:"name,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

type TerminalList struct {
	Items           []Terminal `json:"items"`
	PaginationToken string     `json:"paginationToken,omitempty"`
}

type Terminal struct {
	ID           string            `json:"id,omitempty"`
	SerialNumber string            `json:"serialNumber,omitempty"`
	LocationID   string            `json:"locationId,omitempty"`
	Name         string            `json:"name,omitempty"`
	Status       string            `json:"status,omitempty"`
	Metadata     map[string]string `json:"metadata,omitempty"`
}

type CreateTerminalRequest struct {
	SerialNumber string            `json:"serialNumber"`
	LocationID   string            `json:"locationId"`
	Name         string            `json:"name,omitempty"`
	Metadata     map[string]string `json:"metadata,omitempty"`
}

type UpdateTerminalRequest struct {
	LocationID string            `json:"locationId,omitempty"`
	Name       string            `json:"name,omitempty"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

type RequestedAmounts struct {
	Requested int `json:"requested"`
}

type TerminalPaymentSessionRequest struct {
	PlatformFee     int               `json:"platformFee,omitempty"`
	Metadata        map[string]string `json:"metadata,omitempty"`
	PaymentSettings map[string]any    `json:"paymentSettings,omitempty"`
}

type TerminalPaymentRequest struct {
	Amounts        RequestedAmounts               `json:"amounts"`
	Currency       string                         `json:"currency"`
	PaymentSession *TerminalPaymentSessionRequest `json:"paymentSession,omitempty"`
	Settings       map[string]any                 `json:"settings,omitempty"`
}

type TerminalRefundPaymentSessionReference struct {
	ID string `json:"id"`
}

type TerminalRefundRequest struct {
	PaymentSession    TerminalRefundPaymentSessionReference `json:"paymentSession"`
	Amount            int                                   `json:"amount,omitempty"`
	RefundPlatformFee bool                                  `json:"refundPlatformFee,omitempty"`
	Settings          map[string]any                        `json:"settings,omitempty"`
}

type ReceiptCopyStatus struct {
	Status string `json:"status"`
}

type TerminalConfirmReceiptRequest struct {
	CustomerCopy *ReceiptCopyStatus `json:"customerCopy,omitempty"`
	MerchantCopy *ReceiptCopyStatus `json:"merchantCopy,omitempty"`
}

func (s *InPersonProductsService) List(
	ctx context.Context,
	ascending bool,
	limit int,
	startsAfter string,
) (*InPersonProductList, error) {
	query := buildListQuery(ascending, limit, startsAfter)
	req, err := s.client.newRequestWithQuery(ctx, "GET", "in-person/products", query, nil)
	if err != nil {
		return nil, err
	}

	var products InPersonProductList
	if err := s.client.doJSON(req, &products); err != nil {
		return nil, err
	}

	return &products, nil
}

func (s *InPersonProductsService) Get(ctx context.Context, productID string) (*InPersonProduct, error) {
	req, err := s.client.newRequest(ctx, "GET", "in-person/products/"+productID, nil)
	if err != nil {
		return nil, err
	}

	var product InPersonProduct
	if err := s.client.doJSON(req, &product); err != nil {
		return nil, err
	}

	return &product, nil
}

func (s *InPersonSkusService) List(
	ctx context.Context,
	country string,
	limit int,
	startsAfter string,
	productID string,
) (*InPersonSkuList, error) {
	query := url.Values{}
	query.Set("country", country)
	if limit > 0 {
		query.Set("limit", itoa(limit))
	}
	if startsAfter != "" {
		query.Set("startsAfter", startsAfter)
	}
	if productID != "" {
		query.Set("productId", productID)
	}

	req, err := s.client.newRequestWithQuery(ctx, "GET", "in-person/skus", query, nil)
	if err != nil {
		return nil, err
	}

	var skus InPersonSkuList
	if err := s.client.doJSON(req, &skus); err != nil {
		return nil, err
	}

	return &skus, nil
}

func (s *InPersonSkusService) Get(ctx context.Context, skuID string) (*InPersonSku, error) {
	req, err := s.client.newRequest(ctx, "GET", "in-person/skus/"+skuID, nil)
	if err != nil {
		return nil, err
	}

	var sku InPersonSku
	if err := s.client.doJSON(req, &sku); err != nil {
		return nil, err
	}

	return &sku, nil
}

func (s *InPersonOrdersService) List(
	ctx context.Context,
	ascending bool,
	limit int,
	startsAfter string,
) (*InPersonOrderList, error) {
	return s.ListForAccount(ctx, ascending, limit, startsAfter, "")
}

func (s *InPersonOrdersService) ListForAccount(
	ctx context.Context,
	ascending bool,
	limit int,
	startsAfter string,
	accountID string,
) (*InPersonOrderList, error) {
	query := buildListQuery(ascending, limit, startsAfter)
	req, err := s.client.newRequestWithQuery(ctx, "GET", "in-person/orders", query, nil)
	if err != nil {
		return nil, err
	}
	if accountID != "" {
		req.Header.Set("Account", accountID)
	}

	var orders InPersonOrderList
	if err := s.client.doJSON(req, &orders); err != nil {
		return nil, err
	}

	return &orders, nil
}

func (s *InPersonOrdersService) Get(ctx context.Context, orderID string) (*InPersonOrder, error) {
	return s.GetForAccount(ctx, orderID, "")
}

func (s *InPersonOrdersService) GetForAccount(
	ctx context.Context,
	orderID string,
	accountID string,
) (*InPersonOrder, error) {
	req, err := s.client.newRequest(ctx, "GET", "in-person/orders/"+orderID, nil)
	if err != nil {
		return nil, err
	}
	if accountID != "" {
		req.Header.Set("Account", accountID)
	}

	var order InPersonOrder
	if err := s.client.doJSON(req, &order); err != nil {
		return nil, err
	}

	return &order, nil
}

func (s *InPersonLocationsService) Create(
	ctx context.Context,
	request CreateInPersonLocationRequest,
) (*InPersonLocation, error) {
	return s.CreateForAccount(ctx, request, "")
}

func (s *InPersonLocationsService) CreateForAccount(
	ctx context.Context,
	request CreateInPersonLocationRequest,
	accountID string,
) (*InPersonLocation, error) {
	req, err := s.client.newRequest(ctx, "POST", "in-person/locations", request)
	if err != nil {
		return nil, err
	}
	if accountID != "" {
		req.Header.Set("Account", accountID)
	}

	var location InPersonLocation
	if err := s.client.doJSON(req, &location); err != nil {
		return nil, err
	}

	return &location, nil
}

func (s *InPersonLocationsService) List(
	ctx context.Context,
	ascending bool,
	limit int,
	startsAfter string,
) (*InPersonLocationList, error) {
	return s.ListForAccount(ctx, ascending, limit, startsAfter, "")
}

func (s *InPersonLocationsService) ListForAccount(
	ctx context.Context,
	ascending bool,
	limit int,
	startsAfter string,
	accountID string,
) (*InPersonLocationList, error) {
	query := buildListQuery(ascending, limit, startsAfter)
	req, err := s.client.newRequestWithQuery(ctx, "GET", "in-person/locations", query, nil)
	if err != nil {
		return nil, err
	}
	if accountID != "" {
		req.Header.Set("Account", accountID)
	}

	var locations InPersonLocationList
	if err := s.client.doJSON(req, &locations); err != nil {
		return nil, err
	}

	return &locations, nil
}

func (s *InPersonLocationsService) Get(ctx context.Context, locationID string) (*InPersonLocation, error) {
	return s.GetForAccount(ctx, locationID, "")
}

func (s *InPersonLocationsService) GetForAccount(
	ctx context.Context,
	locationID string,
	accountID string,
) (*InPersonLocation, error) {
	req, err := s.client.newRequest(ctx, "GET", "in-person/locations/"+locationID, nil)
	if err != nil {
		return nil, err
	}
	if accountID != "" {
		req.Header.Set("Account", accountID)
	}

	var location InPersonLocation
	if err := s.client.doJSON(req, &location); err != nil {
		return nil, err
	}

	return &location, nil
}

func (s *InPersonLocationsService) Update(
	ctx context.Context,
	locationID string,
	request UpdateInPersonLocationRequest,
) (*InPersonLocation, error) {
	return s.UpdateForAccount(ctx, locationID, request, "")
}

func (s *InPersonLocationsService) UpdateForAccount(
	ctx context.Context,
	locationID string,
	request UpdateInPersonLocationRequest,
	accountID string,
) (*InPersonLocation, error) {
	req, err := s.client.newRequest(ctx, "PATCH", "in-person/locations/"+locationID, request)
	if err != nil {
		return nil, err
	}
	if accountID != "" {
		req.Header.Set("Account", accountID)
	}

	var location InPersonLocation
	if err := s.client.doJSON(req, &location); err != nil {
		return nil, err
	}

	return &location, nil
}

func (s *InPersonLocationsService) Delete(ctx context.Context, locationID string) (*DeletedResource, error) {
	return s.DeleteForAccount(ctx, locationID, "")
}

func (s *InPersonLocationsService) DeleteForAccount(
	ctx context.Context,
	locationID string,
	accountID string,
) (*DeletedResource, error) {
	req, err := s.client.newRequest(ctx, "DELETE", "in-person/locations/"+locationID, nil)
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

func (s *InPersonTerminalsService) Create(
	ctx context.Context,
	request CreateTerminalRequest,
) (*Terminal, error) {
	return s.CreateForAccount(ctx, request, "")
}

func (s *InPersonTerminalsService) CreateForAccount(
	ctx context.Context,
	request CreateTerminalRequest,
	accountID string,
) (*Terminal, error) {
	req, err := s.client.newRequest(ctx, "POST", "in-person/terminals", request)
	if err != nil {
		return nil, err
	}
	if accountID != "" {
		req.Header.Set("Account", accountID)
	}

	var terminal Terminal
	if err := s.client.doJSON(req, &terminal); err != nil {
		return nil, err
	}

	return &terminal, nil
}

func (s *InPersonTerminalsService) List(
	ctx context.Context,
	ascending bool,
	limit int,
	startsAfter string,
) (*TerminalList, error) {
	return s.ListForAccount(ctx, ascending, limit, startsAfter, "")
}

func (s *InPersonTerminalsService) ListForAccount(
	ctx context.Context,
	ascending bool,
	limit int,
	startsAfter string,
	accountID string,
) (*TerminalList, error) {
	query := buildListQuery(ascending, limit, startsAfter)
	req, err := s.client.newRequestWithQuery(ctx, "GET", "in-person/terminals", query, nil)
	if err != nil {
		return nil, err
	}
	if accountID != "" {
		req.Header.Set("Account", accountID)
	}

	var terminals TerminalList
	if err := s.client.doJSON(req, &terminals); err != nil {
		return nil, err
	}

	return &terminals, nil
}

func (s *InPersonTerminalsService) Get(ctx context.Context, terminalID string) (*Terminal, error) {
	return s.GetForAccount(ctx, terminalID, "")
}

func (s *InPersonTerminalsService) GetForAccount(
	ctx context.Context,
	terminalID string,
	accountID string,
) (*Terminal, error) {
	req, err := s.client.newRequest(ctx, "GET", "in-person/terminals/"+terminalID, nil)
	if err != nil {
		return nil, err
	}
	if accountID != "" {
		req.Header.Set("Account", accountID)
	}

	var terminal Terminal
	if err := s.client.doJSON(req, &terminal); err != nil {
		return nil, err
	}

	return &terminal, nil
}

func (s *InPersonTerminalsService) Update(
	ctx context.Context,
	terminalID string,
	request UpdateTerminalRequest,
) (*Terminal, error) {
	return s.UpdateForAccount(ctx, terminalID, request, "")
}

func (s *InPersonTerminalsService) UpdateForAccount(
	ctx context.Context,
	terminalID string,
	request UpdateTerminalRequest,
	accountID string,
) (*Terminal, error) {
	req, err := s.client.newRequest(ctx, "PATCH", "in-person/terminals/"+terminalID, request)
	if err != nil {
		return nil, err
	}
	if accountID != "" {
		req.Header.Set("Account", accountID)
	}

	var terminal Terminal
	if err := s.client.doJSON(req, &terminal); err != nil {
		return nil, err
	}

	return &terminal, nil
}

func (s *InPersonTerminalsService) Delete(ctx context.Context, terminalID string) (*DeletedResource, error) {
	return s.DeleteForAccount(ctx, terminalID, "")
}

func (s *InPersonTerminalsService) DeleteForAccount(
	ctx context.Context,
	terminalID string,
	accountID string,
) (*DeletedResource, error) {
	req, err := s.client.newRequest(ctx, "DELETE", "in-person/terminals/"+terminalID, nil)
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

func (s *InPersonTerminalsService) InitiatePayment(
	ctx context.Context,
	terminalID string,
	request TerminalPaymentRequest,
) (*Terminal, error) {
	return s.InitiatePaymentForAccount(ctx, terminalID, request, "")
}

func (s *InPersonTerminalsService) InitiatePaymentForAccount(
	ctx context.Context,
	terminalID string,
	request TerminalPaymentRequest,
	accountID string,
) (*Terminal, error) {
	req, err := s.client.newRequest(ctx, "POST", "in-person/terminals/"+terminalID+"/payment", request)
	if err != nil {
		return nil, err
	}
	if accountID != "" {
		req.Header.Set("Account", accountID)
	}

	var terminal Terminal
	if err := s.client.doJSON(req, &terminal); err != nil {
		return nil, err
	}

	return &terminal, nil
}

func (s *InPersonTerminalsService) InitiateRefund(
	ctx context.Context,
	terminalID string,
	request TerminalRefundRequest,
) (*Terminal, error) {
	return s.InitiateRefundForAccount(ctx, terminalID, request, "")
}

func (s *InPersonTerminalsService) InitiateRefundForAccount(
	ctx context.Context,
	terminalID string,
	request TerminalRefundRequest,
	accountID string,
) (*Terminal, error) {
	req, err := s.client.newRequest(ctx, "POST", "in-person/terminals/"+terminalID+"/refund", request)
	if err != nil {
		return nil, err
	}
	if accountID != "" {
		req.Header.Set("Account", accountID)
	}

	var terminal Terminal
	if err := s.client.doJSON(req, &terminal); err != nil {
		return nil, err
	}

	return &terminal, nil
}

func (s *InPersonTerminalsService) CancelAction(ctx context.Context, terminalID string) (*Terminal, error) {
	return s.CancelActionForAccount(ctx, terminalID, "")
}

func (s *InPersonTerminalsService) CancelActionForAccount(
	ctx context.Context,
	terminalID string,
	accountID string,
) (*Terminal, error) {
	req, err := s.client.newRequest(ctx, "POST", "in-person/terminals/"+terminalID+"/cancel-action", nil)
	if err != nil {
		return nil, err
	}
	if accountID != "" {
		req.Header.Set("Account", accountID)
	}

	var terminal Terminal
	if err := s.client.doJSON(req, &terminal); err != nil {
		return nil, err
	}

	return &terminal, nil
}

func (s *InPersonTerminalsService) ConfirmReceipt(
	ctx context.Context,
	terminalID string,
	request TerminalConfirmReceiptRequest,
) (*Terminal, error) {
	return s.ConfirmReceiptForAccount(ctx, terminalID, request, "")
}

func (s *InPersonTerminalsService) ConfirmReceiptForAccount(
	ctx context.Context,
	terminalID string,
	request TerminalConfirmReceiptRequest,
	accountID string,
) (*Terminal, error) {
	req, err := s.client.newRequest(ctx, "POST", "in-person/terminals/"+terminalID+"/confirm-receipt", request)
	if err != nil {
		return nil, err
	}
	if accountID != "" {
		req.Header.Set("Account", accountID)
	}

	var terminal Terminal
	if err := s.client.doJSON(req, &terminal); err != nil {
		return nil, err
	}

	return &terminal, nil
}
