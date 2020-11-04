package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Voucher ...
type Voucher struct {
	Code string
}

// Vouchers ...
type Vouchers struct {
	Voucher []Voucher
}

// Check ...
func (c Vouchers) Check(code string) string {
	for _, item := range c.Voucher {
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

var vouchers Vouchers

func main() {
	voucher := Voucher{
		Code: "123",
	}

	vouchers.Voucher = append(vouchers.Voucher, voucher)

	http.HandleFunc("/", home)
	http.ListenAndServe(":9093", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	voucher := r.PostFormValue("voucher")
	valid := vouchers.Check(voucher)

	result := Result{Status: valid}

	jsonResult, err := json.Marshal(result)
	if err != nil {
		log.Fatal("Error converting json")
	}

	fmt.Fprintf(w, string(jsonResult))
}
