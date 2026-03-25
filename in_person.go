package ryft

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
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
	ID       string            `json:"id,omitempty"`
	Status   string            `json:"status,omitempty"`
	Amount   int               `json:"amount,omitempty"`
	Currency string            `json:"currency,omitempty"`
	Items    []json.RawMessage `json:"items,omitempty"`
	Customer json.RawMessage   `json:"customer,omitempty"`
	Shipping json.RawMessage   `json:"shipping,omitempty"`
	Tracking json.RawMessage   `json:"tracking,omitempty"`
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

type InPersonProductListParams struct {
	ListParams
}

type InPersonSkuListParams struct {
	ListParams
	Country   string
	ProductID string
}

type InPersonOrderListParams struct {
	ListParams
}

type InPersonLocationListParams struct {
	ListParams
}

type TerminalListParams struct {
	ListParams
}

func (s *InPersonProductsService) List(ctx context.Context, params InPersonProductListParams) (*InPersonProductList, error) {
	query := buildListQuery(params.ListParams)
	req, err := s.client.newRequestWithQuery(ctx, http.MethodGet, "in-person/products", query, nil)
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
	req, err := s.client.newRequest(ctx, http.MethodGet, "in-person/products/"+productID, nil)
	if err != nil {
		return nil, err
	}

	var product InPersonProduct
	if err := s.client.doJSON(req, &product); err != nil {
		return nil, err
	}

	return &product, nil
}

