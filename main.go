package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
)

func handlePdfToText(w http.ResponseWriter, r *http.Request) {
	// manual 404 setting
	if r.URL.Path != "/api/pdftotext" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// send cors headers
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Max-Age", "86400")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// check auth token
	query := r.URL.Query()
	token := query.Get("token")
	if token != os.Getenv("SECRET") {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// handle request
	fmt.Println("got /api/pdftotext request")
	PdfToText(w, r)
}

func main() {
	http.HandleFunc("/api/pdftotext", handlePdfToText)

	// get port from environment variable
	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "3333"
	}

	// start http server
	fmt.Printf("server started at localhost:%s\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("server closed")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
