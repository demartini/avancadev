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

// Coupon ...
type Coupon struct {
	Code string
}

// Coupons ...
type Coupons struct {
	Coupon []Coupon
}

// Check ...
func (c Coupons) Check(code string) string {
	for _, item := range c.Coupon {
		if code == item.Code {
			return "valid"
		}
	}
	return "invalid"
}

// Result ...
type Result struct {
	Status  string
	Message string
}

var coupons Coupons

func main() {
	coupon := Coupon{
		Code: "abc",
	}

	coupons.Coupon = append(coupons.Coupon, coupon)

	http.HandleFunc("/", home)
	http.ListenAndServe(":9092", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	coupon := r.PostFormValue("coupon")
	voucher := r.PostFormValue("voucher")
	valid := coupons.Check(coupon)

	resultVoucher := makeHTTPCall("http://localhost:9093", voucher)

	result := Result{Status: valid}

	if resultVoucher.Status == "invalid" && result.Status == "invalid" {
		result.Status = "invalid"
		result.Message = "voucher and coupon invalid"
	}

	if resultVoucher.Status == "invalid" && result.Status == "valid" {
		result.Status = "valid"
		result.Message = "voucher invalid"
	}

	if resultVoucher.Status == "valid" && result.Status == "invalid" {
		result.Status = "invalid"
		result.Message = "coupon invalid"
	}

	if resultVoucher.Status == "valid" && result.Status == "valid" {
		result.Status = "valid"
		result.Message = "voucher and coupon valid"
	}

	jsonResult, err := json.Marshal(result)
	if err != nil {
		log.Fatal("Error converting json")
	}

	fmt.Fprintf(w, string(jsonResult))
}

func makeHTTPCall(urlMicroservice string, voucher string) Result {
	values := url.Values{}
	values.Add("voucher", voucher)

	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 5

	res, err := retryClient.PostForm(urlMicroservice, values)
	if err != nil {
		result := Result{Status: "Servidor fora do ar!"}
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
