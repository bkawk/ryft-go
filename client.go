package ryft

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	sandboxBaseURL = "https://sandbox-api.ryftpay.com/v1"
	liveBaseURL    = "https://api.ryftpay.com/v1"
)

type Config struct {
	SecretKey  string
	BaseURL    string
	HTTPClient *http.Client
}

type Client struct {
	baseURL    string
	secretKey  string
	httpClient *http.Client

	Customers           *CustomersService
	PaymentSessions     *PaymentSessionsService
	Webhooks            *WebhooksService
	Accounts            *AccountsService
	Persons             *PersonsService
	PayoutMethods       *PayoutMethodsService
	Payouts             *PayoutsService
	Transfers           *TransfersService
	Balances            *BalancesService
	BalanceTransactions *BalanceTransactionsService
	PaymentMethods      *PaymentMethodsService
	AccountLinks        *AccountLinksService
	Subscriptions       *SubscriptionsService
	Events              *EventsService
	PlatformFees        *PlatformFeesService
	Files               *FilesService
	Disputes            *DisputesService
}

func NewClient(config Config) (*Client, error) {
	if strings.TrimSpace(config.SecretKey) == "" {
		return nil, errors.New("secret key is required")
	}

	baseURL := strings.TrimRight(strings.TrimSpace(config.BaseURL), "/")
	if baseURL == "" {
		var err error
		baseURL, err = determineBaseURL(config.SecretKey)
		if err != nil {
			return nil, err
		}
	}

	httpClient := config.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 30 * time.Second}
	}

	client := &Client{
		baseURL:    baseURL,
		secretKey:  config.SecretKey,
		httpClient: httpClient,
	}
	client.Customers = &CustomersService{client: client}
	client.PaymentSessions = &PaymentSessionsService{client: client}
	client.Webhooks = &WebhooksService{client: client}
	client.Accounts = &AccountsService{client: client}
	client.Persons = &PersonsService{client: client}
	client.PayoutMethods = &PayoutMethodsService{client: client}
	client.Payouts = &PayoutsService{client: client}
	client.Transfers = &TransfersService{client: client}
	client.Balances = &BalancesService{client: client}
	client.BalanceTransactions = &BalanceTransactionsService{client: client}
	client.PaymentMethods = &PaymentMethodsService{client: client}
	client.AccountLinks = &AccountLinksService{client: client}
	client.Subscriptions = &SubscriptionsService{client: client}
	client.Events = &EventsService{client: client}
	client.PlatformFees = &PlatformFeesService{client: client}
	client.Files = &FilesService{client: client}
	client.Disputes = &DisputesService{client: client}

	return client, nil
}

func determineBaseURL(secretKey string) (string, error) {
	switch {
	case strings.HasPrefix(secretKey, "sk_sandbox"):
		return sandboxBaseURL, nil
	case strings.HasPrefix(secretKey, "sk_"):
		return liveBaseURL, nil
	default:
		return "", errors.New("invalid secret key: expected prefix 'sk_'")
	}
}

func (c *Client) newRequest(ctx context.Context, method string, path string, body any) (*http.Request, error) {
	return c.newRequestWithQuery(ctx, method, path, nil, body)
}

func (c *Client) newRequestWithQuery(
	ctx context.Context,
	method string,
	path string,
	query url.Values,
	body any,
) (*http.Request, error) {
	endpoint, err := url.JoinPath(c.baseURL, path)
	if err != nil {
		return nil, fmt.Errorf("build request url: %w", err)
	}

	if len(query) > 0 {
		endpoint = endpoint + "?" + query.Encode()
	}

	var reader io.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
		reader = bytes.NewReader(payload)
	}

	req, err := http.NewRequestWithContext(ctx, method, endpoint, reader)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", c.secretKey)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

func (c *Client) doJSON(req *http.Request, out any) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return parseAPIError(resp.StatusCode, body)
	}

	if out == nil || len(body) == 0 {
		return nil
	}

	if err := json.Unmarshal(body, out); err != nil {
		return fmt.Errorf("decode response body: %w", err)
	}

	return nil
}

func (c *Client) doMultipartFile(
	ctx context.Context,
	path string,
	accountID string,
	filePath string,
	category string,
	out any,
) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	contentType := mime.TypeByExtension(strings.ToLower(filepath.Ext(filePath)))
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	headers := make(textproto.MIMEHeader)
	headers.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, filepath.Base(filePath)))
	headers.Set("Content-Type", contentType)

	part, err := writer.CreatePart(headers)
	if err != nil {
		return fmt.Errorf("create multipart file part: %w", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return fmt.Errorf("copy multipart file: %w", err)
	}
	if err := writer.WriteField("category", category); err != nil {
		return fmt.Errorf("write category field: %w", err)
	}
	if err := writer.Close(); err != nil {
		return fmt.Errorf("close multipart writer: %w", err)
	}

	endpoint, err := url.JoinPath(c.baseURL, path)
	if err != nil {
		return fmt.Errorf("build request url: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, &body)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", c.secretKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if accountID != "" {
		req.Header.Set("Account", accountID)
	}

	return c.doJSON(req, out)
}

func buildListQuery(ascending bool, limit int, startsAfter string) url.Values {
	query := url.Values{}
	query.Set("ascending", boolString(ascending))
	if limit > 0 {
		query.Set("limit", itoa(limit))
	}
	if startsAfter != "" {
		query.Set("startsAfter", startsAfter)
	}
	return query
}
