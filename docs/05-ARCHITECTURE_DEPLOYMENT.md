# Architecture & Deployment Guide

## System Architecture

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        Client Layer                               │
│              (Postman, Browser, cURL, etc.)                      │
└──────────────────────────────┬──────────────────────────────────┘
                               │ HTTP/REST
┌──────────────────────────────▼──────────────────────────────────┐
│                   API Server (Go + Fiber)                         │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │ Routes                                                    │    │
│  │ - POST /api/users                                        │    │
│  │ - GET  /api/users/:id                                   │    │
│  │ - GET  /api/users                                       │    │
│  │ - POST /api/books                                       │    │
│  │ - GET  /api/books/:id                                  │    │
│  │ - GET  /api/books                                      │    │
│  │ - GET  /api/books/search?query=...                     │    │
│  └─────────────────────────────────────────────────────────┘    │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │ Middleware Stack                                         │    │
│  │ - CORS (AllowOrigins: *)                                 │    │
│  │ - Request ID Generation                                 │    │
│  │ - Logging (structured)                                  │    │
│  │ - Recovery (panic handling)                             │    │
│  │ - Tracing (OpenTelemetry)                              │    │
│  └─────────────────────────────────────────────────────────┘    │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │ Three-Tier Architecture                                 │    │
│  │                                                          │    │
│  │  Handler Layer (HTTP layer)                             │    │
│  │  ├─ user_handler.go                                     │    │
│  │  └─ book_handler.go                                     │    │
│  │           │                                             │    │
│  │           ▼                                             │    │
│  │  Service Layer (Business logic)                         │    │
│  │  ├─ user_service.go                                     │    │
│  │  └─ book_service.go                                     │    │
│  │           │                                             │    │
│  │           ▼                                             │    │
│  │  Repository Layer (Data access)                         │    │
│  │  ├─ user_repository.go                                  │    │
│  │  └─ book_repository.go                                  │    │
│  └─────────────────────────────────────────────────────────┘    │
│                       Port: 3000                                  │
└──────────┬──────────────────────────────────────────┬────────────┘
           │                                          │
           │ BSON                                     │ BSON
           │                                          │
    ┌──────▼──────┐                         ┌─────────▼────────┐
    │   MongoDB   │                         │  Elasticsearch   │
    │  Collections │                         │   Index          │
    │  ├─ users   │                         │  ├─ books        │
    │  └─ books   │                         │  └─ (full-text) │
    │  Port: 27017│                         │  Port: 9200      │
    └─────────────┘                         └──────────────────┘
                                                    │
                                            ┌───────▼────────┐
                                            │    Kibana       │
                                            │  (Visualization)│
                                            │  Port: 5601     │
                                            └─────────────────┘

           ┌─────────────────────────────────────────┐
           │         Observability Stack              │
           │                                          │
           │ OpenTelemetry (OTEL)                    │
           │      │                                  │
           │      ▼                                  │
           │  OTEL HTTP Exporter                    │
           │      │                                  │
           │      ▼                                  │
           │   Jaeger (Distributed Tracing)         │
           │   Port: 6831 (UDP), 16686 (UI)         │
           │                                          │
           └─────────────────────────────────────────┘
```

## Component Details

### API Server (Fiber)

**Port:** 3000

**Health Check:**
```bash
curl http://localhost:3000/health
# Response: {"status":"ok","timestamp":"2024-01-29T..."}
```

**Key Features:**
- RESTful endpoints for User and Book resources
- Automatic request ID generation (X-Request-ID header)
- CORS enabled for browser and Postman requests
- Panic recovery with error responses
- Structured JSON logging
- OpenTelemetry span creation for all requests

### MongoDB

**Port:** 27017

**Collections:**
- `users` - User accounts and profiles
- `books` - Book catalog with metadata

**Features:**
- Document-oriented database
- Automatic ObjectID generation
- Flexible schema
- Support for nested documents and arrays

**Connection String:**
```
mongodb://root:password@localhost:27017
```

### Elasticsearch

**Port:** 9200

**Indices:**
- `books` - Full-text search index for books

**Features:**
- Inverted index for fast full-text search
- Analyzer for title and author fields
- Automatic tokenization and stemming

**Sample Query:**
```bash
curl -X GET "localhost:9200/books/_search" -H "Content-Type: application/json" -d'
{
  "query": {
    "multi_match": {
      "query": "go",
      "fields": ["title", "author"]
    }
  }
}'
```

### Jaeger (Distributed Tracing)

**Ports:**
- 6831/UDP - Jaeger agent (trace receiver)
- 16686/HTTP - Jaeger UI (http://localhost:16686)

**Features:**
- Distributed tracing across services
- Latency visualization
- Error tracking and correlation
- Service dependency graphs

**Key Traces:**
- `CreateUser` - User creation with database write
- `GetUser` - User retrieval with database read
- `CreateBook` - Book creation with dual-write (MongoDB + Elasticsearch)
- `SearchBooks` - Elasticsearch search operation

## Deployment Strategies

### 1. Local Development (Current Setup)

**Start Services:**
```bash
# Terminal 1: Start all Docker containers
docker-compose up -d

