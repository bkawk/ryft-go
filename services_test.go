package ryft

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func newTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	client, err := NewClient(Config{
		SecretKey: "sk_sandbox_123",
		BaseURL:   server.URL,
	})
	if err != nil {
		t.Fatalf("NewClient returned error: %v", err)
	}

	return client
}

func decodeJSONBody(t *testing.T, r *http.Request) map[string]any {
	t.Helper()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("ReadAll returned error: %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}

	return payload
}

func assertQueryValues(t *testing.T, got url.Values, want map[string]string) {
	t.Helper()

	for key, wantValue := range want {
		if gotValue := got.Get(key); gotValue != wantValue {
			t.Fatalf("query[%q] = %q, want %q", key, gotValue, wantValue)
		}
	}
}

func TestBuildListQuery(t *testing.T) {
	t.Parallel()

	query := buildListQuery(ListParams{Ascending: true, Limit: 25, StartsAfter: "cus_123"})
	assertQueryValues(t, query, map[string]string{
		"ascending":   "true",
		"limit":       "25",
		"startsAfter": "cus_123",
	})

	empty := buildListQuery(ListParams{Ascending: false})
	assertQueryValues(t, empty, map[string]string{
		"ascending": "false",
	})
	if empty.Get("limit") != "" {
		t.Fatalf("limit = %q, want empty", empty.Get("limit"))
	}
}

func TestCustomersServiceMethods(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/customers":
			payload := decodeJSONBody(t, r)
			if payload["email"] != "jane@example.com" {
				t.Fatalf("email = %v, want %q", payload["email"], "jane@example.com")
			}
			_, _ = io.WriteString(w, `{"id":"cus_123","email":"jane@example.com"}`)
		case r.Method == http.MethodGet && r.URL.Path == "/customers":
			assertQueryValues(t, r.URL.Query(), map[string]string{
				"email":          "jane@example.com",
				"startTimestamp": "100",
				"endTimestamp":   "200",
				"ascending":      "true",
				"limit":          "5",
				"startsAfter":    "cus_001",
			})
			_, _ = io.WriteString(w, `{"items":[{"id":"cus_123","email":"jane@example.com"}]}`)
		case r.Method == http.MethodGet && r.URL.Path == "/customers/cus_123":
			_, _ = io.WriteString(w, `{"id":"cus_123","email":"jane@example.com"}`)
		case r.Method == http.MethodPatch && r.URL.Path == "/customers/cus_123":
			payload := decodeJSONBody(t, r)
			if payload["firstName"] != "Jane" {
				t.Fatalf("firstName = %v, want %q", payload["firstName"], "Jane")
			}
			_, _ = io.WriteString(w, `{"id":"cus_123","firstName":"Jane"}`)
		case r.Method == http.MethodDelete && r.URL.Path == "/customers/cus_123":
			_, _ = io.WriteString(w, `{"id":"cus_123"}`)
		case r.Method == http.MethodGet && r.URL.Path == "/customers/cus_123/payment-methods":
			_, _ = io.WriteString(w, `{"items":[{"id":"pm_123","customerId":"cus_123"}]}`)
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.String())
		}
	})

	ctx := context.Background()

	customer, err := client.Customers.Create(ctx, CreateCustomerRequest{Email: "jane@example.com"})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if customer.ID != "cus_123" {
		t.Fatalf("customer.ID = %q, want %q", customer.ID, "cus_123")
	}

	customers, err := client.Customers.List(ctx, CustomerListParams{ListParams: ListParams{Ascending: true, Limit: 5, StartsAfter: "cus_001"}, Email: "jane@example.com", StartTimestamp: 100, EndTimestamp: 200})
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(customers.Items) != 1 {
		t.Fatalf("len(customers.Items) = %d, want 1", len(customers.Items))
	}

	gotCustomer, err := client.Customers.Get(ctx, "cus_123")
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if gotCustomer.ID != "cus_123" {
		t.Fatalf("gotCustomer.ID = %q, want %q", gotCustomer.ID, "cus_123")
	}

	updatedCustomer, err := client.Customers.Update(ctx, "cus_123", UpdateCustomerRequest{FirstName: "Jane"})
	if err != nil {
		t.Fatalf("Update returned error: %v", err)
	}
	if updatedCustomer.FirstName != "Jane" {
		t.Fatalf("updatedCustomer.FirstName = %q, want %q", updatedCustomer.FirstName, "Jane")
	}

	deletedCustomer, err := client.Customers.Delete(ctx, "cus_123")
	if err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}
	if deletedCustomer.ID != "cus_123" {
		t.Fatalf("deletedCustomer.ID = %q, want %q", deletedCustomer.ID, "cus_123")
	}

	paymentMethods, err := client.Customers.GetPaymentMethods(ctx, "cus_123")
	if err != nil {
		t.Fatalf("GetPaymentMethods returned error: %v", err)
	}
	if len(paymentMethods.Items) != 1 {
		t.Fatalf("len(paymentMethods.Items) = %d, want 1", len(paymentMethods.Items))
	}
}

func TestAccountsAndAccountLinksMethods(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/accounts":
			payload := decodeJSONBody(t, r)
			if payload["entityType"] != "Business" {
				t.Fatalf("entityType = %v, want %q", payload["entityType"], "Business")
			}
			_, _ = io.WriteString(w, `{"id":"ac_123","entityType":"Business","email":"ops@example.com"}`)
		case r.Method == http.MethodGet && r.URL.Path == "/accounts/ac_123":
			_, _ = io.WriteString(w, `{"id":"ac_123","email":"ops@example.com"}`)
		case r.Method == http.MethodPatch && r.URL.Path == "/accounts/ac_123":
			payload := decodeJSONBody(t, r)
			metadata := payload["metadata"].(map[string]any)
			if metadata["segment"] != "beta" {
				t.Fatalf("metadata.segment = %v, want %q", metadata["segment"], "beta")
			}
			_, _ = io.WriteString(w, `{"id":"ac_123","metadata":{"segment":"beta"}}`)
		case r.Method == http.MethodPost && r.URL.Path == "/accounts/ac_123/verify":
			_, _ = io.WriteString(w, `{"id":"ac_123","entityType":"Business"}`)
		case r.Method == http.MethodPost && r.URL.Path == "/accounts/authorize":
			payload := decodeJSONBody(t, r)
			if payload["redirectUrl"] != "https://example.com/return" {
				t.Fatalf("redirectUrl = %v, want %q", payload["redirectUrl"], "https://example.com/return")
			}
			_, _ = io.WriteString(w, `{"url":"https://connect.example.com/auth"}`)
		case r.Method == http.MethodPost && r.URL.Path == "/account-links":
			payload := decodeJSONBody(t, r)
			if payload["accountId"] != "ac_123" {
				t.Fatalf("accountId = %v, want %q", payload["accountId"], WithAccount("ac_123"))
			}
			_, _ = io.WriteString(w, `{"url":"https://connect.example.com/link","createdTimestamp":1,"expiresTimestamp":2}`)
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.String())
		}
	})

	ctx := context.Background()

	account, err := client.Accounts.Create(ctx, CreateAccountRequest{
		OnboardingFlow: "Hosted",
		EntityType:     "Business",
		Email:          "ops@example.com",
		TermsOfService: TermsOfService{Acceptance: Acceptance{IPAddress: "127.0.0.1"}},
	})
	if err != nil {
		t.Fatalf("Accounts.Create returned error: %v", err)
	}
	if account.ID != "ac_123" {
		t.Fatalf("account.ID = %q, want %q", account.ID, "ac_123")
	}

	gotAccount, err := client.Accounts.Get(ctx, "ac_123")
	if err != nil {
		t.Fatalf("Accounts.Get returned error: %v", err)
	}
	if gotAccount.Email != "ops@example.com" {
		t.Fatalf("gotAccount.Email = %q, want %q", gotAccount.Email, "ops@example.com")
	}

	updatedAccount, err := client.Accounts.Update(ctx, "ac_123", UpdateAccountRequest{
		Metadata: map[string]string{"segment": "beta"},
	})
	if err != nil {
		t.Fatalf("Accounts.Update returned error: %v", err)
	}
	if updatedAccount.Metadata["segment"] != "beta" {
		t.Fatalf("updatedAccount.Metadata[segment] = %q, want %q", updatedAccount.Metadata["segment"], "beta")
	}

	verifiedAccount, err := client.Accounts.Verify(ctx, "ac_123")
	if err != nil {
		t.Fatalf("Accounts.Verify returned error: %v", err)
	}
	if verifiedAccount.ID != "ac_123" {
		t.Fatalf("verifiedAccount.ID = %q, want %q", verifiedAccount.ID, "ac_123")
	}

	authLink, err := client.Accounts.CreateAuthLink(ctx, CreateAccountAuthorizationRequest{
		Email:       "ops@example.com",
		RedirectURL: "https://example.com/return",
	})
	if err != nil {
		t.Fatalf("Accounts.CreateAuthLink returned error: %v", err)
	}
	if authLink.URL == "" {
		t.Fatal("authLink.URL is empty")
	}

	accountLink, err := client.AccountLinks.GenerateTemporaryAccountLink(ctx, CreateTemporaryAccountLinkRequest{
		AccountID:   "ac_123",
		RedirectURL: "https://example.com/return",
	})
	if err != nil {
		t.Fatalf("GenerateTemporaryAccountLink returned error: %v", err)
	}
	if accountLink.URL == "" {
		t.Fatal("accountLink.URL is empty")
	}
}

