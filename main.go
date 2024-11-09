package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/chachacollins/iptracer/spinner"
	"github.com/charmbracelet/huh"
)

type Response struct {
	Country     string `json:"country"`
	CountryCode string `json:"countryCode"`
	City        string `json:"city"`
	Query       string `json:"query"`
}

func main() {
	var inputBody string
	huh.NewInput().
		Title("Enter IP ADDRESS").
		Prompt(">").
		Value(&inputBody).Run()

	spinner.SpinnerClass("Proccessing input")
	baseApiURL := "http://ip-api.com/batch"
	body := []map[string]string{
		{
			"query":  inputBody,
			"fields": "city,country,countryCode,query",
		},
	}
	jsonData, err := json.Marshal(body)
	if err != nil {
		fmt.Println("Could not marshal the JSON")
		return
	}
	spinner.SpinnerClass("Marshalling json")
	res, err := http.Post(baseApiURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Could not make a POST request to IPADDRESS: %s\n", inputBody)
		return
	}
	defer res.Body.Close()
	spinner.SpinnerClass("fetching response")

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Could not read response body")
		return
	}

	var responseData []Response
	if err := json.Unmarshal(resBody, &responseData); err != nil {
		fmt.Println("Could not unmarshal the response body")
		return
	}
	spinner.SpinnerClass("Unmarshalling json")

	if len(responseData) > 0 {
		fmt.Println("Country Code:", responseData[0].CountryCode)
		fmt.Println("Country:", responseData[0].Country)
		fmt.Println("City:", responseData[0].City)
		fmt.Println("Query:", responseData[0].Query)
	} else {
		fmt.Println("No data returned")
	}
}
