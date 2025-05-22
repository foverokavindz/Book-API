# Book API

A simple RESTful API for managing a collection of books, built with Go.

## Features

- CRUD operations for books (Create, Read, Update, Delete)
- JSON data format
- In-memory data store
- RESTful API design
- Docker containerization
- Kubernetes deployment support

## Installation

### Prerequisites

- Go 1.16 or higher

### Setup

Run the application:

```
go run main.go
```

The server will start on port 8080 by default.

### Development with Hot-Reload

You can use [Fresh](https://github.com/gravityblast/fresh), a tool that automatically rebuilds and restarts your Go application when files change.

#### Installing Fresh

```bash
go install github.com/pilu/fresh@latest
```

#### Running the Application with Fresh

Navigate to your project directory and run:

```bash
fresh
```

## Running Tests

The project includes unit tests for all API endpoints.

### Installing Test Dependencies

```bash
go get github.com/stretchr/testify/assert
go get github.com/gorilla/mux
```

### Running All Tests

To run all tests in the project:

```bash
go test ./...
```

### Running Specific Test Files

To run tests in a specific file:

```bash
go test ./tests/book_test.go
```

### Test Structure

The tests are organized in the `tests` directory:

- `book_test.go`: Tests for CRUD operations on books
- `search_test.go`: Tests for book search functionality

## Docker Containerization

### Prerequisites

- Docker installed on your machine

### Steps to Containerize

Build the Docker image:

```bash
docker build -t book-api:latest .
```

Run the Docker container:

```bash
docker run -p 8080:8080 book-api:latest
```

4. Access the API at http://localhost:8080/books

## Kubernetes Deployment

### Environment

- Windows OS
- Minikube running with Docker driver

### Prerequisites

- Docker
- Minikube
- kubectl

### Setup Instructions

1. Install Minikube and start it:

```bash
minikube start
```

2. Configure Docker to use Minikube's Docker daemon:

```bash
minikube docker-env | Invoke-Expression
```

3. Build the Docker image:

```bash
docker build -t book-api:latest .
```

4. Apply Kubernetes configurations:

```bash
kubectl apply -f kubernetes/deployment.yaml
kubectl apply -f kubernetes/service.yaml
```

5. Access the API:

```bash
minikube service book-api-service
```

6. Use the tunnel URL (e.g., http://127.0.0.1:10137) to access the API endpoints

### Kubernetes Objects Used

- **Deployment**: Manages 2 replicas of the book API application
- **Service**: NodePort service that exposes the application

## API Endpoints

| Method | Endpoint    | Description       |
| ------ | ----------- | ----------------- |
| GET    | /books      | Get all books     |
| GET    | /books/{id} | Get a book by ID  |
| POST   | /books      | Create a new book |
| PUT    | /books/{id} | Update a book     |
| DELETE | /books/{id} | Delete a book     |

## Data Format

Books are represented as JSON objects with the following structure:

```json
{
  "bookId": "8k7j6i5h-4g3f-2e1d-0c9b-8a7z6y5x4w3v",
  "authorId": "5t4s3r2q-1p0o-9n8m-7l6k-5j4i3h2g1f0e",
  "publisherId": "2c3d4e5f-6g7h-8i9j-0k1l-2m3n4o5p6q7r",
  "title": "Brave New World",
  "publicationDate": "1932-10-27",
  "isbn": "9780060850524",
  "pages": 288,
  "genre": "Dystopian",
  "description": "A futuristic society that has replaced the pain and drama of human life with social stability at the cost of individuality.",
  "price": 15.25,
  "quantity": 19
}
```

### Example Book Creation Payload

```json
{
  "bookId": "8k7j6i5h-4g3f-2e1d-0c9b-8a7z6y5x4w3v",
  "authorId": "5t4s3r2q-1p0o-9n8m-7l6k-5j4i3h2g1f0e",
  "publisherId": "2c3d4e5f-6g7h-8i9j-0k1l-2m3n4o5p6q7r",
  "title": "Brave New World",
  "publicationDate": "1932-10-27",
  "isbn": "9780060850524",
  "pages": 288,
  "genre": "Dystopian",
  "description": "A futuristic society that has replaced the pain and drama of human life with social stability at the cost of individuality.",
  "price": 15.25,
  "quantity": 19
}
```

## Usage Examples

### Get all books

```bash
curl -X GET http://localhost:8080/books
```

### Get a specific book

```bash
curl -X GET http://localhost:8080/books/1
```

### Create a new book

```bash
curl -X POST http://localhost:8080/books \
  -H "Content-Type: application/json" \
  -d ' {
    "bookId": "8k7j6i5h-4g3f-2e1d-0c9b-8a7z6y5x4w3v",
    "authorId": "5t4s3r2q-1p0o-9n8m-7l6k-5j4i3h2g1f0e",
    "publisherId": "2c3d4e5f-6g7h-8i9j-0k1l-2m3n4o5p6q7r",
    "title": "Brave New World",
    "publicationDate": "1932-10-27",
    "isbn": "9780060850524",
    "pages": 288,
    "genre": "Dystopian",
    "description": "A futuristic society that has replaced the pain and drama of human life with social stability at the cost of individuality.",
    "price": 15.25,
    "quantity": 19
  }'
```

### Update a book

```bash
curl -X PUT http://localhost:8080/books/4 \
  -H "Content-Type: application/json" \
  -d ' {
    "bookId": "8k7j6i5h-4g3f-2e1d-0c9b-8a7z6y5x4w3v",
    "authorId": "5t4s3r2q-1p0o-9n8m-7l6k-5j4i3h2g1f0e",
    "publisherId": "2c3d4e5f-6g7h-8i9j-0k1l-2m3n4o5p6q7r",
    "title": "Brave New World - Updated",
    "publicationDate": "1932-10-27",
    "isbn": "9780060850524",
    "pages": 288,
    "genre": "Dystopian",
    "description": "A futuristic society that has replaced the pain and drama of human life with social stability at the cost of individuality.",
    "price": 15.25,
    "quantity": 19
  }'
```

### Delete a book

```bash
curl -X DELETE http://localhost:8080/books/k7j6i5h-4g3f-2e1d-0c9b-8a7z6y5x4w3v
```

## Project Structure

```
book-api/
├── main.go                 # Entry point with router setup
├── controllers/            # Request handlers for each endpoint
│   └── books.go            # Book controller functions
│   └── search.go           # Book search function
├── models/                 # Data models
│   └── book.go             # Book struct definition
├── Dockerfile              # Docker configuration
├── docker-compose.yml      # Docker Compose configuration
├── kubernetes/             # Kubernetes manifests
│   ├── deployment.yaml     # Deployment configuration
│   └── service.yaml        # Service configuration
└── README.md               # Documentation
```

## Screenshots

### API Testing in Postman

![API Response in Postman](images/Screenshot%202025-03-30%20215028.png)

### Docker Build Process

![Docker Build](images/Screenshot%202025-03-30%20223447.png)

### Docker Container Running

![Docker Container](images/Screenshot%202025-03-30%20223512.png)

### Kubernetes Deployment

![Kubernetes Deployment](images/Screenshot%202025-03-30%20223553.png)

### Minikube Service

![Minikube Service](images/Screenshot%202025-03-30%20223741.png)
