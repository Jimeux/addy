package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/Jimeux/addy/adyen"
	"github.com/gin-gonic/gin"
)

var (
	apiKey          = os.Getenv("ADYEN_API_KEY")
	merchantAccount = os.Getenv("ADYEN_MERCHANT_ACCOUNT")
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/test", handlePaymentSessionRequest)
	router.POST("/result", handleVerifyPaymentResult)
	router.GET("/result", handlePaymentResultRedirect)
	router.GET("/payment", handleRecurringPayment)
	router.POST("/notify", handleNotification)

	router.Run()
}

func handleNotification(c *gin.Context) {
	debugRequest(c.Request)

	var notifications adyen.NotificationRequest
	if err := c.Bind(&notifications); err != nil {
		log.Fatal(err)
	}
	fmt.Println(notifications)

	for _, n := range notifications.NotificationItems {
		fmt.Printf("notification for PSPref %s", n.PSPReference)
	}

	c.Header("Content-Type", "application/json")
	c.Header("Accept", "application/json")
	c.String(http.StatusOK, adyen.NotificationResponse)
}

var lastVerifyResponse adyen.VerifyPaymentResponse

func handleRecurringPayment(c *gin.Context) {
	payload := adyen.NewRecurringPaymentRequest(
		merchantAccount,
		lastVerifyResponse.AdditionalData.RecurringDetailReference,
		lastVerifyResponse.MerchantReference,
		lastVerifyResponse.AdditionalData.ShopperReference,
		adyen.PaymentAmount{
			Currency: "KRW",
			Value:    100,
		},
	)

	body, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}

	resp := doPostRequest(adyen.MakePaymentEndpoint, body)
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	resBody := make(map[string]interface{})
	err = decoder.Decode(&resBody)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode >= 400 {
		log.Fatal(resBody)
	}

	fmt.Println("makePaymentResponse: ", resBody)
	c.JSON(http.StatusOK, resBody)
}

func handlePaymentSessionRequest(c *gin.Context) {
	amount := adyen.PaymentAmount{"USD", 1000}
	ref := fmt.Sprintf("randomId123354asdfasdf%d", time.Now().Unix())
	userRef := fmt.Sprintf("1234565asdfsadf789%d", time.Now().Unix())

	payment := adyen.NewPaymentSessionRequest(amount, ref, userRef, merchantAccount)
	body, err := json.Marshal(payment)
	if err != nil {
		log.Fatal(err)
	}

	resp := doPostRequest(adyen.PaymentSessionEndpoint, body)
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	resBody := make(map[string]interface{})
	err = decoder.Decode(&resBody)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode >= 400 {
		log.Fatal(resBody)
	}

	fmt.Println("sessionResponse: ", resBody)

	c.HTML(http.StatusOK, "index.html", gin.H{
		"session": resBody["paymentSession"],
	})
}

func handlePaymentResultRedirect(c *gin.Context) {
	var query adyen.PaymentResultParams
	if err := c.BindQuery(&query); err != nil {
		log.Fatal(err)
	}

	// HTMLエンコードされているので、デコードしておく
	escaped, err := url.PathUnescape(query.Payload)
	if err != nil {
		log.Fatal("Invalid payload", err)
	}

	payload := adyen.PaymentResultPayload{Payload: escaped}

	if err := verifyReq(payload); err != nil {
		c.String(400, "Error")
	}
	c.String(200, "Success")
}

func handleVerifyPaymentResult(c *gin.Context) {
	var payload adyen.PaymentResultPayload
	if err := c.Bind(&payload); err != nil {
		log.Fatal(err)
	}
	if err := verifyReq(payload); err != nil {
		c.String(400, "Error")
	}
	c.String(200, "Success")
}

func verifyReq(payload adyen.PaymentResultPayload) error {
	body, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}

	resp := doPostRequest(adyen.PaymentVerificationEndpoint, body)
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	var verifyResp adyen.VerifyPaymentResponse
	err = decoder.Decode(&verifyResp)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode >= 400 {
		log.Fatal(verifyResp)
	}

	fmt.Println("verificationResponse: ", verifyResp)
	lastVerifyResponse = verifyResp // TODO DBなんかに入れる

	return nil
}

func doPostRequest(endpoint string, body []byte) *http.Response {
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", apiKey)

	debugRequest(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	return resp
}

func debugRequest(req *http.Request) {
	r, _ := httputil.DumpRequest(req, true)
	fmt.Println()
	fmt.Println(string(r))
	fmt.Println()
}
