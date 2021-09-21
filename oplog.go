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
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Key string
	Url string
}

type Request struct {
	Links struct {
		Activity struct {
			Href string
		}
	}
	WorkPackage struct {
		Href string
	} `json:"workPackage"`
	Hours   string
	Comment struct {
		Raw string
	}
	SpentOn string
}

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

	f := flag.String("file", "status_"+time.Now().Format("20060102")+".dailystatus", "file to read from")
	flag.Parse()

	filename := *f
	datestr := filename[7:15] // extracts date string from filenames of format "status_20210921.dailystatus"
	fmt.Println(filename, datestr)

	// read config data (user key, base-url etc) from a file
	key, url, err := getKeyAndUrl("config.toml.example") // this is an example config file
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(key, url)

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
	}

	input := bufio.NewScanner(file)
	for input.Scan() {
		linecount++
		line := input.Text()
		wp, dur, desc, err := extractData(line)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(desc, wp, dur)

		request := makeRequestStruct(wp, dur, desc)
		fmt.Printf("%+v\n", request)

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

// make a request struct with appropriate params
func makeRequestStruct(wp int, dur string, desc string) Request {
	request := Request{}
	request.Links.Activity.Href = "url/to/href"
	request.WorkPackage.Href = "url/to/workpackage"
	request.Hours = "PT6.5H"
	request.Comment.Raw = "example comment"
	request.SpentOn = "20210905"

	return request
}

// Read config filename and return the API key and url
func getKeyAndUrl(configFileName string) (string, string, error) {
	file, err := os.Open(configFileName)
	if err != nil {
		return "", "", fmt.Errorf("cannot open file: %v", err)
	}

	b, err := io.ReadAll(bufio.NewReader(file))
	if err != nil {
		return "", "", fmt.Errorf("cannot read file: %v", err)
	}

	var conf Config
	if _, err := toml.Decode(string(b[:]), &conf); err != nil {
		return "", "", fmt.Errorf("cannot decode toml: %v", err)
	}

	fmt.Println(conf)

	return conf.Key, conf.Url, nil
}
