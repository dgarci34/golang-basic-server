package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func getApiKey() string {
	// TODO this key is only a demo, register and get a real api key
	return "https://eodhistoricaldata.com/api/eod/MCD.US?from=2023-05-05&to=2023-05-05&period=d&fmt=json&api_token=demo"
}

type StockDayInfo struct {
	Date           string
	Open           float64
	High           float64
	Low            float64
	Close          float64
	Adjusted_close float64
	Volume         int
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/hello" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "Hello!")
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	stock := r.FormValue("Stock")

	if len(stock) == 0 {
		log.Fatalln("Error: stock name not specified")
	}

	fmt.Fprintf(w, "Stock name = %s\n", stock)

	resp, err := http.Get(getApiKey())

	if err != nil {
		log.Fatalln(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var dayInfoArr []StockDayInfo
	json.Unmarshal(body, &dayInfoArr)
	dayInfo := dayInfoArr[0]

	fmt.Fprintf(w, "Date: %s\nOpen: %.2f\nHigh: %.2f\nLow: %.2f\nClose: %.2f\nAdjusted Close: %.2f\nVolume: %d\n",
		dayInfo.Date,
		dayInfo.Open,
		dayInfo.High,
		dayInfo.Low,
		dayInfo.Close,
		dayInfo.Adjusted_close,
		dayInfo.Volume)

	fmt.Println("Recieved form handle request")
}

func main() {
	fileServer := http.FileServer(http.Dir("./static"))

	http.Handle("/", fileServer)
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/form", formHandler)

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
