package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"jeera/decorators"
	"jeera/decorators/foreground"
)

var DEBUGflag = flag.Bool("debug", false, "enable debugging messages for the app")

func main() {
	// Load configuration
	config := LoadConfig()
	
	// Validate configuration
	if !config.Validate() {
		fmt.Println(foreground.RED + "[ERROR] Missing required configuration!" + decorators.RESET_ALL)
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

	var boardID      = flag.Int("board", -1, "JIRA Board ID (required for fetching sprint issues)")
	var boardName    = flag.String("board-name", "", "JIRA Board Name (if Board ID is not known)")
	var create       = flag.Bool("create", false, "Create a new issue")
	var get          = flag.Bool("get", false, "Get an existing issue")
	var update       = flag.Bool("update", false, "Update an existing issue")
	var transition   = flag.Bool("trans", false, "Transition an existing issue")
	var getComments  = flag.String("comments", "", "Manage comments of an existing issue")

	// Create JIRA client
	client := NewJiraClient(config)

	flag.Parse()

	// Start interactive CLI
	if *DEBUGflag {
		fmt.Println("jeera v2.0" + decorators.BOLD + " (Running in Debug Mode)" + decorators.RESET_ALL)
	} else {
		fmt.Println("jeera v2.0")
	}
	fmt.Println("=================================")

	fmt.Printf("Connected to: %s\n", config.BaseURL)
	user, err := client.GetCurrentUser()
	if err != nil {
		log.Printf(foreground.RED + "[ERROR] fetching user details: %v" + decorators.RESET_ALL, err)
	}
	fmt.Printf("Display Name: %s\n", user.DisplayName)
	fmt.Printf("Name: %s\n", user.Name)

	client.config.Username = user.Name    // override username with the actual name fetched from JIRA instance

	if config.UsePAT {
		fmt.Printf("Authentication: Personal Access Token (Bearer)\n\n")
	} else {
		fmt.Printf("Authentication: API Token (Basic)\n\n")
	}

	if *boardID >= 0 || *boardName != "" {
		if *DEBUGflag {
			fmt.Printf(foreground.YELLOW + "[DEBUG] board ID: %d\n" + decorators.RESET_ALL, *boardID)
			fmt.Printf(foreground.YELLOW + "[DEBUG] board name: %s\n" + decorators.RESET_ALL, *boardName)
		}
		fetchDataFromActiveSprint(user, client, boardID, boardName)
	}
	// else if *boardID <= 0 && *boardName == "" {
	// 	fmt.Println("`board ID` or `board name` is required to proceed.")
	// 	return
	// }

	if *create || *get || *update || *transition || *getComments != "" {
		scanner := bufio.NewScanner(os.Stdin)

		if *create {
			createIssueInteractive(client, scanner)
		}
		if *get {
			getIssueInteractive(client, scanner)
		}
		if *update {
			updateIssueInteractive(client, scanner)
		}
		if *transition {
			doTransitionInteractive(client, scanner)
		}
		if *getComments != "" {
			getCommentsInteractive(client, scanner)
		}
	}
}

func fetchDataFromActiveSprint(user *User, client *JiraClient, boardID *int, boardName *string) {
	var projectKeys []string
	var projectKey string

	if *boardID >= 0 {
		var err error
		projectKeys, err = client.GetProjectKeys(*boardID, projectKeys)
		if err != nil {
			log.Printf(foreground.RED + "[ERROR] getting project keys: %v" + decorators.RESET_ALL, err)
			return
		}

		// needed for JQL at the very least if not for anything else...
		fmt.Println("Associated Project Keys: -")
		for i, key := range projectKeys {
			fmt.Printf("Key %d: %s\n", i+1, key)
		}
		// probably select the first entry as the default key for the session...(and maybe ask to change it?)
		projectKey = projectKeys[0]
	} else {
		var err error
		boardIDs, idErr := client.GetBoardID(*boardName)
		if idErr != nil {
			log.Printf(foreground.RED + "[ERROR] getting board ID: %v" + decorators.RESET_ALL, idErr)
			return
		}

		*boardID = boardIDs.ID
		*boardName = boardIDs.Name
		fmt.Println("Board ID:", *boardID)
		fmt.Println("Board Name:", *boardName)

		projectKeys, err = client.GetProjectKeys(*boardID, projectKeys)
		if err != nil {
			log.Printf(foreground.RED + "[ERROR] getting project keys: %v" + decorators.RESET_ALL, err)
			return
		}

		// needed for JQL at the very least if not for anything else...
		fmt.Println("Associated Project Keys: -")
		for i, key := range projectKeys {
			fmt.Printf("Key %d: %s\n", i+1, key)
		}
		// probably select the first entry as the default key for the session...(and maybe ask to change it?)
		projectKey = projectKeys[0]
	}

	activeSprintID, sprintName, err := client.GetActiveSprintID(*boardID)
	if err != nil {
		log.Printf(foreground.RED + "[ERROR] getting active sprint ID: %v" + decorators.RESET_ALL, err)
		return
	}
	fmt.Printf("Active Sprint: %s (ID: %d)\n", sprintName, activeSprintID)

	fmt.Println("Your issues in the active sprint:-")

	// change the first parameter to some decided name later
	myIssues, err := client.GetMyIssuesInActiveSprint(projectKey, sprintName, user.Name)
	if err != nil {
		log.Printf(foreground.RED + "[ERROR] getting issues in active sprint: %v" + decorators.RESET_ALL, err)
		return
	}

	var totalStoryPoints float32 = 0.0

	for _, issue := range myIssues {
		totalStoryPoints += issue.Fields.StoryPoints
		fmt.Printf("%s\n%s : %.f\n\t%s\n", issue.Key, issue.Fields.Summary, issue.Fields.StoryPoints, issue.Fields.Status.Name)
		fmt.Printf("--------------------------------\n\n")
	}
	fmt.Printf("Total Story Points assigned to you in this sprint: %.f\n\n", totalStoryPoints)
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
			Project: &Project{
				Key: projectKey,
			},
			IssueType: &IssueType{
				Name: issueType,
			},
			Summary:     summary,
			Description: description,
		},
	}

	result, err := client.CreateIssue(issue)
	if err != nil {
		log.Printf(foreground.RED + "[ERROR] creating issue: %v" + decorators.RESET_ALL, err)
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
		log.Printf(foreground.RED + "[ERROR] getting issue: %v" + decorators.RESET_ALL, err)
		return
	}

	fmt.Printf("✅ Issue retrieved successfully!\n")
	fmt.Printf("Key: %s\n", issue.Key)
	fmt.Printf("ID: %s\n", issue.ID)
	fmt.Printf("Summary: %s\n", issue.Fields.Summary)
	fmt.Printf("Description: %s\n", issue.Fields.Description)
	fmt.Printf("Issue Type: %s\n", issue.Fields.IssueType.Name)
	fmt.Printf("Assignee: %s\n", issue.Fields.Assignee.DisplayName)
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

	fmt.Print("New Assignee ID (leave empty to keep current): ")
	scanner.Scan()
	assignee := strings.TrimSpace(scanner.Text())

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
			fmt.Printf(foreground.RED + "[ERROR] Invalid story points value: %v\n" + decorators.RESET_ALL, err)
		} else {
			fields.StoryPoints = float32(sp)
		}
	}
	if assignee != "" {
		tmpAssignee := &Assignee{}
		tmpAssignee.Name = assignee

		if err := client.UpdateAssignee(issueIDOrKey, tmpAssignee); err != nil {
			log.Printf(foreground.RED + "[ERROR] updating assignee: %v" + decorators.RESET_ALL, err)
			return
		}
		fmt.Printf("✅ Issue %s assigned to %s successfully!\n", issueIDOrKey, assignee)
	}

	if fields.Summary == "" && fields.Description == "" && fields.AcceptanceCriteria == "" && fields.StoryPoints <= 0.0 {
		fmt.Println("No changes specified.")
		return
	}

	err := client.UpdateIssue(issueIDOrKey, fields)
	if err != nil {
		log.Printf(foreground.RED + "[ERROR] updating issue: %v" + decorators.RESET_ALL, err)
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
		log.Printf(foreground.RED + "[ERROR] fetching transitions: %v" + decorators.RESET_ALL, err)
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
		fmt.Println(foreground.RED + "[ERROR] Invalid choice." + decorators.RESET_ALL)
		return
	}

	selectedTransition := transitions[choice-1]

	err = client.DoTransition(issueIDOrKey, selectedTransition.ID)
	if err != nil {
		log.Printf(foreground.RED + "[ERROR] performing transition: %v" + decorators.RESET_ALL, err)
		return
	}

	fmt.Printf("✅ Issue %s transitioned to '%s' successfully!\n", issueIDOrKey, selectedTransition.Name)
}

