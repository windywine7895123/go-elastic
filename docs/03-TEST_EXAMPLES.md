# API Testing Examples

This guide provides examples for testing the GO-ELASTIC API using different tools.

## Using cURL (Command Line)

### Create a Book
```bash
curl -X POST http://localhost:8080/api/books \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Future Legacy Vol.3",
    "author": "John Doe",
    "isbn": "978-34-9435-339-1",
    "description": "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
    "publisher": "Packt",
    "publish_date": "2022-02-18T11:55:44.8714886+07:00",
    "pages": 493,
    "language": "Thai"
  }'
```

### Get All Books
```bash
curl -X GET http://localhost:8080/api/books \
  -H "Content-Type: application/json"
```

### Get Book by ID
```bash
curl -X GET http://localhost:8080/api/books/507f1f77bcf86cd799439011 \
  -H "Content-Type: application/json"
```

### Search by Title
```bash
curl -X GET "http://localhost:8080/api/books/search?type=title&q=Future" \
  -H "Content-Type: application/json"
```

### Search by Author
```bash
curl -X GET "http://localhost:8080/api/books/search?type=author&q=John" \
  -H "Content-Type: application/json"
```

## Using PowerShell

### Create a Book
```powershell
$body = @{
    title = "Future Legacy Vol.3"
    author = "John Doe"
    isbn = "978-34-9435-339-1"
    description = "Lorem ipsum dolor sit amet, consectetur adipiscing elit."
    publisher = "Packt"
    publish_date = "2022-02-18T11:55:44.8714886+07:00"
    pages = 493
    language = "Thai"
} | ConvertTo-Json

Invoke-WebRequest -Uri "http://localhost:8080/api/books" `
    -Method POST `
    -Headers @{"Content-Type" = "application/json"} `
    -Body $body
```

### Get All Books
```powershell
$response = Invoke-WebRequest -Uri "http://localhost:8080/api/books" `
    -Method GET `
    -Headers @{"Content-Type" = "application/json"}

$response.Content | ConvertFrom-Json | ConvertTo-Json -Depth 10
```

### Get Book by ID
```powershell
$bookId = "507f1f77bcf86cd799439011"

$response = Invoke-WebRequest -Uri "http://localhost:8080/api/books/$bookId" `
    -Method GET `
    -Headers @{"Content-Type" = "application/json"}

$response.Content | ConvertFrom-Json | ConvertTo-Json
```

### Search Books
```powershell
# Search by title
$response = Invoke-WebRequest -Uri "http://localhost:8080/api/books/search?type=title&q=Future" `
    -Method GET `
    -Headers @{"Content-Type" = "application/json"}

$response.Content | ConvertFrom-Json | ConvertTo-Json -Depth 10
```

## Using VS Code REST Client Extension

Create a file named `api.http` in project root:

```http
### Variables
@baseUrl = http://localhost:8080/api/books

### Create a Book
POST {{baseUrl}}
Content-Type: application/json

{
  "title": "Future Legacy Vol.3",
  "author": "John Doe",
  "isbn": "978-34-9435-339-1",
  "description": "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
  "publisher": "Packt",
  "publish_date": "2022-02-18T11:55:44.8714886+07:00",
  "pages": 493,
  "language": "Thai"
}

### Get All Books
GET {{baseUrl}}

### Search by Title
GET {{baseUrl}}/search?type=title&q=Future

### Search by Author
GET {{baseUrl}}/search?type=author&q=John

### Get Book by ID
@bookId = 507f1f77bcf86cd799439011
GET {{baseUrl}}/{{bookId}}
```

**Usage:**
- Click "Send Request" above each request
- Responses show in side panel
- Great for quick testing without Postman

## Using Node.js/JavaScript

### Using fetch API
```javascript
// Create a book
const book = {
  title: "Future Legacy Vol.3",
  author: "John Doe",
  isbn: "978-34-9435-339-1",
  description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
  publisher: "Packt",
  publish_date: "2022-02-18T11:55:44.8714886+07:00",
  pages: 493,
  language: "Thai"
};

fetch('http://localhost:8080/api/books', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify(book)
})
  .then(res => res.json())
  .then(data => console.log('Created:', data))
  .catch(err => console.error('Error:', err));
```

### Using axios
```javascript
const axios = require('axios');

// Create a book
axios.post('http://localhost:8080/api/books', {
  title: "Future Legacy Vol.3",
  author: "John Doe",
  isbn: "978-34-9435-339-1",
  pages: 493,
  language: "Thai"
})
  .then(res => console.log('Created:', res.data))
  .catch(err => console.error('Error:', err.response.data));
```

## Using Python

### Using requests library
```python
import requests
import json

# Create a book
book_data = {
    "title": "Future Legacy Vol.3",
    "author": "John Doe",
    "isbn": "978-34-9435-339-1",
    "description": "Lorem ipsum dolor sit amet.",
    "publisher": "Packt",
    "publish_date": "2022-02-18T11:55:44.8714886+07:00",
    "pages": 493,
    "language": "Thai"
}

response = requests.post(
    'http://localhost:8080/api/books',
    json=book_data,
    headers={'Content-Type': 'application/json'}
)

print(f"Status: {response.status_code}")
print(f"Response: {json.dumps(response.json(), indent=2)}")

# Get all books
response = requests.get('http://localhost:8080/api/books')
books = response.json()
print(f"Books: {json.dumps(books, indent=2)}")

