# MongoDB Migration Guide

## Overview

This project was migrated from PostgreSQL (with GORM ORM) to MongoDB as the primary database. This guide explains the migration, benefits, and usage.

## What Changed

### Before (PostgreSQL + GORM)
```
├── SQL Database: PostgreSQL
├── ORM: GORM (database/sql wrapper)
├── Queries: SQL strings with prepared statements
├── Migrations: GORM AutoMigrate
├── Models: GORM tags (gorm:"...")
└── Transactions: SQL transactions via GORM
```

### After (MongoDB Native)
```
├── NoSQL Database: MongoDB
├── Driver: Official MongoDB Go Driver
├── Queries: BSON filters and aggregations
├── Migrations: Manual index creation
├── Models: BSON tags (bson:"...")
└── Transactions: MongoDB sessions
```

## Key Differences

### Data Model

**PostgreSQL (GORM)**
```go
type User struct {
    ID        uint           `gorm:"primaryKey"`
    Username  string         `gorm:"unique;not null"`
    Email     string         `gorm:"unique;not null"`
    DeletedAt gorm.DeletedAt `gorm:"index"`
}
```

**MongoDB (Native)**
```go
type User struct {
    ID        primitive.ObjectID `bson:"_id,omitempty"`
    Username  string             `bson:"username"`
    Email     string             `bson:"email"`
}
```

### Query Examples

**Get User by ID - PostgreSQL**
```go
var user User
db.First(&user, id)  // SQL: SELECT * FROM users WHERE id = ?
```

**Get User by ID - MongoDB**
```go
var user User
collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
```

**Create User - PostgreSQL**
```go
db.Create(&user)  // SQL: INSERT INTO users ...
```

**Create User - MongoDB**
```go
collection.InsertOne(ctx, user)
```

## Benefits of MongoDB

### Flexibility
- **Flexible Schema**: No predefined schema required
- **Nested Documents**: Store complex relationships naturally
- **Arrays**: Store multiple values in single field

### Scalability
- **Horizontal Scaling**: Built-in sharding support
- **Replication**: Automatic failover with replica sets
- **Performance**: Optimized for document-based queries

### Development Speed
- **Rapid Iteration**: No migrations needed for schema changes
- **Natural Data Model**: JSON-like documents match application objects
- **Easy Nesting**: Related data in single document

## Migration Files Changed

### Database Layer
- **database/database.go** - MongoDB connection instead of PostgreSQL
- Uses MongoDB client and native driver
- Connection string: `mongodb://root:password@localhost:27017`

### Models
- **models/user.go** - BSON tags instead of GORM tags
- Uses `primitive.ObjectID` instead of `uint`
- No soft delete support (removed DeletedAt)

### Repository Layer
- **repository/user_repository.go** - MongoDB queries instead of GORM
- Uses BSON filters: `bson.M{"_id": objectID}`
- Methods now work with string IDs (hex format)

### Service Layer
- **service/user_service.go** - No GORM transactions
- Simplified to use repository directly
- Method signatures updated for MongoDB

### API Handlers
- **handler/user_handler.go** - String IDs instead of uint
- No `strconv.Atoi()` conversion
- IDs passed as hex strings

## Setup Instructions

### 1. Install MongoDB

**Docker (Recommended)**
```bash
docker-compose up -d mongodb
```

**Local Installation**
```bash
# Windows - chocolatey
choco install mongodb

# macOS
brew tap mongodb/brew
brew install mongodb-community

# Linux
sudo apt-get install -y mongodb
```

### 2. Configure Environment

Create `.env` file:
```env
MONGODB_URI=mongodb://root:password@localhost:27017
MONGODB_DB_NAME=go_logger
```

### 3. Start Application

```bash
go run main.go
```

Expected output:
```
MongoDB connection established
```

## Working with MongoDB

### Connect to MongoDB

**Using mongosh (CLI)**
```bash
mongosh "mongodb://root:password@localhost:27017"

# Inside mongosh:
use go_logger
db.users.find()
db.books.find().pretty()
```

**Using MongoDB Compass (GUI)**
1. Download from https://www.mongodb.com/products/tools/compass
2. Connect to: `mongodb://root:password@localhost:27017`
3. Browse collections visually

### Create Indexes

For better performance on frequently searched fields:

