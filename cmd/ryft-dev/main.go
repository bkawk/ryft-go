package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/bkawk/ryft-go"
)

func main() {
	if err := run(context.Background(), os.Args[1:]); err != nil {
		var apiErr *ryft.APIError
		if errors.As(err, &apiErr) {
			printJSON(os.Stderr, apiErr)
			os.Exit(1)
		}

		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return errors.New("usage: ryft-dev <customer-create|customer-update|entity-get> ...")
	}

	client, err := ryft.NewClient(ryft.Config{
		SecretKey: os.Getenv("RYFT_SECRET_KEY"),
	})
	if err != nil {
		return err
	}

	switch args[0] {
	case "customer-create":
		if len(args) < 2 {
			return errors.New("usage: ryft-dev customer-create <email> [first-name] [last-name] [metadata-json]")
		}

		request, err := customerCreateRequest(args[1:])
		if err != nil {
			return err
		}

		customer, err := client.Customers.Create(ctx, request)
		if err != nil {
			return err
		}

		printJSON(os.Stdout, customer)
		return nil
	case "customer-update":
		if len(args) < 2 {
			return errors.New("usage: ryft-dev customer-update <id> [first-name] [last-name] [metadata-json]")
		}

		customerID := args[1]
		request, err := customerUpdateRequest(args[2:])
		if err != nil {
			return err
		}

		customer, err := client.Customers.Update(ctx, customerID, request)
		if err != nil {
			return err
		}

		printJSON(os.Stdout, customer)
		return nil
	case "entity-get":
		if len(args) < 3 {
			return errors.New("usage: ryft-dev entity-get <entity-type> <id> [parent-id]")
		}

		entityType := args[1]
		entityID := args[2]
		parentID := ""
		if len(args) > 3 {
			parentID = args[3]
		}
		return handleEntityGet(ctx, client, entityType, entityID, parentID)
	case "payment-session-create":
		if len(args) < 1 {
			return errors.New("usage: ryft-dev payment-session-create <options-json>")
		}

		request, accountID, err := paymentSessionCreateRequest(args[1:])
		if err != nil {
			return err
		}

		paymentSession, err := client.PaymentSessions.Create(ctx, request, ryft.WithAccount(accountID))
		if err != nil {
			return err
		}

		printJSON(os.Stdout, paymentSession)
		return nil
	case "payment-session-update":
		if len(args) < 2 {
			return errors.New("usage: ryft-dev payment-session-update <id> <options-json>")
		}

		paymentSessionID := args[1]
		request, accountID, err := paymentSessionUpdateRequest(args[2:])
		if err != nil {
			return err
		}

		paymentSession, err := client.PaymentSessions.Update(ctx, paymentSessionID, request, ryft.WithAccount(accountID))
		if err != nil {
			return err
		}

		printJSON(os.Stdout, paymentSession)
		return nil
	case "webhook-create":
		if len(args) < 4 {
			return errors.New("usage: ryft-dev webhook-create <url> <active> <event-types-json>")
		}

		request, err := webhookCreateRequest(args[1:])
		if err != nil {
			return err
		}

		webhook, err := client.Webhooks.Create(ctx, request)
		if err != nil {
			return err
		}

		printJSON(os.Stdout, webhook)
		return nil
	case "webhook-update":
		if len(args) < 2 {
			return errors.New("usage: ryft-dev webhook-update <id> [url] [active] [event-types-json]")
		}

		webhookID := args[1]
		request, err := webhookUpdateRequest(args[2:])
		if err != nil {
			return err
		}

		webhook, err := client.Webhooks.Update(ctx, webhookID, request)
		if err != nil {
			return err
		}

		printJSON(os.Stdout, webhook)
		return nil
	case "webhook-delete":
		if len(args) < 2 {
			return errors.New("usage: ryft-dev webhook-delete <id>")
		}

		webhook, err := client.Webhooks.Delete(ctx, args[1])
		if err != nil {
			return err
		}

		printJSON(os.Stdout, webhook)
		return nil
	case "webhook-list":
		webhooks, err := client.Webhooks.List(ctx)
		if err != nil {
			return err
		}

		printJSON(os.Stdout, webhooks)
		return nil
	case "in-person-location-list":
		if len(args) < 2 {
			return errors.New("usage: ryft-dev in-person-location-list <limit> [account-id]")
		}
		limit, err := parseIntArg(args[1], "limit")
		if err != nil {
			return err
		}
		accountID := ""
		if len(args) > 2 {
			accountID = args[2]
		}
		locations, err := client.InPersonLocations.List(ctx, ryft.InPersonLocationListParams{ListParams: ryft.ListParams{Ascending: false, Limit: limit}}, ryft.WithAccount(accountID))
		if err != nil {
			return err
		}
		printJSON(os.Stdout, locations)
		return nil
	case "account-create":
		if len(args) < 3 {
			return errors.New("usage: ryft-dev account-create <entity-type> <email> [metadata-json] [mode] [onboarding-flow]")
		}
		request, err := accountCreateRequest(args[1:])
		if err != nil {
			return err
		}
		account, err := client.Accounts.Create(ctx, request)
		if err != nil {
			return err
		}
		printJSON(os.Stdout, account)
		return nil
	case "account-verify":
		if len(args) < 2 {
			return errors.New("usage: ryft-dev account-verify <account-id>")
		}
		account, err := client.Accounts.Verify(ctx, args[1])
		if err != nil {
			return err
		}
		printJSON(os.Stdout, account)
		return nil
	case "account-authorize":
		if len(args) < 3 {
			return errors.New("usage: ryft-dev account-authorize <email> <redirect-url>")
		}
		authorization, err := client.Accounts.CreateAuthLink(ctx, ryft.CreateAccountAuthorizationRequest{
			Email:       args[1],
			RedirectURL: args[2],
		})
		if err != nil {
			return err
		}
		printJSON(os.Stdout, authorization)
		return nil
	case "person-create":
		if len(args) < 3 {
			return errors.New("usage: ryft-dev person-create <account-id> <email> [metadata-json]")
		}
		accountID := args[1]
		request, err := personCreateRequest(args[2:])
		if err != nil {
			return err
		}
		person, err := client.Persons.Create(ctx, accountID, request)
		if err != nil {
			return err
		}
		printJSON(os.Stdout, person)
		return nil
	case "person-list":
		if len(args) < 3 {
			return errors.New("usage: ryft-dev person-list <limit> <account-id>")
		}
		limit, err := parseIntArg(args[1], "limit")
		if err != nil {
			return err
		}
		people, err := client.Persons.List(ctx, args[2], ryft.PersonListParams{ListParams: ryft.ListParams{Ascending: true, Limit: limit}})
		if err != nil {
			return err
		}
		printJSON(os.Stdout, people)
		return nil
	case "payout-method-create":
		if len(args) < 3 {
			return errors.New("usage: ryft-dev payout-method-create <account-id> <display-name>")
		}
		accountID := args[1]
		request := payoutMethodCreateRequest(args[2:])
		payoutMethod, err := client.PayoutMethods.Create(ctx, accountID, request)
		if err != nil {
			return err
		}
		printJSON(os.Stdout, payoutMethod)
		return nil
	case "payout-method-list":
		if len(args) < 3 {
			return errors.New("usage: ryft-dev payout-method-list <limit> <account-id>")
		}
		limit, err := parseIntArg(args[1], "limit")
		if err != nil {
			return err
		}
		payoutMethods, err := client.PayoutMethods.List(ctx, args[2], ryft.PayoutMethodListParams{ListParams: ryft.ListParams{Ascending: true, Limit: limit}})
		if err != nil {
			return err
		}
		printJSON(os.Stdout, payoutMethods)
		return nil
	case "payout-create":
		if len(args) < 5 {
			return errors.New("usage: ryft-dev payout-create <account-id> <amount> <currency> <payout-method-id> [metadata-json]")
		}
		accountID := args[1]
		request, err := payoutCreateRequest(args[2:])
		if err != nil {
			return err
		}
		payout, err := client.Payouts.Create(ctx, accountID, request)
		if err != nil {
			return err
		}
		printJSON(os.Stdout, payout)
		return nil
	case "transfer-create":
		if len(args) < 4 {
			return errors.New("usage: ryft-dev transfer-create <destination-account-id> <amount> <currency> [metadata-json]")
		}
		request, err := transferCreateRequest(args[1:])
		if err != nil {
			return err
		}
		transfer, err := client.Transfers.Create(ctx, request)
		if err != nil {
			return err
		}
		printJSON(os.Stdout, transfer)
		return nil
	case "transfer-list":
		if len(args) < 2 {
			return errors.New("usage: ryft-dev transfer-list <limit>")
		}
		limit, err := parseIntArg(args[1], "limit")
		if err != nil {
			return err
		}
		transfers, err := client.Transfers.List(ctx, ryft.TransferListParams{ListParams: ryft.ListParams{Limit: limit}})
		if err != nil {
			return err
		}
		printJSON(os.Stdout, transfers)
		return nil
	case "balance-list":
		if len(args) < 3 {
			return errors.New("usage: ryft-dev balance-list <currency> <account-id>")
		}
		balances, err := client.Balances.List(ctx, ryft.BalanceListParams{Currency: args[1]}, ryft.WithAccount(args[2]))
		if err != nil {
			return err
		}
		printJSON(os.Stdout, balances)
		return nil
	case "balance-transaction-list":
		if len(args) < 3 {
			return errors.New("usage: ryft-dev balance-transaction-list <limit> <account-id>")
		}
		limit, err := parseIntArg(args[1], "limit")
		if err != nil {
			return err
		}
		balanceTransactions, err := client.BalanceTransactions.List(ctx, ryft.BalanceTransactionListParams{ListParams: ryft.ListParams{Limit: limit}}, ryft.WithAccount(args[2]))
		if err != nil {
			return err
		}
		printJSON(os.Stdout, balanceTransactions)
		return nil
	case "event-list":
		accountID := ""
		if len(args) > 1 {
			accountID = args[1]
		}
		events, err := client.Events.List(ctx, ryft.EventListParams{ListParams: ryft.ListParams{Ascending: false, Limit: 50}}, ryft.WithAccount(accountID))
		if err != nil {
			return err
		}
		printJSON(os.Stdout, events)
		return nil
	case "platform-fee-list":
		fees, err := client.PlatformFees.List(ctx, ryft.PlatformFeeListParams{ListParams: ryft.ListParams{Ascending: false, Limit: 50}})
		if err != nil {
			return err
		}
		printJSON(os.Stdout, fees)
		return nil
	case "platform-fee-refund-list":
		if len(args) < 2 {
			return errors.New("usage: ryft-dev platform-fee-refund-list <platform-fee-id>")
		}
		refunds, err := client.PlatformFees.GetRefunds(ctx, args[1])
		if err != nil {
			return err
		}
		printJSON(os.Stdout, refunds)
		return nil
	case "file-list":
		category := ""
		if len(args) > 1 {
			category = args[1]
		}
		files, err := client.Files.List(ctx, ryft.FileListParams{ListParams: ryft.ListParams{Ascending: false, Limit: 50}, Category: category})
		if err != nil {
			return err
		}
		printJSON(os.Stdout, files)
		return nil
	case "file-create":
		if len(args) < 2 {
			return errors.New("usage: ryft-dev file-create <file-path> [category]")
		}
		category := "Evidence"
		if len(args) > 2 && args[2] != "" {
			category = args[2]
		}
		file, err := client.Files.Create(ctx, ryft.CreateFileRequest{
			FilePath: args[1],
			Category: category,
		})
		if err != nil {
			return err
		}
		printJSON(os.Stdout, file)
		return nil
	case "dispute-list":
		disputes, err := client.Disputes.List(ctx, ryft.DisputeListParams{ListParams: ryft.ListParams{Ascending: false, Limit: 50}})
		if err != nil {
			return err
		}
		printJSON(os.Stdout, disputes)
		return nil
	case "dispute-accept":
		if len(args) < 2 {
			return errors.New("usage: ryft-dev dispute-accept <dispute-id>")
		}
		dispute, err := client.Disputes.Accept(ctx, args[1])
		if err != nil {
			return err
		}
		printJSON(os.Stdout, dispute)
		return nil
	case "dispute-challenge":
		if len(args) < 2 {
			return errors.New("usage: ryft-dev dispute-challenge <dispute-id>")
		}
		dispute, err := client.Disputes.Challenge(ctx, args[1])
		if err != nil {
			return err
		}
		printJSON(os.Stdout, dispute)
		return nil
	case "dispute-add-evidence":
		if len(args) < 3 {
			return errors.New("usage: ryft-dev dispute-add-evidence <dispute-id> <file-id>")
		}
		dispute, err := client.Disputes.AddEvidence(ctx, args[1], ryft.AddDisputeEvidenceRequest{
			Files: map[string]any{
				"uncategorised": map[string]any{
					"id": args[2],
				},
			},
		})
		if err != nil {
			return err
		}
		printJSON(os.Stdout, dispute)
		return nil
	case "dispute-delete-evidence":
		if len(args) < 2 {
			return errors.New("usage: ryft-dev dispute-delete-evidence <dispute-id>")
		}
		dispute, err := client.Disputes.DeleteEvidence(ctx, args[1], ryft.DeleteDisputeEvidenceRequest{
			Files: []string{"uncategorised"},
		})
		if err != nil {
			return err
		}
		printJSON(os.Stdout, dispute)
		return nil
	case "customer-payment-method-list":
		if len(args) < 2 {
			return errors.New("usage: ryft-dev customer-payment-method-list <customer-id>")
		}
		paymentMethods, err := client.Customers.GetPaymentMethods(ctx, args[1])
		if err != nil {
			return err
		}
		printJSON(os.Stdout, paymentMethods)
		return nil
	case "payment-method-update":
		if len(args) < 3 {
			return errors.New("usage: ryft-dev payment-method-update <id> <billing-address-json>")
		}
		billingAddress, err := parseMetadata(args[2])
		if err != nil {
			return err
		}
		paymentMethod, err := client.PaymentMethods.Update(ctx, args[1], ryft.UpdatePaymentMethodRequest{
			BillingAddress: billingAddress,
		})
		if err != nil {
			return err
		}
		printJSON(os.Stdout, paymentMethod)
		return nil
	case "payment-method-delete":
		if len(args) < 2 {
			return errors.New("usage: ryft-dev payment-method-delete <id>")
		}
		deleted, err := client.PaymentMethods.Delete(ctx, args[1])
		if err != nil {
			return err
		}
		printJSON(os.Stdout, deleted)
		return nil
	case "account-link-create":
		if len(args) < 3 {
			return errors.New("usage: ryft-dev account-link-create <account-id> <redirect-url>")
		}
		link, err := client.AccountLinks.GenerateTemporaryAccountLink(
			ctx,
			ryft.CreateTemporaryAccountLinkRequest{
				AccountID:   args[1],
				RedirectURL: args[2],
			},
		)
		if err != nil {
			return err
		}
		printJSON(os.Stdout, link)
		return nil
	case "subscription-create":
		if len(args) < 3 {
			return errors.New("usage: ryft-dev subscription-create <customer-id> <payment-method-id> [options-json]")
		}
		request, err := subscriptionCreateRequest(args[1:])
		if err != nil {
			return err
		}
		subscription, err := client.Subscriptions.Create(ctx, request)
		if err != nil {
			return err
		}
		printJSON(os.Stdout, subscription)
		return nil
	case "subscription-update":
		if len(args) < 2 {
			return errors.New("usage: ryft-dev subscription-update <subscription-id> [description] [metadata-json]")
		}
		request, err := subscriptionUpdateRequest(args[2:])
		if err != nil {
			return err
		}
		subscription, err := client.Subscriptions.Update(ctx, args[1], request)
		if err != nil {
			return err
		}
		printJSON(os.Stdout, subscription)
		return nil
	case "subscription-list":
		startTimestamp, endTimestamp := collectionWindowFromEnv()
		subscriptions, err := client.Subscriptions.List(ctx, ryft.SubscriptionListParams{ListParams: ryft.ListParams{Ascending: false, Limit: 10}, StartTimestamp: startTimestamp, EndTimestamp: endTimestamp})
		if err != nil {
			return err
		}
		printJSON(os.Stdout, subscriptions)
		return nil
	case "subscription-payment-session-list":
		if len(args) < 2 {
			return errors.New("usage: ryft-dev subscription-payment-session-list <subscription-id>")
		}
		startTimestamp, endTimestamp := collectionWindowFromEnv()
		paymentSessions, err := client.Subscriptions.GetPaymentSessions(
			ctx,
			args[1],
			ryft.SubscriptionListParams{
				ListParams:     ryft.ListParams{Ascending: false, Limit: 10},
				StartTimestamp: startTimestamp,
				EndTimestamp:   endTimestamp,
			},
		)
		if err != nil {
			return err
		}
		printJSON(os.Stdout, paymentSessions)
		return nil
	case "payment-session-refund":
		if len(args) < 2 {
			return errors.New("usage: ryft-dev payment-session-refund <payment-session-id> [options-json]")
		}
		request, err := refundPaymentSessionRequest(args[2:])
		if err != nil {
			return err
		}
		refund, err := client.PaymentSessions.Refund(ctx, args[1], request)
		if err != nil {
			return err
		}
		printJSON(os.Stdout, refund)
		return nil
	case "subscription-pause":
		if len(args) < 2 {
			return errors.New("usage: ryft-dev subscription-pause <subscription-id> [reason] [resume-timestamp] [unschedule]")
		}
		request, err := subscriptionPauseRequest(args[2:])
		if err != nil {
			return err
		}
		subscription, err := client.Subscriptions.Pause(ctx, args[1], request)
		if err != nil {
			return err
		}
		printJSON(os.Stdout, subscription)
		return nil
	case "subscription-resume":
		if len(args) < 2 {
			return errors.New("usage: ryft-dev subscription-resume <subscription-id>")
		}
		subscription, err := client.Subscriptions.Resume(ctx, args[1])
		if err != nil {
			return err
		}
		printJSON(os.Stdout, subscription)
		return nil
	case "subscription-cancel":
		if len(args) < 2 {
			return errors.New("usage: ryft-dev subscription-cancel <subscription-id>")
		}
		deleted, err := client.Subscriptions.Cancel(ctx, args[1])
		if err != nil {
			return err
		}
		printJSON(os.Stdout, deleted)
		return nil
	default:
		return fmt.Errorf("unknown command: %s", args[0])
	}
}

