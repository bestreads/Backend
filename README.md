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
