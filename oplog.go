package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// this represents a dummy data to post
type Oprequest struct {
	Word string `json:"word"`
}

func main() {

	// TODO: read input from a file and parse it

	// TODO: read config data (user key, base-url etc) from a file

	// TODO: make things run in parallel using goroutines
	// (because this is an embarassingly parallel workload)

	words := []string{
		"apple",
		"banana",
		"grape",
		"mausambi",
		"orange",
		"passion fruit",
		"pineapple",
		"potatoe",
		"strawberry",
		"watermelon",
	}

	for _, word := range words {
		res, statusCode := makeRequest(word)
		fmt.Printf("%s %d\n", res, statusCode)
	}
}

// make request to API and return the response body and status code
func makeRequest(word string) (string, int) {
	// TODO: change the way this URL is determined
	url := "http://127.0.0.1:5000/reverse"

	data, err := json.Marshal(Oprequest{
		Word: word,
	})
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("error reading response", err)
	}
	return string(body), resp.StatusCode
}