# Terminal 2: Start Go application
go run main.go

# Terminal 3: Monitor logs
docker-compose logs -f
```

**Verify:**
```bash
# Check containers running
docker-compose ps

# Check API health
curl http://localhost:3000/health

# Check Jaeger
open http://localhost:16686

# Check Kibana
open http://localhost:5601
```

### 2. Docker Deployment

**Build Docker Image:**
```dockerfile
# Dockerfile for go-elastic
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .

EXPOSE 3000
CMD ["./main"]
```

**Build and Run:**
```bash
# Build image
docker build -t go-elastic:latest .

# Run with docker-compose (recommended)
docker-compose up -d --build

# Or run standalone
docker run -p 3000:3000 \
  -e MONGODB_URI="mongodb://mongo:27017" \
  -e ELASTICSEARCH_URL="http://elasticsearch:9200" \
  go-elastic:latest
```

### 3. Kubernetes Deployment

**Basic Deployment YAML:**
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: go-elastic-config
data:
  MONGODB_URI: "mongodb://mongo-service:27017"
  ELASTICSEARCH_URL: "http://elasticsearch-service:9200"
  JAEGER_ENDPOINT: "http://jaeger-service:14268/api/traces"

---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-elastic
spec:
  replicas: 3
  selector:
    matchLabels:
      app: go-elastic
  template:
    metadata:
      labels:
        app: go-elastic
    spec:
      containers:
      - name: go-elastic
        image: go-elastic:latest
        ports:
        - containerPort: 3000
        envFrom:
        - configMapRef:
            name: go-elastic-config
        livenessProbe:
          httpGet:
            path: /health
            port: 3000
          initialDelaySeconds: 5
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 3000
          initialDelaySeconds: 5
          periodSeconds: 5

---
apiVersion: v1
kind: Service
metadata:
  name: go-elastic-service
spec:
  selector:
    app: go-elastic
  type: LoadBalancer
  ports:
  - protocol: TCP
    port: 80
    targetPort: 3000
```

**Deploy:**
```bash
kubectl apply -f deployment.yaml
```

### 4. Cloud Deployment (Azure)

**Using Azure App Service:**

```bash
# 1. Create resource group
az group create --name go-elastic-rg --location eastus

# 2. Create App Service Plan
az appservice plan create --name go-elastic-plan \
  --resource-group go-elastic-rg \
  --sku B2 --is-linux

# 3. Create App Service
az webapp create --resource-group go-elastic-rg \
  --plan go-elastic-plan \
  --name go-elastic-app \
  --runtime "GO|1.24"

# 4. Deploy from GitHub
az webapp up --resource-group go-elastic-rg \
  --name go-elastic-app \
  --runtime "GO|1.24"

# 5. Configure environment variables
az webapp config appsettings set --resource-group go-elastic-rg \
  --name go-elastic-app \
  --settings MONGODB_URI="..." ELASTICSEARCH_URL="..."
```

## Performance Considerations

### 1. Database Optimization

**MongoDB Indexes:**
```go
// Create compound index for book search
collection.Indexes().CreateOne(ctx, mongo.IndexModel{
    Keys: bson.D{
        {Key: "title", Value: "text"},
        {Key: "author", Value: "text"},
    },
})
```

