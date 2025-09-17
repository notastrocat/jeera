package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
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
	Summary             string      `json:"summary"`
	Description         string      `json:"description"`
	IssueType           *IssueType  `json:"issuetype,omitempty"`
	Project             *Project    `json:"project,omitempty"`
	Priority            *Priority   `json:"priority,omitempty"`
	Status              *Status     `json:"status,omitempty"`
	AcceptanceCriteria  string      `json:"customfield_11028"` // Replace with your actual custom field ID
	StoryPoints         float32     `json:"customfield_10002"` // Replace with your actual custom field ID
	Assignee            *Assignee   `json:"assignee,omitempty"`
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

type Assignee struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName,omitempty"`
}

type Comment struct {
	ID          string `json:"id,omitempty"`
	Body        string `json:"body"`
	Author      string `json:"author,omitempty"`
	Created     string `json:"created,omitempty"`
	LastUpdated string `json:"updated,omitempty"`
	TimeZone    string `json:"timeZone,omitempty"`
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

func (client *JiraClient) UpdateAssignee(issueIDOrKey string, assignee *Assignee) error {
	updateRequest := make(map[string]string)
	if assignee != nil {
		updateRequest["name"] = assignee.Name
	}

	// updateRequest := map[string]interface{}{
	// 	"fields": updateFields,
	// }

	endpoint := fmt.Sprintf("/rest/api/2/issue/%s/assignee", issueIDOrKey)

	resp, err := client.makeRequest("PUT", endpoint, updateRequest)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to update assignee: status %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
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
	if fields.AcceptanceCriteria != "" {
		updateFields["customfield_11028"] = fields.AcceptanceCriteria
	}
	if fields.StoryPoints >= 0.0 {
		updateFields["customfield_10002"] = fields.StoryPoints
	}

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

func (client *JiraClient) DoTransition(issueIDOrKey, transitionID string) error {
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

func (client *JiraClient) GetComments(issueIDOrKey string) ([]Comment, error) {
	endpoint := fmt.Sprintf("/rest/api/2/issue/%s/comment", issueIDOrKey)

	resp, err := client.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get comments: status %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var result []Comment
	var temp struct {
		Comments []map[string]interface{} `json:"comments"`
	}

	if *DEBUGflag {
		fmt.Println("GetComments response:")
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Println(string(bodyBytes))
		// Rewind the response body for decoding
		resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	if err := json.NewDecoder(resp.Body).Decode(&temp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	/*
		    "body": "Change 501887 had a related patch set uploaded by Sharma, Archit:\n[GTJ-691] enable QNX_IPC for ACF bindings\n\n[https://gitgerrit.asux.aptiv.com/c/ux/ispfw/core/+/501887|https://gitgerrit.asux.aptiv.com/c/ux/ispfw/core/+/501887]",
			"updateAuthor": {
				"self": "https://jiraprod.aptiv.com/rest/api/2/user?username=GID_GERRITJIRA",
				"name": "GID_GERRITJIRA",
				"key": "JIRAUSER176187",
				"emailAddress": "",
				"avatarUrls": {
				"48x48": "https://jiraprod.aptiv.com/secure/useravatar?avatarId=11426",
				"24x24": "https://jiraprod.aptiv.com/secure/useravatar?size=small&avatarId=11426",
				"16x16": "https://jiraprod.aptiv.com/secure/useravatar?size=xsmall&avatarId=11426",
				"32x32": "https://jiraprod.aptiv.com/secure/useravatar?size=medium&avatarId=11426"
				},
				"displayName": "GID_GERRITJIRA",
				"active": true,
				"timeZone": "Europe/Amsterdam"
			},
			"created": "2025-09-08T07:49:29.479+0200",
			"updated": "2025-09-08T07:49:29.479+0200"
	*/

	for _, c := range temp.Comments {
		authorMeta := c["updateAuthor"].(map[string]interface{})

		result = append(result, Comment{
			ID:          c["id"].(string),
			Body:        c["body"].(string),
			Author:      authorMeta["displayName"].(string),
			Created:     c["created"].(string),
			LastUpdated: c["updated"].(string),
			TimeZone:    authorMeta["timeZone"].(string),
		})
	}

	return result, nil
}

func (client *JiraClient) GetBoardID(boardName string) (int, error) {
	// "/rest/agile/1.0/board?name=CoreFW_AutoScrum"
	endpoint := fmt.Sprintf("/rest/agile/1.0/board?name=%s", boardName)

	resp, err := client.makeRequest("GET", endpoint, nil)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("failed to get board ID: status %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var temp struct {
		Values []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"values"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&temp); err != nil {
		return 0, fmt.Errorf("failed to decode response: %v", err)
	}

	if len(temp.Values) == 0 {
		return 0, fmt.Errorf("no board found with name %s\nBoard names are case sensitive. Please check your board name", boardName)
	}

	// Return the ID of the first matching board
	return temp.Values[0].ID, nil
}

func (client *JiraClient) GetProjectKeys(boardID int, projectKeys []string) ([]string, error) {
	// interestingly enough, the board ID is an int as opposed to a string like most other IDs in JIRA
	// /rest/agile/1.0/board/14190/project
	endpoint := fmt.Sprintf("/rest/agile/1.0/board/%d/project", boardID)

	resp, err := client.makeRequest("GET", endpoint, nil)
	if err != nil {
		return projectKeys, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return projectKeys, fmt.Errorf("failed to get project key: status %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var temp struct {
		Values []struct {
			Key string `json:"key"`
		} `json:"values"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&temp); err != nil {
		return projectKeys, fmt.Errorf("failed to decode response: %v", err)
	}

	if len(temp.Values) == 0 {
		return projectKeys, fmt.Errorf("no project found for board ID %d", boardID)
	}

	for _, v := range temp.Values {
		projectKeys = append(projectKeys, v.Key)
	}

	// Return the key of the first project
	return projectKeys, nil
}

func (client *JiraClient) GetActiveSprintID(boardID int) (int, string, error) {
	endpoint := fmt.Sprintf("/rest/agile/1.0/board/%d/sprint?state=active", boardID) // amazing, didn't know about this query param.

	resp, err := client.makeRequest("GET", endpoint, nil)
	if err != nil {
		return 0, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return 0, "", fmt.Errorf("failed to get active sprint: status %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var temp struct {
		Values []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"values"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&temp); err != nil {
		return 0, "", fmt.Errorf("failed to decode response: %v", err)
	}

	if len(temp.Values) == 0 {
		return 0, "", fmt.Errorf("no active sprints found for board ID %d", boardID)
	}

	// Return the ID of the first active sprint
	return temp.Values[0].ID, temp.Values[0].Name, nil
}

// func (client *JiraClient) GetIssuesInSprint(sprintID int) ([]Issue, error) {
// 	endpoint := fmt.Sprintf("/rest/agile/1.0/sprint/%d/issue", sprintID)

// 	resp, err := client.makeRequest("GET", endpoint, nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		bodyBytes, _ := io.ReadAll(resp.Body)
// 		return nil, fmt.Errorf("failed to get issues in sprint: status %d, body: %s", resp.StatusCode, string(bodyBytes))
// 	}

// 	var temp struct {
// 		Issues []Issue `json:"issues"`
// 	}

// 	if err := json.NewDecoder(resp.Body).Decode(&temp); err != nil {
// 		return nil, fmt.Errorf("failed to decode response: %v", err)
// 	}

// 	return temp.Issues, nil
// }

func (client *JiraClient) GetMyIssuesInActiveSprint(projectKey, sprintID, assignee string) ([]Issue, error) {
	jql := fmt.Sprintf("project = %s AND sprint = \"%s\" AND assignee in (%s)", projectKey, sprintID, assignee)
	encodedJQL := url.QueryEscape(jql) // Properly encode the JQL

	endpoint := fmt.Sprintf("/rest/api/2/search?jql=%s", encodedJQL)

	resp, err := client.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get issues in sprint: status %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var temp struct {
		Issues []Issue `json:"issues"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&temp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return temp.Issues, nil
}
