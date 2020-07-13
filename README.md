# cGF3ZWwtY3plcndpZW5pZWMK
### fetcher-api

#### Endpoints
- GET /api/fetcher
-Show information about all saved links

- GET /api/fetcher/{id}
-Show information about specific link

- GET /api/fetcher/{id}/history
-Show history of scrape for the link

- DELETE /api/fetcher/{id}
-Deletes saved link

- POST /api/fetcher
`{url":"https://httpbin.org/6/delay","interval":10}`
-Create new or update existing link

#### Worker
###### Scrape saved links with the given interval
- Limit of 5 seconds for the response