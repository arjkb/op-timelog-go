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
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
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

	filename := flag.String("file", "status_"+time.Now().Format("20060102")+".dailystatus", "file to read from") // "samplefile.txt" // TODO: read filename as a command-line arg
	flag.Parse()
	fmt.Println(*filename)

	// TODO: read config data (user key, base-url etc) from a file

	file, err := os.Open(*filename)
	if err != nil {
		log.Fatalln(err)
	}

	input := bufio.NewScanner(file)
	for input.Scan() {
		linecount++
		line := input.Text()
		s := strings.SplitN(line, " ", 3)

		fmt.Println(s[0])
		fmt.Println(s[1])
		fmt.Println(s[2])
		go func(word string) {
			var gr GoroutineResponse
			gr.resp, gr.err = makeRequest(word)
			ch <- gr
		}(line)
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

// Parsed the input line and extracts the work-package code, duration,
// description, and error if any
func extractData(s string) (int, string, string, error) {
	// work-in-progress

	var wp int
	var dur string
	var desc string

	split := strings.SplitN(s, " ", 3)
	if len(split) != 3 {
		return wp, dur, desc, fmt.Errorf("cannot split %q into 3 parts", s)
	}

	wp, err := strconv.Atoi(split[0])
	if err != nil {
		return wp, dur, desc, fmt.Errorf("converting %q to int: %v", split[0], err)
	}

	dur, desc = split[1], split[2]

	return wp, dur, desc, nil
}