func TestApplePayMethods(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/apple-pay/web-domains":
			if got := r.Header.Get("Account"); got != "ac_123" {
				t.Fatalf("Account header = %q, want %q", got, "ac_123")
			}
			payload := decodeJSONBody(t, r)
			if payload["domainName"] != "checkout.example.com" {
				t.Fatalf("domainName = %v, want %q", payload["domainName"], "checkout.example.com")
			}
			_, _ = io.WriteString(w, `{"id":"apwd_123","domainName":"checkout.example.com"}`)
		case r.Method == http.MethodGet && r.URL.Path == "/apple-pay/web-domains":
			if got := r.Header.Get("Account"); got != "ac_123" {
				t.Fatalf("Account header = %q, want %q", got, "ac_123")
			}
			assertQueryValues(t, r.URL.Query(), map[string]string{
				"ascending":   "true",
				"limit":       "10",
				"startsAfter": "apwd_prev",
			})
			_, _ = io.WriteString(w, `{"items":[{"id":"apwd_123","domainName":"checkout.example.com"}]}`)
		case r.Method == http.MethodGet && r.URL.Path == "/apple-pay/web-domains/apwd_123":
			if got := r.Header.Get("Account"); got != "ac_123" {
				t.Fatalf("Account header = %q, want %q", got, "ac_123")
			}
			_, _ = io.WriteString(w, `{"id":"apwd_123","domainName":"checkout.example.com"}`)
		case r.Method == http.MethodDelete && r.URL.Path == "/apple-pay/web-domains/apwd_123":
			if got := r.Header.Get("Account"); got != "ac_123" {
				t.Fatalf("Account header = %q, want %q", got, "ac_123")
			}
			_, _ = io.WriteString(w, `{"id":"apwd_123"}`)
		case r.Method == http.MethodPost && r.URL.Path == "/apple-pay/sessions":
			if got := r.Header.Get("Account"); got != "ac_123" {
				t.Fatalf("Account header = %q, want %q", got, "ac_123")
			}
			payload := decodeJSONBody(t, r)
			if payload["displayName"] != "Ryft Demo" {
				t.Fatalf("displayName = %v, want %q", payload["displayName"], "Ryft Demo")
			}
			_, _ = io.WriteString(w, `{"sessionObject":"apple_pay_session"}`)
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.String())
		}
	})

	ctx := context.Background()

	domain, err := client.ApplePay.RegisterDomain(ctx, RegisterApplePayWebDomainRequest{
		DomainName: "checkout.example.com",
	}, WithAccount("ac_123"))
	if err != nil {
		t.Fatalf("RegisterDomainForAccount returned error: %v", err)
	}
	if domain.ID != "apwd_123" {
		t.Fatalf("domain.ID = %q, want %q", domain.ID, "apwd_123")
	}

	domains, err := client.ApplePay.ListDomains(ctx, ApplePayWebDomainListParams{ListParams: ListParams{Ascending: true, Limit: 10, StartsAfter: "apwd_prev"}}, WithAccount("ac_123"))
	if err != nil {
		t.Fatalf("ListDomainsForAccount returned error: %v", err)
	}
	if len(domains.Items) != 1 {
		t.Fatalf("len(domains.Items) = %d, want 1", len(domains.Items))
	}

	gotDomain, err := client.ApplePay.GetDomain(ctx, "apwd_123", WithAccount("ac_123"))
	if err != nil {
		t.Fatalf("GetDomainForAccount returned error: %v", err)
	}
	if gotDomain.DomainName != "checkout.example.com" {
		t.Fatalf("gotDomain.DomainName = %q, want %q", gotDomain.DomainName, "checkout.example.com")
	}

	deleted, err := client.ApplePay.DeleteDomain(ctx, "apwd_123", WithAccount("ac_123"))
	if err != nil {
		t.Fatalf("DeleteDomainForAccount returned error: %v", err)
	}
	if deleted.ID != "apwd_123" {
		t.Fatalf("deleted.ID = %q, want %q", deleted.ID, "apwd_123")
	}

	session, err := client.ApplePay.CreateSession(ctx, CreateApplePayWebSessionRequest{
		DisplayName: "Ryft Demo",
		DomainName:  "checkout.example.com",
	}, WithAccount("ac_123"))
	if err != nil {
		t.Fatalf("CreateSessionForAccount returned error: %v", err)
	}
	if session.SessionObject != "apple_pay_session" {
		t.Fatalf("session.SessionObject = %q, want %q", session.SessionObject, "apple_pay_session")
	}
}