func handleEntityGet(ctx context.Context, client *ryft.Client, entityType string, entityID string, parentID string) error {
	switch entityType {
	case "customer":
		customer, err := client.Customers.Get(ctx, entityID)
		if err != nil {
			return err
		}
		printJSON(os.Stdout, customer)
		return nil
	case "payment-session":
		paymentSession, err := client.PaymentSessions.Get(ctx, entityID, ryft.WithAccount(parentID))
		if err != nil {
			return err
		}
		printJSON(os.Stdout, paymentSession)
		return nil
	case "webhook":
		webhook, err := client.Webhooks.Get(ctx, entityID)
		if err != nil {
			return err
		}
		printJSON(os.Stdout, webhook)
		return nil
	case "account":
		account, err := client.Accounts.Get(ctx, entityID)
		if err != nil {
			return err
		}
		printJSON(os.Stdout, account)
		return nil
	case "person":
		if parentID == "" {
			return errors.New("person requires parent account id")
		}
		person, err := client.Persons.Get(ctx, parentID, entityID)
		if err != nil {
			return err
		}
		printJSON(os.Stdout, person)
		return nil
	case "payout-method":
		if parentID == "" {
			return errors.New("payout-method requires parent account id")
		}
		payoutMethod, err := client.PayoutMethods.Get(ctx, parentID, entityID)
		if err != nil {
			return err
		}
		printJSON(os.Stdout, payoutMethod)
		return nil
	case "payout":
		if parentID == "" {
			return errors.New("payout requires parent account id")
		}
		payout, err := client.Payouts.Get(ctx, parentID, entityID)
		if err != nil {
			return err
		}
		printJSON(os.Stdout, payout)
		return nil
	case "transfer":
		transfer, err := client.Transfers.Get(ctx, entityID)
		if err != nil {
			return err
		}
		printJSON(os.Stdout, transfer)
		return nil
	case "payment-method":
		paymentMethod, err := client.PaymentMethods.Get(ctx, entityID)
		if err != nil {
			return err
		}
		printJSON(os.Stdout, paymentMethod)
		return nil
	case "subscription":
		subscription, err := client.Subscriptions.Get(ctx, entityID)
		if err != nil {
			return err
		}
		printJSON(os.Stdout, subscription)
		return nil
	case "event":
		event, err := client.Events.Get(ctx, entityID, ryft.WithAccount(parentID))
		if err != nil {
			return err
		}
		printJSON(os.Stdout, event)
		return nil
	case "payment-transaction":
		if parentID == "" {
			return errors.New("payment-transaction requires parent payment-session id")
		}
		transaction, err := client.PaymentSessions.GetTransaction(ctx, parentID, entityID)
		if err != nil {
			return err
		}
		printJSON(os.Stdout, transaction)
		return nil
	case "platform-fee":
		fee, err := client.PlatformFees.Get(ctx, entityID)
		if err != nil {
			return err
		}
		printJSON(os.Stdout, fee)
		return nil
	case "file":
		file, err := client.Files.Get(ctx, entityID)
		if err != nil {
			return err
		}
		printJSON(os.Stdout, file)
		return nil
	case "dispute":
		dispute, err := client.Disputes.Get(ctx, entityID)
		if err != nil {
			return err
		}
		printJSON(os.Stdout, dispute)
		return nil
	case "in-person-location":
		location, err := client.InPersonLocations.Get(ctx, entityID, ryft.WithAccount(parentID))
		if err != nil {
			return err
		}
		printJSON(os.Stdout, location)
		return nil
	default:
		return fmt.Errorf("unsupported entity type: %s", entityType)
	}
}

