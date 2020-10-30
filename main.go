package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type RequestInfo struct {
	Method     string
	URL        string
	Payload    map[string]string
	httpClient *http.Client
}

func main() {
	r := RequestInfo{}
	r.Config()

	reqBody, err := json.Marshal(r.Payload)
	if err != nil {
		print(err)
	}

	req, err := http.NewRequest(r.Method, r.URL, bytes.NewBuffer(reqBody))
	if err != nil {
		print(err)
	}

	if os.Getenv("NEW_RELIC_API_KEY") != "" {
		req.Header.Set("API-Key", os.Getenv("NEW_RELIC_API_KEY"))
	} else {
		log.Fatal("ENV NEW_RELIC_API_KEY is not set")
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.httpClient.Do(req)
	if err != nil {
		log.Println(err)
	}

	f, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	resp.Body.Close()
	fmt.Println(string(f))
}

func (r *RequestInfo) Config() {
	var account_id string
	if os.Getenv("ACCOUNT_ID") != "" {
		account_id = os.Getenv("ACCOUNT_ID")
	} else {
		log.Fatal("ENV ACCOUNT_ID is not set")
	}
	mutation := `
	{
		actor {
		  account(id: ` + account_id + ` ) {
			nrql(query: "SELECT * FROM Transaction") {
			  staticChartUrl(chartType: TABLE, format: PDF)
			}
		  }
		}
	}`
	payload := map[string]string{
		"query":     mutation,
		"variables": "",
	}
	r.Method = "POST"
	r.URL = "https://api.newrelic.com/graphql"
	r.Payload = payload
	r.httpClient = &http.Client{
		Timeout: time.Second * 10,
	}
}
