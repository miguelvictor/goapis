package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
)

func PdfToText(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// get the file from the request
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Printf("Received %s (%d bytes)", handler.Filename, handler.Size)
	defer file.Close()

	// create a temporary file to store the uploaded pdf file
	src, err := os.CreateTemp("", "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer src.Close()

	// copy the contents of the file to the new file
	_, err = io.Copy(src, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// create a temporary file to store the converted text file
	dst, err := os.CreateTemp("", "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// convert pdf to text using the pdftotext binary
	cmd := exec.Command("pdftotext", src.Name(), dst.Name())
	if err := cmd.Run(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// determine the size of the converted text file
	stat, err := dst.Stat()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// send file as the response with the appropriate headers
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Disposition", "attachment; filename=output.txt")
	w.Header().Set("x-src-name", src.Name())
	w.Header().Set("x-dst-name", dst.Name())
	w.Header().Set("Content-Length", strconv.FormatInt(stat.Size(), 10))
	io.Copy(w, dst)
}