func customerCreateRequest(args []string) (ryft.CreateCustomerRequest, error) {
	request := ryft.CreateCustomerRequest{
		Email: args[0],
	}

	if len(args) > 1 {
		request.FirstName = args[1]
	}
	if len(args) > 2 {
		request.LastName = args[2]
	}
	if len(args) > 3 && args[3] != "" {
		metadata, err := parseMetadata(args[3])
		if err != nil {
			return ryft.CreateCustomerRequest{}, err
		}
		request.Metadata = metadata
	}

	return request, nil
}

func customerUpdateRequest(args []string) (ryft.UpdateCustomerRequest, error) {
	var request ryft.UpdateCustomerRequest

	if len(args) > 0 {
		request.FirstName = args[0]
	}
	if len(args) > 1 {
		request.LastName = args[1]
	}
	if len(args) > 2 && args[2] != "" {
		metadata, err := parseMetadata(args[2])
		if err != nil {
			return ryft.UpdateCustomerRequest{}, err
		}
		request.Metadata = metadata
	}

	return request, nil
}

func paymentSessionCreateRequest(args []string) (ryft.CreatePaymentSessionRequest, string, error) {
	options := map[string]any{}
	if len(args) > 0 && args[0] != "" {
		if err := json.Unmarshal([]byte(args[0]), &options); err != nil {
			return ryft.CreatePaymentSessionRequest{}, "", fmt.Errorf("parse payment session options json: %w", err)
		}
	}

	request := ryft.CreatePaymentSessionRequest{
		Amount:        intFromAny(options["amount"], 500),
		Currency:      stringFromAny(options["currency"], "GBP"),
		CustomerEmail: stringFromAny(options["customerEmail"], ""),
		PaymentType:   stringFromAny(options["paymentType"], "Standard"),
		EntryMode:     stringFromAny(options["entryMode"], "Online"),
		CaptureFlow:   stringFromAny(options["captureFlow"], "Automatic"),
		ReturnURL:     stringFromAny(options["returnUrl"], "https://example.com/return"),
		PlatformFee:   intFromAny(options["platformFee"], 0),
	}

	if rawMetadata, ok := options["metadata"].(map[string]any); ok {
		metadata := map[string]string{}
		for key, value := range rawMetadata {
			metadata[key] = fmt.Sprint(value)
		}
		request.Metadata = metadata
	}

	if customerID := stringFromAny(options["customerId"], ""); customerID != "" {
		request.CustomerDetails = &ryft.PaymentSessionCustomer{ID: customerID}
	}

	if previousPaymentID := stringFromAny(options["previousPaymentId"], ""); previousPaymentID != "" {
		request.PreviousPayment = &ryft.PaymentSessionReference{ID: previousPaymentID}
	}

	if rawRebilling, ok := options["rebillingDetail"].(map[string]any); ok {
		request.RebillingDetail = &ryft.RebillingDetail{
			AmountVariance:              stringFromAny(rawRebilling["amountVariance"], ""),
			NumberOfDaysBetweenPayments: intFromAny(rawRebilling["numberOfDaysBetweenPayments"], 0),
			TotalNumberOfPayments:       intFromAny(rawRebilling["totalNumberOfPayments"], 0),
			CurrentPaymentNumber:        intFromAny(rawRebilling["currentPaymentNumber"], 0),
		}
	}

	if rawSplits, ok := options["splits"].(map[string]any); ok {
		request.Splits = buildSplitPaymentRequest(rawSplits)
	}

	if rawAttemptPayment, ok := options["attemptPayment"].(map[string]any); ok {
		payload, err := json.Marshal(rawAttemptPayment)
		if err != nil {
			return ryft.CreatePaymentSessionRequest{}, "", fmt.Errorf("marshal attemptPayment: %w", err)
		}
		request.AttemptPayment = payload
	}

	return request, stringFromAny(options["accountId"], ""), nil
}