func TestInPersonMethods(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/in-person/products":
			assertQueryValues(t, r.URL.Query(), map[string]string{
				"ascending":   "true",
				"limit":       "5",
				"startsAfter": "ippd_prev",
			})
			_, _ = io.WriteString(w, `{"items":[{"id":"ippd_123","name":"Starter Kit"}]}`)
		case r.Method == http.MethodGet && r.URL.Path == "/in-person/products/ippd_123":
			_, _ = io.WriteString(w, `{"id":"ippd_123","name":"Starter Kit"}`)
		case r.Method == http.MethodGet && r.URL.Path == "/in-person/skus":
			assertQueryValues(t, r.URL.Query(), map[string]string{
				"country":   "GB",
				"limit":     "5",
				"productId": "ippd_123",
			})
			_, _ = io.WriteString(w, `{"items":[{"id":"ipsku_123","productId":"ippd_123"}]}`)
		case r.Method == http.MethodGet && r.URL.Path == "/in-person/skus/ipsku_123":
			_, _ = io.WriteString(w, `{"id":"ipsku_123","productId":"ippd_123"}`)
		case r.Method == http.MethodGet && r.URL.Path == "/in-person/orders":
			if got := r.Header.Get("Account"); got != "ac_123" {
				t.Fatalf("Account header = %q, want %q", got, "ac_123")
			}
			assertQueryValues(t, r.URL.Query(), map[string]string{
				"ascending":   "true",
				"limit":       "5",
				"startsAfter": "ipord_prev",
			})
			_, _ = io.WriteString(w, `{"items":[{"id":"ipord_123","status":"ReadyToShip"}]}`)
		case r.Method == http.MethodGet && r.URL.Path == "/in-person/orders/ipord_123":
			if got := r.Header.Get("Account"); got != "ac_123" {
				t.Fatalf("Account header = %q, want %q", got, "ac_123")
			}
			_, _ = io.WriteString(w, `{"id":"ipord_123","status":"ReadyToShip"}`)
		case r.Method == http.MethodPost && r.URL.Path == "/in-person/locations":
			if got := r.Header.Get("Account"); got != "ac_123" {
				t.Fatalf("Account header = %q, want %q", got, "ac_123")
			}
			payload := decodeJSONBody(t, r)
			if payload["name"] != "Store 1" {
				t.Fatalf("name = %v, want %q", payload["name"], "Store 1")
			}
			_, _ = io.WriteString(w, `{"id":"iploc_123","name":"Store 1"}`)
		case r.Method == http.MethodGet && r.URL.Path == "/in-person/locations":
			if got := r.Header.Get("Account"); got != "ac_123" {
				t.Fatalf("Account header = %q, want %q", got, "ac_123")
			}
			assertQueryValues(t, r.URL.Query(), map[string]string{
				"ascending":   "true",
				"limit":       "5",
				"startsAfter": "iploc_prev",
			})
			_, _ = io.WriteString(w, `{"items":[{"id":"iploc_123","name":"Store 1"}]}`)
		case r.Method == http.MethodGet && r.URL.Path == "/in-person/locations/iploc_123":
			if got := r.Header.Get("Account"); got != "ac_123" {
				t.Fatalf("Account header = %q, want %q", got, "ac_123")
			}
			_, _ = io.WriteString(w, `{"id":"iploc_123","name":"Store 1"}`)
		case r.Method == http.MethodPatch && r.URL.Path == "/in-person/locations/iploc_123":
			if got := r.Header.Get("Account"); got != "ac_123" {
				t.Fatalf("Account header = %q, want %q", got, "ac_123")
			}
			payload := decodeJSONBody(t, r)
			if payload["name"] != "Store 1A" {
				t.Fatalf("name = %v, want %q", payload["name"], "Store 1A")
			}
			_, _ = io.WriteString(w, `{"id":"iploc_123","name":"Store 1A"}`)
		case r.Method == http.MethodDelete && r.URL.Path == "/in-person/locations/iploc_123":
			if got := r.Header.Get("Account"); got != "ac_123" {
				t.Fatalf("Account header = %q, want %q", got, "ac_123")
			}
			_, _ = io.WriteString(w, `{"id":"iploc_123"}`)
		case r.Method == http.MethodPost && r.URL.Path == "/in-person/terminals":
			if got := r.Header.Get("Account"); got != "ac_123" {
				t.Fatalf("Account header = %q, want %q", got, "ac_123")
			}
			payload := decodeJSONBody(t, r)
			if payload["serialNumber"] != "SN-001" {
				t.Fatalf("serialNumber = %v, want %q", payload["serialNumber"], "SN-001")
			}
			_, _ = io.WriteString(w, `{"id":"tml_123","serialNumber":"SN-001"}`)
		case r.Method == http.MethodGet && r.URL.Path == "/in-person/terminals":
			if got := r.Header.Get("Account"); got != "ac_123" {
				t.Fatalf("Account header = %q, want %q", got, "ac_123")
			}
			assertQueryValues(t, r.URL.Query(), map[string]string{
				"ascending":   "true",
				"limit":       "5",
				"startsAfter": "tml_prev",
			})
			_, _ = io.WriteString(w, `{"items":[{"id":"tml_123","serialNumber":"SN-001"}]}`)
		case r.Method == http.MethodGet && r.URL.Path == "/in-person/terminals/tml_123":
			if got := r.Header.Get("Account"); got != "ac_123" {
				t.Fatalf("Account header = %q, want %q", got, "ac_123")
			}
			_, _ = io.WriteString(w, `{"id":"tml_123","serialNumber":"SN-001"}`)
		case r.Method == http.MethodPatch && r.URL.Path == "/in-person/terminals/tml_123":
			if got := r.Header.Get("Account"); got != "ac_123" {
				t.Fatalf("Account header = %q, want %q", got, "ac_123")
			}
			payload := decodeJSONBody(t, r)
			if payload["name"] != "Front Desk" {
				t.Fatalf("name = %v, want %q", payload["name"], "Front Desk")
			}
			_, _ = io.WriteString(w, `{"id":"tml_123","name":"Front Desk"}`)
		case r.Method == http.MethodDelete && r.URL.Path == "/in-person/terminals/tml_123":
			if got := r.Header.Get("Account"); got != "ac_123" {
				t.Fatalf("Account header = %q, want %q", got, "ac_123")
			}
			_, _ = io.WriteString(w, `{"id":"tml_123"}`)
		case r.Method == http.MethodPost && r.URL.Path == "/in-person/terminals/tml_123/payment":
			if got := r.Header.Get("Account"); got != "ac_123" {
				t.Fatalf("Account header = %q, want %q", got, "ac_123")
			}
			payload := decodeJSONBody(t, r)
			amounts := payload["amounts"].(map[string]any)
			if amounts["requested"] != float64(1200) {
				t.Fatalf("amounts.requested = %v, want %v", amounts["requested"], float64(1200))
			}
			_, _ = io.WriteString(w, `{"id":"tml_123","status":"AwaitingCard"}`)
		case r.Method == http.MethodPost && r.URL.Path == "/in-person/terminals/tml_123/refund":
			if got := r.Header.Get("Account"); got != "ac_123" {
				t.Fatalf("Account header = %q, want %q", got, "ac_123")
			}
			payload := decodeJSONBody(t, r)
			paymentSession := payload["paymentSession"].(map[string]any)
			if paymentSession["id"] != "ps_123" {
				t.Fatalf("paymentSession.id = %v, want %q", paymentSession["id"], "ps_123")
			}
			_, _ = io.WriteString(w, `{"id":"tml_123","status":"AwaitingCard"}`)
		case r.Method == http.MethodPost && r.URL.Path == "/in-person/terminals/tml_123/cancel-action":
			if got := r.Header.Get("Account"); got != "ac_123" {
				t.Fatalf("Account header = %q, want %q", got, "ac_123")
			}
			_, _ = io.WriteString(w, `{"id":"tml_123","status":"Ready"}`)
		case r.Method == http.MethodPost && r.URL.Path == "/in-person/terminals/tml_123/confirm-receipt":
			if got := r.Header.Get("Account"); got != "ac_123" {
				t.Fatalf("Account header = %q, want %q", got, "ac_123")
			}
			payload := decodeJSONBody(t, r)
			customerCopy := payload["customerCopy"].(map[string]any)
			if customerCopy["status"] != "Succeeded" {
				t.Fatalf("customerCopy.status = %v, want %q", customerCopy["status"], "Succeeded")
			}
			_, _ = io.WriteString(w, `{"id":"tml_123","status":"ReceiptConfirmed"}`)
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.String())
		}
	})

	ctx := context.Background()

	products, err := client.InPersonProducts.List(ctx, InPersonProductListParams{ListParams: ListParams{Ascending: true, Limit: 5, StartsAfter: "ippd_prev"}})
	if err != nil {
		t.Fatalf("InPersonProducts.List returned error: %v", err)
	}
	if len(products.Items) != 1 {
		t.Fatalf("len(products.Items) = %d, want 1", len(products.Items))
	}

	product, err := client.InPersonProducts.Get(ctx, "ippd_123")
	if err != nil {
		t.Fatalf("InPersonProducts.Get returned error: %v", err)
	}
	if product.ID != "ippd_123" {
		t.Fatalf("product.ID = %q, want %q", product.ID, "ippd_123")
	}

	skus, err := client.InPersonSkus.List(ctx, InPersonSkuListParams{ListParams: ListParams{Limit: 5}, Country: "GB", ProductID: "ippd_123"})
	if err != nil {
		t.Fatalf("InPersonSkus.List returned error: %v", err)
	}
	if len(skus.Items) != 1 {
		t.Fatalf("len(skus.Items) = %d, want 1", len(skus.Items))
	}

	sku, err := client.InPersonSkus.Get(ctx, "ipsku_123")
	if err != nil {
		t.Fatalf("InPersonSkus.Get returned error: %v", err)
	}
	if sku.ID != "ipsku_123" {
		t.Fatalf("sku.ID = %q, want %q", sku.ID, "ipsku_123")
	}

	orders, err := client.InPersonOrders.List(ctx, InPersonOrderListParams{ListParams: ListParams{Ascending: true, Limit: 5, StartsAfter: "ipord_prev"}}, WithAccount("ac_123"))
	if err != nil {
		t.Fatalf("InPersonOrders.ListForAccount returned error: %v", err)
	}
	if len(orders.Items) != 1 {
		t.Fatalf("len(orders.Items) = %d, want 1", len(orders.Items))
	}

	order, err := client.InPersonOrders.Get(ctx, "ipord_123", WithAccount("ac_123"))
	if err != nil {
		t.Fatalf("InPersonOrders.GetForAccount returned error: %v", err)
	}
	if order.ID != "ipord_123" {
		t.Fatalf("order.ID = %q, want %q", order.ID, "ipord_123")
	}

	location, err := client.InPersonLocations.Create(ctx, CreateInPersonLocationRequest{
		Name: "Store 1",
		Address: InPersonLocationAddress{
			FirstLine:  "1 Main St",
			City:       "London",
			PostalCode: "SW1A 1AA",
			Country:    "GB",
		},
	}, WithAccount("ac_123"))
	if err != nil {
		t.Fatalf("InPersonLocations.CreateForAccount returned error: %v", err)
	}
	if location.ID != "iploc_123" {
		t.Fatalf("location.ID = %q, want %q", location.ID, "iploc_123")
	}

	locations, err := client.InPersonLocations.List(ctx, InPersonLocationListParams{ListParams: ListParams{Ascending: true, Limit: 5, StartsAfter: "iploc_prev"}}, WithAccount("ac_123"))
	if err != nil {
		t.Fatalf("InPersonLocations.ListForAccount returned error: %v", err)
	}
	if len(locations.Items) != 1 {
		t.Fatalf("len(locations.Items) = %d, want 1", len(locations.Items))
	}

	gotLocation, err := client.InPersonLocations.Get(ctx, "iploc_123", WithAccount("ac_123"))
	if err != nil {
		t.Fatalf("InPersonLocations.GetForAccount returned error: %v", err)
	}
	if gotLocation.ID != "iploc_123" {
		t.Fatalf("gotLocation.ID = %q, want %q", gotLocation.ID, "iploc_123")
	}

	updatedLocation, err := client.InPersonLocations.Update(ctx, "iploc_123", UpdateInPersonLocationRequest{
		Name: "Store 1A",
	}, WithAccount("ac_123"))
	if err != nil {
		t.Fatalf("InPersonLocations.UpdateForAccount returned error: %v", err)
	}
	if updatedLocation.Name != "Store 1A" {
		t.Fatalf("updatedLocation.Name = %q, want %q", updatedLocation.Name, "Store 1A")
	}

	deletedLocation, err := client.InPersonLocations.Delete(ctx, "iploc_123", WithAccount("ac_123"))
	if err != nil {
		t.Fatalf("InPersonLocations.DeleteForAccount returned error: %v", err)
	}
	if deletedLocation.ID != "iploc_123" {
		t.Fatalf("deletedLocation.ID = %q, want %q", deletedLocation.ID, "iploc_123")
	}

	terminal, err := client.InPersonTerminals.Create(ctx, CreateTerminalRequest{
		SerialNumber: "SN-001",
		LocationID:   "iploc_123",
	}, WithAccount("ac_123"))
	if err != nil {
		t.Fatalf("InPersonTerminals.CreateForAccount returned error: %v", err)
	}
	if terminal.ID != "tml_123" {
		t.Fatalf("terminal.ID = %q, want %q", terminal.ID, "tml_123")
	}

	terminals, err := client.InPersonTerminals.List(ctx, TerminalListParams{ListParams: ListParams{Ascending: true, Limit: 5, StartsAfter: "tml_prev"}}, WithAccount("ac_123"))
	if err != nil {
		t.Fatalf("InPersonTerminals.ListForAccount returned error: %v", err)
	}
	if len(terminals.Items) != 1 {
		t.Fatalf("len(terminals.Items) = %d, want 1", len(terminals.Items))
	}

	gotTerminal, err := client.InPersonTerminals.Get(ctx, "tml_123", WithAccount("ac_123"))
	if err != nil {
		t.Fatalf("InPersonTerminals.GetForAccount returned error: %v", err)
	}
	if gotTerminal.ID != "tml_123" {
		t.Fatalf("gotTerminal.ID = %q, want %q", gotTerminal.ID, "tml_123")
	}

	updatedTerminal, err := client.InPersonTerminals.Update(ctx, "tml_123", UpdateTerminalRequest{
		Name: "Front Desk",
	}, WithAccount("ac_123"))
	if err != nil {
		t.Fatalf("InPersonTerminals.UpdateForAccount returned error: %v", err)
	}
	if updatedTerminal.Name != "Front Desk" {
		t.Fatalf("updatedTerminal.Name = %q, want %q", updatedTerminal.Name, "Front Desk")
	}

	paymentTerminal, err := client.InPersonTerminals.InitiatePayment(ctx, "tml_123", TerminalPaymentRequest{
		Amounts:  RequestedAmounts{Requested: 1200},
		Currency: "GBP",
	}, WithAccount("ac_123"))
	if err != nil {
		t.Fatalf("InPersonTerminals.InitiatePaymentForAccount returned error: %v", err)
	}
	if paymentTerminal.Status != "AwaitingCard" {
		t.Fatalf("paymentTerminal.Status = %q, want %q", paymentTerminal.Status, "AwaitingCard")
	}

	refundTerminal, err := client.InPersonTerminals.InitiateRefund(ctx, "tml_123", TerminalRefundRequest{
		PaymentSession: TerminalRefundPaymentSessionReference{ID: "ps_123"},
	}, WithAccount("ac_123"))
	if err != nil {
		t.Fatalf("InPersonTerminals.InitiateRefundForAccount returned error: %v", err)
	}
	if refundTerminal.Status != "AwaitingCard" {
		t.Fatalf("refundTerminal.Status = %q, want %q", refundTerminal.Status, "AwaitingCard")
	}

	cancelledTerminal, err := client.InPersonTerminals.CancelAction(ctx, "tml_123", WithAccount("ac_123"))
	if err != nil {
		t.Fatalf("InPersonTerminals.CancelActionForAccount returned error: %v", err)
	}
	if cancelledTerminal.Status != "Ready" {
		t.Fatalf("cancelledTerminal.Status = %q, want %q", cancelledTerminal.Status, "Ready")
	}

	confirmedTerminal, err := client.InPersonTerminals.ConfirmReceipt(ctx, "tml_123", TerminalConfirmReceiptRequest{
		CustomerCopy: &ReceiptCopyStatus{Status: "Succeeded"},
	}, WithAccount("ac_123"))
	if err != nil {
		t.Fatalf("InPersonTerminals.ConfirmReceiptForAccount returned error: %v", err)
	}
	if confirmedTerminal.Status != "ReceiptConfirmed" {
		t.Fatalf("confirmedTerminal.Status = %q, want %q", confirmedTerminal.Status, "ReceiptConfirmed")
	}

	deletedTerminal, err := client.InPersonTerminals.Delete(ctx, "tml_123", WithAccount("ac_123"))
	if err != nil {
		t.Fatalf("InPersonTerminals.DeleteForAccount returned error: %v", err)
	}
	if deletedTerminal.ID != "tml_123" {
		t.Fatalf("deletedTerminal.ID = %q, want %q", deletedTerminal.ID, "tml_123")
	}
}