func (s *InPersonSkusService) List(ctx context.Context, params InPersonSkuListParams) (*InPersonSkuList, error) {
	query := url.Values{}
	if params.Country != "" {
		query.Set("country", params.Country)
	}
	if params.Limit > 0 {
		query.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.StartsAfter != "" {
		query.Set("startsAfter", params.StartsAfter)
	}
	if params.ProductID != "" {
		query.Set("productId", params.ProductID)
	}

	req, err := s.client.newRequestWithQuery(ctx, http.MethodGet, "in-person/skus", query, nil)
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
	req, err := s.client.newRequest(ctx, http.MethodGet, "in-person/skus/"+skuID, nil)
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
	params InPersonOrderListParams,
	opts ...RequestOption,
) (*InPersonOrderList, error) {
	query := buildListQuery(params.ListParams)
	req, err := s.client.newRequestWithQuery(ctx, http.MethodGet, "in-person/orders", query, nil)
	if err != nil {
		return nil, err
	}
	s.client.applyRequestOptions(req, opts...)

	var orders InPersonOrderList
	if err := s.client.doJSON(req, &orders); err != nil {
		return nil, err
	}

	return &orders, nil
}

func (s *InPersonOrdersService) Get(
	ctx context.Context,
	orderID string,
	opts ...RequestOption,
) (*InPersonOrder, error) {
	req, err := s.client.newRequest(ctx, http.MethodGet, "in-person/orders/"+orderID, nil)
	if err != nil {
		return nil, err
	}
	s.client.applyRequestOptions(req, opts...)

	var order InPersonOrder
	if err := s.client.doJSON(req, &order); err != nil {
		return nil, err
	}

	return &order, nil
}

func (s *InPersonLocationsService) Create(
	ctx context.Context,
	request CreateInPersonLocationRequest,
	opts ...RequestOption,
) (*InPersonLocation, error) {
	req, err := s.client.newRequest(ctx, http.MethodPost, "in-person/locations", request)
	if err != nil {
		return nil, err
	}
	s.client.applyRequestOptions(req, opts...)

	var location InPersonLocation
	if err := s.client.doJSON(req, &location); err != nil {
		return nil, err
	}

	return &location, nil
}

func (s *InPersonLocationsService) List(
	ctx context.Context,
	params InPersonLocationListParams,
	opts ...RequestOption,
) (*InPersonLocationList, error) {
	query := buildListQuery(params.ListParams)
	req, err := s.client.newRequestWithQuery(ctx, http.MethodGet, "in-person/locations", query, nil)
	if err != nil {
		return nil, err
	}
	s.client.applyRequestOptions(req, opts...)

	var locations InPersonLocationList
	if err := s.client.doJSON(req, &locations); err != nil {
		return nil, err
	}

	return &locations, nil
}

func (s *InPersonLocationsService) Get(
	ctx context.Context,
	locationID string,
	opts ...RequestOption,
) (*InPersonLocation, error) {
	req, err := s.client.newRequest(ctx, http.MethodGet, "in-person/locations/"+locationID, nil)
	if err != nil {
		return nil, err
	}
	s.client.applyRequestOptions(req, opts...)

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
	opts ...RequestOption,
) (*InPersonLocation, error) {
	req, err := s.client.newRequest(ctx, http.MethodPatch, "in-person/locations/"+locationID, request)
	if err != nil {
		return nil, err
	}
	s.client.applyRequestOptions(req, opts...)

	var location InPersonLocation
	if err := s.client.doJSON(req, &location); err != nil {
		return nil, err
	}

	return &location, nil
}

func (s *InPersonLocationsService) Delete(
	ctx context.Context,
	locationID string,
	opts ...RequestOption,
) (*DeletedResource, error) {
	req, err := s.client.newRequest(ctx, http.MethodDelete, "in-person/locations/"+locationID, nil)
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

func (s *InPersonTerminalsService) Create(
	ctx context.Context,
	request CreateTerminalRequest,
	opts ...RequestOption,
) (*Terminal, error) {
	req, err := s.client.newRequest(ctx, http.MethodPost, "in-person/terminals", request)
	if err != nil {
		return nil, err
	}
	s.client.applyRequestOptions(req, opts...)

	var terminal Terminal
	if err := s.client.doJSON(req, &terminal); err != nil {
		return nil, err
	}

	return &terminal, nil
}

func (s *InPersonTerminalsService) List(
	ctx context.Context,
	params TerminalListParams,
	opts ...RequestOption,
) (*TerminalList, error) {
	query := buildListQuery(params.ListParams)
	req, err := s.client.newRequestWithQuery(ctx, http.MethodGet, "in-person/terminals", query, nil)
	if err != nil {
		return nil, err
	}
	s.client.applyRequestOptions(req, opts...)

	var terminals TerminalList
	if err := s.client.doJSON(req, &terminals); err != nil {
		return nil, err
	}

	return &terminals, nil
}

func (s *InPersonTerminalsService) Get(
	ctx context.Context,
	terminalID string,
	opts ...RequestOption,
) (*Terminal, error) {
	req, err := s.client.newRequest(ctx, http.MethodGet, "in-person/terminals/"+terminalID, nil)
	if err != nil {
		return nil, err
	}
	s.client.applyRequestOptions(req, opts...)

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
	opts ...RequestOption,
) (*Terminal, error) {
	req, err := s.client.newRequest(ctx, http.MethodPatch, "in-person/terminals/"+terminalID, request)
	if err != nil {
		return nil, err
	}
	s.client.applyRequestOptions(req, opts...)

	var terminal Terminal
	if err := s.client.doJSON(req, &terminal); err != nil {
		return nil, err
	}

	return &terminal, nil
}

func (s *InPersonTerminalsService) Delete(
	ctx context.Context,
	terminalID string,
	opts ...RequestOption,
) (*DeletedResource, error) {
	req, err := s.client.newRequest(ctx, http.MethodDelete, "in-person/terminals/"+terminalID, nil)
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

func (s *InPersonTerminalsService) InitiatePayment(
	ctx context.Context,
	terminalID string,
	request TerminalPaymentRequest,
	opts ...RequestOption,
) (*Terminal, error) {
	req, err := s.client.newRequest(ctx, http.MethodPost, "in-person/terminals/"+terminalID+"/payment", request)
	if err != nil {
		return nil, err
	}
	s.client.applyRequestOptions(req, opts...)

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
	opts ...RequestOption,
) (*Terminal, error) {
	req, err := s.client.newRequest(ctx, http.MethodPost, "in-person/terminals/"+terminalID+"/refund", request)
	if err != nil {
		return nil, err
	}
	s.client.applyRequestOptions(req, opts...)

	var terminal Terminal
	if err := s.client.doJSON(req, &terminal); err != nil {
		return nil, err
	}

	return &terminal, nil
}

func (s *InPersonTerminalsService) CancelAction(
	ctx context.Context,
	terminalID string,
	opts ...RequestOption,
) (*Terminal, error) {
	req, err := s.client.newRequest(ctx, http.MethodPost, "in-person/terminals/"+terminalID+"/cancel-action", nil)
	if err != nil {
		return nil, err
	}
	s.client.applyRequestOptions(req, opts...)

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
	opts ...RequestOption,
) (*Terminal, error) {
	req, err := s.client.newRequest(ctx, http.MethodPost, "in-person/terminals/"+terminalID+"/confirm-receipt", request)
	if err != nil {
		return nil, err
	}
	s.client.applyRequestOptions(req, opts...)

	var terminal Terminal
	if err := s.client.doJSON(req, &terminal); err != nil {
		return nil, err
	}

	return &terminal, nil
}
