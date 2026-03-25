# ryft-go

Prerelease Go SDK for the Ryft API.

This repository currently focuses on a solid core API surface with live parity coverage and Go-native tests. The API may still evolve before a stable `v1` release.

## Installation

```bash
go get github.com/bkawk/ryft-go
```

## Quickstart

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/bkawk/ryft-go"
)

func main() {
	client, err := ryft.NewClient(ryft.Config{
		SecretKey: "sk_sandbox_your_secret_key",
	})
	if err != nil {
		log.Fatal(err)
	}

	customer, err := client.Customers.Create(context.Background(), ryft.CreateCustomerRequest{
		Email:     "sdk-example@ryftpay.test",
		FirstName: "Go",
		LastName:  "Example",
		Metadata: map[string]string{
			"source": "readme",
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(customer.ID)
}
```

Or run the example directly:

```bash
export RYFT_SECRET_KEY=sk_sandbox_your_secret_key
go run ./examples/basic
```

## HTTP Server Example

There is also a small Gin example that shows how to initialize the SDK once at startup, bind an incoming request, call Ryft, and translate `*ryft.APIError` into an HTTP response:

```bash
export RYFT_SECRET_KEY=sk_sandbox_your_secret_key
go run ./examples/http-gin
```

Then create a customer through the local server:

```bash
curl -X POST http://localhost:8080/customers \
  -H "Content-Type: application/json" \
  -d '{
    "email": "sdk-example@ryftpay.test",
    "firstName": "Gin",
    "lastName": "Example",
    "metadata": {
      "source": "http-gin"
    }
  }'
```

## Configuration

`ryft.NewClient` accepts:

- `SecretKey`: required
- `BaseURL`: optional override for custom environments
- `HTTPClient`: optional custom `*http.Client`

If `BaseURL` is omitted, the SDK selects the Ryft sandbox or live API automatically from the secret key prefix.

## Current API Surface

- customers: create, list, get, update, delete
- payment sessions: create, get, update, refund
- payment-session transactions: list, get
- webhooks: create, list, get, update, delete
- accounts: create, get, update, verify, authorization links
- persons: create, list, get, update, delete
- payout methods: create, list, get, update, delete
- payouts: create, list, get
- transfers: create, list, get
- balances: list
- balance transactions: list
- payment methods: get, update, delete
- account links: create
- subscriptions: create, list, get, update, pause, resume, cancel
- subscription payment sessions: list
- events: list, get
- platform fees: list, get, refunds list
- files: create, list, get
- disputes: list, get, accept, challenge, add evidence, delete evidence

## Error Handling

API failures return `*ryft.APIError`, which includes:

- `Status`
- `Code`
- `Message`
- `RequestID`
- `Errors`

Example:

```go
import (
	"errors"
	"log"

	"github.com/bkawk/ryft-go"
)

if err != nil {
	var apiErr *ryft.APIError
	if errors.As(err, &apiErr) {
		log.Printf("ryft error: status=%d code=%s message=%s", apiErr.Status, apiErr.Code, apiErr.Message)
		return
	}

	log.Fatal(err)
}
```

## Development

The repo includes Go tests plus command-line tooling for local development and validation.

Run the local test suite with:

```bash
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out
```

GitHub Actions runs formatting, tests, and coverage summary on pushes to `main` and on pull requests.