func TestPersonsAndPayoutMethodsMethods(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/accounts/ac_123/persons":
			payload := decodeJSONBody(t, r)
			if payload["firstName"] != "Jane" {
				t.Fatalf("firstName = %v, want %q", payload["firstName"], "Jane")
			}
			_, _ = io.WriteString(w, `{"id":"per_123","firstName":"Jane"}`)
		case r.Method == http.MethodGet && r.URL.Path == "/accounts/ac_123/persons/per_123":
			_, _ = io.WriteString(w, `{"id":"per_123","email":"jane@example.com"}`)
		case r.Method == http.MethodGet && r.URL.Path == "/accounts/ac_123/persons":
			assertQueryValues(t, r.URL.Query(), map[string]string{
				"ascending":   "false",
				"limit":       "10",
				"startsAfter": "per_001",
			})
			_, _ = io.WriteString(w, `{"items":[{"id":"per_123"}]}`)
		case r.Method == http.MethodPatch && r.URL.Path == "/accounts/ac_123/persons/per_123":
			payload := decodeJSONBody(t, r)
			if payload["middleNames"] != "Q" {
				t.Fatalf("middleNames = %v, want %q", payload["middleNames"], "Q")
			}
			_, _ = io.WriteString(w, `{"id":"per_123","email":"updated@example.com"}`)
		case r.Method == http.MethodDelete && r.URL.Path == "/accounts/ac_123/persons/per_123":
			_, _ = io.WriteString(w, `{"id":"per_123"}`)
		case r.Method == http.MethodPost && r.URL.Path == "/accounts/ac_123/payout-methods":
			payload := decodeJSONBody(t, r)
			if payload["displayName"] != "Main Account" {
				t.Fatalf("displayName = %v, want %q", payload["displayName"], "Main Account")
			}
			_, _ = io.WriteString(w, `{"id":"pmo_123","displayName":"Main Account"}`)
		case r.Method == http.MethodGet && r.URL.Path == "/accounts/ac_123/payout-methods/pmo_123":
			_, _ = io.WriteString(w, `{"id":"pmo_123","currency":"GBP"}`)
		case r.Method == http.MethodGet && r.URL.Path == "/accounts/ac_123/payout-methods":
			assertQueryValues(t, r.URL.Query(), map[string]string{
				"ascending":   "true",
				"limit":       "20",
				"startsAfter": "pmo_001",
			})
			_, _ = io.WriteString(w, `{"items":[{"id":"pmo_123"}]}`)
		case r.Method == http.MethodPatch && r.URL.Path == "/accounts/ac_123/payout-methods/pmo_123":
			payload := decodeJSONBody(t, r)
			if payload["displayName"] != "Reserve Account" {
				t.Fatalf("displayName = %v, want %q", payload["displayName"], "Reserve Account")
			}
			_, _ = io.WriteString(w, `{"id":"pmo_123","displayName":"Reserve Account"}`)
		case r.Method == http.MethodDelete && r.URL.Path == "/accounts/ac_123/payout-methods/pmo_123":
			_, _ = io.WriteString(w, `{"id":"pmo_123"}`)
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.String())
		}
	})

	ctx := context.Background()

	person, err := client.Persons.Create(ctx, "ac_123", CreatePersonRequest{
		FirstName:     "Jane",
		LastName:      "Doe",
		Email:         "jane@example.com",
		DateOfBirth:   "1990-01-01",
		Gender:        "Female",
		Nationalities: []string{"GB"},
		Address: Address{
			LineOne:    "1 High Street",
			City:       "London",
			Country:    "GB",
			PostalCode: "SW1A 1AA",
		},
		PhoneNumber:   "+447700900000",
		BusinessRoles: []string{"Director"},
	})
	if err != nil {
		t.Fatalf("Persons.Create returned error: %v", err)
	}
	if person.ID != "per_123" {
		t.Fatalf("person.ID = %q, want %q", person.ID, "per_123")
	}

	gotPerson, err := client.Persons.Get(ctx, "ac_123", "per_123")
	if err != nil {
		t.Fatalf("Persons.Get returned error: %v", err)
	}
	if gotPerson.Email != "jane@example.com" {
		t.Fatalf("gotPerson.Email = %q, want %q", gotPerson.Email, "jane@example.com")
	}

	persons, err := client.Persons.List(ctx, "ac_123", PersonListParams{ListParams: ListParams{Ascending: false, Limit: 10, StartsAfter: "per_001"}})
	if err != nil {
		t.Fatalf("Persons.List returned error: %v", err)
	}
	if len(persons.Items) != 1 {
		t.Fatalf("len(persons.Items) = %d, want 1", len(persons.Items))
	}

	updatedPerson, err := client.Persons.Update(ctx, "ac_123", "per_123", UpdatePersonRequest{
		MiddleNames: "Q",
	})
	if err != nil {
		t.Fatalf("Persons.Update returned error: %v", err)
	}
	if updatedPerson.Email != "updated@example.com" {
		t.Fatalf("updatedPerson.Email = %q, want %q", updatedPerson.Email, "updated@example.com")
	}

	deletedPerson, err := client.Persons.Delete(ctx, "ac_123", "per_123")
	if err != nil {
		t.Fatalf("Persons.Delete returned error: %v", err)
	}
	if deletedPerson.ID != "per_123" {
		t.Fatalf("deletedPerson.ID = %q, want %q", deletedPerson.ID, "per_123")
	}

	payoutMethod, err := client.PayoutMethods.Create(ctx, "ac_123", CreatePayoutMethodRequest{
		Type:        "BankAccount",
		DisplayName: "Main Account",
		Currency:    "GBP",
		Country:     "GB",
		BankAccount: BankAccountDetails{
			AccountNumberType: "SortCodeAndAccountNumber",
			AccountNumber:     "12345678",
			BankIDType:        "SortCode",
			BankID:            "112233",
		},
	})
	if err != nil {
		t.Fatalf("PayoutMethods.Create returned error: %v", err)
	}
	if payoutMethod.ID != "pmo_123" {
		t.Fatalf("payoutMethod.ID = %q, want %q", payoutMethod.ID, "pmo_123")
	}

	gotPayoutMethod, err := client.PayoutMethods.Get(ctx, "ac_123", "pmo_123")
	if err != nil {
		t.Fatalf("PayoutMethods.Get returned error: %v", err)
	}
	if gotPayoutMethod.Currency != "GBP" {
		t.Fatalf("gotPayoutMethod.Currency = %q, want %q", gotPayoutMethod.Currency, "GBP")
	}

	payoutMethods, err := client.PayoutMethods.List(ctx, "ac_123", PayoutMethodListParams{ListParams: ListParams{Ascending: true, Limit: 20, StartsAfter: "pmo_001"}})
	if err != nil {
		t.Fatalf("PayoutMethods.List returned error: %v", err)
	}
	if len(payoutMethods.Items) != 1 {
		t.Fatalf("len(payoutMethods.Items) = %d, want 1", len(payoutMethods.Items))
	}

	updatedPayoutMethod, err := client.PayoutMethods.Update(ctx, "ac_123", "pmo_123", UpdatePayoutMethodRequest{
		DisplayName: "Reserve Account",
		BankAccount: BankAccountDetails{},
	})
	if err != nil {
		t.Fatalf("PayoutMethods.Update returned error: %v", err)
	}
	if updatedPayoutMethod.DisplayName != "Reserve Account" {
		t.Fatalf("updatedPayoutMethod.DisplayName = %q, want %q", updatedPayoutMethod.DisplayName, "Reserve Account")
	}

	deletedPayoutMethod, err := client.PayoutMethods.Delete(ctx, "ac_123", "pmo_123")
	if err != nil {
		t.Fatalf("PayoutMethods.Delete returned error: %v", err)
	}
	if deletedPayoutMethod.ID != "pmo_123" {
		t.Fatalf("deletedPayoutMethod.ID = %q, want %q", deletedPayoutMethod.ID, "pmo_123")
	}
}

