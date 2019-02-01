package adyen

// SDKによる決済がリダイレクトする際のクエリパラム
type PaymentResultParams struct {
	Payload    string `form:"payload" json:"payload"`
	Type       string `form:"type" json:"type"`
	ResultCode string `form:"resultCode" json:"resultCode"`
}

// SDKによる決済の結果
type PaymentResultPayload struct {
	Payload string `form:"payload" json:"payload"`
}

type PaymentAmount struct {
	Currency string `json:"currency"` // "USD"
	Value    int    `json:"value"`    // 100 -- セント単位（小数点はない）
}

// SDKのPayloadを検証（verify）するエンドポイントのレスポンス
// https://docs.adyen.com/api-explorer/#/PaymentSetupAndVerificationService/payments/result
type VerifyPaymentResponse struct {
	AdditionalData    AdditionalData `json:"additionalData"`
	PSPReference      string         `json:"pspReference"`
	ResultCode        string         `json:"resultCode"`
	MerchantReference string         `json:"merchantReference"`
	PaymentMethod     string         `json:"paymentMethod"`
	ShopperLocale     string         `json:"shopperLocale"`
}

type AdditionalData struct {
	RecurringDetailReference string `json:"recurring.recurringDetailReference"`
	ShopperReference         string `json:"recurring.shopperReference"`
}

// 保存済みのカードから決済を取るリクエスト
// https://docs.adyen.com/developers/features/tokenization/making-payments-with-tokens
type RecurringPaymentRequest struct {
	MerchantAccount    string                 `json:"merchantAccount"` // https://ca-test.adyen.com/ca/ca/accounts/show.shtml?accountTypeCode=MerchantAccount
	Amount             PaymentAmount          `json:"amount"`
	PaymentMethod      RecurringPaymentMethod `json:"paymentMethod"`
	Reference          string                 `json:"reference"`          // 決済を特定する値（ID）
	ReturnURL          string                 `json:"returnUrl"`          // リダイレクト時にアクセスされるURL
	ShopperReference   string                 `json:"shopperReference"`   // ユーザIDなど
	ShopperInteraction string                 `json:"shopperInteraction"` // ContAuth（サブ系の決済）Ecommerce（客による）
}

type RecurringPaymentMethod struct {
	RecurringDetailReference string `json:"recurringDetailReference"`
}

func NewRecurringPaymentRequest(account, recurRef, ref, userRef string, amount PaymentAmount) *RecurringPaymentRequest {
	return &RecurringPaymentRequest{
		MerchantAccount:    account,
		Amount:             amount,
		PaymentMethod:      RecurringPaymentMethod{recurRef},
		Reference:          ref,
		ShopperReference:   userRef,
		ShopperInteraction: shopperInteraction,
		ReturnURL:          recurringReturnURL,
	}
}
