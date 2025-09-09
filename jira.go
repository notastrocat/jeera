package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// JiraClient represents a JIRA API client
type JiraClient struct {
	config     *Config
	httpClient *http.Client
}

// NewJiraClient creates a new JIRA client
func NewJiraClient(config *Config) *JiraClient {
	return &JiraClient{
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Issue represents a JIRA issue structure
type Issue struct {
	ID     string      `json:"id,omitempty"`
	Key    string      `json:"key,omitempty"`
	Fields IssueFields `json:"fields"`
}

// IssueFields represents the fields of a JIRA issue
type IssueFields struct {
	Summary     string      `json:"summary"`
	Description string      `json:"description,omitempty"`
	IssueType   IssueType   `json:"issuetype"`
	Project     Project     `json:"project"`
	Priority    *Priority   `json:"priority,omitempty"`
	Status      *Status     `json:"status,omitempty"`
}

// IssueType represents a JIRA issue type
type IssueType struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name"`
}

// Project represents a JIRA project
type Project struct {
	Key string `json:"key"`
	ID  string `json:"id,omitempty"`
}

// Priority represents a JIRA priority
type Priority struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name"`
}

// Status represents a JIRA status
type Status struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name"`
}

// CreateIssueRequest represents the request structure for creating an issue
type CreateIssueRequest struct {
	Fields IssueFields `json:"fields"`
}

// CreateIssueResponse represents the response from creating an issue
type CreateIssueResponse struct {
	ID   string `json:"id"`
	Key  string `json:"key"`
	Self string `json:"self"`
}

// makeRequest performs an HTTP request with authentication
func (client *JiraClient) makeRequest(method, endpoint string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %v", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	url := fmt.Sprintf("%s%s", client.config.BaseURL, endpoint)
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
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

// CreateIssue creates a new JIRA issue
func (client *JiraClient) CreateIssue(issue *Issue) (*CreateIssueResponse, error) {
	request := CreateIssueRequest{Fields: issue.Fields}
	
	resp, err := client.makeRequest("POST", "/rest/api/2/issue", request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to create issue: status %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var result CreateIssueResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &result, nil
}

// GetIssue retrieves a JIRA issue by ID or key
func (client *JiraClient) GetIssue(issueIDOrKey string) (*Issue, error) {
	endpoint := fmt.Sprintf("/rest/api/2/issue/%s", issueIDOrKey)
	
	resp, err := client.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get issue: status %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var issue Issue
	if err := json.NewDecoder(resp.Body).Decode(&issue); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &issue, nil
}

// UpdateIssue updates an existing JIRA issue
func (client *JiraClient) UpdateIssue(issueIDOrKey string, fields IssueFields) error {
	updateRequest := map[string]interface{}{
		"fields": fields,
	}
	
	endpoint := fmt.Sprintf("/rest/api/2/issue/%s", issueIDOrKey)
	
	resp, err := client.makeRequest("PUT", endpoint, updateRequest)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to update issue: status %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}
