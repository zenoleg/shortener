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
