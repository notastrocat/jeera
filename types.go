package main

import (
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

type User struct {
	Name         string  `json:"name"`
	DisplayName  string  `json:"displayName,omitempty"`
	Email        string  `json:"emailAddress,omitempty"`
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
