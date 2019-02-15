package adyen

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	apiKey          string
	merchantAccount string
}

func NewClient(apiKey string, merchantAccount string) *Client {
	return &Client{apiKey: apiKey, merchantAccount: merchantAccount}
}

func (c *Client) CreatePaymentSession(amount PaymentAmount, ref, shopperRef, origin, returnURL string) (*SessionResponse, error) {
	req := NewPaymentSessionRequest(amount, c.merchantAccount, ref, shopperRef, origin, returnURL)
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	var session SessionResponse
	if err := c.doPost(body, PaymentSessionEndpoint, &session); err != nil {
		return nil, err
	}

	return &session, nil
}

func (c *Client) VerifyPayment(payload string) (*VerifyPaymentResponse, error) {
	if payload == "" {
		return nil, fmt.Errorf("payload required for VerifyPayment request")
	}

	req := PaymentResultPayload{payload}
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	var data VerifyPaymentResponse
	if err := c.doPost(body, VerifyPaymentEndpoint, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

func (c *Client) MakeRecurringPayment(amount PaymentAmount, recurringRef, shopperRef, ref string) (*RecurringPaymentResponse, error) {
	req := NewRecurringPaymentRequest(c.merchantAccount, recurringRef, ref, shopperRef, amount)
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	var data RecurringPaymentResponse
	if err := c.doPost(body, MakePaymentEndpoint, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

func (c *Client) doPost(body []byte, endpoint string, decodedResp interface{}) error {
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(body))
	req.Header.Set("X-API-Key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("bad Adyen request %s: status %d", endpoint, resp.StatusCode)
	}
	if err := json.NewDecoder(resp.Body).Decode(&decodedResp); err != nil {
		return err
	}
	return nil
}