func TestPaymentSessionsAndPaymentMethodsMethods(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/payment-sessions":
			if got := r.Header.Get("Account"); got != "ac_123" {
				t.Fatalf("Account header = %q, want %q", got, "ac_123")
			}
			payload := decodeJSONBody(t, r)
			if payload["amount"] != float64(2000) {
				t.Fatalf("amount = %v, want %v", payload["amount"], 2000)
			}
			if payload["platformFee"] != float64(50) {
				t.Fatalf("platformFee = %v, want %v", payload["platformFee"], 50)
			}
			customerDetails := payload["customerDetails"].(map[string]any)
			if customerDetails["id"] != "cus_123" {
				t.Fatalf("customerDetails.id = %v, want %q", customerDetails["id"], "cus_123")
			}
			previousPayment := payload["previousPayment"].(map[string]any)
			if previousPayment["id"] != "ps_prev" {
				t.Fatalf("previousPayment.id = %v, want %q", previousPayment["id"], "ps_prev")
			}
			rebillingDetail := payload["rebillingDetail"].(map[string]any)
			if rebillingDetail["currentPaymentNumber"] != float64(2) {
				t.Fatalf("rebillingDetail.currentPaymentNumber = %v, want %v", rebillingDetail["currentPaymentNumber"], 2)
			}
			splits := payload["splits"].(map[string]any)
			items := splits["items"].([]any)
			if len(items) != 1 {
				t.Fatalf("len(splits.items) = %d, want 1", len(items))
			}
			item := items[0].(map[string]any)
			if item["accountId"] != "ac_split" {
				t.Fatalf("splits.items[0].accountId = %v, want %q", item["accountId"], "ac_split")
			}
			_, _ = io.WriteString(w, `{"id":"ps_123","amount":2000,"currency":"GBP"}`)
		case r.Method == http.MethodGet && r.URL.Path == "/payment-sessions/ps_123":
			if got := r.Header.Get("Account"); got != "ac_123" {
				t.Fatalf("Account header = %q, want %q", got, "ac_123")
			}
			_, _ = io.WriteString(w, `{"id":"ps_123","currency":"GBP"}`)
		case r.Method == http.MethodPatch && r.URL.Path == "/payment-sessions/ps_123":
			if got := r.Header.Get("Account"); got != "ac_123" {
				t.Fatalf("Account header = %q, want %q", got, "ac_123")
			}
			payload := decodeJSONBody(t, r)
			if payload["customerEmail"] != "buyer@example.com" {
				t.Fatalf("customerEmail = %v, want %q", payload["customerEmail"], "buyer@example.com")
			}
			_, _ = io.WriteString(w, `{"id":"ps_123","customerEmail":"buyer@example.com"}`)
		case r.Method == http.MethodPost && r.URL.Path == "/payment-sessions/ps_123/refunds":
			if got := r.Header.Get("Account"); got != "ac_123" {
				t.Fatalf("Account header = %q, want %q", got, "ac_123")
			}
			payload := decodeJSONBody(t, r)
			if payload["refundPlatformFee"] != true {
				t.Fatalf("refundPlatformFee = %v, want true", payload["refundPlatformFee"])
			}
			_, _ = io.WriteString(w, `{"id":"txn_123","paymentSessionId":"ps_123","type":"Refund","refundedAmount":2000}`)
		case r.Method == http.MethodGet && r.URL.Path == "/payment-sessions/ps_123/transactions":
			if got := r.Header.Get("Account"); got != "ac_123" {
				t.Fatalf("Account header = %q, want %q", got, "ac_123")
			}
			assertQueryValues(t, r.URL.Query(), map[string]string{
				"startTimestamp": "100",
				"endTimestamp":   "200",
				"ascending":      "false",
				"limit":          "3",
			})
			_, _ = io.WriteString(w, `{"items":[{"id":"txn_123","paymentSessionId":"ps_123"}]}`)
		case r.Method == http.MethodGet && r.URL.Path == "/payment-sessions/ps_123/transactions/txn_123":
			if got := r.Header.Get("Account"); got != "ac_123" {
				t.Fatalf("Account header = %q, want %q", got, "ac_123")
			}
			_, _ = io.WriteString(w, `{"id":"txn_123","paymentSessionId":"ps_123","status":"Approved"}`)
		case r.Method == http.MethodGet && r.URL.Path == "/payment-methods/pm_123":
			_, _ = io.WriteString(w, `{"id":"pm_123","type":"Card"}`)
		case r.Method == http.MethodPatch && r.URL.Path == "/payment-methods/pm_123":
			payload := decodeJSONBody(t, r)
			billingAddress := payload["billingAddress"].(map[string]any)
			if billingAddress["city"] != "London" {
				t.Fatalf("billingAddress.city = %v, want %q", billingAddress["city"], "London")
			}
			_, _ = io.WriteString(w, `{"id":"pm_123","billingAddress":{"city":"London"}}`)
		case r.Method == http.MethodDelete && r.URL.Path == "/payment-methods/pm_123":
			_, _ = io.WriteString(w, `{"id":"pm_123"}`)
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.String())
		}
	})

	ctx := context.Background()

	paymentSession, err := client.PaymentSessions.Create(ctx, CreatePaymentSessionRequest{
		Amount:          2000,
		Currency:        "GBP",
		PlatformFee:     50,
		CustomerDetails: &PaymentSessionCustomer{ID: "cus_123"},
		PreviousPayment: &PaymentSessionReference{ID: "ps_prev"},
		RebillingDetail: &RebillingDetail{
			AmountVariance:              "Fixed",
			NumberOfDaysBetweenPayments: 30,
			TotalNumberOfPayments:       12,
			CurrentPaymentNumber:        2,
		},
		Splits: &CreateSplitPaymentRequest{
			Items: []SplitItem{
				{
					AccountID:   "ac_split",
					Amount:      2000,
					Description: "sub account split",
					Fee:         &SplitFee{Amount: 25},
					Metadata:    map[string]string{"scenario": "platform"},
				},
			},
		},
	}, WithAccount("ac_123"))
	if err != nil {
		t.Fatalf("PaymentSessions.CreateForAccount returned error: %v", err)
	}
	if paymentSession.ID != "ps_123" {
		t.Fatalf("paymentSession.ID = %q, want %q", paymentSession.ID, "ps_123")
	}

	gotPaymentSession, err := client.PaymentSessions.Get(ctx, "ps_123", WithAccount("ac_123"))
	if err != nil {
		t.Fatalf("PaymentSessions.GetForAccount returned error: %v", err)
	}
	if gotPaymentSession.Currency != "GBP" {
		t.Fatalf("gotPaymentSession.Currency = %q, want %q", gotPaymentSession.Currency, "GBP")
	}

	updatedPaymentSession, err := client.PaymentSessions.Update(ctx, "ps_123", UpdatePaymentSessionRequest{
		CustomerEmail: "buyer@example.com",
	}, WithAccount("ac_123"))
	if err != nil {
		t.Fatalf("PaymentSessions.UpdateForAccount returned error: %v", err)
	}
	if updatedPaymentSession.CustomerEmail != "buyer@example.com" {
		t.Fatalf("updatedPaymentSession.CustomerEmail = %q, want %q", updatedPaymentSession.CustomerEmail, "buyer@example.com")
	}

	refund, err := client.PaymentSessions.Refund(ctx, "ps_123", RefundPaymentSessionRequest{
		Amount:            2000,
		Reason:            "requested_by_customer",
		RefundPlatformFee: true,
	}, WithAccount("ac_123"))
	if err != nil {
		t.Fatalf("PaymentSessions.RefundForAccount returned error: %v", err)
	}
	if refund.ID != "txn_123" {
		t.Fatalf("refund.ID = %q, want %q", refund.ID, "txn_123")
	}

	transactions, err := client.PaymentSessions.ListTransactions(ctx, "ps_123", PaymentSessionTransactionListParams{ListParams: ListParams{Ascending: false, Limit: 3}, StartTimestamp: 100, EndTimestamp: 200}, WithAccount("ac_123"))
	if err != nil {
		t.Fatalf("PaymentSessions.ListTransactionsForAccount returned error: %v", err)
	}
	if len(transactions.Items) != 1 {
		t.Fatalf("len(transactions.Items) = %d, want 1", len(transactions.Items))
	}

	transaction, err := client.PaymentSessions.GetTransaction(ctx, "ps_123", "txn_123", WithAccount("ac_123"))
	if err != nil {
		t.Fatalf("PaymentSessions.GetTransactionForAccount returned error: %v", err)
	}
	if transaction.Status != "Approved" {
		t.Fatalf("transaction.Status = %q, want %q", transaction.Status, "Approved")
	}

	paymentMethod, err := client.PaymentMethods.Get(ctx, "pm_123")
	if err != nil {
		t.Fatalf("PaymentMethods.Get returned error: %v", err)
	}
	if paymentMethod.Type != "Card" {
		t.Fatalf("paymentMethod.Type = %q, want %q", paymentMethod.Type, "Card")
	}

	updatedPaymentMethod, err := client.PaymentMethods.Update(ctx, "pm_123", UpdatePaymentMethodRequest{
		BillingAddress: map[string]string{"city": "London"},
	})
	if err != nil {
		t.Fatalf("PaymentMethods.Update returned error: %v", err)
	}
	if updatedPaymentMethod.BillingAddress["city"] != "London" {
		t.Fatalf("updatedPaymentMethod.BillingAddress[city] = %v, want %q", updatedPaymentMethod.BillingAddress["city"], "London")
	}

	deletedPaymentMethod, err := client.PaymentMethods.Delete(ctx, "pm_123")
	if err != nil {
		t.Fatalf("PaymentMethods.Delete returned error: %v", err)
	}
	if deletedPaymentMethod.ID != "pm_123" {
		t.Fatalf("deletedPaymentMethod.ID = %q, want %q", deletedPaymentMethod.ID, "pm_123")
	}
}

