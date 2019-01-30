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

	router.Run()
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

	req, err := http.NewRequest(http.MethodPost, adyen.MakePaymentEndpoint, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", apiKey)

	debugRequest(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
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
	payment := adyen.NewPaymentSessionRequest(merchantAccount)
	body, err := json.Marshal(payment)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPost, adyen.PaymentSessionEndpoint, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", apiKey)

	debugRequest(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
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
	var query adyen.PayloadQueryString
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

	req, err := http.NewRequest(http.MethodPost, adyen.PaymentVerificationEndpoint, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	debugRequest(req)

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

func debugRequest(req *http.Request) {
	r, _ := httputil.DumpRequest(req, true)
	fmt.Println()
	fmt.Println(string(r))
	fmt.Println()
}
