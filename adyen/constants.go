package adyen

const (
	webSDKVersion      = "1.6.5"
	channelWeb         = "Web"
	recurringReturnURL = "http://localhost:8080/result"
	countryCode        = "KR"
	shopperLocale      = "ko_KR"
	shopperInteraction = "ContAuth" // ユーザがデバイスを動作していない間に取られる決済
)

// 決済結果コード
const (
	ResultCodeAuthorised      = "Authorised"      // Indicates the payment authorisation was successfully completed. This state serves as an indicator to proceed with the delivery of goods and services. This is a final state.
	ResultCodeRefused         = "Refused"         // Indicates the payment was refused. The reason is given in the refusalReason field. This is a final state.
	ResultCodeRedirectShopper = "RedirectShopper" // Indicates the shopper should be redirected to an external web page or app to complete the authorisation.
	ResultCodeReceived        = "Received"        // Indicates the payment has successfully been received by Adyen, and will be processed. This is the initial state for all payments.
	ResultCodeCancelled       = "Cancelled"       // Indicates the payment has been cancelled (either by the shopper or the merchant) before processing was completed. This is a final state.
	ResultCodePending         = "Pending"         // Indicates that it is not possible to obtain the final status of the payment. This can happen if the systems providing final status information for the payment are unavailable, or if the shopper needs to take further action to complete the payment. For more information on handling a pending payment, refer to Payments with pending status.
	ResultCodeError           = "Error"           // Indicates an error occurred during processing of the payment. The reason is given in the refusalReason field. This is a final state.
)

const (
	apiURL = "https://checkout-test.adyen.com/v40"
	// https://docs.adyen.com/api-explorer/#/PaymentSetupAndVerificationService/v40/paymentSession
	PaymentSessionEndpoint = apiURL + "/paymentSession"
	// https://docs.adyen.com/api-explorer/#/PaymentSetupAndVerificationService/v40/payments/result
	VerifyPaymentEndpoint = apiURL + "/payments/result"
	// https://docs.adyen.com/api-explorer/#/PaymentSetupAndVerificationService/v40/payments
	MakePaymentEndpoint = "https://checkout-test.adyen.com/v40/payments"
)
