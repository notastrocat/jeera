package main

import (
	"fmt"
	"net/http"
	"io"
	"encoding/json"

	"jeera/decorators"
	"jeera/decorators/foreground"
)

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
		return 0, fmt.Errorf(foreground.RED + "[ERROR] failed to get board ID: status %d, body: %s" + decorators.RESET_ALL, resp.StatusCode, string(bodyBytes))
	}

	var temp struct {
		Values []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"values"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&temp); err != nil {
		return 0, fmt.Errorf(foreground.RED + "[ERROR] failed to decode response: %v" + decorators.RESET_ALL, err)
	}

	if len(temp.Values) == 0 {
		return 0, fmt.Errorf(foreground.RED + "[ERROR] no board found with name %s\nBoard names are case sensitive. Please check your board name" + decorators.RESET_ALL, boardName)
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
		return projectKeys, fmt.Errorf(foreground.RED + "[ERROR] failed to get project key: status %d, body: %s" + decorators.RESET_ALL, resp.StatusCode, string(bodyBytes))
	}

	var temp struct {
		Values []struct {
			Key string `json:"key"`
		} `json:"values"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&temp); err != nil {
		return projectKeys, fmt.Errorf(foreground.RED + "[ERROR] failed to decode response: %v" + decorators.RESET_ALL, err)
	}

	if len(temp.Values) == 0 {
		return projectKeys, fmt.Errorf(foreground.RED + "[ERROR] no project found for board ID %d" + decorators.RESET_ALL, boardID)
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
		return 0, "", fmt.Errorf(foreground.RED + "[ERROR] failed to get active sprint: status %d, body: %s" + decorators.RESET_ALL, resp.StatusCode, string(bodyBytes))
	}

	var temp struct {
		Values []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"values"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&temp); err != nil {
		return 0, "", fmt.Errorf(foreground.RED + "[ERROR] failed to decode response: %v" + decorators.RESET_ALL, err)
	}

	if len(temp.Values) == 0 {
		return 0, "", fmt.Errorf(foreground.RED + "[ERROR] no active sprints found for board ID %d" + decorators.RESET_ALL, boardID)
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
