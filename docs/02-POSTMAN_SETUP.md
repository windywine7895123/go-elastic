# Postman Setup Guide

## Overview
This guide explains how to use Postman to test the GO-ELASTIC API.

## What is Postman?
Postman is a popular API client that lets you test HTTP requests without writing code. It's perfect for:
- Testing API endpoints
- Debugging requests
- Sharing API collections
- Automating test runs

## Installation

1. Download from https://www.postman.com/downloads/
2. Install for your OS (Windows/Mac/Linux)
3. Create a free account or use without account
4. Launch Postman

## Quick Start (Import Collection)

### Step 1: Get Collection File
The project includes `postman_collection.json` with all endpoints pre-configured.

### Step 2: Import into Postman
1. Open Postman
2. Click **Import** (top-left)
3. Select **Upload Files** tab
4. Browse to project root and select `postman_collection.json`
5. Click **Import**

All endpoints are now ready to use!

## Manual Setup (If not importing)

### Creating Your First Request

#### Create a Book (POST)
1. **New Request**
   - Click **+** tab or File → New → Request
   
2. **Configure Request**
   - Method: `POST` (dropdown)
   - URL: `http://localhost:8080/api/books`
   
3. **Add Headers**
   - Go to **Headers** tab
   - Add `Content-Type: application/json`
   
4. **Add Body**
   - Go to **Body** tab
   - Select **raw**
   - Select **JSON** from dropdown
   - Paste JSON:
   ```json
   {
     "title": "The Go Programming Language",
     "author": "Alan Donovan, Brian Kernighan",
     "isbn": "978-0134190440",
     "pages": 400,
     "language": "English"
   }
   ```
   
5. **Send Request**
   - Click **Send**
   - View response below

#### Get All Books (GET)
1. Method: `GET`
2. URL: `http://localhost:8080/api/books`
3. Headers: `Content-Type: application/json`
4. Body: (none for GET)
5. Click **Send**

#### Search Books (GET)
1. Method: `GET`
2. URL: `http://localhost:8080/api/books/search?type=title&q=Go`
3. Headers: `Content-Type: application/json`
4. Body: (none)
5. Click **Send**

## Environment Variables

### Create Custom Environment
1. Click **Environments** (left sidebar)
2. Click **+** to create new
3. Name it: `go-elastic-dev`
4. Add variables:
   ```
   base_url     : http://localhost:8080
   book_id      : (leave empty)
   ```
5. Click **Save**

### Use Variables in Requests
Instead of hardcoding URLs:

```
{{base_url}}/api/books              instead of    http://localhost:8080/api/books
{{base_url}}/api/books/{{book_id}}  instead of    http://localhost:8080/api/books/507f1f77bcf86cd799439011
```

## Running Collections (Automated Testing)

### Runner Feature
Perfect for running all requests in sequence.

**Steps:**
1. Click **Runner** (top menu)
2. Select your collection (e.g., "Go-Elastic Book API")
3. Select environment: `go-elastic-dev` (if created)
4. Set iterations: `1`
5. Click **Start Run**

Results show:
- ✅ Passed requests
- ❌ Failed requests
- Response times
- Status codes

### Tests in Collections
Add assertions to validate responses automatically.

**Example Test:**
```javascript
pm.test("Status is 200", function () {
    pm.response.to.have.status(200);
});

pm.test("Response has books array", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData).to.be.an('array');
});
```

## Pre-request Scripts

Run scripts before each request (e.g., authentication).

**Example:**
```javascript
// Generate timestamp for created_at
pm.environment.set("created_at", new Date().toISOString());
```

Then use `{{created_at}}` in request body.

## Common Tasks

### Extract Value from Response
**In Tests tab:**
```javascript
// After successful POST, extract the book ID
var jsonData = pm.response.json();
pm.environment.set("book_id", jsonData.id);
```

Now use `{{book_id}}` in next request (Get Book by ID).

### Validate Response Structure
**In Tests tab:**
```javascript
pm.test("Book has required fields", function () {
    var book = pm.response.json();
    pm.expect(book).to.have.property('title');
    pm.expect(book).to.have.property('author');
    pm.expect(book).to.have.property('id');
});
```

### Check Response Time
**In Tests tab:**
```javascript
pm.test("Response time is less than 200ms", function () {
    pm.expect(pm.response.responseTime).to.be.below(200);
});
```

## Troubleshooting

### Connection Refused
```
Error: Cannot GET http://localhost:8080/api/books
```

**Solution:**
1. Make sure app is running: `go run main.go`
2. Check port is correct (default: 8080)
3. Check firewall isn't blocking

### "Cannot parse JSON" Error
```json
{
  "error": "Cannot parse JSON",
  "details": "...",
  "hint": "Ensure Content-Type header is 'application/json' and request body is valid JSON"
}
```

**Solution:**
1. Go to **Headers** tab
2. Add `Content-Type: application/json`
3. Verify JSON is valid (use jsonlint.com)
4. Try again

### Empty or Malformed Response
**Solution:**
1. Open **Postman Console** (View → Show Postman Console)
2. See raw request/response
3. Check for hidden characters or encoding issues
4. Verify content-type is `application/json`

## Pro Tips

### Use Collections Folder Structure
Organize requests in folders:
```
Go-Elastic API
├── Books
│   ├── Create Book
│   ├── Get All Books
│   ├── Get Book by ID
│   └── Search Books
└── Users
    ├── Create User
    └── Get All Users
```

### Save Common Tests
Create a collection-level script (Tests tab at collection level) that runs for all requests:

```javascript
pm.test("Response is valid JSON", function () {
    pm.response.to.be.json;
});
```

### Mock Servers
Use Postman's mock server feature to test without real API:
1. Right-click collection
2. Select **Mock collection**
3. Define mock responses
4. Run requests against mock

### Collaboration
Share collections with team:
1. Click **Share** on collection
2. Generate invite link
3. Team members can import and collaborate

## API Workflow Example

**Complete workflow for creating and retrieving a book:**

1. **POST /api/books** - Create book
   - Response: `201` with book ID
   - Extract ID with script

2. **GET /api/books/{{book_id}}** - Get specific book
   - Response: `200` with full book data
   
3. **GET /api/books/search** - Search for book
   - Response: `200` with search results

4. **GET /api/books** - List all books
   - Response: `200` with all books

## Next Steps

- Review [01-BOOK_API.md](01-BOOK_API.md) for all endpoints
- Check [03-TEST_EXAMPLES.md](03-TEST_EXAMPLES.md) for cURL examples
- Explore [Postman Learning Center](https://learning.postman.com/)

---

**Tip:** Use **Postman Console** (Ctrl+Alt+C) to debug requests in detail!

**Last Updated:** January 29, 2026
