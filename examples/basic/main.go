package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bkawk/ryft-go"
)

func main() {
	secretKey := os.Getenv("RYFT_SECRET_KEY")
	if secretKey == "" {
		log.Fatal("RYFT_SECRET_KEY is required")
	}

	client, err := ryft.NewClient(ryft.Config{
		SecretKey: secretKey,
	})
	if err != nil {
		log.Fatal(err)
	}

	customer, err := client.Customers.Create(context.Background(), ryft.CreateCustomerRequest{
		Email:     fmt.Sprintf("example-%d@ryftpay.test", time.Now().Unix()),
		FirstName: "Go",
		LastName:  "Example",
		Metadata: map[string]string{
			"source": "examples/basic",
		},
	})
	if err != nil {
		var apiErr *ryft.APIError
		if errors.As(err, &apiErr) {
			log.Fatalf("ryft api error: status=%d code=%s message=%s", apiErr.Status, apiErr.Code, apiErr.Message)
		}
		log.Fatal(err)
	}

	fmt.Printf("Created customer %s for %s\n", customer.ID, customer.Email)
}
