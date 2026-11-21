# jaik-io-api
Used for educational demonstrations for courses being taught at Spokane Colleges.

Published at **https://api.jaik.io**. One Go service; each endpoint is a route in `cmd/api/routes.go`.

| Route | Methods |
|---|---|
| `/healthz` | GET |
| `/books`, `/books/{id}` | GET, POST / GET, PUT, DELETE (PostgreSQL) |

Run locally after creating the `books` table and setting the connection string:

```powershell
$env:READINGLIST_DB_DSN = 'postgres://USER:PASSWORD@HOST:5432/readinglist?sslmode=require'
psql "$env:READINGLIST_DB_DSN" -c "CREATE TABLE books (
  id bigserial PRIMARY KEY,
  created_at timestamptz NOT NULL DEFAULT now(),
  title text NOT NULL,
  published integer NOT NULL DEFAULT 0,
  pages integer NOT NULL DEFAULT 0,
  genres text[] NOT NULL DEFAULT '{}',
  rating real NOT NULL DEFAULT 0,
  version integer NOT NULL DEFAULT 1
);"
go run ./cmd/api
```

The API reads and writes books in PostgreSQL. Keep the DSN in an environment variable or secret
manager, not in source control. After starting the API, run `curl.exe http://localhost:8080/books`.
Tests: `go test ./...`.

Deploy: merge to `develop` → CI builds and pushes the image to ECR → Deploy applies `k8s/` to the `jaik-llc` EKS cluster
behind the shared ALB. Platform pieces (ECR repo, TLS certificate, deploy-role trust) live in `jaik-llc-infra`.
