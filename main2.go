package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

func test(c web.C, w http.ResponseWriter, r *http.Request) {
	type SessionRequest struct {
		UserID    int64  `json:"user_id"`
		Origin    string `json:"origin"`
		ReturnURL string `json:"return_url"`
		Amount    struct {
			Currency string `json:"currency"`
			Value    int    `json:"value"`
		}
	}

	decoder := json.NewDecoder(r.Body)
	var s SessionRequest
	if err := decoder.Decode(&s); err != nil {
		log.Fatal(err)
	}
	fmt.Println(s)

	fmt.Fprintf(w, "testing")
}

func main() {
	goji.Post("/test/:id", test)
	goji.Serve()
}
