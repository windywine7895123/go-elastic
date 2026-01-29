# Book API Documentation

## Overview
The Book API allows you to create, retrieve, and search books. All book data is stored in MongoDB and automatically indexed in Elasticsearch for fast full-text search capabilities.

## Base URL
```
http://localhost:8080/api/books
```

## Endpoints

### 1. Create a Book
**Endpoint:** `POST /api/books`

Creates a new book and stores it in both MongoDB and Elasticsearch.

**Request Body:**
```json
{
  "title": "The Go Programming Language",
  "author": "Alan Donovan, Brian Kernighan",
  "isbn": "978-0134190440",
  "description": "A comprehensive guide to Go programming language",
  "publisher": "Addison-Wesley",
  "publish_date": "2015-10-26T00:00:00Z",
  "pages": 400,
  "language": "English"
}
```

**Required Fields:**
- `title` (string) - The title of the book
- `author` (string) - The author of the book

**Optional Fields:**
- `isbn` (string) - International Standard Book Number
- `description` (string) - Book description
- `publisher` (string) - Publisher name
- `publish_date` (ISO 8601 datetime) - Publication date
- `pages` (integer) - Number of pages
- `language` (string) - Language of the book

**Response:** `201 Created`
```json
{
  "id": "507f1f77bcf86cd799439011",
  "title": "The Go Programming Language",
  "author": "Alan Donovan, Brian Kernighan",
  "isbn": "978-0134190440",
  "description": "A comprehensive guide to Go programming language",
  "publisher": "Addison-Wesley",
  "publish_date": "2015-10-26T00:00:00Z",
  "pages": 400,
  "language": "English",
  "created_at": "2024-01-29T10:30:00Z",
  "updated_at": "2024-01-29T10:30:00Z"
}
```

**Error Response:** `400 Bad Request`
```json
{
  "error": "Title and Author are required"
}
```

### 2. Get All Books
**Endpoint:** `GET /api/books`

Retrieves all books from MongoDB.

**Response:** `200 OK`
```json
[
  {
    "id": "507f1f77bcf86cd799439011",
    "title": "The Go Programming Language",
    "author": "Alan Donovan, Brian Kernighan",
    "isbn": "978-0134190440",
    "description": "A comprehensive guide to Go programming language",
    "publisher": "Addison-Wesley",
    "publish_date": "2015-10-26T00:00:00Z",
    "pages": 400,
    "language": "English",
    "created_at": "2024-01-29T10:30:00Z",
    "updated_at": "2024-01-29T10:30:00Z"
  }
]
```

### 3. Get Book by ID
**Endpoint:** `GET /api/books/:id`

Retrieves a specific book by its MongoDB ObjectID.

**URL Parameters:**
- `id` (string) - MongoDB ObjectID in hex format

**Example:**
```
GET http://localhost:8080/api/books/507f1f77bcf86cd799439011
```

**Response:** `200 OK`
```json
{
  "id": "507f1f77bcf86cd799439011",
  "title": "The Go Programming Language",
  "author": "Alan Donovan, Brian Kernighan",
  "isbn": "978-0134190440",
  "description": "A comprehensive guide to Go programming language",
  "publisher": "Addison-Wesley",
  "publish_date": "2015-10-26T00:00:00Z",
  "pages": 400,
  "language": "English",
  "created_at": "2024-01-29T10:30:00Z",
  "updated_at": "2024-01-29T10:30:00Z"
}
```

**Error Response:** `404 Not Found`
```json
{
  "error": "Book not found"
}
```

### 4. Search Books
**Endpoint:** `GET /api/books/search?type=<type>&q=<query>`

Searches books using Elasticsearch full-text search.

**Query Parameters:**
- `type` (string) - Search field: `title` or `author`
- `q` (string) - Search query

**Examples:**

Search by title:
```
GET http://localhost:8080/api/books/search?type=title&q=Go+Programming
```

Search by author:
```
GET http://localhost:8080/api/books/search?type=author&q=Donovan
```

**Response:** `200 OK`
```json
[
  {
    "id": "507f1f77bcf86cd799439011",
    "title": "The Go Programming Language",
    "author": "Alan Donovan, Brian Kernighan",
    "isbn": "978-0134190440",
    "description": "A comprehensive guide to Go programming language",
    "publisher": "Addison-Wesley",
    "publish_date": "2015-10-26T00:00:00Z",
    "pages": 400,
    "language": "English",
    "created_at": "2024-01-29T10:30:00Z",
    "updated_at": "2024-01-29T10:30:00Z"
  }
]
```

## Error Codes

| Code | Message | Cause |
|------|---------|-------|
| 400 | Cannot parse JSON | Invalid JSON format or missing Content-Type header |
| 400 | Title and Author are required | Missing required fields |
| 404 | Book not found | Invalid book ID or book doesn't exist |
| 500 | Internal Server Error | Server error (check logs) |

## Data Flow

When you create a book:

```
1. POST /api/books with JSON body
         ↓
2. Handler validates JSON and required fields
         ↓
3. Service creates timestamps if not provided
         ↓
4. Repository saves to MongoDB
         ↓
5. Repository indexes in Elasticsearch
         ↓
6. Response with created book (201)
```

## Best Practices

### Required Timestamps
The API automatically sets `created_at` and `updated_at` if not provided. To specify custom values:

```json
{
  "title": "My Book",
  "author": "Author Name",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### Search Optimization
For better search results:
- Use complete words ("Programming Language" not "Prog")
- Elasticsearch handles case-insensitivity
- Special characters are handled automatically
- Spaces in queries are handled as-is

### ObjectID Format
MongoDB ObjectIDs are 24-character hex strings:
- `507f1f77bcf86cd799439011` ✅ Valid
- `invalid-id` ❌ Invalid (must be 24 hex chars)

## Testing

See the project root [README.md](../README.md#testing) for testing instructions.

---

**Last Updated:** January 29, 2026
