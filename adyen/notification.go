package adyen

// Notificationを正常に受信したらHTTP200で返すレスポンス
var NotificationResponse = `"{notificationResponse": "[accepted]"}`

// https://docs.adyen.com/developers/api-reference/notifications-api
type NotificationRequest struct {
	Live              string `json:"live"`
	NotificationItems []struct {
		NotificationItem `json:"NotificationRequestItem"`
	} `json:"notificationItems"`
}

// 何故か配列の要素（NotificationItem）がJSONオブジェクトに包まれている
/*type NotificationWrapper struct {
	Item NotificationItem `json:"NotificationRequestItem"`
}*/

type NotificationItem struct {
	Data                AdditionalNotificationData `json:"additionalData"`
	Amount              PaymentAmount              `json:"amount"`
	EventCode           string                     `json:"eventCode"`
	EventDate           string                     `json:"eventDate"`
	MerchantAccountCode string                     `json:"merchantAccountCode"`
	MerchantReference   string                     `json:"merchantReference"`
	Operations          []string                   `json:"operations"`
	PaymentMethod       string                     `json:"paymentMethod"`
	PSPReference        string                     `json:"pspReference"`
	Reason              string                     `json:"reason"`
	Success             string                     `json:"success"`
}

type AdditionalNotificationData struct {
	ExpiryDate               string `json:"expiryDate"`
	AuthCode                 string `json:"authCode"`
	RecurringDetailReference string `json:"recurring.recurringDetailReference"`
	CardSummary              string `json:"cardSummary"`
	ShopperReference         string `json:"recurring.shopperReference"`
}