func getCommentsInteractive(client *JiraClient, scanner *bufio.Scanner) {
	fmt.Println("\n--- Get Comments ---")

	fmt.Print("Issue ID or Key: ")
	scanner.Scan()
	issueIDOrKey := strings.TrimSpace(scanner.Text())

	comments, err := client.GetComments(issueIDOrKey)
	if err != nil {
		log.Printf(foreground.RED + "[ERROR] getting comments: %v" + decorators.RESET_ALL, err)
		return
	}

	if len(comments) == 0 {
		fmt.Println("No comments found for this issue.")
		return
	}

	fmt.Printf("✅ Comments retrieved successfully! Total: %d\n", len(comments))
	for _, c := range comments {
		fmt.Printf("\nCommentID %s\n", c.ID)
		fmt.Printf("Author: %s\n", c.Author)
		fmt.Printf("Created: %s\n", c.Created)
		fmt.Printf("Last Updated: %s\n", c.LastUpdated)
		//          Last Updated: 2025-09-08T11:18:04.666+0200 -> longest field
		fmt.Printf("------------------------------------------\n%s\n\n", c.Body)
	}
}

// I need to later move it out of the options menu, create a flag while starting the application...this would populate the cache as well as get the active sprint tasks for only me (or a specific user: [TODO] later)

