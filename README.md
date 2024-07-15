# Project Title

This is a simple Go project that provides URL shortening services. The project uses the Echo framework and an K-V storage for URL mapping.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

- Go (>= v1.22.1)
- Docker (latest version)

### Running the App Locally

```bash
make run
```

### Running the App Locally using Docker

```bash
make dc
```

### Running App tests

```bash
make test
```

### Running Linter

```bash
make lint
```

### Environment Variables
The application uses the following environment variables, which can be set in the .env file:  

**HTTP_ADDRESS**: The address where the HTTP server will listen. **Default is :0 (random)**.

**LOGGER_LEVEL**: The level of logging. Can be debug, info, warn, error, fatal, panic. **Default is debug**.

**LOGGER_FORMAT**: The format of the logs. Can be console or json. **Default is console**.

**LEVEL_DB_PATH**: The path where LevelDB will store its data. **Default is /tmp**.

## Endpoints

| Method | Endpoint               | Description                                                                                  | Request Parameters                                                                 | Response                                                                                               |
|--------|-------------------------|----------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------|
| POST   | `/api/v1/shorten`       | Accepts a long URL as input and returns a shortened URL.                                      | **Body:** `{ "url": "https://www.example.com" }`                                    | **201:** `{ "destination": "http://localhost:8080/link/abc123" }` <br> **400:** `{ "message": "error" }`|
| GET    | `/api/v1/shorten`       | Accepts a long URL as a query parameter and returns a shortened URL.                          | **Query:** `url=https://www.example.com`                                            | **200:** `{ "destination": "http://localhost:8080/link/abc123" }` <br> **400:** `{ "message": "error" }` <br> **404:** `{ "message": "error" }`|
| GET    | `/api/v1/original`      | Accepts a shortened URL as a query parameter and returns the original long URL.               | **Query:** `url=http://localhost:8080/link/abc123`                                  | **200:** `{ "destination": "https://www.example.com" }` <br> **400:** `{ "message": "error" }` <br> **404:** `{ "message": "error" }`|
| GET    | `/link/{shortID}`       | Accepts a short ID as a path parameter and redirects to the original long URL.                | **Path:** `{shortID}`                                                               | **301:** Redirect with `Location` header <br> **400:** `{ "message": "error" }` <br> **404:** `{ "message": "error" }`|
