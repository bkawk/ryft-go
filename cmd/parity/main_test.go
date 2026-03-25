package main

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/bkawk/ryft-go"
)

func newParityClient(t *testing.T, handler http.HandlerFunc) *ryft.Client {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	client, err := ryft.NewClient(ryft.Config{
		SecretKey: "sk_sandbox_123",
		BaseURL:   server.URL,
	})
	if err != nil {
		t.Fatalf("NewClient returned error: %v", err)
	}

	return client
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	originalStdout := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe returned error: %v", err)
	}

	os.Stdout = writer
	defer func() {
		os.Stdout = originalStdout
	}()

	fn()

	if err := writer.Close(); err != nil {
		t.Fatalf("writer.Close returned error: %v", err)
	}

	output, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("ReadAll returned error: %v", err)
	}

	return string(output)
}

func TestCustomerCreateRequest(t *testing.T) {
	t.Parallel()

	request, err := customerCreateRequest([]string{"jane@example.com", "Jane", "Doe", `{"source":"test","attempt":2}`})
	if err != nil {
		t.Fatalf("customerCreateRequest returned error: %v", err)
	}

	if request.Email != "jane@example.com" {
		t.Fatalf("Email = %q, want %q", request.Email, "jane@example.com")
	}
	if request.FirstName != "Jane" {
		t.Fatalf("FirstName = %q, want %q", request.FirstName, "Jane")
	}
	if request.LastName != "Doe" {
		t.Fatalf("LastName = %q, want %q", request.LastName, "Doe")
	}
	if got := request.Metadata["attempt"]; got != "2" {
		t.Fatalf("Metadata[attempt] = %q, want %q", got, "2")
	}
}

func TestCustomerUpdateRequest(t *testing.T) {
	t.Parallel()

	request, err := customerUpdateRequest([]string{"Jane", "Doe", `{"segment":"beta"}`})
	if err != nil {
		t.Fatalf("customerUpdateRequest returned error: %v", err)
	}

	if request.FirstName != "Jane" || request.LastName != "Doe" {
		t.Fatalf("request = %#v, want first/last names set", request)
	}
	if got := request.Metadata["segment"]; got != "beta" {
		t.Fatalf("Metadata[segment] = %q, want %q", got, "beta")
	}
}

func TestPaymentSessionCreateRequest(t *testing.T) {
	t.Parallel()

	request, err := paymentSessionCreateRequest([]string{"2500", "GBP", "buyer@example.com", "Ecommerce", "CustomerEntry", "Automatic", `{"flow":"checkout"}`})
	if err != nil {
		t.Fatalf("paymentSessionCreateRequest returned error: %v", err)
	}

	if request.Amount != 2500 {
		t.Fatalf("Amount = %d, want 2500", request.Amount)
	}
	if request.ReturnURL != "https://example.com/return" {
		t.Fatalf("ReturnURL = %q, want default return URL", request.ReturnURL)
	}
	if got := request.Metadata["flow"]; got != "checkout" {
		t.Fatalf("Metadata[flow] = %q, want %q", got, "checkout")
	}
}

func TestParseMetadata(t *testing.T) {
	t.Parallel()

	metadata, err := parseMetadata(`{"attempt":1,"owner":"ryft-go"}`)
	if err != nil {
		t.Fatalf("parseMetadata returned error: %v", err)
	}

	want := map[string]string{
		"attempt": "1",
		"owner":   "ryft-go",
	}
	if !reflect.DeepEqual(metadata, want) {
		t.Fatalf("parseMetadata = %#v, want %#v", metadata, want)
	}
}

