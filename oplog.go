package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// this represents a dummy data to post
type Oprequest struct {
	Word string `json:"word"`
}

type Opresponse struct {
	Reversed string `json:"reversed"`
}

type Respp struct {
	Original   string
	Answer     string
	StatusCode int
}

func main() {

	var linecount int

	filename := "samplefile.txt" // TODO: read filename as a command-line arg

	// TODO: read config data (user key, base-url etc) from a file

	// TODO: make things run in parallel using goroutines
	// (because this is an embarassingly parallel workload)

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
	}

	input := bufio.NewScanner(file)
	for input.Scan() {
		linecount++
		res, err := makeRequest(input.Text())
		if err != nil {
			log.Println(err)
			continue
		}
		fmt.Println(res)
	}
}

// make request to API and return the response body and status code
func makeRequest(word string) (Respp, error) {
	var response Respp

	// TODO: change the way this URL is determined
	url := "http://127.0.0.1:5000/reverse"

	data, err := json.Marshal(Oprequest{
		Word: word,
	})
	if err != nil {
		return response, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return response, fmt.Errorf("got status %d", resp.StatusCode)
	}

	var opresponse Opresponse
	if err := json.NewDecoder(resp.Body).Decode(&opresponse); err != nil {
		fmt.Println(err)
		return response, err
	}

	response = Respp{
		Original:   word,
		Answer:     opresponse.Reversed,
		StatusCode: resp.StatusCode,
	}

	return response, nil
}