func paymentSessionUpdateRequest(args []string) (ryft.UpdatePaymentSessionRequest, string, error) {
	var request ryft.UpdatePaymentSessionRequest
	accountID := ""

	if len(args) == 0 || args[0] == "" {
		return request, accountID, nil
	}

	options := map[string]any{}
	if err := json.Unmarshal([]byte(args[0]), &options); err != nil {
		return ryft.UpdatePaymentSessionRequest{}, "", fmt.Errorf("parse payment session update options json: %w", err)
	}

	if value, ok := options["amount"]; ok {
		amount := intFromAny(value, 0)
		request.Amount = &amount
	}
	request.CustomerEmail = stringFromAny(options["customerEmail"], "")
	request.CaptureFlow = stringFromAny(options["captureFlow"], "")
	accountID = stringFromAny(options["accountId"], "")

	if rawMetadata, ok := options["metadata"].(map[string]any); ok {
		metadata := map[string]string{}
		for key, value := range rawMetadata {
			metadata[key] = fmt.Sprint(value)
		}
		request.Metadata = metadata
	}

	return request, accountID, nil
}

func buildSplitPaymentRequest(rawSplits map[string]any) *ryft.CreateSplitPaymentRequest {
	rawItems, ok := rawSplits["items"].([]any)
	if !ok || len(rawItems) == 0 {
		return nil
	}

	items := make([]ryft.SplitItem, 0, len(rawItems))
	for _, rawItem := range rawItems {
		itemMap, ok := rawItem.(map[string]any)
		if !ok {
			continue
		}

		item := ryft.SplitItem{
			AccountID:   stringFromAny(itemMap["accountId"], ""),
			Amount:      intFromAny(itemMap["amount"], 0),
			Description: stringFromAny(itemMap["description"], ""),
		}

		if rawFee, ok := itemMap["fee"].(map[string]any); ok {
			item.Fee = &ryft.SplitFee{
				Amount: intFromAny(rawFee["amount"], 0),
			}
		}

		if rawMetadata, ok := itemMap["metadata"].(map[string]any); ok {
			metadata := map[string]string{}
			for key, value := range rawMetadata {
				metadata[key] = fmt.Sprint(value)
			}
			item.Metadata = metadata
		}

		items = append(items, item)
	}

	if len(items) == 0 {
		return nil
	}

	return &ryft.CreateSplitPaymentRequest{Items: items}
}

