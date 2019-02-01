package adyen

// SDKの初期化に必要なデータを返却するリクエスト
// https://docs.adyen.com/api-explorer/#/PaymentSetupAndVerificationService/v40/paymentSession
type PaymentSessionRequest struct {
	MerchantAccount  string        `json:"merchantAccount"` // /ca/ca/accounts/show.shtml?accountTypeCode=MerchantAccount
	SDKVersion       string        `json:"sdkVersion"`      // 1.6.4
	Channel          string        `json:"channel"`         // プラットフォーム：Web / iOS / Android
	Amount           PaymentAmount `json:"amount"`
	Reference        string        `json:"reference"`        // ペイメントID,
	CountryCode      string        `json:"countryCode"`      // US
	ShopperLocale    string        `json:"shopperLocale"`    // en_US
	ShopperReference string        `json:"shopperReference"` // ユーザを特定するID
	Origin           string        `json:"origin"`           // http://localhost:8080,
	ReturnURL        string        `json:"returnUrl"`        // http://localhost:8080/completed
	EnableRecurring  bool          `json:"enableRecurring"`
	EnableOneClick   bool          `json:"enableOneClick"`
}

func NewPaymentSessionRequest(amount PaymentAmount, ref, shopperRef, merchantAccount string) *PaymentSessionRequest {
	return &PaymentSessionRequest{
		MerchantAccount:  merchantAccount,
		SDKVersion:       webSDKVersion,
		Channel:          channelWeb,
		Amount:           amount,
		Reference:        ref,
		CountryCode:      countryCode,
		ShopperLocale:    shopperLocale,
		ShopperReference: shopperRef,
		Origin:           origin,
		ReturnURL:        initialReturnURL,
		EnableRecurring:  true,
		EnableOneClick:   true,
	}
}
