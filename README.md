# Phone Book API

## Overview
The Phone Book API is a simple web server API that allows users to manage their contacts. It provides endpoints to add, edit, delete, search, and retrieve contacts. The API is built using Golang and follows best practices for clean code, SOLID principles, and error handling.

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
│   │   ├── repository.go      # Data access layer for contacts
│   │   └── service.go        # Business logic for handling contacts
│   ├── database
│   │   └── db.go             # Database connection and setup
│   ├── metrics
│   │   └── metrics.go        # Metrics collection for monitoring
│   └── router
│       └── router.go         # API routing setup
├── test
│   └── contacts_test.go      # Unit tests for contact functionality
├── Dockerfile                 # Instructions for building the Docker image
├── docker-compose.yml         # Docker Compose configuration
├── go.mod                     # Go module definition
└── go.sum                     # Checksums for module dependencies
```

## Setup Instructions
1. **Clone the repository:**
   ```
   git clone <repository-url>
   cd phone-book-api
   ```

2. **Build and run the application using Docker:**
   ```
   docker-compose up --build
   ```

3. **Access the API:**
   The API will be available at `http://localhost:8080`.

## API Documentation
### Endpoints
- **GET /contacts**: Retrieve a list of contacts (supports pagination).
- **POST /contacts**: Add a new contact.
- **PUT /contacts/{id}**: Edit an existing contact.
- **DELETE /contacts/{id}**: Delete a contact.
- **GET /contacts/search**: Search for a contact by name or phone number.

## Testing
To run the tests, use the following command:
```
go test ./test
```

## Metrics
The application includes metrics collection to monitor API usage and performance. Metrics can be accessed through the designated endpoint.

## License
This project is licensed under the MIT License. See the LICENSE file for details.