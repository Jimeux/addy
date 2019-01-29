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

	router.Run()
}

func handlePaymentResultRedirect(c *gin.Context) {
	type Query struct {
		Payload    string `form:"payload" json:"payload"`
		Type       string `form:"type" json:"type"`
		ResultCode string `form:"resultCode" json:"resultCode"`
	}

	var query Query
	if err := c.BindQuery(&query); err != nil {
		log.Fatal(err)
	}

	escaped, err := url.PathUnescape(query.Payload)
	if err != nil {
		log.Fatal("Invalid payload", err)
	}

	payload := adyen.PaymentResultPayload{Payload: escaped}

	fmt.Println(payload)

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

func handlePaymentSessionRequest(c *gin.Context) {
	payment := adyen.NewPaymentSessionRequest(merchantAccount)
	body, err := json.Marshal(payment)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPost, adyen.PaymentSessionEndpoint, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", apiKey)

	r, err := httputil.DumpRequest(req, true)
	fmt.Println()
	fmt.Println(string(r))
	fmt.Println()

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

	// fmt.Println("sessionResponse: ", resBody)

	c.HTML(http.StatusOK, "index.html", gin.H{
		"session": resBody["paymentSession"],
	})
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

	r, err := httputil.DumpResponse(resp, true)
	fmt.Println()
	fmt.Println(string(r))
	fmt.Println()

	decoder := json.NewDecoder(resp.Body)
	resBody := make(map[string]interface{})
	err = decoder.Decode(&resBody)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode >= 400 {
		log.Fatal(resBody)
	}
	return nil
}
