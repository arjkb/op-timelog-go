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
	Key      string
	Url      string
	Activity ConfigActivity
}

type ConfigActivity struct {
	Default    int
	Meeting    int
	MeetingWPS []int `toml:"meeting_wps"`
}

type Request struct {
	Links struct {
		Activity struct {
			Href string `json:"href"`
		} `json:"activity"`
		WorkPackage struct {
			Href string `json:"href"`
		} `json:"workPackage"`
	} `json:"_links"`
	Hours   string `json:"hours"`
	Comment struct {
		Raw string `json:"raw"`
	} `json:"comment"`
	SpentOn string `json:"spentOn"`
}

type GoroutineResponse struct {
	StatusCode int
	Payload    []byte
	Err        error
}

func main() {
	var linecount int
	ch := make(chan GoroutineResponse)

	f := flag.String("file", "status_"+time.Now().Format("20060102")+".dailystatus", "file to read from")
	flag.Parse()

	filename := *f
	datestr := filename[7:15] // extracts date string from filenames of format "status_20210921.dailystatus"
	// TODO: make the above line failsafe

	// read config data
	config, err := getConfig("config.toml")
	if err != nil {
		log.Fatal(err)
	}

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

		jsonMarshalled, err := makePostDataJSON(wp, dur, desc, datestr, config.Activity)
		if err != nil {
			log.Println(err)
			continue
		}

		// incrementing linecount here (instead of at the top) to avoid
		// counting lines for the cases where parsing it resulted in an error.
		linecount++
		go func(url string, key string, payload []byte) {
			var gr GoroutineResponse
			gr.StatusCode, gr.Payload, gr.Err = makeRequest(url, key, payload)
			ch <- gr
		}(config.Url, config.Key, jsonMarshalled)
	}

	for i := 0; i < linecount; i++ {
		gr := <-ch
		if gr.Err != nil {
			log.Println(gr.Err)
			continue
		}
		fmt.Printf("%d %s\n", gr.StatusCode, gr.Payload)
	}
}

// make request to API and return the status code and original payload
func makeRequest(url string, key string, payload []byte) (int, []byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewReader(payload))
	if err != nil {
		return 0, nil, fmt.Errorf("error creating new request: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth("apikey", key)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("error posting data: %v", err)
	}
	defer resp.Body.Close()

	return resp.StatusCode, payload, nil
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
func makePostDataJSON(wp int, dur string, desc string, datestr string, configActivity ConfigActivity) ([]byte, error) {

	// determine activity to be used
	activityCode := configActivity.Default
	for _, meetingActivity := range configActivity.MeetingWPS {
		if wp == meetingActivity {
			activityCode = configActivity.Meeting
			break
		}
	}

	request := Request{}
	request.Links.Activity.Href = "/api/v3/time_entries/activities/" + strconv.Itoa(activityCode)
	request.Links.WorkPackage.Href = "/api/v3/work_packages/" + strconv.Itoa(wp)
	request.Hours = "PT" + dur + "H"
	request.Comment.Raw = desc
	request.SpentOn = datestr

	data, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %v", err)
	}

	return data, nil
}

// Read config file and return all its contents.
func getConfig(configFileName string) (Config, error) {
	file, err := os.Open(configFileName)
	if err != nil {
		return Config{}, fmt.Errorf("cannot open file: %v", err)
	}

	b, err := io.ReadAll(bufio.NewReader(file))
	if err != nil {
		return Config{}, fmt.Errorf("cannot read file: %v", err)
	}

	var conf Config
	if _, err := toml.Decode(string(b[:]), &conf); err != nil {
		return Config{}, fmt.Errorf("cannot decode toml: %v", err)
	}

	return conf, nil
}
