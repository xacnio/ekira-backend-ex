package utils

import (
	"errors"
	"fmt"
	"github.com/dongri/phonenumber"
	req2 "github.com/imroc/req/v3"
	"os"
	"strings"
)

type VerifyKitStartReference struct {
	Reference string `json:"reference"`
	Timeout   int    `json:"timeout"`
}

func StartVerifyKitOTP(ip string, phone string) (*VerifyKitStartReference, error) {
	type VerifyKitResponse struct {
		Meta struct {
			RequestId      string `json:"requestId"`
			HttpStatusCode int    `json:"httpStatusCode"`
		} `json:"meta"`
		Result struct {
			Reference string `json:"reference"`
			Timeout   int    `json:"timeout"`
		} `json:"result"`
	}
	type VerifyKitRequest struct {
		PhoneNumber string `json:"phoneNumber"`
		CountryCode string `json:"countryCode"`
	}
	vfReq := VerifyKitRequest{}

	phoneCountry := phonenumber.GetISO3166ByNumber(phone, true)
	if phoneCountry.Alpha2 == "" {
		return nil, errors.New("invalid phone number")
	}
	numberNoCountry := strings.Replace(phone, phoneCountry.CountryCode, "", 1)
	vfReq.PhoneNumber = numberNoCountry
	vfReq.CountryCode = phoneCountry.Alpha2

	apiKey := os.Getenv("VERIFY_KIT_WEB_KEY")
	client := req2.C()
	resp, err := client.R().SetHeaders(map[string]string{
		"Content-Type":        "application/json",
		"X-Vfk-Server-Key":    apiKey,
		"X-Vfk-Forwarded-For": ip,
	}).SetBody(vfReq).Post("https://web-rest.verifykit.com/v1.0/send-otp")
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("verify kit request error")
	}
	if resp.StatusCode != 201 {
		fmt.Println(resp.StatusCode)
		fmt.Println(resp.String())
		return nil, errors.New("verify kit otp error")
	}
	vfRes := VerifyKitResponse{}
	e := resp.Unmarshal(&vfRes)
	if e != nil || vfRes.Result.Reference == "" {
		fmt.Println(e)
		return nil, errors.New("verify kit error")
	}
	vfReference := VerifyKitStartReference{
		Reference: vfRes.Result.Reference,
		Timeout:   vfRes.Result.Timeout,
	}
	return &vfReference, nil
}

func CheckVerifyKitOTP(ip, phone, reference, code string) error {
	type VerifyKitResponse struct {
		Meta struct {
			RequestId      string `json:"requestId"`
			HttpStatusCode int    `json:"httpStatusCode"`
		} `json:"meta"`
		Result struct {
			ValidationStatus string `json:"validationStatus"`
			SessionID        string `json:"sessionId"`
		} `json:"result"`
	}
	type VerifyKitRequest struct {
		Reference   string `json:"reference"`
		Code        string `json:"code"`
		PhoneNumber string `json:"phoneNumber"`
		CountryCode string `json:"countryCode"`
	}
	vfReq := VerifyKitRequest{}
	vfReq.Reference = reference
	vfReq.Code = code

	phoneCountry := phonenumber.GetISO3166ByNumber(phone, true)
	if phoneCountry.Alpha2 == "" {
		return errors.New("invalid phone number")
	}
	numberNoCountry := strings.Replace(phone, phoneCountry.CountryCode, "", 1)
	vfReq.PhoneNumber = numberNoCountry
	vfReq.CountryCode = phoneCountry.Alpha2

	apiKey := os.Getenv("VERIFY_KIT_WEB_KEY")
	client := req2.C()
	resp, err := client.R().SetHeaders(map[string]string{
		"Content-Type":        "application/json",
		"X-Vfk-Server-Key":    apiKey,
		"X-Vfk-Forwarded-For": ip,
	}).SetBody(vfReq).Post("https://web-rest.verifykit.com/v1.0/check-otp")
	if err != nil {
		fmt.Println(err)
		return errors.New("verify kit request error")
	}
	if resp.StatusCode != 200 {
		fmt.Println(resp.StatusCode)
		fmt.Println(resp.String())
		return errors.New("verify kit otp error")
	}
	return nil
}
