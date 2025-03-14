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

## License
This project is licensed under the MIT License. See the LICENSE file for details.