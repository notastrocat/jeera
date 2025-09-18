## ToDo checklist

- [x] the app displays the username - which in my case at least is dumb since I use PAT always. But, this could be used to fetch the actual name of the person in JIRA instance & then display it there. This would then serve as a basic check whether the JIRA instance is responsive or not.
- [x] an argument to only GET the issues in active sprint for a particular user.
  - [x] need to fetch the board & some related data.
    - [x] this is part of the /rest/agile/1.0 API as opposed to the traditional /rest/api/2/ we've been using till now.
    - [x] I'm now able to fetch the board & it's data. Parse it. Even get the board ID from board name supplied by user (searching using the same APIs different endpoint)
  - [ ] This changes everything, I am now thinking of caching a few fields especially the board name, ID, active sprint name, ID, project name, ID a few other fields. This would in turn make the queries *somewhat* faster at the very least.
    - [ ] fetch this data during the initialization of the application...
- [ ] functionality to move newly created issues to (active) sprint.

---

- [x] need to fix transition based on *required* fields.
    - [x] handle updates to acceptance criteria + story points
- [x] a function for assignee
- [ ] function(s) for comments - these are part of each issue.
  - [x] getting the comments
  - [ ] updating the comments
  - [ ] deleting the comments 
- [ ] handle linked issues
- [ ] argument: filename for a CSV file. this file would then be parsed to create issues in bulk; the general flow of application won't start wherein it asks for user input. This will only be used to "create" new issues in bulk.
- [ ] Refactor the code to break it into multiple files - it's easy to manage that way.
- [ ] update README.md
- [ ] Tests
