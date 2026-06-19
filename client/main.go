package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"
)

type DollarQuotation struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

type ApiResponse struct {
	USDBRL DollarQuotation `json:"USDBRL"`
}

func main() {
	url := "http://localhost:8080/cotacao"
	ctx, cancel := context.WithTimeout(context.Background(), 300 * time.Millisecond)

	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		log.Printf("Error creating request: %v\n", err)
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error fetching cotacao: %v\n", err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error reading response body: %v\n", err)
		return
	}

	var response ApiResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("Error parsing response: %v\n", err)
		return
	}

	if err := WriteFileWithBid(response); err != nil {
		log.Printf("Error writing file: %v\n", err)
	}
}

func WriteFileWithBid(apiResponse ApiResponse) error{
	file, err := os.OpenFile("cotacao.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	const templateText = "Dólar: {{.USDBRL.Bid }}\n"

	tmpl, err := template.New("cotacao").Parse(templateText)
	if err != nil {
		return err
	}

	tmpl.Execute(file, apiResponse)

	defer file.Close()
	return nil
}