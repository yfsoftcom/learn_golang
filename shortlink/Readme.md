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