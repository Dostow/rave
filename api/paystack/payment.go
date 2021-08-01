package paystack

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/dostow/rave/api/models"
	"github.com/go-resty/resty/v2"
)

// PaymentRequest payment request
type PaymentRequest struct {
	TxRef             string   `json:"referennce"`
	Amount            string   `json:"amount"`
	Currency          string   `json:"currency"`
	CallbackURL       string   `json:"callbask_url"`
	Channels          []string `json:"channels"`
	Plan              string   `json:"plan"`
	Email             string   `json:"email"`
	TransactionCharge string   `json:"transaction_charge"`
	Subaccount        string   `json:"subaccount"`
	SplitCode         string   `json:"split_code"`
	Metadata          string   `json:"metadata,omitempty"`
}

// PaymentLinkData payment link payload
type PaymentLinkData struct {
	Link string `json:"link"`
}

// InitializePaymentResultData   result
type InitializePaymentResultData struct {
	AuthorizationURL string `json:"authorization_url"`
	AccessCode       string `json:"acecess_code"`
	Reference        string `json:"referennce"`
}

// InitializePaymentResult delete result
type InitializePaymentResult struct {
	Status  bool             `json:"status"`
	Message string           `json:"message"`
	Data    *json.RawMessage `json:"data"`
}

// InitializePayment initialize a payment
func (p *Paystack) InitializePayment(ctx context.Context, seckey string, req *models.PaymentRequest) (*models.PaymentResponse, error) {
	client := resty.New()
	metadata, _ := json.Marshal(&req.Meta)
	preq := PaymentRequest{
		TxRef:       req.TxRef,
		Amount:      req.Amount,
		Currency:    req.Currency,
		CallbackURL: req.RedirectURL,
		Channels:    strings.Split(req.PaymentOptions, ","),
		Email:       req.Customer.Email,
		Metadata:    string(metadata),
	}
	resp, err := client.R().
		EnableTrace().
		SetHeader("Authorization", "Bearer "+seckey).
		SetResult(&InitializePaymentResult{}).
		SetError(&errorResponse{}).
		SetBody(preq).
		Post(fmt.Sprintf("%s/transaction/initialize", url))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() == 200 {
		result := resp.Result().(*InitializePaymentResult)
		if result.Status {
			rsp := InitializePaymentResultData{}
			if err := json.Unmarshal(*result.Data, &rsp); err != nil {
				return nil, err
			}
			return &models.PaymentResponse{Link: rsp.AuthorizationURL, Original: result.Data}, nil
		}
		return nil, errors.New(result.Message)
	}
	respBody := resp.Error().(*errorResponse)
	return nil, errors.New(respBody.Message)
}
