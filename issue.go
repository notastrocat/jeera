package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"jeera/decorators"
	"jeera/decorators/foreground"
)

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
		return nil, fmt.Errorf(foreground.RED + "[ERROR] failed to create issue: status %d, body: %s" + decorators.RESET_ALL, resp.StatusCode, string(bodyBytes))
	}

	var result CreateIssueResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf(foreground.RED + "[ERROR] failed to decode response: %v" + decorators.RESET_ALL, err)
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
		return nil, fmt.Errorf(foreground.RED + "[ERROR] failed to get issue: status %d, body: %s" + decorators.RESET_ALL, resp.StatusCode, string(bodyBytes))
	}

	var issue Issue
	if err := json.NewDecoder(resp.Body).Decode(&issue); err != nil {
		return nil, fmt.Errorf(foreground.RED + "[ERROR] failed to decode response: %v" + decorators.RESET_ALL, err)
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
		return fmt.Errorf(foreground.RED + "[ERROR] failed to update assignee: status %d, body: %s" + decorators.RESET_ALL, resp.StatusCode, string(bodyBytes))
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
		fmt.Printf(foreground.YELLOW + "[DEBUG] UpdateIssue request: %+v\n" + decorators.RESET_ALL, updateRequest)
	}

	endpoint := fmt.Sprintf("/rest/api/2/issue/%s", issueIDOrKey)

	resp, err := client.makeRequest("PUT", endpoint, updateRequest)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf(foreground.RED + "[ERROR] failed to update issue: status %d, body: %s" + decorators.RESET_ALL, resp.StatusCode, string(bodyBytes))
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
		return nil, fmt.Errorf(foreground.RED + "[ERROR] failed to get transitions: status %d, body: %s" + decorators.RESET_ALL, resp.StatusCode, string(bodyBytes))
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
		fmt.Println(foreground.YELLOW + "[DEBUG] GetTransitions response:" + decorators.RESET_ALL)
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Println(string(bodyBytes))
		// Rewind the response body for decoding
		resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	if err := json.NewDecoder(resp.Body).Decode(&temp); err != nil {
		return nil, fmt.Errorf(foreground.RED + "[ERROR] failed to decode response: %v" + decorators.RESET_ALL, err)
	}

	if *DEBUGflag {
		fmt.Printf(foreground.YELLOW + "[DEBUG] Decoded transitions: %+v\n" + decorators.RESET_ALL, temp.RawTransitions)
		fmt.Printf(foreground.YELLOW + "[DEBUG] \n\nchecking value: %s\n\n" + decorators.RESET_ALL, temp.RawTransitions[0]["name"])
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
		return fmt.Errorf(foreground.RED + "[ERROR] failed to do transition: status %d, body: %s" + decorators.RESET_ALL, resp.StatusCode, string(bodyBytes))
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
		return nil, fmt.Errorf(foreground.RED + "[ERROR] failed to get comments: status %d, body: %s" + decorators.RESET_ALL, resp.StatusCode, string(bodyBytes))
	}

	var result []Comment
	var temp struct {
		Comments []map[string]interface{} `json:"comments"`
	}

	if *DEBUGflag {
		fmt.Println(foreground.YELLOW + "[DEBUG] GetComments response:" + decorators.RESET_ALL)
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Println(string(bodyBytes))
		// Rewind the response body for decoding
		resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	if err := json.NewDecoder(resp.Body).Decode(&temp); err != nil {
		return nil, fmt.Errorf(foreground.RED + "[ERROR] failed to decode response: %v" + decorators.RESET_ALL, err)
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
