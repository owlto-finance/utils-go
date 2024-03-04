package network

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func Request(url string, data interface{}, result interface{}) error {

	// Create a new HTTP client
	client := http.Client{}
	var method = "GET"
	var dataBytes []byte
	if data != nil {
		method = "POST"
		var err error
		dataBytes, err = json.Marshal(data)
		if err != nil {
			return fmt.Errorf("failed to marshal request body %v : %v", url, err)
		}
	}

	// Create a new HTTP request
	req, err := http.NewRequest(method, url, bytes.NewReader(dataBytes))
	if err != nil {
		return fmt.Errorf("error creating request %v : %v", url, err)
	}

	// Set request headers
	req.Header.Set("Accept", "application/json")

	// Send the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request %v : %v", url, err)
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code %v : %v", url, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error read body %v : %v - %v", url, err, body)
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return fmt.Errorf("error read body %v : %v - %v", url, err, body)
	}
	return nil
}
