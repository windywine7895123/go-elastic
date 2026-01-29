# 🚀 GO-ELASTIC: Distributed Book Management API

[![Go](https://img.shields.io/badge/Go-1.24-blue.svg)](https://golang.org/)
[![MongoDB](https://img.shields.io/badge/MongoDB-7.0-green.svg)](https://www.mongodb.com/)
[![Elasticsearch](https://img.shields.io/badge/Elasticsearch-8.12-yellow.svg)](https://www.elastic.co/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](#license)

A modern Go REST API for managing books with distributed storage, full-text search, and distributed tracing capabilities.

## 📋 Table of Contents

- [Project Purpose](#project-purpose)
- [Architecture](#architecture)
- [Quick Start](#quick-start)
- [Documentation](#documentation)
- [Debugging Guide](#debugging-guide)
- [API Endpoints](#api-endpoints)
- [Contributing](#contributing)

## 🎯 Project Purpose

**GO-ELASTIC** is a production-ready book management API that demonstrates:

### Core Features
- ✅ **RESTful API** - Complete CRUD operations for books and users
- ✅ **MongoDB Storage** - Primary database with flexible document model
- ✅ **Elasticsearch Integration** - Full-text search capabilities with indexed data
- ✅ **Distributed Tracing** - OpenTelemetry + Jaeger for performance monitoring
- ✅ **CORS Support** - Works seamlessly with Postman, browsers, and frontend apps
- ✅ **Error Handling** - Graceful degradation when services are unavailable
- ✅ **Request Logging** - Complete HTTP request tracking with request IDs
- ✅ **Recovery Middleware** - Panic recovery to prevent crashes

### Use Cases
1. **Learning Go** - See best practices for Go web applications
2. **Microservices Pattern** - Three-tier architecture (Handler→Service→Repository)
3. **Database Integration** - MongoDB and Elasticsearch integration
4. **Observability** - Distributed tracing setup and monitoring
5. **API Development** - RESTful API design with proper error handling
6. **Unit test-writing** - Practice unit test for use in realworld

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────┐
│         GO-ELASTIC API Server (Fiber)               │
│                                                     │
│  ┌─────────────────────────────────────────────┐  │
│  │         HTTP Handler Layer                  │  │
│  │  (CreateBook, GetBook, SearchBooks, etc)   │  │
│  └────────────────────┬────────────────────────┘  │
│                       │                            │
│  ┌────────────────────▼────────────────────────┐  │
│  │         Business Logic Layer                │  │
│  │  (Validation, Tracing, Orchestration)      │  │
│  └────────────────────┬────────────────────────┘  │
│                       │                            │
│  ┌────────────────────▼────────────────────────┐  │
│  │      Repository Layer                       │  │
│  │  (MongoDB + Elasticsearch Operations)       │  │
│  └────────────────────┬────────────────────────┘  │
│                       │                            │
└───────────────────────┼────────────────────────────┘
                        │
        ┌───────────────┼───────────────┐
        │               │               │
        ▼               ▼               ▼
    ┌────────┐    ┌──────────────┐  ┌─────────────┐
    │MongoDB │    │Elasticsearch │  │   Jaeger    │
    │:27017  │    │   :9200      │  │  :16686     │
    └────────┘    └──────────────┘  └─────────────┘
                        │
                        ▼
                   ┌─────────┐
                   │  Kibana │
                   │ :5601   │
                   └─────────┘
```

### Technology Stack

| Component | Technology | Purpose |
|-----------|-----------|---------|
| **Framework** | Fiber v2 | Fast, Gin-like Go web framework |
| **Database** | MongoDB 7.0 | Document-based primary storage |
| **Search** | Elasticsearch 8.12 | Full-text search & indexing |
| **Tracing** | OpenTelemetry + Jaeger | Distributed request tracing |
| **Visualization** | Kibana | Elasticsearch data visualization |
| **Language** | Go 1.24 | Backend language |

## 🚀 Quick Start

### Prerequisites
- Go 1.24+ installed
- Docker & Docker Compose
- Postman or cURL (for API testing)

### 1. Clone and Setup

```bash
cd c:\Users\pearl\lab\go-elastic
cp .env.example .env
```

### 2. Start Services

```bash
# Start all backend services (MongoDB, Elasticsearch, Jaeger)
docker-compose up -d

# Verify services are running
docker-compose ps
```

Expected output:
```
NAME         STATUS              PORTS
mongodb      Up 2 minutes        0.0.0.0:27017->27017/tcp
elasticsearch Up 2 minutes       0.0.0.0:9200->9200/tcp
kibana       Up 2 minutes        0.0.0.0:5601->5601/tcp
jaeger       Up 2 minutes        0.0.0.0:4318->4318/tcp, 0.0.0.0:16686->16686/tcp
```

### 3. Start Application

```bash
go run main.go
```

Expected output:
```
server starting on :8080
OpenTelemetry tracer initialized (endpoint: localhost:4318, service: go-elastic-api)
```

### 4. Test API

```bash
curl -X POST http://localhost:8080/api/books \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Future Legacy Vol.3",
    "author": "John Doe",
    "isbn": "978-34-9435-339-1",
    "pages": 493,
    "language": "Thai"
  }'
```

## 📚 Documentation

Comprehensive guides are available in the `/docs` folder:

### API Documentation
- [**Book API Guide**](docs/01-BOOK_API.md) - Complete Book endpoint reference
- [**Postman Setup**](docs/02-POSTMAN_SETUP.md) - How to use Postman with the API
- [**API Examples**](docs/03-TEST_EXAMPLES.md) - cURL, PowerShell, and REST client examples

### Setup & Configuration
- [**MongoDB Migration**](docs/04-MONGODB_MIGRATION.md) - Database setup and migration guide
- [**Tracing Setup**](docs/05-TRACING_SETUP.md) - OpenTelemetry and Jaeger configuration
- [**Environment Variables**](docs/.env.example) - Configuration reference

### Architecture & Design
- [**Architecture Overview**](docs/06-ARCHITECTURE.md) - System design and patterns
- [**Debugging Guide**](#debugging-guide) - Below in this README

## 🐛 Debugging Guide

### 1. Application Won't Start

#### Error: "Failed to connect to database"

```
Error: failed to connect to MongoDB
```

**Solution:**
```bash
# Check if MongoDB container is running
docker-compose ps mongodb

# If not running, start it
docker-compose up -d mongodb

# Check MongoDB logs
docker-compose logs mongodb

# Verify connection manually
mongosh "mongodb://root:password@localhost:27017"
```

#### Error: "Port 8080 already in use"

```
Error: listen tcp :8080: bind: An attempt was made to reuse a socket...
```

**Solution:**
```bash
# Find process using port 8080 (Windows)
netstat -ano | findstr ":8080"

# Kill the process (replace PID)
taskkill /PID <PID> /F

# Or use different port
go run main.go  # Then modify main.go to use different port
```

### 2. API Requests Failing

#### Error: "Cannot parse JSON"

```json
{
  "error": "Cannot parse JSON",
  "details": "EOF",
  "hint": "Ensure Content-Type header is 'application/json' and request body is valid JSON"
}
```

**Solutions:**
1. **Check Content-Type header:**
   ```bash
   curl -X POST http://localhost:8080/api/books \
     -H "Content-Type: application/json" \  # ← MUST have this
     -d '{"title":"Test","author":"John"}'
   ```

2. **Validate JSON syntax:**
   - Use online JSON validator: https://jsonlint.com
   - Missing commas, quotes, or brackets

3. **Check request body:**
   - Empty body for POST?
   - Binary/image data instead of JSON?

#### Error: "Status is 400 - Title and Author are required"

**Solution:** Both fields are mandatory:
```json
{
  "title": "Book Title",      // Required
  "author": "Author Name"     // Required
}
```

#### Error: "Cannot connect to http://localhost:8080"

**Solution:**
```bash
# Check if app is running
# You should see in terminal: "server starting on :8080"

# If not, start it
go run main.go

# Test connection
curl http://localhost:8080/hello
```

### 3. Database Issues

#### MongoDB: "Connection refused"

```
Error: connection refused (127.0.0.1:27017)
```

**Solution:**
```bash
# Start MongoDB container
docker-compose up -d mongodb

# Wait a few seconds for MongoDB to start
Start-Sleep -Seconds 5

# Test connection
mongosh "mongodb://root:password@localhost:27017"
```

#### Elasticsearch: "No such host"

**Solution:**
```bash
# Start Elasticsearch
docker-compose up -d elasticsearch

# Wait for startup
Start-Sleep -Seconds 10

# Test connection
curl http://localhost:9200

# Check logs
docker-compose logs elasticsearch
```

### 4. Tracing Issues

#### Error: "traces export: parse ... invalid URL escape"

**Solution:** Check your `.env` file:

```env
# ✅ Correct
JAEGER_ENDPOINT=localhost:4318

# ✅ Also correct
OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4318

# ❌ Wrong (will cause URL parsing error)
OTEL_EXPORTER_OTLP_ENDPOINT=http://http://localhost:4318
```

#### Traces not showing in Jaeger UI

**Solution:**
```bash
# 1. Verify Jaeger is running
docker-compose ps jaeger

# 2. Check endpoint in .env
cat .env | grep JAEGER_ENDPOINT

# 3. Restart app with correct endpoint
go run main.go

# 4. Make an API request
curl -X GET http://localhost:8080/api/books \
  -H "Content-Type: application/json"

# 5. Open Jaeger UI
# http://localhost:16686
# Select service: go-elastic-api
# Click: Find Traces
```

### 5. Elasticsearch Integration

#### Books not appearing in Elasticsearch

**Solution:**
```bash
# 1. Check if index was created
curl http://localhost:9200/_cat/indices

# You should see: "books"

# 2. Check if documents were indexed
curl http://localhost:9200/books/_search

# 3. View in Kibana
# http://localhost:5601
# Dev Tools → Console
# GET /books/_search

# 4. Check repository logs for errors
# App should show: "Indexed book in Elasticsearch"
```

#### Search returns no results

**Solution:**
```bash
# 1. Verify books exist in MongoDB
mongosh << EOF
use go_logger
db.books.find().pretty()
EOF

# 2. Check search query format
# Correct:
curl "http://localhost:8080/api/books/search?type=title&q=Go"

# Wrong:
curl "http://localhost:8080/api/books/search?title=Go"

# 3. Check Elasticsearch index mapping
curl http://localhost:9200/books/_mapping
```

### 6. Performance Debugging

#### Slow API responses

**Using Jaeger:**
1. Open http://localhost:16686
2. Select service: "go-elastic-api"
3. Choose an operation (POST /api/books, GET /api/books, etc)
4. View duration breakdown:
   ```
   Total: 500ms
   ├── Request Parse: 10ms
   ├── Validation: 5ms
   ├── DB Insert: 400ms ← Slow!
   ├── ES Index: 70ms
   └── Response: 15ms
   ```

**Optimize slow components:**
```go
// In service/book_service.go
tr := otel.Tracer("book-service")
ctx, span := tr.Start(ctx, "CreateBook")
defer span.End()

// Add timing information
span.AddEvent("Starting database insert")
start := time.Now()

// Your code here
err := s.repo.Create(ctx, book)

span.AddEvent(fmt.Sprintf("DB insert took %dms", time.Since(start).Milliseconds()))
```

### 7. Common Docker Issues

#### Docker daemon not running

```
Error: Cannot connect to Docker daemon
```

**Solution:**
- Windows: Open Docker Desktop application
- Linux: `sudo systemctl start docker`
- Mac: Open Docker.app from Applications

#### Container exit code 1

```bash
# Check what went wrong
docker-compose logs mongodb

# Common issue: port already in use
# Change port in docker-compose.yaml
```

#### Out of disk space

```bash
# Clean up unused Docker images and containers
docker-compose down -v  # Remove volumes too
docker system prune -a

# Then restart
docker-compose up -d
```

## 📊 Monitoring & Observability

### Jaeger UI (Distributed Tracing)
- **URL:** http://localhost:16686
- **Service:** go-elastic-api
- **View:** Request traces, latency, errors, dependencies

### Kibana (Elasticsearch Visualization)
- **URL:** http://localhost:5601
- **Dashboards:** Create dashboards for book data
- **Queries:** Run complex queries on indexed data

### Application Logs
```bash
# Real-time logs
go run main.go

# Structured logging to analyze
go run main.go | tee app.log
```

## 📝 Testing

### Unit Tests
```bash
go test ./...
```

### Integration Tests
```bash
# Requires all services running
docker-compose up -d
go test -v ./tests/integration
```

### Load Testing
```bash
# Using Apache Bench
ab -n 1000 -c 10 http://localhost:8080/api/books

# Using hey
go install github.com/rakyll/hey@latest
hey -n 1000 -c 10 http://localhost:8080/api/books
```

## 🔧 Troubleshooting Checklist

When something doesn't work:

- [ ] Services running? `docker-compose ps`
- [ ] App started? Check terminal output
- [ ] Correct endpoint? `http://localhost:8080`
- [ ] Content-Type header? `application/json`
- [ ] Valid JSON body? Validate at jsonlint.com
- [ ] Required fields? title + author
- [ ] MongoDB accessible? Try mongosh
- [ ] Elasticsearch accessible? Try `curl http://localhost:9200`
- [ ] Jaeger accessible? Try http://localhost:16686
- [ ] Firewall blocking? Check Windows Firewall
- [ ] Port conflicts? Check with netstat

## 🚢 Deployment

For production deployment:

1. **Update .env** with real endpoints
2. **Set APP_ENV=production**
3. **Use environment-specific docker-compose**
4. **Enable HTTPS** for API
5. **Configure Jaeger persistence**
6. **Set up database backups**
7. **Monitor resource usage**

See [docs/DEPLOYMENT.md](docs/DEPLOYMENT.md) for details.

## 📖 Additional Resources

### Official Documentation
- [Fiber Framework](https://docs.gofiber.io/)
- [MongoDB Go Driver](https://pkg.go.dev/go.mongodb.org/mongo-driver)
- [Elasticsearch Go Client](https://www.elastic.co/guide/en/elasticsearch/client/go-api/current/index.html)
- [OpenTelemetry Go](https://opentelemetry.io/docs/instrumentation/go/)

### Example Requests
See [docs/API_EXAMPLES.md](docs/03-TEST_EXAMPLES.md) for:
- cURL examples
- PowerShell examples
- REST Client examples
- Postman collection

## �️ Frontend HTML - Elasticsearch Book Search

The project includes a mock frontend HTML file (`mock_frontend.html`) that demonstrates real-time book search functionality with Elasticsearch integration.

### Features

**Real-Time Search with Debouncing**
```javascript
// Debounce timer (300ms)
// Waits for user to stop typing before searching
// Prevents excessive API calls during typing
debounceTimer = setTimeout(() => {
    performSearch(query);
}, 300);
```

**Race Condition Prevention**
```javascript
// Abort controller to cancel previous requests
// When user types a new character, old requests are cancelled
// Only the latest search result is rendered
if (abortController) {
    abortController.abort();
}
abortController = new AbortController();
```

**CSS-in-JS Styling**
```css
/* Important: Set box-sizing to border-box */
/* Ensures padding is included in width calculation */
box-sizing: border-box;
```

**Responsive Design**
- Smooth hover animations (card lift effect)
- Mobile-friendly input with focus states
- Status indicator showing search state (Ready, Typing, Searching, Results)
- Badge system for language display

### HTML Structure

**Search Input**
```html
<input type="text" id="searchInput" 
       placeholder="Type to search books (e.g. Martin, Clean Code)...">
```

**Result Cards**
```html
<div class="card">
    <span class="badge">${book.language || 'EN'}</span>
    <h3>${book.title}</h3>
    <p>✍️ <strong>Author:</strong> ${book.author}</p>
    <p>📖 <strong>Publisher:</strong> ${book.publisher} | 📄 ${book.pages} pages</p>
    <p style="...">ISBN: ${book.isbn}</p>
</div>
```

### How It Works

1. **User Types** - Input event triggers debounce timer
2. **Wait for Input** - 300ms delay to batch keystrokes
3. **Search Request** - Calls `/api/books/search?q=query` endpoint
4. **Race Prevention** - Cancels previous requests if new search starts
5. **Render Results** - Displays books in styled cards

### Running the Frontend

**Option 1: Open Directly in Browser**
```bash
# Simply open the HTML file
open mock_frontend.html
```

**Option 2: Serve via Python**
```bash
# Python 3
python -m http.server 8000

# Then visit: http://localhost:8000/mock_frontend.html
```

**Option 3: Use Go HTTP Server**
```bash
# Add to main.go
app.Static("/", "./mock_frontend.html")

# Then visit: http://localhost:3000/mock_frontend.html
```

### Frontend-API Integration

**API Endpoint Called**
```
GET /api/books/search?type=title&q={query}
```

**Expected Response Format**
```json
[
  {
    "_id": "507f1f77bcf86cd799439011",
    "title": "Clean Code",
    "author": "Robert Martin",
    "publisher": "Prentice Hall",
    "pages": 464,
    "isbn": "978-0132350884",
    "language": "EN"
  }
]
```

**CORS Requirements**
- Frontend must run on different port than API (e.g., 8000 vs 3000)
- API CORS middleware handles cross-origin requests
- No additional frontend configuration needed

### Performance Optimizations

**1. Debouncing (300ms)**
- Reduces API calls during rapid typing
- Example: User types "Clean Code" (10 characters)
  - Without debounce: 10 API calls
  - With debounce: 1 API call

**2. Request Cancellation**
- Cancels previous request if new search starts
- Prevents race conditions where old results override new ones
- Uses AbortController API

**3. Local Rendering**
- Results rendered on frontend (no server-side template)
- Faster UI updates
- Reduces server load

### Styling Reference

**Color Scheme**
| Element | Color | Purpose |
|---------|-------|---------|
| Background | #f4f6f8 | Neutral gray |
| Cards | white | Content containers |
| Links | #007bff | Interactive elements |
| Text | #333 | Primary text |
| Placeholder | #999 | Hints |
| Badge | #e3f2fd | Language indicator |

**Interactive Elements**
```css
/* Input focus state */
input:focus {
    border-color: #007bff;
    box-shadow: 0 4px 12px rgba(0,123,255,0.2);
}

/* Card hover effect */
.card:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 8px rgba(0,0,0,0.1);
}
```

### Troubleshooting Frontend

**Issue: CORS Errors**
```
Access to XMLHttpRequest blocked by CORS policy
```
**Solution:** Ensure API_URL in HTML matches running server
```javascript
const API_URL = "http://localhost:3000/api/books/search?type=title&q=";
```

**Issue: No Results**
```
Status shows "Searching..." but no results
```
**Solutions:**
1. Check API is running: `curl http://localhost:3000/health`
2. Check MongoDB has books: MongoDB Compass or mongosh
3. Check browser console for error messages

**Issue: Slow Search Response**
```
Searching takes >1 second
```
**Solutions:**
1. Ensure Elasticsearch is indexed: Check Kibana
2. Create database indexes for title/author fields
3. Check network latency: Browser DevTools > Network tab

---

## �👥 Support

For issues and questions:
1. Check [Debugging Guide](#debugging-guide) above
2. Review [docs/](docs/) folder for detailed guides
3. Check application logs
4. Review Jaeger traces for performance insights

---

**Last Updated:** January 29, 2026  
**Go Version:** 1.24  
