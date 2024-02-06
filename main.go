package main

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Result struct {
	StatusCode int    `json:"statusCode"`
	Time       int64  `json:"time"`
	Data       string `json:"data"`
}

func parseLine(line string) string {
	lines := strings.Split(line, " ")

	result := 0

	for _, num := range lines {
		intNum, err := strconv.Atoi(num)

		if err != nil {
			continue
		}

		result += intNum
	}

	return strconv.Itoa(result)
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20)

	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")

	if err != nil {
		http.Error(w, "Error while retreiving file", http.StatusInternalServerError)
		io.WriteString(w, "0")
		return
	}

	defer file.Close()

	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")

	scanner := bufio.NewScanner(file)

	var result bytes.Buffer

	for scanner.Scan() {
		result.WriteString(parseLine(scanner.Text()))
		result.WriteString(" ")
	}

	io.WriteString(w, result.String())
	return
}

func main() {
	http.HandleFunc("/upload", handleUpload)

	port, found := os.LookupEnv("PORT")

	if !found {
		port = "8080"
	}

	err := http.ListenAndServe(":"+port, nil)

	if err != nil {
		panic(err)
	}
}
