package chttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// Sends API request to Pterodactyl API (or perhaps technically any endpoint with /api/ :)).
func SendHTTPReq(endpoint string, auth string, request_type string, form_data map[string]interface{}) (string, int, error) {
	// Initialize data and return code (status code).
	d := ""
	rc := -1

	var post_body io.Reader

	// Check to see if we need to send post data.
	if request_type == "POST" || request_type == "PUT" {
		// Convert to JSON and use as body.
		j, err := json.Marshal(form_data)

		if err != nil {
			return d, rc, err
		}

		// Read byte array into IO reader.
		post_body = bytes.NewBuffer(j)

		if err != nil {
			return d, rc, err
		}
	}

	// Setup HTTP GET request.
	client := &http.Client{Timeout: time.Second * 5}
	req, err := http.NewRequest(request_type, endpoint, post_body)

	if err != nil {
		fmt.Println(err)

		return d, rc, err
	}

	// Set Application API token.
	req.Header.Set("Authorization", auth)

	// Accept only JSON.
	req.Header.Set("Accept", "application/json")

	// Set content stype.
	req.Header.Set("Content-Type", "application/json")

	// Perform HTTP request and check for errors.
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)

		return d, rc, err
	}

	// Set return code.
	rc = resp.StatusCode

	// Read body.
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err)

		return d, rc, err
	}

	// Return data as a string.
	d = string(body)

	return d, rc, nil
}
