package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
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
	src, err := os.CreateTemp("", "*"+filepath.Ext(handler.Filename))
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

	// convert pdf file
	if strings.HasSuffix(handler.Filename, ".pdf") {
		output, err := _PdfToText(src)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer output.Close()

		// determine the size of the converted text file
		stat, err := output.Stat()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// send file as the response with the appropriate headers
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Disposition", "attachment; filename=output.txt")
		w.Header().Set("Content-Length", strconv.FormatInt(stat.Size(), 10))
		io.Copy(w, output)
		return
	}

	// convert doc/docx file
	if strings.HasSuffix(handler.Filename, ".doc") || strings.HasSuffix(handler.Filename, ".docx") {
		output, err := _DocToText(src)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// delete the file in the current directory
		defer os.Remove(output.Name())

		// determine the size of the converted text file
		stat, err := output.Stat()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// send file as the response with the appropriate headers
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Disposition", "attachment; filename=output.txt")
		w.Header().Set("Content-Length", strconv.FormatInt(stat.Size(), 10))
		io.Copy(w, output)
		return
	}

	// unknown document type
	http.Error(w, "Unknown document extension", http.StatusBadRequest)
}

func _PdfToText(src *os.File) (*os.File, error) {
	// create a temporary file to store the converted text file
	dst, err := os.CreateTemp("", "*.txt")
	if err != nil {
		return nil, err

	}

	// convert pdf to text using the pdftotext binary
	cmd := exec.Command("pdftotext", src.Name(), dst.Name())
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	return dst, nil
}

func _DocToText(src *os.File) (*os.File, error) {
	// convert pdf to text using the pdftotext binary
	path := src.Name()
	cmd := exec.Command("libreoffice", "--headless", "--convert-to", "txt", path, "--outdir", "./")
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	// get the converted file
	nameWithExt := filepath.Base(path)
	ext := filepath.Ext(path)
	outputPath := strings.TrimSuffix(nameWithExt, ext) + ".txt"

	return os.Open(outputPath)
}
