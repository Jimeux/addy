package adyen

// redirectUrlに渡されるのGETパラム。
// Web SDKのフォームが外部ページへ遷移する時に必要。
type PayloadQueryString struct {
	Payload    string `form:"payload" json:"payload"`
	Type       string `form:"type" json:"type"`
	ResultCode string `form:"resultCode" json:"resultCode"`
}

type AdditionalData struct {
	RecurringDetailReference string `json:"recurring.recurringDetailReference"`
	ShopperReference         string `json:"recurring.shopperReference"`
}

type VerifyPaymentResponse struct {
	AdditionalData    AdditionalData `json:"additionalData"`
	PSPReference      string         `json:"pspReference"`
	ResultCode        string         `json:"resultCode"`
	MerchantReference string         `json:"merchantReference"`
	PaymentMethod     string         `json:"paymentMethod"`
	ShopperLocale     string         `json:"shopperLocale"`
}
