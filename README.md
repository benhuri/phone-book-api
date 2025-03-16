# Phone Book API

## Overview
The Phone Book API is a simple web server API that allows users to manage their contacts. It provides endpoints to add, edit, delete, search, and retrieve contacts. The API is built using Golang.

## Features
- Add a new contact
- Edit an existing contact
- Delete a contact
- Search for a contact
- Retrieve a list of contacts with pagination (maximum of 10 contacts per request)

## Project Structure
```
phone-book-api
├── cmd
│   └── main.go               # Entry point of the application
├── internal
│   ├── contacts
│   │   ├── handler.go        # HTTP handlers for contact-related API endpoints
│   │   ├── model.go          # Defines the Contact struct
│   │   ├── repository.go     # Data access layer for contacts
│   │   └── service.go        # Business logic for handling contacts
│   ├── config
│   │   └── config.go         # Configuration setup
│   ├── database
│   │   └── db.go             # Database connection and setup
│   ├── metrics
│   │   └── metrics.go        # Metrics collection for monitoring
│   └── router
│       └── router.go         # API routing setup
├── test
│   └── contacts_test.go      # Unit tests for contact functionality
├── Dockerfile                # Instructions for building the Docker image
├── docker-compose.yml        # Docker Compose configuration
├── go.mod                    # Go module definition
└── go.sum                    # Checksums for module dependencies
```

## Prerequisites
- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [Go](https://golang.org/dl/) (if you want to run the application locally without Docker)
- PostgreSQL database (if running locally without Docker)

## Setup Instructions
1. **Clone the repository:**
   ```sh
   git clone https://github.com/benhuri/phone-book-api
   cd phone-book-api
   ```

2. **Build and run the application using Docker:**
   ```sh
   docker-compose up --build
   ```

3. **Access the API:**
   The API will be available at `http://localhost:8080`.

4. **Run the application locally:**
   ```sh
   go run cmd/main.go
   ```

## Configuration

### User Credentials
User credentials and other configuration settings can be specified in the `config.yaml` file or as environment variables.

#### `config.yaml` File
The `config.yaml` file should be located in the root of the project directory. Here is an example configuration:

```yaml
DB_USER: postgres
DB_PASSWORD: "123"
DB_NAME: phonebook
DB_HOST: localhost
DB_PORT: 5432
```

#### Environment Variables
Alternatively, you can set the following environment variables:

- `DB_USER`: The database user (e.g., `postgres`)
- `DB_PASSWORD`: The database password (e.g., `123`)
- `DB_NAME`: The database name (e.g., `phonebook`)
- `DB_HOST`: The database host (e.g., `localhost`)
- `DB_PORT`: The database port (e.g., `5432`)

### Example of Setting Environment Variables

#### On Windows (Command Prompt)
```sh
set DB_USER=postgres
set DB_PASSWORD=123
set DB_NAME=phonebook
set DB_HOST=localhost
set DB_PORT=5432
go run cmd/main.go
```

#### On Windows (PowerShell)
```sh
$env:DB_USER="postgres"
$env:DB_PASSWORD="123"
$env:DB_NAME="phonebook"
$env:DB_HOST="localhost"
$env:DB_PORT="5432"
go run cmd/main.go
```

#### On Unix-based Systems (Linux, macOS)
```sh
export DB_USER=postgres
export DB_PASSWORD=123
export DB_NAME=phonebook
export DB_HOST=localhost
export DB_PORT=5432
go run cmd/main.go
```

## API Documentation

### Endpoints
- **GET /contacts**: Retrieve a list of contacts (supports pagination).
- **POST /contacts**: Add a new contact.
- **PUT /contacts/{id}**: Edit an existing contact.
- **DELETE /contacts/{id}**: Delete a contact.
- **GET /contacts/search**: Search for a contact by name or phone number.

### Validations
The following validations are applied to the contact fields:
- `first_name`: Required, minimum length of 1, maximum length of 50.
- `last_name`: Required, minimum length of 1, maximum length of 50.
- `phone_number`: Required, exactly 10 characters, numeric.
- `address`: Required, minimum length of 2, maximum length of 100.

### Example Requests

#### Add a New Contact
**Endpoint:** `POST /contacts`

**Request Body:**
```json
{
  "first_name": "John",
  "last_name": "Doe",
  "phone_number": "1234567890",
  "address": "123 Main St"
}
```

**Example Request:**
```sh
curl -X POST http://localhost:8080/contacts \
     -H "Content-Type: application/json" \
     -d '{
           "first_name": "John",
           "last_name": "Doe",
           "phone_number": "1234567890",
           "address": "123 Main St"
         }'
```

#### Retrieve Contacts
**Endpoint:** `GET /contacts`

**Query Parameters:**
- `page`: The page number to retrieve (default is 1).
- `limit`: The number of contacts per page (default is 10).

**Example Request:**
```sh
curl -X GET http://localhost:8080/contacts?page=1&limit=10
```

#### Edit a Contact
**Endpoint:** `PUT /contacts/{id}`

**Request Body:**
```json
{
  "first_name": "Jane",
  "last_name": "Doe",
  "phone_number": "0987654321",
  "address": "456 Elm St"
}
```

**Example Request:**
```sh
curl -X PUT http://localhost:8080/contacts/1 \
     -H "Content-Type: application/json" \
     -d '{
           "first_name": "Jane",
           "last_name": "Doe",
           "phone_number": "0987654321",
           "address": "456 Elm St"
         }'
```

#### Delete a Contact
**Endpoint:** `DELETE /contacts/{id}`

**Example Request:**
```sh
curl -X DELETE http://localhost:8080/contacts/1
```

#### Search for a Contact
**Endpoint:** `GET /contacts/search`

**Query Parameters:**
- `query`: The search query (e.g., name or phone number).

**Example Request:**
```sh
curl -X GET http://localhost:8080/contacts/search?query=John
```

## Testing
To run the tests, use the following command:
```sh
go test ./test
```

## Metrics
The application includes metrics collection to monitor API usage and performance. Metrics can be accessed through the designated endpoint.
http://localhost:8080/metrics