## IssueKey/editmeta endpoint
```
"customfield_10002": {
  "required": false,
  "schema": {
    "type": "number",
    "custom": "com.atlassian.jira.plugin.system.customfieldtypes:float",
    "customId": 10002
  },
  "name": "Story Points",
  "fieldId": "customfield_10002",
  "operations": [
    "set"
  ]
}
```

```
"customfield_11028": {
      "required": false,
      "schema": {
        "type": "string",
        "custom": "com.atlassian.jira.plugin.system.customfieldtypes:textarea",
        "customId": 11028
      },
      "name": "Acceptance Criteria",
      "fieldId": "customfield_11028",
      "operations": [
        "set"
      ]
    }
```

`components.allowedValues -> ` this has a member called `self` which is another URL (but also has a *name* member) - so maybe cache it?

`versions.allowedValues -> ` this has a list of all the release versions planned - this one grows dynamically I believe. [self, name]

`fixVersions.allowedValues -> ` similar to the one above. not sure which one is actually being used? [self, name]

`customfield_15400.allowedValues -> ` a list of PIs. each PI further has a `children`: array of it's sprints. [self, value]

`issuelinks`: probably an array?

`assignee`: has something called *autoCompleteUrl*. Can it be used as a query to assign it to some user?

`issuetype`: has it's *operations* array empty. maybe I cannot change the issue type via this endpoint?

## IssueKey/comments

## Board

* this is part of the agile/1.0 API as opposed to the traditional api/2.

```
{
  "id": 14190,
  "self": "https://jiraprod.aptiv.com/rest/agile/1.0/board/14190",
  "name": "CoreFW_AutoScrum",
  "type": "scrum"
}
```

## `agile/1.0/board/{boardID}/sprint`

* This returns all the sprints for a (scrum/kanban) board.
* Which can then be used to find out the active sprint.