**Elasticsearch Optimization:**
- Use appropriate analyzers
- Set refresh interval appropriately
- Monitor disk space

### 2. Connection Pooling

**MongoDB (Default):**
- 128 connections in pool
- Configurable via connection string

**Elasticsearch:**
- HTTP connection reuse
- Default pool size: 256

### 3. Caching Strategy

Consider implementing caching for:
- Frequently accessed books
- User profile data
- Search results

**Redis Integration Example:**
```go
import "github.com/redis/go-redis/v9"

// Cache book
cache.Set(ctx, "book:"+id.Hex(), bookJSON, 1*time.Hour)

// Check cache first
cached := cache.Get(ctx, "book:"+id.Hex())
```

## Monitoring & Observability

### Jaeger Traces

View at http://localhost:16686

**Key Metrics:**
- Operation duration
- Error rates per service
- Service dependency graph

### Elasticsearch Monitoring

```bash
# Cluster health
curl http://localhost:9200/_cluster/health

# Indices status
curl http://localhost:9200/_cat/indices?v

# Search performance
curl http://localhost:9200/books/_stats
```

### MongoDB Monitoring

```javascript
// In mongosh
db.currentOp()  // Running operations
db.serverStatus()  // Server statistics
db.stats()  // Database stats
```

### Application Logs

**Log Levels:**
- DEBUG: Detailed tracing information
- INFO: General information
- WARN: Warning messages
- ERROR: Error messages

**Log Format:**
```json
{
  "timestamp": "2024-01-29T12:34:56Z",
  "level": "INFO",
  "message": "Book created successfully",
  "book_id": "507f1f77bcf86cd799439011",
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

## Scaling Strategies

### Horizontal Scaling

**Load Balancing:**
```yaml
# Docker Compose with multiple replicas
services:
  api-1:
    build: .
    ports: ["3001:3000"]
  api-2:
    build: .
    ports: ["3002:3000"]
  api-3:
    build: .
    ports: ["3003:3000"]
  
  nginx:
    image: nginx:latest
    ports: ["80:80"]
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
```

### Database Scaling

**MongoDB Sharding:**
```javascript
// Enable sharding on database
sh.enableSharding("go_elastic")

// Shard books collection
sh.shardCollection("go_elastic.books", { "_id": "hashed" })
```

**Elasticsearch Sharding:**
- Automatically handles sharding
- Default: 1 shard, 1 replica
- Adjust based on data size

## Environment Variables

**Essential:**
```env
# Server
PORT=3000

# MongoDB
MONGODB_URI=mongodb://root:password@localhost:27017
MONGODB_DB_NAME=go_elastic

# Elasticsearch
ELASTICSEARCH_URL=http://localhost:9200

# Tracing
OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4318
JAEGER_ENDPOINT=http://localhost:14268/api/traces
```

**Optional:**
```env
# Logging
LOG_LEVEL=info

# Performance
MONGODB_POOL_SIZE=128
REQUEST_TIMEOUT=30s
```

## Disaster Recovery

### Backup Strategy

**MongoDB Backup:**
```bash
# Full backup
mongodump --uri="mongodb://root:password@localhost:27017" \
  --out=/backups/mongodb-$(date +%Y%m%d)

# Restore
mongorestore --uri="mongodb://root:password@localhost:27017" \
  /backups/mongodb-20240129
```

**Elasticsearch Backup:**
```bash
# Register snapshot repository
curl -X PUT "localhost:9200/_snapshot/backup" -H "Content-Type: application/json" -d'
{
  "type": "fs",
  "settings": {
    "location": "/mnt/elasticsearch-backup"
  }
}'

# Create snapshot
curl -X PUT "localhost:9200/_snapshot/backup/snapshot_1"

# Restore
curl -X POST "localhost:9200/_snapshot/backup/snapshot_1/_restore"
```

### Failover Configuration

**MongoDB Replica Set:**
```bash
# Initialize replica set in mongosh
rs.initiate({
  _id: "rs0",
  members: [
    {_id: 0, host: "mongodb-1:27017"},
    {_id: 1, host: "mongodb-2:27017"},
    {_id: 2, host: "mongodb-3:27017"}
  ]
})
```

---

**Last Updated:** January 29, 2026
