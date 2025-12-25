Role Based Authentication System (Go)
====================================

A role-based authentication demo written in Go with minimal dependencies.

---

## Project Overview

This project demonstrates:
- Role-based authentication
- Clean Go project layout (`cmd/`, `internal/`, `pkg/`)
- Containerized deployment using Docker
- Environment-based configuration

This version replaces the earlier Node.js implementation (See tag: `v1.0-nodejs`).

---

## API Documentation
- Hoppscotch collection: `docs/AuthDemo_Hoppscotch_Import.json`
- OpenAPI spec: `docs/openapi-specs.yaml`

### Swagger UI
Run this command to see the specs in [Swagger UI](http://localhost:8081)
```shell
docker run -p 8081:8080 -e SWAGGER_JSON=/docs/openapi-specs.yaml -v ./docs/openapi-specs.yaml:/docs/openapi-specs.yaml swaggerapi/swagger-ui
```

## Requirements

- Go 1.25+ (for local development)
- Docker (recommended)

---

## Running locally (without Docker)

```bash
go mod download
go run ./cmd/server
```

## Running with Docker (recommended)
Follow the documentation as given in https://nodejs.org/en/docs/guides/nodejs-docker-webapp/

 + Install docker application (https://www.docker.com/products/docker-desktop)
 + Create the Dockerfile and place it in the root of the project
 + `docker build . -t <your username>/auth-demo-go`
 + To see your images created by the previous step run `docker images`

 This will show like the following

|      REPOSITORY       |  TAG   |   IMAGE ID   |    CREATED    |  SIZE   |
|:---------------------:|:------:|:------------:|:-------------:|:-------:|
| kararnab/auth-demo-go | latest | c93113ff6c5c | 4 minutes ago | 66.77MB |

+ Run the docker image created using the following cmd
+ `docker run -p 49160:8080 -d --env-file .env <your username>/auth-demo-go`
+ Open the browser link http://localhost:49160

OR 
+ `docker run --publish 8080:8080 kararnab/role_auth_1.0`
+ Open the browser link http://localhost:8080

