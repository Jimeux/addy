package adyen

import (
	"fmt"
	"time"
)

const (
	apiURL = "https://checkout-test.adyen.com/v40"
	// https://docs.adyen.com/api-explorer/#/PaymentSetupAndVerificationService/v40/paymentSession
	PaymentSessionEndpoint = apiURL + "/paymentSession"
	// https://docs.adyen.com/api-explorer/#/PaymentSetupAndVerificationService/v40/payments/result
	PaymentVerificationEndpoint = apiURL + "/payments/result"
	// https://docs.adyen.com/api-explorer/#/PaymentSetupAndVerificationService/v40/payments
	MakePaymentEndpoint = "https://checkout-test.adyen.com/v40/payments"
)

type PaymentResultPayload struct {
	Payload string `form:"payload" json:"payload"`
}

type RecurringPaymentMethod struct {
	RecurringDetailReference string `json:"recurringDetailReference"`
}

type PaymentAmount struct {
	Currency string `json:"currency"` // "EUR"
	Value    int    `json:"value"`    // 10 -- 小数点は？
}

type RecurringPaymentRequest struct {
	MerchantAccount    string                 `json:"merchantAccount"` // https://ca-test.adyen.com/ca/ca/accounts/show.shtml?accountTypeCode=MerchantAccount
	Amount             PaymentAmount          `json:"amount"`
	PaymentMethod      RecurringPaymentMethod `json:"paymentMethod"`
	Reference          string                 `json:"reference"`          // ペイメントID,
	ReturnUrl          string                 `json:"returnUrl"`          // http://localhost:8080/return
	ShopperReference   string                 `json:"shopperReference"`   // ユーザIDなど,
	ShopperInteraction string                 `json:"shopperInteraction"` // ContAuth
}

func NewRecurringPaymentRequest(account, recurRef, ref, userRef string, amount PaymentAmount) *RecurringPaymentRequest {
	return &RecurringPaymentRequest{
		MerchantAccount:    account,
		Amount:             amount,
		PaymentMethod:      RecurringPaymentMethod{recurRef},
		Reference:          ref,
		ShopperReference:   userRef,
		ShopperInteraction: "ContAuth", // ユーザなしの決済（フォームなどを動作していなくて、サーバだけで実行する）
		ReturnUrl:          "http://localhost:8080/return",
	}
}

type PaymentSessionRequest struct {
	MerchantAccount  string        `json:"merchantAccount"` // https://ca-test.adyen.com/ca/ca/accounts/show.shtml?accountTypeCode=MerchantAccount
	SdkVersion       string        `json:"sdkVersion"`      // 1.6.4
	Channel          string        `json:"channel"`         // Web
	Amount           PaymentAmount `json:"amount"`
	Reference        string        `json:"reference"`        // ペイメントID,
	CountryCode      string        `json:"countryCode"`      // "NL",
	ShopperLocale    string        `json:"shopperLocale"`    // "nl_NL",
	ShopperReference string        `json:"shopperReference"` // ユーザIDなど,
	Origin           string        `json:"origin"`           // http://localhost:8080,
	ReturnUrl        string        `json:"returnUrl"`        // http://localhost:8080/completed
	// Html             bool          `json:"html"`
	EnableRecurring bool `json:"enableRecurring"`
	EnableOneClick  bool `json:"enableOneClick"`
}

func NewPaymentSessionRequest(merchantAccount string) *PaymentSessionRequest {
	return &PaymentSessionRequest{
		merchantAccount,
		"1.6.4",
		"Web",
		PaymentAmount{"USD", 1000},
		fmt.Sprintf("randomId123354asdfasdf%d", time.Now().Unix()),
		"US",
		"en_US",
		fmt.Sprintf("1234565asdfsadf789%d", time.Now().Unix()),
		"http://localhost:8080",
		"http://localhost:8080/result",
		// true,
		true,
		true,
	}
}
