// Oplog -- Push timelog updates to OpenProject
// Copyright (C) 2021  Arjun Krishna Babu

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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
	Reversed string
}

type Respp struct {
	Original   string
	Answer     string
	StatusCode int
}

type GoroutineResponse struct {
	resp Respp
	err  error
}

func main() {
	var linecount int
	ch := make(chan GoroutineResponse)

	filename := "samplefile.txt" // TODO: read filename as a command-line arg

	// TODO: read config data (user key, base-url etc) from a file

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
	}

	input := bufio.NewScanner(file)
	for input.Scan() {
		linecount++
		go func(word string) {
			var gr GoroutineResponse
			gr.resp, gr.err = makeRequest(word)
			ch <- gr
		}(input.Text())
	}

	for i := 0; i < linecount; i++ {
		gr := <-ch
		if gr.err != nil {
			log.Println(err)
			continue
		}
		fmt.Println(gr.resp)
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
