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
      "Uid": "2",
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
      "ImageUrl":"https://<base>/api/v1/<todo>"
   }
]
```

wenn man es ohne body aufruft, kriegt man die neuesten posts im allgemeinen

--- 

POST:

**DAS WIRD NOCH GEÄNDERT**

```bash
id=1
curl -v -X POST "http://localhost:3000/api/v1/user/$id/post" -H "Content-Type: application/json" -d "{\"bid\": 1,\"content\": \"awdawdawd\",\"b64image\":\"AAAAAAAAAAAAA\"}"
```
returned nichts bis auf 200 wenn es funktioniert hat


### library

**/user/\<id\>/lib**

GET:

holt die bücher von einem nutzer
```bash
id=1
curl -v -X GET "http://localhost:3000/api/v1/user/$id/lib?limit=10" -H "Content-Type: application/json"
```
man hat den parameter `limit`, womit man die anzahl kontrollieren kann. default `1` (wenn es leer gelassen wird)

returned ein array:
```json
[
   {
      "Uid":2,
      "Book":{
         "ID":1,
         "ISBN":"978-0-439-02340-2",
         "Title":"test0",
         "Author":"test0",
         "CoverURL":"",
         "RatingAvg":0,
         "Description":"test0",
         "ReleaseDate":0
      },
      "State":0,
      "Rating":0
   }
]
```

---

POST:

fügt ein buch dem nutzer hinzu

state geht von 0 - 2: 0 = Want to read, 1 = Reading, 2 = Read
```bash
id=1
curl -v -X POST "http://localhost:3000/api/v1/user/$id/lib" -H "Content-Type: application/json" -d "{\"bid\": 1,\"state\": 0}"
```

returned ok wenn alles in ordnung war


**/user/\<id\>/lib/\<id\>**

PUT:

zustand updaten
```bash
id=1
bid=1
curl -v -X PUT "http://localhost:3000/api/v1/user/$id/lib/$bid" -H "Content-Type: application/json" -d "{\"state\": 1}"
```

returned ok wenn alles in ordnung war

---

DELETE:

buch aus der bibliothek löschen
```bash
id=1
bid=1
curl -v -X DELETE "http://localhost:3000/api/v1/user/$id/lib/$bid" -H "Content-Type: application/json" 
```

returned ok wenn alles in ordnung war
