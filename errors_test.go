package ryft

import "testing"

func TestParseAPIError(t *testing.T) {
	t.Parallel()

	err := parseAPIError(409, []byte(`{"code":"conflict","message":"Already exists","requestId":"req_123","errors":[{"code":"duplicate","message":"Duplicate customer"}]}`))
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}

	if apiErr.Status != 409 {
		t.Fatalf("Status = %d, want 409", apiErr.Status)
	}
	if apiErr.Code != "conflict" {
		t.Fatalf("Code = %q, want %q", apiErr.Code, "conflict")
	}
	if apiErr.Message != "Already exists" {
		t.Fatalf("Message = %q, want %q", apiErr.Message, "Already exists")
	}
	if apiErr.RequestID != "req_123" {
		t.Fatalf("RequestID = %q, want %q", apiErr.RequestID, "req_123")
	}
	if len(apiErr.Errors) != 1 {
		t.Fatalf("Errors length = %d, want 1", len(apiErr.Errors))
	}
}

func TestParseAPIErrorFallsBackToRawBody(t *testing.T) {
	t.Parallel()

	err := parseAPIError(500, []byte("server exploded"))
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}

	if apiErr.Message != "server exploded" {
		t.Fatalf("Message = %q, want %q", apiErr.Message, "server exploded")
	}
}