func TestPaymentSessionUpdateRequest(t *testing.T) {
	t.Parallel()

	request, err := paymentSessionUpdateRequest([]string{"775", "updated@example.com", "Automatic", `{"scenario":"update"}`})
	if err != nil {
		t.Fatalf("paymentSessionUpdateRequest returned error: %v", err)
	}

	if request.Amount == nil || *request.Amount != 775 {
		t.Fatalf("Amount = %#v, want 775", request.Amount)
	}
	if request.CustomerEmail != "updated@example.com" {
		t.Fatalf("CustomerEmail = %q, want %q", request.CustomerEmail, "updated@example.com")
	}
	if request.CaptureFlow != "Automatic" {
		t.Fatalf("CaptureFlow = %q, want %q", request.CaptureFlow, "Automatic")
	}
	if got := request.Metadata["scenario"]; got != "update" {
		t.Fatalf("Metadata[scenario] = %q, want %q", got, "update")
	}
}

func TestWebhookUpdateRequest(t *testing.T) {
	t.Parallel()

	request, err := webhookUpdateRequest([]string{"https://example.com/updated", "false", `["PaymentSession.captured"]`})
	if err != nil {
		t.Fatalf("webhookUpdateRequest returned error: %v", err)
	}

	if request.URL != "https://example.com/updated" {
		t.Fatalf("URL = %q, want %q", request.URL, "https://example.com/updated")
	}
	if request.Active == nil || *request.Active != false {
		t.Fatalf("Active = %#v, want false", request.Active)
	}
	if len(request.EventTypes) != 1 || request.EventTypes[0] != "PaymentSession.captured" {
		t.Fatalf("EventTypes = %#v, want [PaymentSession.captured]", request.EventTypes)
	}
}

func TestRefundPaymentSessionRequest(t *testing.T) {
	t.Parallel()

	request, err := refundPaymentSessionRequest([]string{`{"amount":300,"reason":"RequestedByCustomer","refundPlatformFee":true}`})
	if err != nil {
		t.Fatalf("refundPaymentSessionRequest returned error: %v", err)
	}

	if request.Amount != 300 {
		t.Fatalf("Amount = %d, want 300", request.Amount)
	}
	if request.Reason != "RequestedByCustomer" {
		t.Fatalf("Reason = %q, want %q", request.Reason, "RequestedByCustomer")
	}
	if request.RefundPlatformFee != true {
		t.Fatal("RefundPlatformFee = false, want true")
	}
}

func TestParseIntArgError(t *testing.T) {
	t.Parallel()

	if _, err := parseIntArg("oops", "amount"); err == nil {
		t.Fatal("parseIntArg error = nil, want error")
	}
}

func TestWebhookCreateRequest(t *testing.T) {
	t.Parallel()

	request, err := webhookCreateRequest([]string{"https://example.com/webhook", "True", `["PaymentSession.captured","PaymentSession.failed"]`})
	if err != nil {
		t.Fatalf("webhookCreateRequest returned error: %v", err)
	}

	if !request.Active {
		t.Fatal("Active = false, want true")
	}
	if len(request.EventTypes) != 2 {
		t.Fatalf("EventTypes = %#v, want 2 values", request.EventTypes)
	}
}

func TestAccountCreateRequestIndividualAndInvalidBusiness(t *testing.T) {
	t.Parallel()

	individual, err := accountCreateRequest([]string{"Individual", "person@example.com", `{"team":"onboarding"}`, "", "Hosted"})
	if err != nil {
		t.Fatalf("accountCreateRequest individual returned error: %v", err)
	}
	if individual.OnboardingFlow != "Hosted" {
		t.Fatalf("OnboardingFlow = %q, want %q", individual.OnboardingFlow, "Hosted")
	}
	if individual.Individual == nil || individual.Business != nil {
		t.Fatalf("individual request = %#v, want Individual details only", individual)
	}
	if got := individual.Metadata["team"]; got != "onboarding" {
		t.Fatalf("Metadata[team] = %q, want %q", got, "onboarding")
	}

	invalidBusiness, err := accountCreateRequest([]string{"Individual", "person@example.com", "", "invalid-individual-business"})
	if err != nil {
		t.Fatalf("accountCreateRequest invalid business returned error: %v", err)
	}
	if invalidBusiness.Business == nil {
		t.Fatal("Business = nil, want business details for invalid-individual-business mode")
	}
}

