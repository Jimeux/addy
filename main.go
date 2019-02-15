package main

import (
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

	client := adyen.NewClient(apiKey, merchantAccount)
	controller := Controller{client}

	router.GET("/test", controller.handlePaymentSessionRequest)
	router.POST("/result", controller.handleVerifyPaymentResult)
	router.GET("/result", controller.handlePaymentResultRedirect)
	router.GET("/payment", controller.handleRecurringPayment)
	router.POST("/notify", controller.handleNotification)

	router.Run()
}

type Controller struct {
	client *adyen.Client
}

func (ct *Controller) handleNotification(c *gin.Context) {
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

var lastVerifyResponse *adyen.VerifyPaymentResponse

func (ct *Controller) handleRecurringPayment(c *gin.Context) {
	resp, err := ct.client.MakeRecurringPayment(
		adyen.PaymentAmount{
			Currency: "KRW",
			Value:    100,
		},
		lastVerifyResponse.AdditionalData.RecurringDetailReference,
		lastVerifyResponse.AdditionalData.ShopperReference,
		lastVerifyResponse.MerchantReference,
	)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("makePaymentResponse: ", resp)
	c.JSON(http.StatusOK, resp)
}

func (ct *Controller) handlePaymentSessionRequest(c *gin.Context) {
	amount := adyen.PaymentAmount{"USD", 1000}
	ref := fmt.Sprintf("randomId123354asdfasdf%d", time.Now().Unix())
	userRef := fmt.Sprintf("1234565asdfsadf789%d", time.Now().Unix())

	resp, err := ct.client.CreatePaymentSession(amount, ref, userRef, "http://localhost:8080", "http://localhost:8080/result")
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("sessionResponse: ", resp)

	c.HTML(http.StatusOK, "index.html", gin.H{
		"session": resp.PaymentSession,
	})
}

func (ct *Controller) handlePaymentResultRedirect(c *gin.Context) {
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

	if err := ct.verifyReq(payload); err != nil {
		c.String(400, "Error")
	}
	c.String(200, "Success")
}

func (ct *Controller) handleVerifyPaymentResult(c *gin.Context) {
	var payload adyen.PaymentResultPayload
	if err := c.Bind(&payload); err != nil {
		log.Fatal(err)
	}
	if err := ct.verifyReq(payload); err != nil {
		c.String(400, "Error")
	}
	c.String(200, "Success")
}

func (ct *Controller) verifyReq(payload adyen.PaymentResultPayload) error {
	resp, err := ct.client.VerifyPayment(payload.Payload)
	if err != nil {
		log.Println(err)
		return err
	}

	fmt.Println("verificationResponse: ", resp)
	lastVerifyResponse = resp // TODO DBなんかに入れる

	return nil
}

func debugRequest(req *http.Request) {
	r, _ := httputil.DumpRequest(req, true)
	fmt.Println()
	fmt.Println(string(r))
	fmt.Println()
}
