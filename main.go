package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
)

func handlePdfToText(w http.ResponseWriter, r *http.Request) {
	// add cors headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

	// manual 404 setting
	if r.URL.Path != "/api/pdftotext" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// check auth token
	query := r.URL.Query()
	token := query.Get("token")
	if token != os.Getenv("SECRET") {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// don't send anything if it's an OPTIONS request
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusNoContent)
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
