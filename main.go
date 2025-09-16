package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var DEBUGflag = flag.Bool("debug", false, "enable debugging messages for the app")

func main() {
	// Load configuration
	config := LoadConfig()
	
	// Validate configuration
	if !config.Validate() {
		fmt.Println("Error: Missing required configuration!")
		fmt.Println("Please create a .env file or set environment variables:")
		fmt.Println("  JIRA_BASE_URL - Your JIRA instance URL (e.g., https://yourcompany.atlassian.net)")
		fmt.Println("  JIRA_USERNAME - Your JIRA username/email (e.g., your.email@company.com)")
		fmt.Println("  JIRA_PAT - Your JIRA Personal Access Token (recommended)")
		fmt.Println("    OR")
		fmt.Println("  JIRA_API_TOKEN - Your JIRA API token (legacy)")
		fmt.Println("")
		fmt.Println("Option 1 - Create .env file:")
		fmt.Println("  cp .env.example .env")
		fmt.Println("  # Edit .env file with your actual values")
		fmt.Println("")
		fmt.Println("Option 2 - Use environment variables:")
		fmt.Println("  export JIRA_BASE_URL=https://yourcompany.atlassian.net")
		fmt.Println("  export JIRA_USERNAME=your.email@company.com")
		fmt.Println("  export JIRA_PAT=your-personal-access-token")
		os.Exit(1)
	}

	// Create JIRA client
	client := NewJiraClient(config)

	flag.Parse()
	// Start interactive CLI
	if *DEBUGflag {
		fmt.Println("JIRA Auto - Issue Management Tool (Running in Debug Mode)")
	} else {
		fmt.Println("JIRA Auto - Issue Management Tool")
	}
	fmt.Println("=================================")
	fmt.Printf("Connected to: %s\n", config.BaseURL)
	fmt.Printf("Username: %s\n", config.Username)
	if config.UsePAT {
		fmt.Printf("Authentication: Personal Access Token (Bearer)\n\n")
	} else {
		fmt.Printf("Authentication: API Token (Basic)\n\n")
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("Available commands:")
		fmt.Println("  1. Create issue")
		fmt.Println("  2. Get issue")
		fmt.Println("  3. Update issue")
		fmt.Println("  4. Transition issue")
		fmt.Println("  5. Delete issue")
		fmt.Println("  6. Bulk create issues")
		fmt.Println("  7. Exit")
		fmt.Print("\nEnter your choice (1-7): ")

		scanner.Scan()
		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			createIssueInteractive(client, scanner)
		case "2":
			getIssueInteractive(client, scanner)
		case "3":
			updateIssueInteractive(client, scanner)
		case "4":
			doTransitionInteractive(client, scanner)
		case "5":
			fmt.Println("not implemented yet")
			// updateIssueInteractive(client, scanner)
		case "6":
			fmt.Println("not implemented yet")
			// updateIssueInteractive(client, scanner)
		case "7":
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid choice. Please enter 1-4.")
		}

		fmt.Println()
	}
}

func createIssueInteractive(client *JiraClient, scanner *bufio.Scanner) {
	fmt.Println("\n--- Create New Issue ---")
	
	fmt.Print("Project Key: ")
	scanner.Scan()
	projectKey := strings.TrimSpace(scanner.Text())
	
	fmt.Print("Issue Type (e.g., Bug, Task, Story): ")
	scanner.Scan()
	issueType := strings.TrimSpace(scanner.Text())
	
	fmt.Print("Summary: ")
	scanner.Scan()
	summary := strings.TrimSpace(scanner.Text())
	
	fmt.Print("Description (optional): ")
	scanner.Scan()
	description := strings.TrimSpace(scanner.Text())

	issue := &Issue{
		Fields: IssueFields{
			Project: Project{
				Key: projectKey,
			},
			IssueType: IssueType{
				Name: issueType,
			},
			Summary:     summary,
			Description: description,
		},
	}

	result, err := client.CreateIssue(issue)
	if err != nil {
		log.Printf("Error creating issue: %v", err)
		return
	}

	fmt.Printf("✅ Issue created successfully!\n")
	fmt.Printf("Key: %s\n", result.Key)
	fmt.Printf("ID: %s\n", result.ID)
}

