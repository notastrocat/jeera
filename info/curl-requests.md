# Sample cURL requests I used to verify the endpoints

```
curl -H "Authorization: Bearer PAT_TOKEN" https://jiraprod.aptiv.com/rest/api/2/user?username=d472pb | jq . > user.json
```

### Request for getting an issue & saving it to a JSON file
`curl -H "Authorization: Bearer PAT_TOKEN" https://jiraprod.aptiv.com/rest/api/2/issue/GTJ-687 | jq . > issue-687.json`

### Request for JQL searches & pretty print it to the console
```
curl -G \
    -H "Authorization: Bearer PAT_TOKEN"  \
    --data-urlencode 'jql=project = GTJ AND sprint = "25PI3 S6" AND assignee in (d472pb)' \
    'https://jiraprod.aptiv.com/rest/api/2/search' \
    -d 'fields=key,summary,status' -d 'maxResults=100' | jq .
```

### cURL request using the /agile/1.0/ API

- now, there are two approaches, either get the boardID from the board.values.name (will have to search using the below endpoint) OR get the boardID directly, which I highly doubt that people would know of beforehand...I too struggled to find boardID in JIRA GUI.

```
curl -H "Authorization: Bearer PAT_TOKEN" https://jiraprod.aptiv.com/rest/agile/1.0/board?name=CoreFW_AutoScrum
 | jq .
```

### cURL request to move newly created issues into any sprints...
```
curl --request POST \
  --url 'https://jiraprod.aptiv.com/rest/agile/1.0/sprint/102884/issue' \
  --header 'Authorization: Bearer PAT_TOKEN' \
  --header 'Content-Type: application/json' \
  --data '{
  "issues": [
    "GTJ-701"
  ]
}' | jq .
```