func TestPersonCreateRequest(t *testing.T) {
	t.Parallel()

	request, err := personCreateRequest([]string{"jane@example.com", `{"region":"uk"}`})
	if err != nil {
		t.Fatalf("personCreateRequest returned error: %v", err)
	}

	if request.Email != "jane@example.com" {
		t.Fatalf("Email = %q, want %q", request.Email, "jane@example.com")
	}
	if got := request.Metadata["region"]; got != "uk" {
		t.Fatalf("Metadata[region] = %q, want %q", got, "uk")
	}
}

func TestPayoutMethodCreateRequest(t *testing.T) {
	t.Parallel()

	request := payoutMethodCreateRequest([]string{"Primary Account"})
	if request.DisplayName != "Primary Account" {
		t.Fatalf("DisplayName = %q, want %q", request.DisplayName, "Primary Account")
	}
	if request.BankAccount.BankID != "601613" {
		t.Fatalf("BankID = %q, want %q", request.BankAccount.BankID, "601613")
	}
}

func TestPayoutCreateRequest(t *testing.T) {
	t.Parallel()

	request, err := payoutCreateRequest([]string{"700", "GBP", "pmo_123", `{"batch":"march"}`})
	if err != nil {
		t.Fatalf("payoutCreateRequest returned error: %v", err)
	}

	if request.Amount != 700 {
		t.Fatalf("Amount = %d, want 700", request.Amount)
	}
	if got := request.Metadata["batch"]; got != "march" {
		t.Fatalf("Metadata[batch] = %q, want %q", got, "march")
	}
}

func TestTransferCreateRequest(t *testing.T) {
	t.Parallel()

	request, err := transferCreateRequest([]string{"ac_123", "900", "GBP", `{"flow":"split"}`})
	if err != nil {
		t.Fatalf("transferCreateRequest returned error: %v", err)
	}

	if request.Destination.AccountID != "ac_123" {
		t.Fatalf("Destination.AccountID = %q, want %q", request.Destination.AccountID, "ac_123")
	}
	if request.Reason == "" {
		t.Fatal("Reason is empty")
	}
	if got := request.Metadata["flow"]; got != "split" {
		t.Fatalf("Metadata[flow] = %q, want %q", got, "split")
	}
}

func TestSubscriptionCreateRequest(t *testing.T) {
	t.Parallel()

	request, err := subscriptionCreateRequest([]string{
		"cus_123",
		"pm_123",
		`{"description":"Gold","billingCycleTimestamp":1710000000,"metadata":{"tier":"gold"},"price":{"amount":999,"currency":"EUR","interval":{"unit":"Weeks","count":2,"times":6}},"paymentSettings":{"captureFlow":"Automatic"},"shippingDetails":{"address":{"country":"GB"}}}`,
	})
	if err != nil {
		t.Fatalf("subscriptionCreateRequest returned error: %v", err)
	}

	if request.Customer.ID != "cus_123" || request.PaymentMethod.ID != "pm_123" {
		t.Fatalf("request = %#v, want customer/payment method ids set", request)
	}
	if request.Description != "Gold" {
		t.Fatalf("Description = %q, want %q", request.Description, "Gold")
	}
	if request.Price.Amount != 999 || request.Price.Currency != "EUR" {
		t.Fatalf("Price = %#v, want EUR 999", request.Price)
	}
	if request.Price.Interval.Unit != "Weeks" || request.Price.Interval.Count != 2 || request.Price.Interval.Times != 6 {
		t.Fatalf("Interval = %#v, want Weeks/2/6", request.Price.Interval)
	}
	if request.BillingCycleTimestamp != 1710000000 {
		t.Fatalf("BillingCycleTimestamp = %d, want 1710000000", request.BillingCycleTimestamp)
	}
	if got := request.Metadata["tier"]; got != "gold" {
		t.Fatalf("Metadata[tier] = %q, want %q", got, "gold")
	}
	if value, ok := request.PaymentSettings["captureFlow"].(string); !ok || value != "Automatic" {
		t.Fatalf("PaymentSettings = %#v, want captureFlow Automatic", request.PaymentSettings)
	}
}

