package adyen

const (
	webSDKVersion      = "1.6.5"
	channelWeb         = "Web"
	origin             = "http://localhost:8080"
	initialReturnURL   = "http://localhost:8080/result"
	recurringReturnURL = "http://localhost:8080/result"
	countryCode        = "KR"
	shopperLocale      = "ko_KR"
	shopperInteraction = "ContAuth" // ユーザがデバイスを動作していない間に取られる決済
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
