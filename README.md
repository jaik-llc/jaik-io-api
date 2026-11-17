# jaik-io-api
Used for educational demonstrations for courses being taught at Spokane Colleges.

Published at **https://api.jaik.io**. One Go service; each endpoint is a route in `cmd/api/routes.go`.

| Route | Methods |
|---|---|
| `/healthz` | GET |
| `/books`, `/books/{id}` | GET, POST / GET, PUT, DELETE (PostgreSQL-backed CRUD operations)|

Run locally: `go run ./cmd/api` then `curl localhost:8080/books`. Tests: `go test ./...`.

Deploy: merge to `develop` → CI builds and pushes the image to ECR → Deploy applies `k8s/` to the `jaik-llc` EKS cluster
behind the shared ALB. Platform pieces (ECR repo, TLS certificate, deploy-role trust) live in `jaik-llc-infra`.
## Local PostgreSQL Database Setup

This application uses PostgreSQL to store and manage books.

### Database Setup

1. Install and start PostgreSQL.
2. Create a database named `jaik_api_dev`.
3. Apply the SQL migration:

```powershell
psql -h localhost -U postgres -d jaik_api_dev -f migrations/0001_create_books.sql
```

4. Configure the `READINGLIST_DB_DSN` environment variable with your local PostgreSQL connection string. Keep database credentials out of Git.
5. Start the application:

```powershell
go run ./cmd/api
```

The API runs on port 8080 by default.

### API Endpoints

| Method | Endpoint | Description |
|---|---|---|
| GET | /healthz | Application health |
| GET | /healthcheck/db | Database connectivity |
| GET | /books | List all books |
| POST | /books | Create a book |
| GET | /books/{id} | Retrieve a book |
| PUT | /books/{id} | Update a book |
| DELETE | /books/{id} | Delete a book |

Book records are stored in PostgreSQL rather than hardcoded sample data.
