package main

import (
	"bytes"
	"encoding/json"
	"github.com/Jimeux/addy/adyen"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

var (
	apiKey          = os.Getenv("ADYEN_API_KEY")
	merchantAccount = os.Getenv("ADYEN_MERCHANT_ACCOUNT")
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/", handlePaymentSessionRequest)

	router.Run()
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

	c.HTML(http.StatusOK, "index.html", gin.H{
		"session": resBody["paymentSession"],
	})
}
