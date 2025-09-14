# Library Management System

A comprehensive library management backend API built with Go, Gin, and PostgreSQL.

## Features

- **Book Management**: CRUD operations for books with ISBN validation
- **Author Management**: Manage authors and their biographies
- **Borrower Management**: Library member management with email validation
- **Borrowing System**: Track book borrowings, returns, and overdue books
- **Search & Pagination**: Search functionality across all entities with pagination
- **RESTful API**: Clean REST API design with proper HTTP status codes
- **Database Migrations**: Automatic database schema management

## Tech Stack

- **Go 1.21+**
- **Gin** - HTTP web framework
- **GORM** - ORM library
- **PostgreSQL** - Database
- **UUID** - Unique identifiers
- **Environment Variables** - Configuration management

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- Git

## Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd library-management-go
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Set up environment variables**
   ```bash
   cp env.example .env
   ```
   Edit `.env` file with your database configuration:
   ```
   DATABASE_URL=postgres://username:password@localhost:5432/library_management?sslmode=disable
   PORT=8080
   GIN_MODE=debug
   ```

4. **Set up PostgreSQL database**
   ```sql
   CREATE DATABASE library_management;
   ```

5. **Run the application**
   ```bash
   go run main.go
   ```

The API will be available at `http://localhost:8080`

## API Endpoints

### Health Check
- `GET /api/v1/health` - Check API status

### Authors
- `POST /api/v1/authors` - Create author
- `GET /api/v1/authors` - Get all authors (with pagination and search)
- `GET /api/v1/authors/:id` - Get author by ID
- `PUT /api/v1/authors/:id` - Update author
- `DELETE /api/v1/authors/:id` - Delete author

### Books
- `POST /api/v1/books` - Create book
- `GET /api/v1/books` - Get all books (with pagination and search)
- `GET /api/v1/books/:id` - Get book by ID
- `PUT /api/v1/books/:id` - Update book
- `DELETE /api/v1/books/:id` - Delete book

### Borrowers
- `POST /api/v1/borrowers` - Create borrower
- `GET /api/v1/borrowers` - Get all borrowers (with pagination and search)
- `GET /api/v1/borrowers/:id` - Get borrower by ID
- `PUT /api/v1/borrowers/:id` - Update borrower
- `DELETE /api/v1/borrowers/:id` - Delete borrower

### Borrowings
- `POST /api/v1/borrowings/borrow` - Borrow a book
- `POST /api/v1/borrowings/return` - Return a book
- `GET /api/v1/borrowings` - Get all borrowings (with pagination)
- `GET /api/v1/borrowings/:id` - Get borrowing by ID
- `GET /api/v1/borrowings/borrower/:borrowerId` - Get borrowings by borrower
- `GET /api/v1/borrowings/overdue` - Get overdue borrowings
- `PUT /api/v1/borrowings/update-overdue` - Update overdue status

## Request/Response Examples

### Create Author
```json
POST /api/v1/authors
{
  "name": "J.K. Rowling",
  "biography": "British author, best known for the Harry Potter series"
}
```

### Create Book
```json
POST /api/v1/books
{
  "title": "Harry Potter and the Philosopher's Stone",
  "isbn": "978-0747532699",
  "description": "The first book in the Harry Potter series",
  "author_id": "author-uuid-here",
  "published_at": "1997-06-26T00:00:00Z"
}
```

### Create Borrower
```json
POST /api/v1/borrowers
{
  "name": "John Doe",
  "email": "john.doe@example.com",
  "phone": "+1234567890",
  "address": "123 Main St, City, Country"
}
```

### Borrow Book
```json
POST /api/v1/borrowings/borrow
{
  "book_id": "book-uuid-here",
  "borrower_id": "borrower-uuid-here",
  "due_date": "2024-02-15T00:00:00Z"
}
```

## Query Parameters

### Pagination
- `page` - Page number (default: 1)
- `limit` - Items per page (default: 10, max: 100)

### Search
- `search` - Search query for title, ISBN, author name, etc.

### Example
```
GET /api/v1/books?page=1&limit=20&search=harry potter
```

## Business Rules

1. **Books**: ISBN must be unique, cannot delete books that are currently borrowed
2. **Authors**: Cannot delete authors with existing books
3. **Borrowers**: Email must be unique, cannot delete borrowers with active borrowings
4. **Borrowings**: 
   - Maximum 5 books per borrower
   - Cannot borrow if borrower has overdue books
   - Books become unavailable when borrowed
   - Books become available when returned

## Database Schema

The application uses the following main entities:
- **Authors**: id, name, biography, timestamps
- **Books**: id, title, isbn, description, author_id, published_at, available, timestamps
- **Borrowers**: id, name, email, phone, address, timestamps
- **Borrowings**: id, book_id, borrower_id, borrowed_at, due_date, returned_at, status, timestamps

## Development

### Project Structure
```
library-management-go/
├── main.go
├── go.mod
├── env.example
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── database/
│   │   └── database.go
│   ├── models/
│   │   └── models.go
│   ├── services/
│   │   ├── author_service.go
│   │   ├── book_service.go
│   │   ├── borrower_service.go
│   │   └── borrowing_service.go
│   ├── handlers/
│   │   ├── author_handler.go
│   │   ├── book_handler.go
│   │   ├── borrower_handler.go
│   │   └── borrowing_handler.go
│   └── routes/
│       └── routes.go
└── README.md
```

### Running Tests
```bash
go test ./...
```

### Building
```bash
go build -o library-management main.go
```

## License

This project is licensed under the MIT License.
