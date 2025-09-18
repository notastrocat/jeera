package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"jeera/decorators"
	"jeera/decorators/foreground"
)

// makeRequest performs an HTTP request with authentication
func (client *JiraClient) makeRequest(method, endpoint string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader

	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf(foreground.RED + "[ERROR] failed to marshal request body: %v" + decorators.RESET_ALL, err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	url := fmt.Sprintf("%s%s", client.config.BaseURL, endpoint)
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf(foreground.RED + "[ERROR] failed to create request: %v" + decorators.RESET_ALL, err)
	}

	// Add authentication header
	if client.config.UsePAT {
		// Use Bearer authentication for Personal Access Tokens
		req.Header.Set("Authorization", "Bearer "+client.config.APIToken)
	} else {
		// Use Basic authentication for API tokens
		auth := base64.StdEncoding.EncodeToString([]byte(client.config.Username + ":" + client.config.APIToken))
		req.Header.Set("Authorization", "Basic "+auth)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}

	return resp, nil
}

func (client *JiraClient) GetCurrentUser() (*User, error) {
	resp, err := client.makeRequest("GET", "/rest/api/2/myself", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf(foreground.RED + "[ERROR] failed to get current user: status %d, body: %s" + decorators.RESET_ALL, resp.StatusCode, string(bodyBytes))
	}

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf(foreground.RED + "[ERROR] failed to decode response: %v" + decorators.RESET_ALL, err)
	}

	return &user, nil
}