func TestSubscriptionsAndWebhooksMethods(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/subscriptions":
			payload := decodeJSONBody(t, r)
			if payload["description"] != "Gold plan" {
				t.Fatalf("description = %v, want %q", payload["description"], "Gold plan")
			}
			_, _ = io.WriteString(w, `{"id":"sub_123","description":"Gold plan"}`)
		case r.Method == http.MethodGet && r.URL.Path == "/subscriptions/sub_123":
			_, _ = io.WriteString(w, `{"id":"sub_123","status":"Active"}`)
		case r.Method == http.MethodPatch && r.URL.Path == "/subscriptions/sub_123":
			payload := decodeJSONBody(t, r)
			if payload["description"] != "Updated plan" {
				t.Fatalf("description = %v, want %q", payload["description"], "Updated plan")
			}
			_, _ = io.WriteString(w, `{"id":"sub_123","description":"Updated plan"}`)
		case r.Method == http.MethodGet && r.URL.Path == "/subscriptions":
			assertQueryValues(t, r.URL.Query(), map[string]string{
				"startTimestamp": "100",
				"endTimestamp":   "200",
				"ascending":      "true",
				"limit":          "10",
				"startsAfter":    "sub_001",
			})
			_, _ = io.WriteString(w, `{"items":[{"id":"sub_123"}]}`)
		case r.Method == http.MethodGet && r.URL.Path == "/subscriptions/sub_123/payment-sessions":
			assertQueryValues(t, r.URL.Query(), map[string]string{
				"startTimestamp": "100",
				"endTimestamp":   "200",
				"ascending":      "false",
				"limit":          "4",
				"startsAfter":    "ps_001",
			})
			_, _ = io.WriteString(w, `{"items":[{"id":"ps_123"}]}`)
		case r.Method == http.MethodPatch && r.URL.Path == "/subscriptions/sub_123/pause":
			payload := decodeJSONBody(t, r)
			if payload["reason"] != "customer_request" {
				t.Fatalf("reason = %v, want %q", payload["reason"], "customer_request")
			}
			_, _ = io.WriteString(w, `{"id":"sub_123","status":"Paused"}`)
		case r.Method == http.MethodPatch && r.URL.Path == "/subscriptions/sub_123/resume":
			_, _ = io.WriteString(w, `{"id":"sub_123","status":"Active"}`)
		case r.Method == http.MethodDelete && r.URL.Path == "/subscriptions/sub_123/cancel":
			_, _ = io.WriteString(w, `{"id":"sub_123"}`)
		case r.Method == http.MethodPost && r.URL.Path == "/webhooks":
			payload := decodeJSONBody(t, r)
			if payload["url"] != "https://example.com/webhooks" {
				t.Fatalf("url = %v, want %q", payload["url"], "https://example.com/webhooks")
			}
			_, _ = io.WriteString(w, `{"id":"wh_123","url":"https://example.com/webhooks","active":true}`)
		case r.Method == http.MethodGet && r.URL.Path == "/webhooks/wh_123":
			_, _ = io.WriteString(w, `{"id":"wh_123","active":true}`)
		case r.Method == http.MethodGet && r.URL.Path == "/webhooks":
			_, _ = io.WriteString(w, `{"items":[{"id":"wh_123"}]}`)
		case r.Method == http.MethodPatch && r.URL.Path == "/webhooks/wh_123":
			payload := decodeJSONBody(t, r)
			if payload["active"] != false {
				t.Fatalf("active = %v, want false", payload["active"])
			}
			_, _ = io.WriteString(w, `{"id":"wh_123","active":false}`)
		case r.Method == http.MethodDelete && r.URL.Path == "/webhooks/wh_123":
			_, _ = io.WriteString(w, `{"id":"wh_123"}`)
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.String())
		}
	})

	ctx := context.Background()

	subscription, err := client.Subscriptions.Create(ctx, CreateSubscriptionRequest{
		Customer:    SubscriptionCustomerReference{ID: "cus_123"},
		Description: "Gold plan",
		Price: SubscriptionPrice{
			Amount:   1000,
			Currency: "GBP",
			Interval: SubscriptionInterval{Unit: "Month", Count: 1},
		},
	})
	if err != nil {
		t.Fatalf("Subscriptions.Create returned error: %v", err)
	}
	if subscription.ID != "sub_123" {
		t.Fatalf("subscription.ID = %q, want %q", subscription.ID, "sub_123")
	}

	gotSubscription, err := client.Subscriptions.Get(ctx, "sub_123")
	if err != nil {
		t.Fatalf("Subscriptions.Get returned error: %v", err)
	}
	if gotSubscription.Status != "Active" {
		t.Fatalf("gotSubscription.Status = %q, want %q", gotSubscription.Status, "Active")
	}

	updatedSubscription, err := client.Subscriptions.Update(ctx, "sub_123", UpdateSubscriptionRequest{
		Description: "Updated plan",
	})
	if err != nil {
		t.Fatalf("Subscriptions.Update returned error: %v", err)
	}
	if updatedSubscription.Description != "Updated plan" {
		t.Fatalf("updatedSubscription.Description = %q, want %q", updatedSubscription.Description, "Updated plan")
	}

	subscriptions, err := client.Subscriptions.List(ctx, SubscriptionListParams{ListParams: ListParams{Ascending: true, Limit: 10, StartsAfter: "sub_001"}, StartTimestamp: 100, EndTimestamp: 200})
	if err != nil {
		t.Fatalf("Subscriptions.List returned error: %v", err)
	}
	if len(subscriptions.Items) != 1 {
		t.Fatalf("len(subscriptions.Items) = %d, want 1", len(subscriptions.Items))
	}

	paymentSessions, err := client.Subscriptions.GetPaymentSessions(ctx, "sub_123", SubscriptionListParams{
		ListParams:     ListParams{Ascending: false, Limit: 4, StartsAfter: "ps_001"},
		StartTimestamp: 100,
		EndTimestamp:   200,
	})
	if err != nil {
		t.Fatalf("Subscriptions.GetPaymentSessions returned error: %v", err)
	}
	if len(paymentSessions.Items) != 1 {
		t.Fatalf("len(paymentSessions.Items) = %d, want 1", len(paymentSessions.Items))
	}

	pausedSubscription, err := client.Subscriptions.Pause(ctx, "sub_123", PauseSubscriptionRequest{
		Reason: "customer_request",
	})
	if err != nil {
		t.Fatalf("Subscriptions.Pause returned error: %v", err)
	}
	if pausedSubscription.Status != "Paused" {
		t.Fatalf("pausedSubscription.Status = %q, want %q", pausedSubscription.Status, "Paused")
	}

	resumedSubscription, err := client.Subscriptions.Resume(ctx, "sub_123")
	if err != nil {
		t.Fatalf("Subscriptions.Resume returned error: %v", err)
	}
	if resumedSubscription.Status != "Active" {
		t.Fatalf("resumedSubscription.Status = %q, want %q", resumedSubscription.Status, "Active")
	}

	cancelledSubscription, err := client.Subscriptions.Cancel(ctx, "sub_123")
	if err != nil {
		t.Fatalf("Subscriptions.Cancel returned error: %v", err)
	}
	if cancelledSubscription.ID != "sub_123" {
		t.Fatalf("cancelledSubscription.ID = %q, want %q", cancelledSubscription.ID, "sub_123")
	}

	webhook, err := client.Webhooks.Create(ctx, CreateWebhookRequest{
		URL:        "https://example.com/webhooks",
		Active:     true,
		EventTypes: []string{"PaymentSessionCaptured"},
	})
	if err != nil {
		t.Fatalf("Webhooks.Create returned error: %v", err)
	}
	if webhook.ID != "wh_123" {
		t.Fatalf("webhook.ID = %q, want %q", webhook.ID, "wh_123")
	}

	gotWebhook, err := client.Webhooks.Get(ctx, "wh_123")
	if err != nil {
		t.Fatalf("Webhooks.Get returned error: %v", err)
	}
	if !gotWebhook.Active {
		t.Fatal("gotWebhook.Active = false, want true")
	}

	webhooks, err := client.Webhooks.List(ctx)
	if err != nil {
		t.Fatalf("Webhooks.List returned error: %v", err)
	}
	if len(webhooks.Items) != 1 {
		t.Fatalf("len(webhooks.Items) = %d, want 1", len(webhooks.Items))
	}

	active := false
	updatedWebhook, err := client.Webhooks.Update(ctx, "wh_123", UpdateWebhookRequest{Active: &active})
	if err != nil {
		t.Fatalf("Webhooks.Update returned error: %v", err)
	}
	if updatedWebhook.Active {
		t.Fatal("updatedWebhook.Active = true, want false")
	}

	deletedWebhook, err := client.Webhooks.Delete(ctx, "wh_123")
	if err != nil {
		t.Fatalf("Webhooks.Delete returned error: %v", err)
	}
	if deletedWebhook.ID != "wh_123" {
		t.Fatalf("deletedWebhook.ID = %q, want %q", deletedWebhook.ID, "wh_123")
	}
}

