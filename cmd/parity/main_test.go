package main

import (
	"reflect"
	"testing"
)

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
