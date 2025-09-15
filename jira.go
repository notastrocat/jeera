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
	IssueType   IssueType   `json:"issuetype,omitempty"`
	Project     Project     `json:"project,omitempty"`
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

// Transition represents a JIRA issue transition
type Transition struct {
	ID   string `json:"id"`
	Name string `json:"name"`
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
	updateFields := make(map[string]interface{})

	if fields.Summary != "" {
        updateFields["summary"] = fields.Summary
    }
    if fields.Description != "" {
        updateFields["description"] = fields.Description
    }
    if fields.IssueType.ID != "" || fields.IssueType.Name != "" {
        issueType := make(map[string]interface{})
        if fields.IssueType.ID != "" {
            issueType["id"] = fields.IssueType.ID
        }
        if fields.IssueType.Name != "" {
            issueType["name"] = fields.IssueType.Name
        }
        updateFields["issuetype"] = issueType
    }
    if fields.Project.Key != "" || fields.Project.ID != "" {
        project := make(map[string]interface{})
        if fields.Project.Key != "" {
            project["key"] = fields.Project.Key
        }
        if fields.Project.ID != "" {
            project["id"] = fields.Project.ID
        }
        updateFields["project"] = project
    }
	if fields.Priority != nil {
        priority := make(map[string]interface{})
        if fields.Priority.ID != "" {
            priority["id"] = fields.Priority.ID
        }
        if fields.Priority.Name != "" {
            priority["name"] = fields.Priority.Name
        }
        updateFields["priority"] = priority
    }
	// if fields.Status != nil {
    //     status := make(map[string]interface{})
    //     if fields.Status.ID != "" {
    //         status["id"] = fields.Status.ID
    //     }
    //     if fields.Status.Name != "" {
    //         status["name"] = fields.Status.Name
    //     }
    //     updateFields["statuscategory"] = status
    // }

	updateRequest := map[string]interface{}{
        "fields": updateFields,
    }

	if *DEBUGflag {
		fmt.Printf("UpdateIssue request: %+v\n", updateRequest)
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

func (client *JiraClient) GetTransitions(issueIDOrKey string) ([]Transition, error) {
	endpoint := fmt.Sprintf("/rest/api/2/issue/%s/transitions", issueIDOrKey)

	resp, err := client.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get transitions: status %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var result []Transition
	var temp struct {
		RawTransitions []map[string]interface{} `json:"transitions"`
	}

	// so this approach is a bit hacky but it works to get around the dynamic nature of the transitions
	// ideally we would define a proper struct but the fields can vary widely
	// so we just use a map[string]interface{} and extract the fields we care about
	// see https://developer.atlassian.com/cloud/jira/platform/rest/v2/api-group

	if *DEBUGflag {
		fmt.Println("GetTransitions response:")
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Println(string(bodyBytes))
		// Rewind the response body for decoding
		resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	if err := json.NewDecoder(resp.Body).Decode(&temp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	if *DEBUGflag {
		fmt.Printf("Decoded transitions: %+v\n", temp.RawTransitions)
		fmt.Printf("\n\nchecking value: %s\n\n", temp.RawTransitions[0]["name"])
	}

	for _, t := range temp.RawTransitions {
		result = append(result, Transition{
			ID:   t["id"].(string),
			Name: t["name"].(string),
		})
	}

	return result, nil
}

func (client *JiraClient) DoTransition(issueIDOrKey , transitionID string) error {
	endpoint := fmt.Sprintf("/rest/api/2/issue/%s/transitions", issueIDOrKey)

	transitionRequest := map[string]interface{}{
		"transition": map[string]string{
			"id": transitionID,
		},
	}

	resp, err := client.makeRequest("POST", endpoint, transitionRequest)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to do transition: status %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}