func TestFinancialResourcesMethods(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/accounts/ac_123/payouts":
			payload := decodeJSONBody(t, r)
			if payload["payoutMethodId"] != "pmo_123" {
				t.Fatalf("payoutMethodId = %v, want %q", payload["payoutMethodId"], "pmo_123")
			}
			_, _ = io.WriteString(w, `{"id":"po_123","amount":5000,"currency":"GBP"}`)
		case r.Method == http.MethodGet && r.URL.Path == "/accounts/ac_123/payouts/po_123":
			_, _ = io.WriteString(w, `{"id":"po_123","status":"Pending"}`)
		case r.Method == http.MethodGet && r.URL.Path == "/accounts/ac_123/payouts":
			assertQueryValues(t, r.URL.Query(), map[string]string{
				"startTimestamp": "100",
				"endTimestamp":   "200",
				"ascending":      "false",
				"limit":          "8",
				"startsAfter":    "po_001",
			})
			_, _ = io.WriteString(w, `{"items":[{"id":"po_123"}]}`)
		case r.Method == http.MethodPost && r.URL.Path == "/transfers":
			payload := decodeJSONBody(t, r)
			destination := payload["destination"].(map[string]any)
			if destination["accountId"] != "ac_dest_123" {
				t.Fatalf("destination.accountId = %v, want %q", destination["accountId"], "ac_dest_123")
			}
			_, _ = io.WriteString(w, `{"id":"tr_123","status":"Approved"}`)
		case r.Method == http.MethodGet && r.URL.Path == "/transfers/tr_123":
			_, _ = io.WriteString(w, `{"id":"tr_123","currency":"GBP"}`)
		case r.Method == http.MethodGet && r.URL.Path == "/transfers":
			assertQueryValues(t, r.URL.Query(), map[string]string{
				"ascending": "false",
				"limit":     "6",
			})
			_, _ = io.WriteString(w, `{"items":[{"id":"tr_123"}]}`)
		case r.Method == http.MethodGet && r.URL.Path == "/balances":
			if got := r.Header.Get("Account"); got != "ac_123" {
				t.Fatalf("Account header = %q, want %q", got, "ac_123")
			}
			assertQueryValues(t, r.URL.Query(), map[string]string{
				"currency": "GBP",
			})
			_, _ = io.WriteString(w, `{"items":[{"currency":"GBP","available":{"amount":1000}}]}`)
		case r.Method == http.MethodGet && r.URL.Path == "/balance-transactions":
			if got := r.Header.Get("Account"); got != "ac_123" {
				t.Fatalf("Account header = %q, want %q", got, "ac_123")
			}
			assertQueryValues(t, r.URL.Query(), map[string]string{
				"limit":       "7",
				"startsAfter": "bt_001",
				"payoutId":    "po_123",
			})
			_, _ = io.WriteString(w, `{"items":[{"id":"bt_123","type":"Payout"}]}`)
		case r.Method == http.MethodGet && r.URL.Path == "/platform-fees":
			assertQueryValues(t, r.URL.Query(), map[string]string{
				"ascending": "true",
				"limit":     "9",
			})
			_, _ = io.WriteString(w, `{"items":[{"id":"fee_123"}]}`)
		case r.Method == http.MethodGet && r.URL.Path == "/platform-fees/fee_123":
			_, _ = io.WriteString(w, `{"id":"fee_123","currency":"GBP"}`)
		case r.Method == http.MethodGet && r.URL.Path == "/platform-fees/fee_123/refunds":
			_, _ = io.WriteString(w, `{"items":[{"id":"fr_123","platformFeeId":"fee_123"}]}`)
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.String())
		}
	})

	ctx := context.Background()

	payout, err := client.Payouts.Create(ctx, "ac_123", CreatePayoutRequest{
		Amount:         5000,
		Currency:       "GBP",
		PayoutMethodID: "pmo_123",
	})
	if err != nil {
		t.Fatalf("Payouts.Create returned error: %v", err)
	}
	if payout.ID != "po_123" {
		t.Fatalf("payout.ID = %q, want %q", payout.ID, "po_123")
	}

	gotPayout, err := client.Payouts.Get(ctx, "ac_123", "po_123")
	if err != nil {
		t.Fatalf("Payouts.Get returned error: %v", err)
	}
	if gotPayout.Status != "Pending" {
		t.Fatalf("gotPayout.Status = %q, want %q", gotPayout.Status, "Pending")
	}

	payouts, err := client.Payouts.List(ctx, "ac_123", PayoutListParams{ListParams: ListParams{Ascending: false, Limit: 8, StartsAfter: "po_001"}, StartTimestamp: 100, EndTimestamp: 200})
	if err != nil {
		t.Fatalf("Payouts.List returned error: %v", err)
	}
	if len(payouts.Items) != 1 {
		t.Fatalf("len(payouts.Items) = %d, want 1", len(payouts.Items))
	}

	transfer, err := client.Transfers.Create(ctx, CreateTransferRequest{
		Amount:   900,
		Currency: "GBP",
		Destination: TransferDestination{
			AccountID: "ac_dest_123",
		},
		Reason: "marketplace_split",
	})
	if err != nil {
		t.Fatalf("Transfers.Create returned error: %v", err)
	}
	if transfer.ID != "tr_123" {
		t.Fatalf("transfer.ID = %q, want %q", transfer.ID, "tr_123")
	}

	gotTransfer, err := client.Transfers.Get(ctx, "tr_123")
	if err != nil {
		t.Fatalf("Transfers.Get returned error: %v", err)
	}
	if gotTransfer.Currency != "GBP" {
		t.Fatalf("gotTransfer.Currency = %q, want %q", gotTransfer.Currency, "GBP")
	}

	transfers, err := client.Transfers.List(ctx, TransferListParams{ListParams: ListParams{Limit: 6}})
	if err != nil {
		t.Fatalf("Transfers.List returned error: %v", err)
	}
	if len(transfers.Items) != 1 {
		t.Fatalf("len(transfers.Items) = %d, want 1", len(transfers.Items))
	}

	balances, err := client.Balances.List(ctx, BalanceListParams{Currency: "GBP"}, WithAccount("ac_123"))
	if err != nil {
		t.Fatalf("Balances.List returned error: %v", err)
	}
	if len(balances.Items) != 1 {
		t.Fatalf("len(balances.Items) = %d, want 1", len(balances.Items))
	}

	balanceTransactions, err := client.BalanceTransactions.List(ctx, BalanceTransactionListParams{ListParams: ListParams{Limit: 7, StartsAfter: "bt_001"}, PayoutID: "po_123"}, WithAccount("ac_123"))
	if err != nil {
		t.Fatalf("BalanceTransactions.List returned error: %v", err)
	}
	if len(balanceTransactions.Items) != 1 {
		t.Fatalf("len(balanceTransactions.Items) = %d, want 1", len(balanceTransactions.Items))
	}

	platformFees, err := client.PlatformFees.List(ctx, PlatformFeeListParams{ListParams: ListParams{Ascending: true, Limit: 9}})
	if err != nil {
		t.Fatalf("PlatformFees.List returned error: %v", err)
	}
	if len(platformFees.Items) != 1 {
		t.Fatalf("len(platformFees.Items) = %d, want 1", len(platformFees.Items))
	}

	platformFee, err := client.PlatformFees.Get(ctx, "fee_123")
	if err != nil {
		t.Fatalf("PlatformFees.Get returned error: %v", err)
	}
	if platformFee.Currency != "GBP" {
		t.Fatalf("platformFee.Currency = %q, want %q", platformFee.Currency, "GBP")
	}

	platformFeeRefunds, err := client.PlatformFees.GetRefunds(ctx, "fee_123")
	if err != nil {
		t.Fatalf("PlatformFees.GetRefunds returned error: %v", err)
	}
	if len(platformFeeRefunds.Items) != 1 {
		t.Fatalf("len(platformFeeRefunds.Items) = %d, want 1", len(platformFeeRefunds.Items))
	}
}

