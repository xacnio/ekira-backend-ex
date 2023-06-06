package utils

import (
	"bytes"
	"encoding/xml"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"unicode"
)

// Quoted from: github.com/hakanersu/tcvalidate
func ValidateTcNumber(tcnumber string) bool {
	runes := []rune(tcnumber)
	if isSame(runes) {
		return false
	}
	if len(runes) != 11 {
		return false
	}

	odd, even, sum, rebuild := 0, 0, 0, ""

	for i := 0; i < len(runes)-2; i++ {

		a, _ := strconv.Atoi(string(runes[i]))

		if string(runes[0]) == "0" {
			return false
		}

		if (i+1)%2 == 0 {
			odd += a
		} else {
			even += a
		}

		rebuild += string(runes[i])

		sum += a
	}

	ten := (even*7 - odd) % 10

	indexTen, _ := strconv.Atoi(string(runes[9]))

	eleven := (sum + indexTen) % 10

	build := string(rebuild) + strconv.Itoa(ten) + strconv.Itoa(eleven)

	if build == tcnumber {
		return true
	}

	return false
}

func isSame(a []rune) bool {
	b := a[0:10]
	for i := 1; i < len(b); i++ {
		if b[i] != b[0] {
			return false
		}
	}
	return true
}

// Quoted from: github.com/barisesen/tcverify
type Response struct {
	TCKimlikNoDogrulaResult string `xml:"Body>TCKimlikNoDogrulaResponse>TCKimlikNoDogrulaResult"`
}

func CheckTcNumberNVI(ID string, name string, surname string, birthYear string) (bool, error) {
	rawXml := []byte(`<?xml version="1.0" encoding="utf-8"?>
		<soap:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
			<soap:Body>
				<TCKimlikNoDogrula xmlns="http://tckimlik.nvi.gov.tr/WS">
					<TCKimlikNo>` + ID + `</TCKimlikNo>
					<Ad>` + strings.ToUpperSpecial(unicode.TurkishCase, name) + `</Ad>
					<Soyad>` + strings.ToUpperSpecial(unicode.TurkishCase, surname) + `</Soyad>
					<DogumYili>` + birthYear + `</DogumYili>
				</TCKimlikNoDogrula>
			</soap:Body>
		</soap:Envelope>`)

	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://tckimlik.nvi.gov.tr/Service/KPSPublic.asmx", bytes.NewBuffer(rawXml))
	if err != nil {
		return false, err
	}
	req.Header.Add("Content-Type", "text/xml; charset=utf-8")
	req.Header.Add("SOAPAction", "http://tckimlik.nvi.gov.tr/WS/TCKimlikNoDogrula")
	req.Header.Add("Host", "tckimlik.nvi.gov.tr")
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	response := Response{}
	if err := xml.NewDecoder(resp.Body).Decode(&response); err != nil {
		return false, err
	}

	status, err := strconv.ParseBool(response.TCKimlikNoDogrulaResult)
	if err != nil {
		return false, err
	}

	if !status {
		error := errors.New("Bu bilgileri ait vatandaşlık doğrulanamadı.")
		return status, error
	}
	return status, nil
}

// IsValidTcInfo checks if the given tcNo, name, surname and year are valid.
func IsValidTcInfo(tcNo string, name string, surname string, year string) (bool, error) {
	if len(tcNo) != 11 {
		return false, nil
	}

	if ValidateTcNumber(tcNo) == false {
		return false, nil
	}

	resp, err := CheckTcNumberNVI(tcNo, name, surname, year)
	return resp, err
}