# Search books
response = requests.get(
    'http://localhost:8080/api/books/search',
    params={'type': 'title', 'q': 'Future'}
)
results = response.json()
print(f"Search results: {json.dumps(results, indent=2)}")
```

## Test Data (JSON)

### Minimal Book (Only Required Fields)
```json
{
  "title": "My Book",
  "author": "Author Name"
}
```

### Full Book (All Fields)
```json
{
  "title": "Future Legacy Vol.3",
  "author": "John Doe",
  "isbn": "978-34-9435-339-1",
  "description": "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
  "publisher": "Packt",
  "publish_date": "2022-02-18T11:55:44.8714886+07:00",
  "pages": 493,
  "language": "Thai"
}
```

### Various Languages
```json
{
  "title": "Libro Español",
  "author": "Autor Español",
  "language": "Spanish"
}

{
  "title": "本の日本語",
  "author": "著者",
  "language": "Japanese"
}

{
  "title": "書籍中文",
  "author": "作者",
  "language": "Chinese"
}

{
  "title": "หนังสือไทย",
  "author": "ผู้เขียน",
  "language": "Thai"
}
```

## Response Examples

### Successful Create (201 Created)
```json
{
  "id": "507f1f77bcf86cd799439011",
  "title": "Future Legacy Vol.3",
  "author": "John Doe",
  "isbn": "978-34-9435-339-1",
  "description": "Lorem ipsum dolor sit amet...",
  "publisher": "Packt",
  "publish_date": "2022-02-18T11:55:44.8714886+07:00",
  "pages": 493,
  "language": "Thai",
  "created_at": "2024-01-29T14:30:00Z",
  "updated_at": "2024-01-29T14:30:00Z"
}
```

### Get All (200 OK)
```json
[
  {
    "id": "507f1f77bcf86cd799439011",
    "title": "Future Legacy Vol.3",
    "author": "John Doe",
    ...
  },
  {
    "id": "507f1f77bcf86cd799439012",
    "title": "Another Book",
    "author": "Another Author",
    ...
  }
]
```

### Not Found (404)
```json
{
  "error": "Book not found"
}
```

### Invalid Request (400)
```json
{
  "error": "Title and Author are required"
}
```

## Performance Testing

### Load Testing with Apache Bench
```bash
# Test GET /api/books with 1000 requests, 10 concurrent
ab -n 1000 -c 10 http://localhost:8080/api/books

# Results show:
# - Total time
# - Requests per second
# - Average response time
# - Min/Max response times
```

### Load Testing with Hey
```bash
# Install first
go install github.com/rakyll/hey@latest

# Run test
hey -n 1000 -c 10 http://localhost:8080/api/books

# More detailed output than ab
```

### Stress Testing
```bash
# Increase concurrency for stress test
ab -n 5000 -c 50 http://localhost:8080/api/books

# Watch for:
# - Connection errors
# - Timeout errors
# - Response time degradation
```

## Automated Testing Script

### PowerShell Test Script
```powershell
# test-api.ps1

$baseUrl = "http://localhost:8080/api/books"
$headers = @{"Content-Type" = "application/json"}
$testsPassed = 0
$testsFailed = 0

function Test-Endpoint {
    param(
        [string]$Name,
        [string]$Method,
        [string]$Url,
        [object]$Body,
        [int]$ExpectedStatus
    )
    
    try {
        if ($Body) {
            $response = Invoke-WebRequest -Uri $Url -Method $Method -Headers $headers -Body ($Body | ConvertTo-Json)
        } else {
            $response = Invoke-WebRequest -Uri $Url -Method $Method -Headers $headers
        }
        
        if ($response.StatusCode -eq $ExpectedStatus) {
            Write-Host "✅ $Name" -ForegroundColor Green
            $script:testsPassed++
        } else {
            Write-Host "❌ $Name (Expected $ExpectedStatus, got $($response.StatusCode))" -ForegroundColor Red
            $script:testsFailed++
        }
    } catch {
        Write-Host "❌ $Name (Error: $($_.Exception.Message))" -ForegroundColor Red
        $script:testsFailed++
    }
}

# Run tests
Test-Endpoint -Name "Create Book" -Method POST -Url $baseUrl `
    -Body @{title="Test";author="Test"} -ExpectedStatus 201

Test-Endpoint -Name "Get All Books" -Method GET -Url $baseUrl -ExpectedStatus 200

Test-Endpoint -Name "Search Books" -Method GET `
    -Url "$baseUrl/search?type=title&q=Test" -ExpectedStatus 200

# Summary
Write-Host "`nResults: $testsPassed passed, $testsFailed failed" -ForegroundColor Yellow
```

Run with: `.\test-api.ps1`

## Tips & Tricks

### Preserve Previous Responses
In Postman, click **Timeline** to see history of all requests and responses.

### Use Collection Runner for Bulk Testing
Perfect for running multiple requests in sequence and validating each response.

### Add Delays Between Requests
Use `setNextRequest()` in tests:
```javascript
// Run this request 5 times with 1 second delay
pm.setTimeout(1000);
pm.execution.setNextRequest("Your Request Name");
```

### Monitor Network Performance
- Use Postman Console to see response times
- Check Jaeger UI (http://localhost:16686) for detailed trace analysis
- Use browser DevTools for frontend API calls

---

**Last Updated:** January 29, 2026

