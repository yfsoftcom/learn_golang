## Short Link Project

it generates short link for url

### API

- GET  `/:shortlink`
- POST `/api/shortlink` `{url: string}` -> `{shortlink: string, exp: int}`
- GET `/api/info?shortlink=xxxxxxxx` -> `{url: string, createAt: datetime, exp: int}`

### Store

Redis

### Redis Keys

- `url.next.id` int
- `shortlink:%s:url` string
- `shortlink:%s:detail` map

### Run Redis Server

`$ docker run -p 6379:6379 --name redisserver -d redis:alpine3.11`