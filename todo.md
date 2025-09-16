## ToDo checklist

- [ ] the app displays the username - which in my case at least is dumb since I use PAT always. But, this could be used to fetch the actual name of the in JIRA instance & then display it there. This would then serve as a basic check whether the JIRA instance is responsive or not.
- [ ] an argument to only GET the issues in active sprint for a particular user.
- [x] need to fix transition based on *required* fields.
    - [x] handle updates to acceptance criteria + story points
- [x] a function for assignee
- [ ] function(s) for comments - these are part of each issue.
  - [x] getting the comments
  - [ ] updating the comments
  - [ ] deleting the comments 
- [ ] handle linked issues
- [ ] arugument: filename for a CSV file. this file would then be parsed to create issues in bulk; the general flow of application won't start wherein it asks for user input. This will only be used to "create" new issues in bulk.