func TestSubscriptionUpdateAndPauseRequests(t *testing.T) {
	t.Parallel()

	updateRequest, err := subscriptionUpdateRequest([]string{"Updated plan", `{"phase":"trial"}`})
	if err != nil {
		t.Fatalf("subscriptionUpdateRequest returned error: %v", err)
	}
	if updateRequest.Description != "Updated plan" {
		t.Fatalf("Description = %q, want %q", updateRequest.Description, "Updated plan")
	}
	if got := updateRequest.Metadata["phase"]; got != "trial" {
		t.Fatalf("Metadata[phase] = %q, want %q", got, "trial")
	}

	pauseRequest, err := subscriptionPauseRequest([]string{"customer_request", "1710000000", "TRUE"})
	if err != nil {
		t.Fatalf("subscriptionPauseRequest returned error: %v", err)
	}
	if pauseRequest.Reason != "customer_request" || pauseRequest.ResumeTimestamp != 1710000000 || !pauseRequest.Unschedule {
		t.Fatalf("pauseRequest = %#v, want full pause request", pauseRequest)
	}
}

func TestCollectionWindowAndHelpers(t *testing.T) {
	t.Setenv("RYFT_COLLECTION_START_TIMESTAMP", "100")
	t.Setenv("RYFT_COLLECTION_END_TIMESTAMP", "200")

	start, end := collectionWindowFromEnv()
	if start != 100 || end != 200 {
		t.Fatalf("collectionWindowFromEnv = (%d, %d), want (100, 200)", start, end)
	}

	if got := intFromAny(float64(12), 0); got != 12 {
		t.Fatalf("intFromAny(float64) = %d, want 12", got)
	}
	if got := intFromAny(int64(13), 0); got != 13 {
		t.Fatalf("intFromAny(int64) = %d, want 13", got)
	}
	if got := intFromAny("bad", 9); got != 9 {
		t.Fatalf("intFromAny(fallback) = %d, want 9", got)
	}

	if got := stringFromAny("value", "fallback"); got != "value" {
		t.Fatalf("stringFromAny(value) = %q, want %q", got, "value")
	}
	if got := stringFromAny("", "fallback"); got != "fallback" {
		t.Fatalf("stringFromAny(fallback) = %q, want %q", got, "fallback")
	}

	raw := map[string]any{
		"interval": map[string]any{
			"count": 3,
		},
	}
	if got := mapValue(raw, "interval", "count"); got != 3 {
		t.Fatalf("mapValue = %#v, want 3", got)
	}
	if got := mapValue(raw, "missing", "count"); got != nil {
		t.Fatalf("mapValue missing = %#v, want nil", got)
	}
}

func TestPrintJSON(t *testing.T) {
	t.Parallel()

	tempFile, err := os.CreateTemp(t.TempDir(), "print-json-*.json")
	if err != nil {
		t.Fatalf("CreateTemp returned error: %v", err)
	}
	defer tempFile.Close()

	printJSON(tempFile, map[string]string{"id": "cus_123"})

	if _, err := tempFile.Seek(0, 0); err != nil {
		t.Fatalf("Seek returned error: %v", err)
	}
	var buffer bytes.Buffer
	if _, err := buffer.ReadFrom(tempFile); err != nil {
		t.Fatalf("ReadFrom returned error: %v", err)
	}
	if got := buffer.String(); got == "" || got[0] != '{' {
		t.Fatalf("printJSON output = %q, want json content", got)
	}
}

