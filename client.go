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

// This file implements a generic REST client with JSON support to make it
// easier to test web services. Example usage:
//
// go run client.go \
//	--request_url=http://127.0.0.1:3000/webservice \
//	--request_method=post \
//	--request_data='{"Id":0,"Email":"bga@bug-br.org.br"}'

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Command line flags.
var requestUrl = flag.String("request_url", "",
	"Host and port to send the request to")
var requestMethod = flag.String("request_method", "GET",
	"HTTP request method")
var requestData = flag.String("request_data", "",
	"Marshalled JSON data")

// doRequest executes an HTTP request to the given requestUrl using the given
// requestMethod and requestData.
func doRequest(requestMethod, requestUrl,
	requestData string) (*http.Response, error) {
	// These will hold the return value.
	var res *http.Response
	var err error

	// Convert method to uppercase for easier checking.
	upperRequestMethod := strings.ToUpper(requestMethod)
	switch upperRequestMethod {
	case "DELETE", "PATCH", "PUT":
		// All these methods have no shortcuts in Go's HTTP library, so
		// we have to do them manually.
		if len(requestData) == 0 && upperRequestMethod != "DELETE" {
			// All methods (except for DELETE) require data.
			return nil, fmt.Errorf(
				"--request_data must be provided")
		}

		// NewRequest requires a Reader, so we create a byte buffer
		// for our string data.
		contentBuffer := bytes.NewBufferString(requestData)

		req, err := http.NewRequest(upperRequestMethod, requestUrl,
			contentBuffer)
		if err != nil {
			// Failed creating HTTP request.
			return nil, fmt.Errorf(
				"error creating new HTTP request")
		}

		// Use the default HTTP client to execute the request.
		res, err = http.DefaultClient.Do(req)
	case "GET":
		// Use the HTTP library Get() method.
		res, err = http.Get(requestUrl)
	case "POST":
		// Use the HTTP library Post() method.
		if len(requestData) == 0 {
			// Post requires data.
			return nil, fmt.Errorf(
				"--request_data must be provided")
		}

		// Create Reader for Post data.
		contentBuffer := bytes.NewBufferString(requestData)

		res, err = http.Post (requestUrl, "application/json",
			contentBuffer)
	default:
		// We doń't know how to handle this request.
		return nil, fmt.Errorf(
			"invalid --request_method provided : %s",
				requestMethod)
	}

	return res, err
}

func main() {
	// Parse command line flags.
	flag.Parse()

	if len(*requestUrl) == 0 {
		// We need a request URL.
		fmt.Println("--request_url must be provided")
		return
	}

	if len(*requestMethod) == 0 {
		// And we also need a method.
		fmt.Println("--request_method must be provided")
	}

	// Execute request (if possible).
	res, err := doRequest(*requestMethod, *requestUrl, *requestData)
	if err != nil {
		// Request failed.
		fmt.Println("Error executing HTTP request :", err)
		return
	}

	// Make sure res.Body is closed whwn we are done.
	defer res.Body.Close()

	// Read body data (i.e. our response).
	responseData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		// Could not read data.
		fmt.Println("Error reading response data :", err)
		return
	}

	// Print received data.
	fmt.Println(string(responseData))
}