func getIssueInteractive(client *JiraClient, scanner *bufio.Scanner) {
	fmt.Println("\n--- Get Issue ---")
	
	fmt.Print("Issue ID or Key: ")
	scanner.Scan()
	issueIDOrKey := strings.TrimSpace(scanner.Text())

	issue, err := client.GetIssue(issueIDOrKey)
	if err != nil {
		log.Printf("Error getting issue: %v", err)
		return
	}

	fmt.Printf("✅ Issue retrieved successfully!\n")
	fmt.Printf("Key: %s\n", issue.Key)
	fmt.Printf("ID: %s\n", issue.ID)
	fmt.Printf("Summary: %s\n", issue.Fields.Summary)
	fmt.Printf("Description: %s\n", issue.Fields.Description)
	fmt.Printf("Issue Type: %s\n", issue.Fields.IssueType.Name)
	fmt.Printf("Project: %s\n", issue.Fields.Project.Key)
	if issue.Fields.Status != nil {
		fmt.Printf("Status: %s\n", issue.Fields.Status.Name)
	}
	if issue.Fields.Priority != nil {
		fmt.Printf("Priority: %s\n", issue.Fields.Priority.Name)
	}
}

func updateIssueInteractive(client *JiraClient, scanner *bufio.Scanner) {
	fmt.Println("\n--- Update Issue ---")

	fmt.Print("Issue ID or Key: ")
	scanner.Scan()
	issueIDOrKey := strings.TrimSpace(scanner.Text())

	fmt.Print("New Summary (leave empty to keep current): ")
	scanner.Scan()
	summary := strings.TrimSpace(scanner.Text())

	fmt.Print("New Description (leave empty to keep current): ")
	scanner.Scan()
	description := strings.TrimSpace(scanner.Text())

	fmt.Print("New Acceptance Criteria (leave empty to keep current): ")
	scanner.Scan()
	acceptanceCriteria := strings.TrimSpace(scanner.Text())

	fmt.Print("New Story Points (leave empty to keep current): ")
	scanner.Scan()
	// storyPoints should be an integer
	storyPoints := strings.TrimSpace(scanner.Text())

	// Build update fields
	fields := IssueFields{}
	if summary != "" {
		fields.Summary = summary
	}
	if description != "" {
		fields.Description = description
	}
	if acceptanceCriteria != "" {
		fields.AcceptanceCriteria = acceptanceCriteria
	}
	if storyPoints != "" {
		if sp, err := strconv.ParseFloat(storyPoints, 32); err != nil {
			fmt.Printf("Invalid story points value: %v\n", err)
		} else {
			fields.StoryPoints = float32(sp)
		}
	}

	if fields.Summary == "" && fields.Description == "" && fields.AcceptanceCriteria == "" && fields.StoryPoints <= 0.0 {
		fmt.Println("No changes specified.")
		return
	}

	err := client.UpdateIssue(issueIDOrKey, fields)
	if err != nil {
		log.Printf("Error updating issue: %v", err)
		return
	}

	fmt.Printf("✅ Issue %s updated successfully!\n", issueIDOrKey)
}

func doTransitionInteractive(client *JiraClient, scanner *bufio.Scanner) {
	fmt.Println("\n--- Transition Issue ---")

	fmt.Print("Issue ID or Key: ")
	scanner.Scan()
	issueIDOrKey := strings.TrimSpace(scanner.Text())

	// Fetch available transitions
	transitions, err := client.GetTransitions(issueIDOrKey)
	if err != nil {
		log.Printf("Error fetching transitions: %v", err)
		return
	}

	if len(transitions) == 0 {
		fmt.Println("No transitions available for this issue.")
		return
	}

	fmt.Println("Available Transitions:")
	for i, t := range transitions {
		fmt.Printf("  %d. %s (ID: %s)\n", i+1, t.Name, t.ID)
	}

	fmt.Print("\nSelect transition number: ")
	scanner.Scan()
	choiceStr := strings.TrimSpace(scanner.Text())
	var choice int
	_, err = fmt.Sscanf(choiceStr, "%d", &choice)
	if err != nil || choice < 1 || choice > len(transitions) {
		fmt.Println("Invalid choice.")
		return
	}

	selectedTransition := transitions[choice-1]

	err = client.DoTransition(issueIDOrKey, selectedTransition.ID)
	if err != nil {
		log.Printf("Error performing transition: %v", err)
		return
	}

	fmt.Printf("✅ Issue %s transitioned to '%s' successfully!\n", issueIDOrKey, selectedTransition.Name)
}
