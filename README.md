# Microgate

This project is a **scalable user and task management system** built using a **microservices architecture**.  
Each service is designed to be developed, tested, and deployed independently, ensuring clear separation of concerns and maintainability.  

The system emphasises two key aspects:  
1. It is a **microservices-based design**, where each component functions independently.  
2. All **inter-service communication** is carried out over **gRPC using Protobuf**, providing efficient, type-safe, and reliable communication between services.  

Users interact with the system via HTTP requests in JSON format through the API Gateway, which forwards authentication and task management requests to the relevant microservices over gRPC.  

## Architecture Overview
![alt text](architecture-migrogate.png)
The system consists of three main microservices:

1. **API Gateway Microservice**  
   - Receives user requests in HTTP/JSON format.  
   - Communicates with the **Auth Microservice** via gRPC for user authentication and registration.
   - Communicates with the **Task Microservice** via gRPC to create, list, and manage tasks.  

2. **Auth Microservice**  
   - Stores user information in PostgreSQL.  
   - Maintains access tokens in Redis.  
   - Handles user registration, login, token generation, and validation.  

3. **Task Microservice**  
   - Stores user tasks in MongoDB.  
   - Handles task CRUD.  



## Table of Contents

- [Features](#features)
- [File Structure](#structure)
- [Run Project Without Docker](#bare)
- [Run Project With Docker](#docker)
- [Usage Examples](#usage)
- [License](#license)


## Features <a name="features"></a>
    System Monitoring:
        Structured Logging
        Metrics
        Alerts
        Healthcheck
        Tracing
        Error Handling

    Test:
        Unit Tests
        Integration Tests
        CI

    Security:
        CORS avoidance
        Rate limiting

    Authentication/Authorization:
        Stateful token-based approach for secure authentication and authorization

    Other:
        Continuous integration (CI)
        Docker (Containerization)


## Project Layout <a name="structure"></a>
``` shell
.
├── Makefile
├── README.md
├── cmd
│   ├── auth
│   │   └── main.go
│   ├── gateway
│   │   └── main.go
│   └── task
│       └── main.go
├── deployments
│   ├── compose
│   │   ├── docker-compose-test.yaml
│   │   └── docker-compose.yaml
│   └── docker
│       ├── auth-service.Dockerfile
│       ├── gateway-service.Dockerfile
│       ├── integration-test.Dockerfile
│       └── task-service.Dockerfile
├── go.mod
├── go.sum
├── integration_tests
│   └── user_flow_test.go
├── internal
│   ├── auth
│   │   ├── auth.go
│   │   ├── cache
│   │   │   ├── contract
│   │   │   │   └── interface.go
│   │   │   ├── factory.go
│   │   │   ├── mock
│   │   │   │   ├── mock_cache.go
│   │   │   │   └── mock_cache_test.go
│   │   │   └── redis
│   │   │       └── redis.go
│   │   ├── db
│   │   │   ├── contract
│   │   │   │   └── interface.go
│   │   │   ├── factory.go
│   │   │   ├── migrations
│   │   │   │   ├── 000001_create_user_table.down.sql
│   │   │   │   └── 000001_create_user_table.up.sql
│   │   │   ├── mock
│   │   │   │   ├── mock_db.go
│   │   │   │   └── mock_db_test.go
│   │   │   └── postgres
│   │   │       └── postgres.go
│   │   ├── dto.go
│   │   ├── endpoints.go
│   │   ├── grpc.go
│   │   ├── pb
│   │   │   ├── auth-service.pb.go
│   │   │   ├── auth-service.proto
│   │   │   └── auth-service_grpc.pb.go
│   │   ├── repo
│   │   │   ├── contract
│   │   │   │   └── interface.go
│   │   │   ├── entity
│   │   │   │   ├── anonymous_user.go
│   │   │   │   ├── err.go
│   │   │   │   └── user.go
│   │   │   ├── factory.go
│   │   │   ├── mock
│   │   │   │   ├── mock_user.go
│   │   │   │   └── mock_user_test.go
│   │   │   └── user
│   │   │       ├── user.go
│   │   │       └── user_validation.go
│   │   └── service.go
│   ├── gateway
│   │   ├── client
│   │   │   ├── auth.go
│   │   │   └── task.go
│   │   ├── endpoints
│   │   │   ├── dto.go
│   │   │   └── endpoints.go
│   │   ├── gateway.go
│   │   ├── service.go
│   │   └── transport
│   │       ├── http.go
│   │       └── middleware.go
│   └── task
│       ├── db
│       │   ├── contract
│       │   │   └── interface.go
│       │   ├── factory.go
│       │   ├── mock
│       │   │   ├── mock_db.go
│       │   │   └── mock_db_test.go
│       │   └── mongo
│       │       └── mongo.go
│       ├── dto.go
│       ├── endpoint.go
│       ├── grpc.go
│       ├── pb
│       │   ├── task-service.pb.go
│       │   ├── task-service.proto
│       │   └── task-service_grpc.pb.go
│       ├── repo
│       │   ├── contract
│       │   │   └── interface.go
│       │   ├── entity
│       │   │   └── task.go
│       │   ├── factory.go
│       │   ├── mock
│       │   │   ├── mock_task.go
│       │   │   └── mock_task_test.go
│       │   └── task
│       │       ├── task.go
│       │       └── validate_task.go
│       ├── service.go
│       └── task.go
├── pkg
│   ├── config
│   │   ├── config.go
│   │   ├── contract
│   │   │   └── interface.go
│   │   ├── factory.go
│   │   ├── loader
│   │   │   └── loader.go
│   │   └── mock
│   │       ├── mock.go
│   │       └── mock_test.go
│   ├── ctx
│   │   └── ctx.go
│   ├── err
│   │   └── err.go
│   ├── logger
│   │   ├── contract
│   │   │   └── interface.go
│   │   ├── factory.go
│   │   └── zaplogger
│   │       └── zaplogger.go
│   ├── token
│   │   └── token.go
│   └── validator
│       └── validator.go
└── service.dev.yaml

```

## Instructions
There are two options:
1) You can run the project without docker.
2) You can run the project with docker.

### 1. Run Project Without Docker<a name="bare"></a>

#### Requirements
- [Go](https://go.dev/dl/)
- [MongoDB](https://www.mongodb.com/docs/manual/installation/)
- [Redis](https://redis.io/download/)
- [PostgreSQL](https://www.postgresql.org/download/)
- [make](https://www.gnu.org/software/make/)
- [go-migrate](https://github.com/golang-migrate/migrate)
- [git](https://git-scm.com/downloads)
- Code Editor


#### Steps
- First, clone the repository:
```shell
git clone https://github.com/barantoraman/microgate.git
```
- Navigate to the microgate directory:
```shell
cd microgate
```
- Create a Postgres database named auth.
- Create a Postgres superuser named microgate.
(If you change any names, make sure to update the configurations in the service.dev.yaml file.)
- Update the service.dev.yaml configurations:
```shell
make config/local
```
- Clean up and synchronize the Go module dependencies:
```shell
go mod tidy 
```
- Define the environment variables for your database connection:
```shell
export AUTH_DB_DSN=YOUR_DATABASE_USER:YOUR_USER_PASS@/YOUR_DATABASE_NAME
```
- Apply the necessary database migrations:
```shell
make db/migrate/up/auth
```
- Run services:
```shell
make local/run/task
make local/run/auth
make local/run/gateway
```

- Also, if you want to test gRPC services (Auth and Task Service):
```shell
make test/task
make test/auth
make test/all
```

### 1. Run Project With Docker <a name="docker"></a>

#### Requirements

- [Docker](https://docs.docker.com/get-docker/)
- [make](https://www.gnu.org/software/make/)
- [make](https://sourceforge.net/projects/gnuwin32/files/make/3.81/make-3.81.exe/download?use_mirror=nav&download=) (for windows)
- [git](https://git-scm.com/downloads)

#### Steps

- First, clone the repository:
```shell
git clone https://github.com/barantoraman/microgate.git
```
- Navigate to the microgate directory:
```shell
cd microgate
```
- Build the containers based on the docker-compose.yaml configuration:
```shell
make docker/prod/build
```
- Start the containers:
```shell
make docker/prod/run
```

- If you want to run the Integration tests:
  - Begin by building the test container:
  ```shell
  make docker/test/build
  ```
  - Then, run the test container to execute the tests::
  ```shell
  make docker/test/run
  ```

### Usage Examples <a name="usage"></a>

- Signup Endpoint:
```bash
curl -X POST -d '{
    "user": {
        "email": "wallacegromit@mail.com",
        "password": "wallacegromitpass"
    }
}' localhost:8081/v1/signup
```
- Login Endpoint:
```bash
curl -X POST -d '{
    "user": {
        "email": "wallacegromit@mail.com",
        "password": "wallacegromitpass"
    }
}' localhost:8081/v1/login
```
- Logout Endpoint:
```bash
curl -X POST -d '{
    "token": {
        "plaintext": "...TOKEN..."
    }
}' localhost:8081/v1/logout
```
- Add Task Endpoint:
```bash
curl -X POST -d '{
    "task": {
        "title": "a grand day out",
        "description": "to the moon!",
        "status": "crackers"
    }
}' --header "Authorization: Bearer ...TOKEN..." localhost:8081/v1/task
```
- List Task Endpoint:
```bash
curl --header "Authorization: Bearer ...TOKEN..." \
localhost:8081/v1/task
```
- Delete Task Endpoint:
```bash
curl -X DELETE --header "Authorization: Bearer ...TOKEN..." localhost:8081/v1/task/{task_id}
```
