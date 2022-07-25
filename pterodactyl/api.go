package pterodactyl

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gamemann/Rust-Auto-Wipe/config"
	"github.com/gamemann/Rust-Auto-Wipe/processor"
)

func SendAPIRequest(cfg *config.Config, wipedata *processor.WipeData, request_type string, request_endpoint string, post_data string) (string, int, error) {
	// Initialize data and return code (status code).
	d := ""
	rc := -1

	// Compile our URL.
	urlstr := wipedata.APIURL + "/api/" + request_endpoint

	// Setup HTTP GET request.
	client := &http.Client{Timeout: time.Second * 5}
	req, err := http.NewRequest(request_type, urlstr, nil)

	if err != nil {
		fmt.Println(err)

		return d, rc, err
	}

	// Set Application API token.
	req.Header.Set("Authorization", "Bearer "+wipedata.APIToken)

	// Accept only JSON.
	req.Header.Set("Accept", "application/json")

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
