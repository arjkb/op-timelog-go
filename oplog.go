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
			Href string `json:"href"`
		} `json:"activity"`
	} `json:"_links"`
	WorkPackage struct {
		Href string `json:"href"`
	} `json:"workPackage"`
	Hours   string `json:"hours"`
	Comment struct {
		Raw string `json:"raw"`
	} `json:"comment"`
	SpentOn string `json:"spentOn"`
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

	// read config data
	key, url, err := getKeyAndUrl("config.toml")
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
		line := input.Text()
		wp, dur, desc, err := extractData(line)
		if err != nil {
			log.Println(err)
			continue
		}

		jsonMarshalled, err := makePostDataJSON(wp, dur, desc, datestr)
		if err != nil {
			log.Println(err)
			continue
		}
		fmt.Printf("%+s\n", jsonMarshalled)

		// incrementing linecount here (instead of at the top) to avoid
		// counting lines for the cases where parsing it resulted in an error.
		linecount++
		go func(url string, key string, payload []byte) {
			var gr GoroutineResponse
			gr.resp, gr.err = makeRequest(url, key, payload)
			ch <- gr
		}(url, key, jsonMarshalled)
	}

	for i := 0; i < linecount; i++ {
		gr := <-ch
		if gr.err != nil {
			log.Println(gr.err)
			continue
		}
		fmt.Println(gr.resp)
	}
}

// make request to API and return the response body and status code
func makeRequest(url string, key string, payload []byte) (Respp, error) {
	// TODO: pass key as an HTTP Basic Auth header data
	resp, err := http.Post(url, "application/json", bytes.NewReader(payload))
	if err != nil {
		return Respp{}, fmt.Errorf("error posting data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Respp{}, fmt.Errorf("got status %d", resp.StatusCode)
	}

	var opresponse Opresponse
	if err := json.NewDecoder(resp.Body).Decode(&opresponse); err != nil {
		return Respp{}, fmt.Errorf("error decoding response: %v", err)
	}

	response := Respp{
		// Original:   word, // we do not have a word string anymore because the param has changed
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
func makePostDataJSON(wp int, dur string, desc string, datestr string) ([]byte, error) {
	request := Request{}
	request.Links.Activity.Href = "api/v3/time_entries/activities/3"
	request.WorkPackage.Href = "api/v3/work_package/" + strconv.Itoa(wp)
	request.Hours = "PT" + dur + "H"
	request.Comment.Raw = desc
	request.SpentOn = datestr

	data, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %v", err)
	}

	return data, nil
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
