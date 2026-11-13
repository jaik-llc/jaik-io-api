# jaik-io-api
Used for educational demonstrations for courses being taught at Spokane Colleges.

Published at **https://api.jaik.io**. One Go service; each endpoint is a route in `cmd/api/routes.go`.

| Route | Methods |
|---|---|
| `/healthz` | GET |
| `/books`, `/books/{id}` | GET, POST / GET, PUT, DELETE (sample data, no database yet) |

Run locally: `go run ./cmd/api` then `curl localhost:8080/books`. Tests: `go test ./...`.

Deploy: merge to `develop` → CI builds and pushes the image to ECR → Deploy applies `k8s/` to the `jaik-llc` EKS cluster
behind the shared ALB. Platform pieces (ECR repo, TLS certificate, deploy-role trust) live in `jaik-llc-infra`.