// func getActiveSprintTasksInteractive(client *JiraClient, scanner *bufio.Scanner) {
// 	fmt.Println("\n--- Get Active Sprint Tasks ---")

// 	fmt.Print("Board ID: ")
// 	scanner.Scan()
// 	boardIDStr := strings.TrimSpace(scanner.Text())
// 	boardID, err := strconv.Atoi(boardIDStr)
// 	if err != nil {
// 		fmt.Printf("Invalid Board ID: %v\n", err)
// 		return
// 	}

// 	sprintID, sprintName, err := client.GetActiveSprintID(boardID)
// 	if err != nil {
// 		log.Printf("Error getting active sprint: %v", err)
// 		return
// 	}
// 	if sprintID == 0 {
// 		fmt.Println("No active sprint found for this board.")
// 		return
// 	}

// 	fmt.Printf("Active Sprint: %s (ID: %d)\n", sprintName, sprintID)

// 	projectKey := "dummy"
// 	sprintID := "dummy"
// 	assignee := "dummy"
// 	issues, err := client.GetMyIssuesInActiveSprint(projectKey, sprintID, assignee)
// 	if err != nil {
// 		log.Printf("Error getting issues in sprint: %v", err)
// 		return
// 	}

// 	if len(issues) == 0 {
// 		fmt.Println("No issues found in the active sprint.")
// 		return
// 	}

// 	fmt.Printf("✅ Issues in Active Sprint retrieved successfully! Total: %d\n", len(issues))
// 	for _, issue := range issues {
// 		fmt.Printf("\nKey: %s\n", issue.Key)
// 		fmt.Printf("ID: %s\n", issue.ID)
// 		fmt.Printf("Summary: %s\n", issue.Fields.Summary)
// 		if issue.Fields.Status != nil {
// 			fmt.Printf("Status: %s\n", issue.Fields.Status.Name)
// 		}
// 		if issue.Fields.Assignee != nil {
// 			fmt.Printf("Assignee: %s\n", issue.Fields.Assignee.DisplayName)
// 		} else {
// 			fmt.Printf("Assignee: Unassigned\n")
// 		}
// 	}
// }