func TestRunUsageAndValidationErrors(t *testing.T) {
	t.Setenv("RYFT_SECRET_KEY", "sk_sandbox_123")

	testCases := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{
			name:    "no args",
			args:    nil,
			wantErr: "usage: parity <customer-create|customer-update|entity-get> ...",
		},
		{
			name:    "unknown command",
			args:    []string{"unknown"},
			wantErr: "unknown command: unknown",
		},
		{
			name:    "customer create usage",
			args:    []string{"customer-create"},
			wantErr: "usage: parity customer-create <email> [first-name] [last-name] [metadata-json]",
		},
		{
			name:    "payment session update invalid amount",
			args:    []string{"payment-session-update", "ps_123", "oops"},
			wantErr: "parse amount:",
		},
		{
			name:    "subscription pause invalid resume timestamp",
			args:    []string{"subscription-pause", "sub_123", "reason", "oops"},
			wantErr: "parse resume-timestamp:",
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			err := run(context.Background(), testCase.args)
			if err == nil {
				t.Fatal("run error = nil, want error")
			}
			if !strings.Contains(err.Error(), testCase.wantErr) {
				t.Fatalf("run error = %q, want substring %q", err.Error(), testCase.wantErr)
			}
		})
	}
}

func TestHandleEntityGetValidationErrors(t *testing.T) {
	t.Parallel()

	client := newParityClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected request: %s %s", r.Method, r.URL.String())
	})

	testCases := []struct {
		name       string
		entityType string
		entityID   string
		parentID   string
		wantErr    string
	}{
		{
			name:       "person requires parent",
			entityType: "person",
			entityID:   "per_123",
			wantErr:    "person requires parent account id",
		},
		{
			name:       "payout method requires parent",
			entityType: "payout-method",
			entityID:   "pmo_123",
			wantErr:    "payout-method requires parent account id",
		},
		{
			name:       "payment transaction requires parent",
			entityType: "payment-transaction",
			entityID:   "txn_123",
			wantErr:    "payment-transaction requires parent payment-session id",
		},
		{
			name:       "unsupported type",
			entityType: "mystery",
			entityID:   "id_123",
			wantErr:    "unsupported entity type: mystery",
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			err := handleEntityGet(context.Background(), client, testCase.entityType, testCase.entityID, testCase.parentID)
			if err == nil {
				t.Fatal("handleEntityGet error = nil, want error")
			}
			if err.Error() != testCase.wantErr {
				t.Fatalf("handleEntityGet error = %q, want %q", err.Error(), testCase.wantErr)
			}
		})
	}
}

func TestHandleEntityGetSuccess(t *testing.T) {
	t.Parallel()

	client := newParityClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		case "/customers/cus_123":
			_, _ = io.WriteString(w, `{"id":"cus_123","email":"jane@example.com"}`)
		case "/events/evt_123":
			if got := r.Header.Get("Account"); got != "ac_123" {
				t.Fatalf("Account header = %q, want %q", got, "ac_123")
			}
			_, _ = io.WriteString(w, `{"id":"evt_123","accountId":"ac_123"}`)
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	})

	customerOutput := captureStdout(t, func() {
		if err := handleEntityGet(context.Background(), client, "customer", "cus_123", ""); err != nil {
			t.Fatalf("handleEntityGet customer returned error: %v", err)
		}
	})
	if !strings.Contains(customerOutput, `"id": "cus_123"`) {
		t.Fatalf("customer output = %q, want customer id", customerOutput)
	}

	eventOutput := captureStdout(t, func() {
		if err := handleEntityGet(context.Background(), client, "event", "evt_123", "ac_123"); err != nil {
			t.Fatalf("handleEntityGet event returned error: %v", err)
		}
	})
	if !strings.Contains(eventOutput, `"accountId": "ac_123"`) {
		t.Fatalf("event output = %q, want account id", eventOutput)
	}
}
