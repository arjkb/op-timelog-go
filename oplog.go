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
	// TODO: change the way this URL is determined
	url := "http://127.0.0.1:5000/reverse"

	data, err := json.Marshal(Oprequest{
		Word: "sandwich",
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
	fmt.Printf("%s", body)
}