func parseIntArg(raw string, fieldName string) (int, error) {
	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("parse %s: %w", fieldName, err)
	}
	return value, nil
}

func webhookCreateRequest(args []string) (ryft.CreateWebhookRequest, error) {
	eventTypes, err := parseStringSlice(args[2], "event-types")
	if err != nil {
		return ryft.CreateWebhookRequest{}, err
	}

	return ryft.CreateWebhookRequest{
		URL:        args[0],
		Active:     parseBoolArg(args[1]),
		EventTypes: eventTypes,
	}, nil
}

func webhookUpdateRequest(args []string) (ryft.UpdateWebhookRequest, error) {
	var request ryft.UpdateWebhookRequest

	if len(args) > 0 {
		request.URL = args[0]
	}
	if len(args) > 1 && args[1] != "" {
		active := parseBoolArg(args[1])
		request.Active = &active
	}
	if len(args) > 2 && args[2] != "" {
		eventTypes, err := parseStringSlice(args[2], "event-types")
		if err != nil {
			return ryft.UpdateWebhookRequest{}, err
		}
		request.EventTypes = eventTypes
	}

	return request, nil
}

func parseBoolArg(raw string) bool {
	return raw == "true" || raw == "TRUE" || raw == "True"
}

func parseStringSlice(raw string, fieldName string) ([]string, error) {
	var values []string
	if err := json.Unmarshal([]byte(raw), &values); err != nil {
		return nil, fmt.Errorf("parse %s json: %w", fieldName, err)
	}
	return values, nil
}

