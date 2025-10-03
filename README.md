# Moviepin

Moviepin is a microservices-based application for managing movies and user reviews. It consists of two services: a movie management service and a token management service.

I made this project to learn building microservices with Go. This was my first time building a microservices application with Go, Gin, gRPC, and PostgreSQL.

## Architecture

The application is designed with a microservices architecture. The [`movie-management-service`](./movie-management-service/) acts as the main entry point for clients and communicates with the [`token-management-service`](./token-management-service/) for user authentication via gRPC. The `movie-management-service` is a stateful service that connects to a PostgreSQL database, while the `token-management-service` is a stateless service.

## Services

### Movie Management Service

The movie management service is a RESTful API. It provides endpoints for managing movies, including creating, reading, updating, and deleting movies. It also handles user registration and login.

### Token Management Service

The token management service is a gRPC service. It is responsible for generating and verifying JSON Web Tokens (JWTs) for user authentication.

## Getting Started

To run the application, you need to have Go and PostgreSQL installed.

1. **Clone the repository:**

    ```bash
    git clone https://github.com/ranmerc/moviepin
    ```

2. **Set up the database:**

    ```bash
    # Create the database
    make create-db

    # Migrate the database (create tables)
    make migrate-up
    ```

3. Export environment variables:

    Check the `.env.example` file for the list of environment variables that need to be set.

4. **Run the services:**

    * **Token Management Service:**

        ```bash
        cd token-management-service
        go run main.go
        ```

    * **Movie Management Service:**

        ```bash
        cd movie-management-service
        go run main.go
        ```

## API Endpoints

The API endpoints are defined in the Bruno collection [present here](./movie-management-service/Moviepin_Collection/). You can use the Bruno client to test the API.
