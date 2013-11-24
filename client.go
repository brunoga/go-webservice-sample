// Copyright 2013 Bruno Albuquerque (bga@bug-br.org.br).
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not
// use this file except in compliance with the License. You may obtain a copy of
// the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations under
// the License.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var requestUrl = flag.String("request_url", "",
	"Host and port to send the request to")
var requestMethod = flag.String("request_method", "GET",
	"HTTP request method")
var requestData = flag.String("request_data", "",
	"Marshalled JSON data")

func doRequest(requestMethod, requestUrl,
	requestData string) (*http.Response, error) {
	var res *http.Response
	var err error

	upperRequestMethod := strings.ToUpper(requestMethod)
	switch upperRequestMethod {
	case "DELETE", "PATCH", "PUT":
		if len(requestData) == 0 && upperRequestMethod != "DELETE" {
			return nil, fmt.Errorf(
				"--request_data must be provided")
		}

		contentBuffer := bytes.NewBufferString(requestData)
		req, err := http.NewRequest(upperRequestMethod, requestUrl,
			contentBuffer)
		if err != nil {
			return nil, fmt.Errorf(
				"error creating new HTTP request")
		}

		res, err = http.DefaultClient.Do(req)
	case "GET":
		res, err = http.Get(requestUrl)
	case "POST":
		if len(requestData) == 0 {
			return nil, fmt.Errorf(
				"--request_data must be provided")
		}

		contentBuffer := bytes.NewBufferString(requestData)
		res, err = http.Post (requestUrl, "application/json",
			contentBuffer)
	default:
		return nil, fmt.Errorf(
			"invalid --request_method provided : %s",
				requestMethod)
	}

	return res, err
}

func main() {
	flag.Parse()

	if len(*requestUrl) == 0 {
		fmt.Println("--request_url must be provided")
		return
	}

	if len(*requestMethod) == 0 {
		fmt.Println("--request_method must be provided")
	}

	res, err := doRequest(*requestMethod, *requestUrl, *requestData)
	if err != nil {
		fmt.Println("Error executing HTTP request :", err)
		return
	}

	defer res.Body.Close()
	responseData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response data :", err)
		return
	}

	fmt.Println(string(responseData))
}