func accountCreateRequest(args []string) (ryft.CreateAccountRequest, error) {
	entityType := args[0]
	email := args[1]
	mode := "standard"
	if len(args) > 3 && args[3] != "" {
		mode = args[3]
	}
	onboardingFlow := "NonHosted"
	if len(args) > 4 && args[4] != "" {
		onboardingFlow = args[4]
	}

	metadata := map[string]string{
		"harness":    "sdk-write-platform-resources",
		"scenario":   "account-create",
		"sdk":        "ryft-go",
		"entityType": entityType,
	}
	if len(args) > 2 && args[2] != "" {
		parsed, err := parseMetadata(args[2])
		if err != nil {
			return ryft.CreateAccountRequest{}, err
		}
		for key, value := range parsed {
			metadata[key] = value
		}
	}

	request := ryft.CreateAccountRequest{
		OnboardingFlow: onboardingFlow,
		EntityType:     entityType,
		Email:          email,
		Metadata:       metadata,
		TermsOfService: ryft.TermsOfService{
			Acceptance: ryft.Acceptance{
				IPAddress: "127.0.0.1",
			},
		},
	}

	address := ryft.Address{
		LineOne:    "1 Test Street",
		City:       "London",
		Country:    "GB",
		PostalCode: "EC1A 1BB",
	}
	if mode == "invalid-individual-business" {
		request.Business = &ryft.BusinessAccountDetails{
			Name:               fmt.Sprintf("Invalid Business %s", strings.Split(email, "@")[0]),
			Type:               "PrivateCompany",
			RegistrationNumber: "12345678",
			RegisteredAddress:  address,
			ContactEmail:       email,
		}
	} else if entityType == "Individual" {
		request.Individual = &ryft.IndividualAccountDetails{
			FirstName:     "Sdk",
			LastName:      "Individual",
			Email:         email,
			DateOfBirth:   "1990-01-20",
			Gender:        "Male",
			Nationalities: []string{"GB"},
			Address:       address,
		}
	} else {
		request.Business = &ryft.BusinessAccountDetails{
			Name:               fmt.Sprintf("SDK Account Business %s", strings.Split(email, "@")[0]),
			Type:               "PrivateCompany",
			RegistrationNumber: "12345678",
			RegisteredAddress:  address,
			ContactEmail:       email,
		}
	}

	return request, nil
}

