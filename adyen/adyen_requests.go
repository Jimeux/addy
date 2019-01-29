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
)

type PaymentAmount struct {
	Currency string `json:"currency"` // "EUR"
	Value    int    `json:"value"`    // 10 -- 小数点は？
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
		PaymentAmount{"JPY", 100},
		fmt.Sprintf("randomId123354asdfasdf%d", time.Now().Unix()),
		"JP",
		"ja_JP",
		fmt.Sprintf("1234565asdfsadf789%d", time.Now().Unix()),
		"http://localhost:8080",
		"http://localhost:8080/result",
		// true,
		true,
		true,
	}
}

type PaymentResultPayload struct {
	Payload string `form:"payload" json:"payload"`
}
