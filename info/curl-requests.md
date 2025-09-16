# Sample cURL requests I used to verify the endpoints

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
