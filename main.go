package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Response struct {
	Country     string `json:"country"`
	CountryCode string `json:"countryCode"`
	City        string `json:"city"`
	Query       string `json:"query"`
}

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("Not enough arguments passed to the program")
		help()
		os.Exit(1)
	}

	baseApiURL := "http://ip-api.com/batch"
	body := []map[string]string{
		{
			"query":  args[1],
			"fields": "city,country,countryCode,query",
		},
	}
	jsonData, err := json.Marshal(body)
	if err != nil {
		fmt.Println("Could not marshal the JSON")
		return
	}

	res, err := http.Post(baseApiURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Could not make a POST request to IPADDRESS: %s\n", args[1])
		return
	}
	defer res.Body.Close()

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

	if len(responseData) > 0 {
		fmt.Println("Country Code:", responseData[0].CountryCode)
		fmt.Println("Country:", responseData[0].Country)
		fmt.Println("City:", responseData[0].City)
		fmt.Println("Query:", responseData[0].Query)
	} else {
		fmt.Println("No data returned")
	}
}

func help() {
	fmt.Println("Usage: iptrace <IPADDRESS>")
}
