package ryft

import "context"

type WebhooksService struct {
	client *Client
}

type Webhook struct {
	ID         string   `json:"id"`
	URL        string   `json:"url,omitempty"`
	Active     bool     `json:"active"`
	EventTypes []string `json:"eventTypes,omitempty"`
}

type WebhookList struct {
	Items []Webhook `json:"items"`
}

type CreateWebhookRequest struct {
	URL        string   `json:"url"`
	Active     bool     `json:"active"`
	EventTypes []string `json:"eventTypes"`
}

type UpdateWebhookRequest struct {
	URL        string   `json:"url,omitempty"`
	Active     *bool    `json:"active,omitempty"`
	EventTypes []string `json:"eventTypes,omitempty"`
}

func (s *WebhooksService) Create(ctx context.Context, request CreateWebhookRequest) (*Webhook, error) {
	req, err := s.client.newRequest(ctx, "POST", "webhooks", request)
	if err != nil {
		return nil, err
	}

	var webhook Webhook
	if err := s.client.doJSON(req, &webhook); err != nil {
		return nil, err
	}

	return &webhook, nil
}

func (s *WebhooksService) Get(ctx context.Context, webhookID string) (*Webhook, error) {
	req, err := s.client.newRequest(ctx, "GET", "webhooks/"+webhookID, nil)
	if err != nil {
		return nil, err
	}

	var webhook Webhook
	if err := s.client.doJSON(req, &webhook); err != nil {
		return nil, err
	}

	return &webhook, nil
}

func (s *WebhooksService) List(ctx context.Context) (*WebhookList, error) {
	req, err := s.client.newRequest(ctx, "GET", "webhooks", nil)
	if err != nil {
		return nil, err
	}

	var webhooks WebhookList
	if err := s.client.doJSON(req, &webhooks); err != nil {
		return nil, err
	}

	return &webhooks, nil
}

func (s *WebhooksService) Update(ctx context.Context, webhookID string, request UpdateWebhookRequest) (*Webhook, error) {
	req, err := s.client.newRequest(ctx, "PATCH", "webhooks/"+webhookID, request)
	if err != nil {
		return nil, err
	}

	var webhook Webhook
	if err := s.client.doJSON(req, &webhook); err != nil {
		return nil, err
	}

	return &webhook, nil
}

func (s *WebhooksService) Delete(ctx context.Context, webhookID string) (*DeletedResource, error) {
	req, err := s.client.newRequest(ctx, "DELETE", "webhooks/"+webhookID, nil)
	if err != nil {
		return nil, err
	}

	var deleted DeletedResource
	if err := s.client.doJSON(req, &deleted); err != nil {
		return nil, err
	}

	return &deleted, nil
}