func personCreateRequest(args []string) (ryft.CreatePersonRequest, error) {
	email := args[0]
	metadata := map[string]string{
		"harness":  "sdk-write-platform-resources",
		"scenario": "person-create",
		"sdk":      "ryft-go",
	}
	if len(args) > 1 && args[1] != "" {
		parsed, err := parseMetadata(args[1])
		if err != nil {
			return ryft.CreatePersonRequest{}, err
		}
		for key, value := range parsed {
			metadata[key] = value
		}
	}

	return ryft.CreatePersonRequest{
		FirstName:     "Sdk",
		LastName:      "Person",
		Email:         email,
		DateOfBirth:   "1990-01-15",
		Gender:        "Male",
		Nationalities: []string{"GB"},
		Address: ryft.Address{
			LineOne:    "1 Test Street",
			City:       "London",
			Country:    "GB",
			PostalCode: "EC1A 1BB",
		},
		PhoneNumber:   "+447000000000",
		BusinessRoles: []string{"Director"},
		Documents:     []map[string]any{},
		Metadata:      metadata,
	}, nil
}

func payoutMethodCreateRequest(args []string) ryft.CreatePayoutMethodRequest {
	displayName := args[0]
	return ryft.CreatePayoutMethodRequest{
		Type:        "BankAccount",
		DisplayName: displayName,
		Currency:    "GBP",
		Country:     "GB",
		BankAccount: ryft.BankAccountDetails{
			AccountNumberType: "UnitedKingdom",
			AccountNumber:     "31926819",
			BankIDType:        "SortCode",
			BankID:            "601613",
		},
	}
}

func payoutCreateRequest(args []string) (ryft.CreatePayoutRequest, error) {
	amount, err := parseIntArg(args[0], "amount")
	if err != nil {
		return ryft.CreatePayoutRequest{}, err
	}

	request := ryft.CreatePayoutRequest{
		Amount:         amount,
		Currency:       args[1],
		PayoutMethodID: args[2],
	}
	if len(args) > 3 && args[3] != "" {
		metadata, err := parseMetadata(args[3])
		if err != nil {
			return ryft.CreatePayoutRequest{}, err
		}
		request.Metadata = metadata
	}

	return request, nil
}

func transferCreateRequest(args []string) (ryft.CreateTransferRequest, error) {
	amount, err := parseIntArg(args[1], "amount")
	if err != nil {
		return ryft.CreateTransferRequest{}, err
	}

	request := ryft.CreateTransferRequest{
		Amount:   amount,
		Currency: args[2],
		Destination: ryft.TransferDestination{
			AccountID: args[0],
		},
		Reason: "webhook-events test transfer",
	}
	if len(args) > 3 && args[3] != "" {
		metadata, err := parseMetadata(args[3])
		if err != nil {
			return ryft.CreateTransferRequest{}, err
		}
		request.Metadata = metadata
	}

	return request, nil
}

