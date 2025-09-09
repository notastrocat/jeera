# JIRA Auto - Go Application Project

## Project Overview
This is a Go application for JIRA API integration with functionality to create, retrieve, and update issues.

## Project Structure
- `main.go`: Entry point with CLI interface
- `jira.go`: JIRA API interaction functions
- `config.go`: Configuration management
- `go.mod`: Go module file

## Progress Checklist
- [x] Verify copilot-instructions.md file creation
- [x] Clarify Project Requirements - Go JIRA API application
- [x] Scaffold the Project - Created go.mod, main.go, jira.go, config.go
- [x] Customize the Project - Implemented JIRA API client with create/get/update functionality
- [x] Install Required Extensions - No extensions required for Go development
- [x] Compile the Project - Successfully built with `go build`
- [x] Create and Run Task - Created VS Code tasks for building and running
- [x] Launch the Project - Ready to launch with proper configuration
- [x] Ensure Documentation is Complete - README.md and examples created

## Authentication
Uses Basic Auth with username/API token for JIRA authentication.

## Main Features
- Create JIRA issues (POST /rest/api/2/issue)
- Retrieve JIRA issues (GET /rest/api/2/issue/{issueIdOrKey})
- Update JIRA issues (PUT /rest/api/2/issue/{issueIdOrKey})