func TestEventsDisputesAndFilesMethods(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/events":
			if got := r.Header.Get("Account"); got != "ac_123" {
				t.Fatalf("Account header = %q, want %q", got, "ac_123")
			}
			assertQueryValues(t, r.URL.Query(), map[string]string{
				"ascending": "true",
				"limit":     "11",
			})
			_, _ = io.WriteString(w, `{"items":[{"id":"evt_123","eventType":"PaymentSessionCaptured"}]}`)
		case r.Method == http.MethodGet && r.URL.Path == "/events/evt_123":
			if got := r.Header.Get("Account"); got != "ac_123" {
				t.Fatalf("Account header = %q, want %q", got, "ac_123")
			}
			_, _ = io.WriteString(w, `{"id":"evt_123","accountId":"ac_123"}`)
		case r.Method == http.MethodGet && r.URL.Path == "/files":
			assertQueryValues(t, r.URL.Query(), map[string]string{
				"category":    "Evidence",
				"ascending":   "false",
				"limit":       "4",
				"startsAfter": "file_001",
			})
			_, _ = io.WriteString(w, `{"items":[{"id":"file_123"}]}`)
		case r.Method == http.MethodGet && r.URL.Path == "/files/file_123":
			_, _ = io.WriteString(w, `{"id":"file_123","category":"Evidence"}`)
		case r.Method == http.MethodGet && r.URL.Path == "/disputes":
			assertQueryValues(t, r.URL.Query(), map[string]string{
				"startTimestamp": "100",
				"endTimestamp":   "200",
				"ascending":      "false",
				"limit":          "12",
				"startsAfter":    "dp_001",
			})
			_, _ = io.WriteString(w, `{"items":[{"id":"dp_123"}]}`)
		case r.Method == http.MethodGet && r.URL.Path == "/disputes/dp_123":
			_, _ = io.WriteString(w, `{"id":"dp_123","status":"RequiresAction"}`)
		case r.Method == http.MethodPost && r.URL.Path == "/disputes/dp_123/accept":
			_, _ = io.WriteString(w, `{"id":"dp_123","status":"Accepted"}`)
		case r.Method == http.MethodPost && r.URL.Path == "/disputes/dp_123/challenge":
			_, _ = io.WriteString(w, `{"id":"dp_123","status":"Challenged"}`)
		case r.Method == http.MethodPatch && r.URL.Path == "/disputes/dp_123/evidence":
			payload := decodeJSONBody(t, r)
			files := payload["files"].(map[string]any)
			if files["receipt"] != "file_123" {
				t.Fatalf("files.receipt = %v, want %q", files["receipt"], "file_123")
			}
			_, _ = io.WriteString(w, `{"id":"dp_123","status":"Submitted"}`)
		case r.Method == http.MethodDelete && r.URL.Path == "/disputes/dp_123/evidence":
			payload := decodeJSONBody(t, r)
			fileList := payload["files"].([]any)
			if len(fileList) != 1 || fileList[0] != "receipt" {
				t.Fatalf("files = %v, want [receipt]", fileList)
			}
			_, _ = io.WriteString(w, `{"id":"dp_123","status":"RequiresAction"}`)
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.String())
		}
	})

	ctx := context.Background()

	events, err := client.Events.List(ctx, EventListParams{ListParams: ListParams{Ascending: true, Limit: 11}}, WithAccount("ac_123"))
	if err != nil {
		t.Fatalf("Events.List returned error: %v", err)
	}
	if len(events.Items) != 1 {
		t.Fatalf("len(events.Items) = %d, want 1", len(events.Items))
	}

	event, err := client.Events.Get(ctx, "evt_123", WithAccount("ac_123"))
	if err != nil {
		t.Fatalf("Events.Get returned error: %v", err)
	}
	if event.AccountID != "ac_123" {
		t.Fatalf("event.AccountID = %q, want %q", event.AccountID, WithAccount("ac_123"))
	}

	files, err := client.Files.List(ctx, FileListParams{ListParams: ListParams{Ascending: false, Limit: 4, StartsAfter: "file_001"}, Category: "Evidence"})
	if err != nil {
		t.Fatalf("Files.List returned error: %v", err)
	}
	if len(files.Items) != 1 {
		t.Fatalf("len(files.Items) = %d, want 1", len(files.Items))
	}

	file, err := client.Files.Get(ctx, "file_123")
	if err != nil {
		t.Fatalf("Files.Get returned error: %v", err)
	}
	if file.Category != "Evidence" {
		t.Fatalf("file.Category = %q, want %q", file.Category, "Evidence")
	}

	disputes, err := client.Disputes.List(ctx, DisputeListParams{ListParams: ListParams{Ascending: false, Limit: 12, StartsAfter: "dp_001"}, StartTimestamp: 100, EndTimestamp: 200})
	if err != nil {
		t.Fatalf("Disputes.List returned error: %v", err)
	}
	if len(disputes.Items) != 1 {
		t.Fatalf("len(disputes.Items) = %d, want 1", len(disputes.Items))
	}

	dispute, err := client.Disputes.Get(ctx, "dp_123")
	if err != nil {
		t.Fatalf("Disputes.Get returned error: %v", err)
	}
	if dispute.Status != "RequiresAction" {
		t.Fatalf("dispute.Status = %q, want %q", dispute.Status, "RequiresAction")
	}

	acceptedDispute, err := client.Disputes.Accept(ctx, "dp_123")
	if err != nil {
		t.Fatalf("Disputes.Accept returned error: %v", err)
	}
	if acceptedDispute.Status != "Accepted" {
		t.Fatalf("acceptedDispute.Status = %q, want %q", acceptedDispute.Status, "Accepted")
	}

	challengedDispute, err := client.Disputes.Challenge(ctx, "dp_123")
	if err != nil {
		t.Fatalf("Disputes.Challenge returned error: %v", err)
	}
	if challengedDispute.Status != "Challenged" {
		t.Fatalf("challengedDispute.Status = %q, want %q", challengedDispute.Status, "Challenged")
	}

	submittedDispute, err := client.Disputes.AddEvidence(ctx, "dp_123", AddDisputeEvidenceRequest{
		Files: map[string]any{"receipt": "file_123"},
	})
	if err != nil {
		t.Fatalf("Disputes.AddEvidence returned error: %v", err)
	}
	if submittedDispute.Status != "Submitted" {
		t.Fatalf("submittedDispute.Status = %q, want %q", submittedDispute.Status, "Submitted")
	}

	updatedDispute, err := client.Disputes.DeleteEvidence(ctx, "dp_123", DeleteDisputeEvidenceRequest{
		Files: []string{"receipt"},
	})
	if err != nil {
		t.Fatalf("Disputes.DeleteEvidence returned error: %v", err)
	}
	if updatedDispute.Status != "RequiresAction" {
		t.Fatalf("updatedDispute.Status = %q, want %q", updatedDispute.Status, "RequiresAction")
	}
}

func TestNewRequestAndDoJSONErrorPaths(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/ok":
			if got := r.Header.Get("Authorization"); got != "sk_sandbox_123" {
				t.Fatalf("Authorization header = %q, want %q", got, "sk_sandbox_123")
			}
			if got := r.Header.Get("Accept"); got != "application/json" {
				t.Fatalf("Accept header = %q, want %q", got, "application/json")
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = io.WriteString(w, `{"id":"ok_123"}`)
		case "/api-error":
			w.WriteHeader(http.StatusConflict)
			_, _ = io.WriteString(w, `{"code":"duplicate_value","message":"already exists"}`)
		case "/bad-json":
			w.Header().Set("Content-Type", "application/json")
			_, _ = io.WriteString(w, `{"id":`)
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	})

	ctx := context.Background()

	req, err := client.newRequest(ctx, http.MethodPost, "ok", map[string]string{"name": "test"})
	if err != nil {
		t.Fatalf("newRequest returned error: %v", err)
	}
	if got := req.Header.Get("Content-Type"); got != "application/json" {
		t.Fatalf("Content-Type = %q, want %q", got, "application/json")
	}
	if !strings.HasSuffix(req.URL.Path, "/ok") {
		t.Fatalf("req.URL.Path = %q, want suffix /ok", req.URL.Path)
	}

	var okResponse DeletedResource
	if err := client.doJSON(req, &okResponse); err != nil {
		t.Fatalf("doJSON returned error: %v", err)
	}
	if okResponse.ID != "ok_123" {
		t.Fatalf("okResponse.ID = %q, want %q", okResponse.ID, "ok_123")
	}

	apiErrReq, err := client.newRequest(ctx, http.MethodGet, "api-error", nil)
	if err != nil {
		t.Fatalf("newRequest returned error: %v", err)
	}
	apiErr := client.doJSON(apiErrReq, &DeletedResource{})
	if apiErr == nil {
		t.Fatal("doJSON error = nil, want APIError")
	}
	typedErr, ok := apiErr.(*APIError)
	if !ok {
		t.Fatalf("error type = %T, want *APIError", apiErr)
	}
	if typedErr.Status != http.StatusConflict || typedErr.Code != "duplicate_value" {
		t.Fatalf("APIError = %+v, want status=%d code=%q", typedErr, http.StatusConflict, "duplicate_value")
	}

	badJSONReq, err := client.newRequest(ctx, http.MethodGet, "bad-json", nil)
	if err != nil {
		t.Fatalf("newRequest returned error: %v", err)
	}
	if err := client.doJSON(badJSONReq, &DeletedResource{}); err == nil || !strings.Contains(err.Error(), "decode response body") {
		t.Fatalf("doJSON error = %v, want decode response body error", err)
	}
}

func TestAPIErrorError(t *testing.T) {
	t.Parallel()

	var nilErr *APIError
	if got := nilErr.Error(); got != "" {
		t.Fatalf("nil APIError.Error() = %q, want empty", got)
	}

	if got := (&APIError{Message: "duplicate"}).Error(); got != "duplicate" {
		t.Fatalf("message APIError.Error() = %q, want %q", got, "duplicate")
	}
	if got := (&APIError{Code: "duplicate_value"}).Error(); got != "duplicate_value" {
		t.Fatalf("code APIError.Error() = %q, want %q", got, "duplicate_value")
	}
	if got := (&APIError{}).Error(); got != "ryft api error" {
		t.Fatalf("default APIError.Error() = %q, want %q", got, "ryft api error")
	}
}
