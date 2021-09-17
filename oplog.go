package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	// TODO: change the way this URL is determined
	url := "http://127.0.0.1:5000/"

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	fmt.Printf("%s", body)
}
