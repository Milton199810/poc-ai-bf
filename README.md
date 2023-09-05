# improved-octo-memory

## Installation

Install packages:

```sh
go mod download
```

Create a `.env` file:
```sh
cat << EOF > config/.env.local
SQL_SERVER_HOST=
SQL_SERVER_USER=
SQL_SERVER_PASSWORD=
SQL_SERVER_PORT=
SQL_SERVER_DATABASE=
OPENAI_API_KEY=
AZURE_OPENAI_KEY=
AZURE_OPENAI_ENDPOINT=
GCP_CREDENTIALS_PATH=
GCP_VERTEX_API_ENDPOINT=
GCP_VERTEX_PROJECT_ID=
GCP_VERTEX_MODEL_ID=
EOF
```

Add users:

```go
// internal/infra/web/middleware/basic_auth.go
...
func BasicAuthMiddleware(next http.Handler) http.Handler {
	users := map[string]string{
		"email": "password",
	}
...
```

## Usage

Run:
```sh
go run cmd/app/main.go
```

Access to:

- [localhost:8000/summary](localhost:8000/summary)

## Project status

In development