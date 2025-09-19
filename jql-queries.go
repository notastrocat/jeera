package main

import (
	"fmt"
	"net/http"
	"io"
	"encoding/json"
	"net/url"

	"jeera/decorators"
	"jeera/decorators/foreground"
)

func (client *JiraClient) GetMyIssuesInActiveSprint(projectKey, sprintID, assignee string) ([]Issue, error) {
	jql := fmt.Sprintf("project = %s AND sprint = \"%s\" AND assignee in (%s)", projectKey, sprintID, assignee)
	encodedJQL := url.QueryEscape(jql) // Properly encode the JQL

	if *DEBUGflag {
		fmt.Printf(foreground.YELLOW + "[DEBUG] JQL Query: %s\n" + decorators.RESET_ALL, jql)
		fmt.Printf(foreground.YELLOW + "[DEBUG] Encoded JQL: %s\n" + decorators.RESET_ALL, encodedJQL)
	}

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
