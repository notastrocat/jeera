# JIRA Auto - Go Application

A Go application for interacting with JIRA issues through the REST API. This tool allows you to create, retrieve, and update JIRA issues from the command line.

## Features

- **Create Issues**: Create new JIRA issues with summary, description, project, and issue type
- **Retrieve Issues**: Get detailed information about existing issues by ID or key
- **Update Issues**: Modify existing issues (summary and description)
- **Interactive CLI**: User-friendly command-line interface
- **Secure Authentication**: Uses Basic Auth with API tokens

## Prerequisites

- Go 1.16 or later
- JIRA instance with API access
- JIRA API token (recommended over password)

## Configuration

The application supports configuration through both `.env` files and environment variables:

### Option 1: Using .env file (Recommended)

1. Copy the example file:
```bash
cp .env.example .env
```

2. Edit the `.env` file with your actual values:
```bash
JIRA_BASE_URL=https://yourcompany.atlassian.net
JIRA_USERNAME=your.email@company.com
JIRA_API_TOKEN=your-api-token-here
```

### Option 2: Using Environment Variables

```bash
export JIRA_BASE_URL=https://yourcompany.atlassian.net
export JIRA_USERNAME=your.email@company.com
export JIRA_API_TOKEN=your-api-token
```

The application will first try to load from `.env` file, then fall back to environment variables.

### Getting a JIRA API Token

1. Go to your JIRA account settings
2. Navigate to Security → API tokens
3. Create a new API token
4. Copy the token and use it as `JIRA_API_TOKEN`

## Installation

1. Clone or download this repository
2. Navigate to the project directory
3. Build the application:

```bash
go build -o jira-auto
```

## Usage

1. Set your environment variables (see Configuration section)
2. Run the application:

```bash
./jira-auto
```

3. Follow the interactive prompts to:
   - Create new issues
   - Retrieve existing issues
   - Update issues

## Project Structure

```
jira-auto/
├── main.go      # Entry point with interactive CLI
├── jira.go      # JIRA API client and functions
├── config.go    # Configuration management
├── go.mod       # Go module file
└── README.md    # This file
```

## API Functions

### CreateIssue
- **Endpoint**: POST `/rest/api/2/issue`
- **Purpose**: Creates a new JIRA issue
- **Required fields**: Project key, issue type, summary

### GetIssue
- **Endpoint**: GET `/rest/api/2/issue/{issueIdOrKey}`
- **Purpose**: Retrieves an existing issue
- **Input**: Issue ID or key

### UpdateIssue
- **Endpoint**: PUT `/rest/api/2/issue/{issueIdOrKey}`
- **Purpose**: Updates an existing issue
- **Updatable fields**: Summary, description

## Error Handling

The application includes comprehensive error handling for:
- Missing configuration
- Network errors
- API authentication failures
- Invalid issue keys/IDs
- Malformed requests

## Example Session

```
JIRA Auto - Issue Management Tool
=================================
Connected to: https://yourcompany.atlassian.net
Username: your.email@company.com

Available commands:
  1. Create issue
  2. Get issue
  3. Update issue
  4. Exit

Enter your choice (1-4): 1

--- Create New Issue ---
Project Key: PROJ
Issue Type: Bug
Summary: Fix login issue
Description: Users cannot login with SSO
✅ Issue created successfully!
Key: PROJ-123
ID: 10001
```

## Security Notes

- Never hardcode credentials in the source code
- Use environment variables or secure configuration files
- API tokens are preferred over passwords
- Consider using OAuth for production deployments

## Contributing

Feel free to submit issues and enhancement requests!

## License

This project is open source and available under the MIT License.