```go
opts := options.Index().SetUnique(true)
collection.Indexes().CreateOne(ctx, mongo.IndexModel{
    Keys:    bson.D{{Key: "username", Value: 1}},
    Options: opts,
})
```

### Query Examples

**Find One**
```go
var user User
err := collection.FindOne(ctx, bson.M{"username": "john"}).Decode(&user)
```

**Find Many**
```go
cursor, err := collection.Find(ctx, bson.M{"language": "English"})
var books []Book
err = cursor.All(ctx, &books)
```

**Find with Filter**
```go
filter := bson.M{
    "author": bson.M{"$regex": "John", "$options": "i"},
    "pages": bson.M{"$gte": 400},
}
cursor, err := collection.Find(ctx, filter)
```

**Update**
```go
update := bson.M{"$set": bson.M{"title": "New Title"}}
result, err := collection.UpdateOne(ctx, bson.M{"_id": id}, update)
```

**Delete**
```go
result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
```

## Common Migration Issues

### Issue: "Unique constraint violation"

MongoDB doesn't enforce UNIQUE automatically. Create unique indexes:

```go
opts := options.Index().SetUnique(true)
collection.Indexes().CreateOne(ctx, mongo.IndexModel{
    Keys:    bson.D{{Key: "email", Value: 1}},
    Options: opts,
})
```

### Issue: "No automatic timestamps"

Create manually in code:

```go
book.CreatedAt = time.Now()
book.UpdatedAt = time.Now()
collection.InsertOne(ctx, book)
```

### Issue: "Soft deletes don't work"

Instead of soft delete:

**Option 1: Add deleted flag**
```go
type Book struct {
    ID        primitive.ObjectID `bson:"_id,omitempty"`
    Title     string             `bson:"title"`
    Deleted   bool               `bson:"deleted"`
    DeletedAt time.Time          `bson:"deleted_at"`
}

// Query non-deleted
filter := bson.M{"deleted": false}
```

**Option 2: Use separate archive collection**
```go
// Move to archive collection before deleting
archiveCollection.InsertOne(ctx, book)
booksCollection.DeleteOne(ctx, bson.M{"_id": id})
```

## Performance Tips

### 1. Create Appropriate Indexes

```go
// Index for search
collection.Indexes().CreateOne(ctx, mongo.IndexModel{
    Keys: bson.D{
        {Key: "title", Value: 1},
        {Key: "author", Value: 1},
    },
})
```

### 2. Use Projection to Limit Fields

```go
opts := options.FindOne().SetProjection(bson.M{"title": 1, "author": 1})
collection.FindOne(ctx, filter, opts)
```

### 3. Batch Operations

```go
models := []mongo.WriteModel{
    mongo.NewInsertOneModel().SetDocument(book1),
    mongo.NewInsertOneModel().SetDocument(book2),
}
result, err := collection.BulkWrite(ctx, models)
```

### 4. Connection Pooling

Already configured in MongoDB driver (default 128 connections).

## Troubleshooting

### MongoDB Connection Failed

```bash
# Check if MongoDB is running
docker-compose ps mongodb

# Check logs
docker-compose logs mongodb

# Test connection manually
mongosh "mongodb://root:password@localhost:27017"
```

### ObjectID Format Issues

MongoDB ObjectIDs are 24-character hex strings:
- ✅ `507f1f77bcf86cd799439011`
- ❌ `12345` (too short)
- ❌ `not-a-valid-id`

Convert properly:
```go
// String to ObjectID
id, err := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")

// ObjectID to String
idString := id.Hex()
```

### Memory Usage High

Ensure connections are closed:

```go
defer client.Disconnect(context.Background())
defer cursor.Close(ctx)
```

## Next Steps

1. **Create Indexes** for frequently searched fields
2. **Set up Backups** for production MongoDB
3. **Monitor Performance** using MongoDB Charts
4. **Implement Transactions** if needed (requires replica set)
5. **Plan for Sharding** as data grows

## References

- [MongoDB Go Driver Docs](https://pkg.go.dev/go.mongodb.org/mongo-driver)
- [MongoDB Manual](https://docs.mongodb.com/manual/)
- [Query and Projection Operators](https://docs.mongodb.com/manual/reference/operator/query/)

---

**Last Updated:** January 29, 2026
