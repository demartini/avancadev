package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/hashicorp/go-retryablehttp"
)

// Result ...
type Result struct {
	Status  string
	Message string
}

func main() {
	http.HandleFunc("/", home)
	http.ListenAndServe(":9091", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	coupon := r.PostFormValue("coupon")
	ccNumber := r.PostFormValue("ccNumber")
	voucher := r.PostFormValue("voucher")

	resultCoupon := makeHTTPCall("http://localhost:9092", coupon, voucher)

	result := Result{Status: "declined"}

	if ccNumber == "1" {
		result.Status = "valid"
		result.Message = "approvad"
	}

	if resultCoupon.Status == "invalid" {
		result = resultCoupon
	}

	if resultCoupon.Status == "valid" {
		result.Message = resultCoupon.Message
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		log.Fatal("Error processing json")
	}

	fmt.Fprintf(w, string(jsonData))
}

func makeHTTPCall(urlMicroservice string, coupon string, voucher string) Result {
	values := url.Values{}
	values.Add("coupon", coupon)
	values.Add("voucher", voucher)

	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 5

	res, err := retryClient.PostForm(urlMicroservice, values)
	if err != nil {
		result := Result{Message: "Servidor fora do ar!"}
		return result
	}

	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Error processing result")
	}

	result := Result{}

	json.Unmarshal(data, &result)

	return result
}
