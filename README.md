# Backend
Backend Repo

## Local testing
### Postgres config:
```cmd
podman run -d \
  --name postgres-api \
  -e POSTGRES_USER=api_user \
  -e POSTGRES_PASSWORD=test \
  -e POSTGRES_DB=bestread \
  -p 5432:5432 \
  postgres:17.5-alpine3.22
```

### Example config.yaml
```yaml
DB_HOST: "localhost"
DB_USERNAME: "api_user"
DB_PASSWORD: "test"
DB_NAME: "bestread"
DB_PORT: "5432"
DB_SSL_MODE: false
API_PORT: "3000"
DEBUG_LEVEL: "debug"
API_BASE_PATH: "/api"
```

## doc

### Posts

GET:

```bash
curl -v -X GET "http://localhost:3000/api/v1/post" -H "Content-Type: application/json" -d "{\"uid\":1,\"bid\":1}"
```

returned json objs mit den reviews
```json
[
   {
      "Pfp":"",
      "Username":"",
      "Book":{
         "ID":2,
         "ISBN":"978-0-439-02341-2",
         "Title":"test1",
         "Author":"test1",
         "CoverURL":"",
         "RatingAvg":0,
         "Description":"test1",
         "ReleaseDate":1
      },
      "Content":"awdawdawd",
      "Image":"AAAAAAAAAAAAA"
   }
]
```

--- 

POST:

```bash
id=1
curl -v -X POST "http://localhost:3000/api/v1/user/$id/post" -H "Content-Type: application/json" -d "{\"bid\": 1,\"content\": \"awdawdawd\",\"b64image\":\"AAAAAAAAAAAAA\"}"
```
returned nichts bis auf 200 wenn es funktioniert hat
