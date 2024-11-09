package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"io"
	"log"
	"net/http"

	"github.com/chachacollins/iptracer/spinner"
	"github.com/charmbracelet/huh"
	"os"
)

const (
	purple = lipgloss.Color("#ceaef3")
	pink   = lipgloss.Color("#FFC0CB")
)

type Response struct {
	Country     string `json:"country"`
	CountryCode string `json:"countryCode"`
	City        string `json:"city"`
	Query       string `json:"query"`
}

func main() {
	re := lipgloss.NewRenderer(os.Stdout)
	var (
		// HeaderStyle is the lipgloss style used for the table headers.
		HeaderStyle = re.NewStyle().Foreground(purple).Bold(true).Align(lipgloss.Center)
		// CellStyle is the base lipgloss style used for the table rows.
		CellStyle = re.NewStyle().Padding(0, 1).Width(14)
		// OddRowStyle is the lipgloss style used for odd-numbered table rows.
		OddRowStyle = CellStyle.Foreground(pink)
		// EvenRowStyle is the lipgloss style used for even-numbered table rows.
		EvenRowStyle = CellStyle.Foreground(pink)
		// BorderStyle is the lipgloss style used for the table border.
		BorderStyle = lipgloss.NewStyle().Foreground(purple)
	)
	var inputBody string
	huh.NewInput().
		Title("Enter IP ADDRESS").
		Prompt(">").
		Value(&inputBody).Run()

	spinner.SpinnerClass("Proccessing input")
	var reqForm []string
	form := huh.NewForm(
		huh.NewGroup(

			huh.NewMultiSelect[string]().
				Title("Choose response").
				Options(
					huh.NewOption("CountryCode", "countryCode").Selected(true),
					huh.NewOption("Country", "country").Selected(true),
					huh.NewOption("City", "city"),
					huh.NewOption("Query", "query"),
				).
				Value(&reqForm),
		),
	)
	err := form.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(reqForm)
	baseApiURL := "http://ip-api.com/batch"
	body := []map[string]interface{}{
		{
			"query":  inputBody,
			"fields": reqForm,
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
	t := table.New().
		Border(lipgloss.ThickBorder()).
		BorderStyle(BorderStyle).
		StyleFunc(func(row, col int) lipgloss.Style {
			var style lipgloss.Style

			switch {
			case row == table.HeaderRow:
				return HeaderStyle
			case row%2 == 0:
				style = EvenRowStyle
			default:
				style = OddRowStyle
			}

			// Make the second column a little wider.
			if col == 1 {
				style = style.Width(22)
			}

			return style
		}).
		Headers("Request", "Response").
		Row("Country Code", responseData[0].CountryCode)
	t.Row("Country", responseData[0].Country)
	t.Row("City", responseData[0].City)
	t.Row("Query", responseData[0].Query)
	fmt.Println(t)

}