func subscriptionCreateRequest(args []string) (ryft.CreateSubscriptionRequest, error) {
	customerID := args[0]
	paymentMethodID := args[1]
	options := map[string]any{}
	if len(args) > 2 && args[2] != "" {
		if err := json.Unmarshal([]byte(args[2]), &options); err != nil {
			return ryft.CreateSubscriptionRequest{}, fmt.Errorf("parse options json: %w", err)
		}
	}

	metadata := map[string]string{
		"harness":  "sdk-subscription-readiness",
		"scenario": "subscription-create",
		"sdk":      "ryft-go",
	}
	if rawMetadata, ok := options["metadata"].(map[string]any); ok {
		for key, value := range rawMetadata {
			metadata[key] = fmt.Sprint(value)
		}
	}

	request := ryft.CreateSubscriptionRequest{
		Customer:      ryft.SubscriptionCustomerReference{ID: customerID},
		PaymentMethod: ryft.SubscriptionPaymentMethodReference{ID: paymentMethodID},
		Description:   stringFromAny(options["description"], "SDK subscription readiness"),
		Metadata:      metadata,
		Price: ryft.SubscriptionPrice{
			Amount:   100,
			Currency: "GBP",
			Interval: ryft.SubscriptionInterval{
				Unit:  "Months",
				Count: 1,
				Times: 12,
			},
		},
		PaymentSettings: map[string]any{
			"statementDescriptor": map[string]any{
				"descriptor": "Ryft Ltd",
				"city":       "London",
			},
		},
	}

	if rawPrice, ok := options["price"].(map[string]any); ok {
		request.Price = ryft.SubscriptionPrice{
			Amount:   intFromAny(rawPrice["amount"], 100),
			Currency: stringFromAny(rawPrice["currency"], "GBP"),
			Interval: ryft.SubscriptionInterval{
				Unit:  stringFromAny(mapValue(rawPrice, "interval", "unit"), "Months"),
				Count: intFromAny(mapValue(rawPrice, "interval", "count"), 1),
				Times: intFromAny(mapValue(rawPrice, "interval", "times"), 12),
			},
		}
	}

	if value, ok := options["billingCycleTimestamp"]; ok {
		request.BillingCycleTimestamp = intFromAny(value, 0)
	}
	if rawSettings, ok := options["paymentSettings"].(map[string]any); ok {
		request.PaymentSettings = rawSettings
	}
	if rawShipping, ok := options["shippingDetails"].(map[string]any); ok {
		request.ShippingDetails = rawShipping
	}

	return request, nil
}

func subscriptionUpdateRequest(args []string) (ryft.UpdateSubscriptionRequest, error) {
	var request ryft.UpdateSubscriptionRequest
	if len(args) > 0 {
		request.Description = args[0]
	}
	if len(args) > 1 && args[1] != "" {
		metadata, err := parseMetadata(args[1])
		if err != nil {
			return ryft.UpdateSubscriptionRequest{}, err
		}
		request.Metadata = metadata
	}
	return request, nil
}

func subscriptionPauseRequest(args []string) (ryft.PauseSubscriptionRequest, error) {
	var request ryft.PauseSubscriptionRequest
	if len(args) > 0 {
		request.Reason = args[0]
	}
	if len(args) > 1 && args[1] != "" {
		resumeTimestamp, err := parseIntArg(args[1], "resume-timestamp")
		if err != nil {
			return ryft.PauseSubscriptionRequest{}, err
		}
		request.ResumeTimestamp = resumeTimestamp
	}
	if len(args) > 2 && args[2] != "" {
		request.Unschedule = parseBoolArg(args[2])
	}
	return request, nil
}

func refundPaymentSessionRequest(args []string) (ryft.RefundPaymentSessionRequest, error) {
	var request ryft.RefundPaymentSessionRequest
	if len(args) == 0 || args[0] == "" {
		return request, nil
	}

	var options map[string]any
	if err := json.Unmarshal([]byte(args[0]), &options); err != nil {
		return ryft.RefundPaymentSessionRequest{}, fmt.Errorf("parse refund options json: %w", err)
	}

	request.Amount = intFromAny(options["amount"], 0)
	request.Reason = stringFromAny(options["reason"], "")
	if value, ok := options["refundPlatformFee"].(bool); ok {
		request.RefundPlatformFee = value
	}
	return request, nil
}

func collectionWindowFromEnv() (int, int) {
	return intFromEnv("RYFT_COLLECTION_START_TIMESTAMP"), intFromEnv("RYFT_COLLECTION_END_TIMESTAMP")
}

func intFromEnv(name string) int {
	raw := os.Getenv(name)
	if raw == "" {
		return 0
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0
	}
	return value
}

func intFromAny(value any, fallback int) int {
	switch typed := value.(type) {
	case float64:
		return int(typed)
	case int:
		return typed
	case int64:
		return int(typed)
	default:
		return fallback
	}
}

func stringFromAny(value any, fallback string) string {
	if typed, ok := value.(string); ok && typed != "" {
		return typed
	}
	return fallback
}

func mapValue(raw map[string]any, parent string, child string) any {
	parentMap, ok := raw[parent].(map[string]any)
	if !ok {
		return nil
	}
	return parentMap[child]
}

func parseMetadata(raw string) (map[string]string, error) {
	var metadata map[string]any
	if err := json.Unmarshal([]byte(raw), &metadata); err != nil {
		return nil, fmt.Errorf("parse metadata json: %w", err)
	}

	normalized := make(map[string]string, len(metadata))
	for key, value := range metadata {
		normalized[key] = fmt.Sprint(value)
	}

	return normalized, nil
}

func printJSON(file *os.File, value any) {
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(value)
}
