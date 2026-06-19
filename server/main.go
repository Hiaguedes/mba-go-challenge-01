package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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

type DollarQuotationDB struct {
	gorm.Model
	DollarQuotation `gorm:"embedded"`
}


func main() {
	log.Println("Rodando servidor na porta :8080")

	http.HandleFunc("/", RootHandler)
	http.HandleFunc("/cotacao", CotationHandler)
	http.ListenAndServe(":8080", nil)
}

func RootHandler(writer http.ResponseWriter, req *http.Request){
	writer.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(writer, "Rota nao existente, rodar /cotacao")
	log.Println("Cliente tentou rodar rota nao existente /")
}

func CotationHandler(writer http.ResponseWriter, req *http.Request){
	url := "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	dsn := "app.db"
	log.Println("Request iniciada")
	originalCtx := req.Context()

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Error on gorm sqlite open: %v\n", err)
	}

	db.AutoMigrate(&DollarQuotationDB{})

	apiCtx, cancel := context.WithTimeout(originalCtx, 200*time.Millisecond)
	defer cancel()

	apiReq, err := http.NewRequestWithContext(apiCtx, http.MethodGet, url, nil)
	if err != nil {
		log.Printf("Error request with context: %v\n", err)
		http.Error(writer, "internal error", http.StatusInternalServerError)
		return
	}

	resp, err := http.DefaultClient.Do(apiReq)
	if err != nil {
		log.Printf("Error Default Client: %v\n", err)
		http.Error(writer, "error fetching cotacao", http.StatusGatewayTimeout)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error io read all: %v\n", err)
		http.Error(writer, "error reading response", http.StatusInternalServerError)
		return
	}

	defer log.Println("Request finalizada")

	var apiReponse ApiResponse
	err = json.Unmarshal(body, &apiReponse)
	if err != nil {
		log.Printf("Error JSON Unmarshal: %v\n", err)
		http.Error(writer, "Error Parsing Response: ", http.StatusInternalServerError)
		return
	}

	ctxDb, cancelDb := context.WithTimeout(originalCtx, 10*time.Millisecond)
	defer cancelDb()
	db.WithContext(ctxDb).Create(&DollarQuotationDB{DollarQuotation: apiReponse.USDBRL})

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(body)
}

