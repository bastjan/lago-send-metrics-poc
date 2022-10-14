package lago

import (
	"context"
	"fmt"
)

type EventRequest struct {
	client *Client
}

type EventParams struct {
	Event *EventInput `json:"event"`
}

type EventInput struct {
	TransactionID          string            `json:"transaction_id,omitempty"`
	ExternalCustomerID     string            `json:"external_customer_id,omitempty"`
	ExternalSubscriptionID string            `json:"external_subscription_id,omitempty"`
	Code                   string            `json:"code,omitempty"`
	Timestamp              int64             `json:"timestamp,omitempty"`
	Properties             map[string]string `json:"properties,omitempty"`
}

func (c *Client) Event() *EventRequest {
	return &EventRequest{
		client: c,
	}
}

func (er *EventRequest) Create(ctx context.Context, eventInput *EventInput) *Error {
	eventParams := &EventParams{
		Event: eventInput,
	}

	clientRequest := &ClientRequest{
		Path: "events",
		Body: eventParams,
	}

	err := er.client.PostWithoutResult(ctx, clientRequest)
	if err != nil {
		return err
	}

	return nil
}

func (cr *EventRequest) Get(ctx context.Context, transactionID string) (map[string]any, *Error) {
	subPath := fmt.Sprintf("%s/%s", "events", transactionID)
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &map[string]any{},
	}

	result, err := cr.client.Get(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	r := *(result.(*map[string]any))
	return r, nil
}
