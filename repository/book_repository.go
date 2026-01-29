package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go-elastic/database"
	"go-elastic/models"
	"time"

	"github.com/elastic/go-elasticsearch/v8/esapi"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookRepository interface {
	Create(ctx context.Context, book *models.Book) error
	FindByID(ctx context.Context, id string) (*models.Book, error)
	FindAll(ctx context.Context) ([]models.Book, error)
	SearchByTitle(ctx context.Context, title string) ([]models.Book, error)
	SearchByAuthor(ctx context.Context, author string) ([]models.Book, error)
}

type bookRepository struct {
	mongoCollection *mongo.Collection
	esClient        interface{}
}

func NewBookRepository(mongoCollection *mongo.Collection) BookRepository {
	return &bookRepository{
		mongoCollection: mongoCollection,
		esClient:        database.ESClient,
	}
}

// Create saves book to MongoDB and indexes it in Elasticsearch
func (r *bookRepository) Create(ctx context.Context, book *models.Book) error {
	if book.ID.IsZero() {
		book.ID = primitive.NewObjectID()
	}
	if book.CreatedAt.IsZero() {
		book.CreatedAt = time.Now()
	}
	if book.UpdatedAt.IsZero() {
		book.UpdatedAt = time.Now()
	}

	// Insert into MongoDB
	res, err := r.mongoCollection.InsertOne(ctx, book)
	if err != nil {
		return err
	}

	// Index in Elasticsearch
	esClient := database.ESClient
	bookJSON, err := json.Marshal(book)
	if err != nil {
		return err
	}

	req := esapi.IndexRequest{
		Index:      "books",
		DocumentID: book.ID.Hex(),
		Body:       bytes.NewReader(bookJSON),
	}

	resp, err := req.Do(ctx, esClient)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.IsError() {
		return err
	}

	book.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

// FindByID retrieves a book from MongoDB
func (r *bookRepository) FindByID(ctx context.Context, id string) (*models.Book, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var book models.Book
	err = r.mongoCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&book)
	if err != nil {
		return nil, err
	}
	return &book, nil
}

// FindAll retrieves all books from MongoDB
func (r *bookRepository) FindAll(ctx context.Context) ([]models.Book, error) {
	cursor, err := r.mongoCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var books []models.Book
	err = cursor.All(ctx, &books)
	return books, err
}

// SearchByTitle searches for books by title in Elasticsearch
func (r *bookRepository) SearchByTitle(ctx context.Context, title string) ([]models.Book, error) {
	// 1. สร้าง Query (เหมือนเดิม)
	query := map[string]interface{}{
		"size": 100,
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"title": title, // หรือใช้ "fuzziness": "AUTO" เพื่อกันพิมพ์ผิด
			},
		},
	}

	// Tip: ใช้ strings.Builder หรือ library สร้าง JSON จะเร็วกว่า map แต่ map อ่านง่ายกว่า (รับได้)
	queryJSON, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	req := esapi.SearchRequest{
		Index: []string{"books"},
		Body:  bytes.NewReader(queryJSON),
	}

	res, err := req.Do(ctx, database.ESClient) // ควรใช้ client จาก struct (Dependency Injection)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// 2. เช็ค Error จาก Server (แก้จุดที่ 1)
	if res.IsError() {
		return nil, fmt.Errorf("error searching documents: %s", res.String())
	}

	// 3. สร้าง Struct มารอรับ Response (แก้จุดที่ 2)
	// วิธีนี้เร็วที่สุด เพราะ Decode ทีเดียวลง Struct เลย ไม่ต้องแปลงไปมา
	type ESResponse struct {
		Hits struct {
			Hits []struct {
				Source models.Book `json:"_source"` // Mapping ตรงนี้สำคัญ
			} `json:"hits"`
		} `json:"hits"`
	}

	var response ESResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	// 4. ดึงข้อมูลออกมา
	var books []models.Book
	for _, hit := range response.Hits.Hits {
		books = append(books, hit.Source)
	}

	return books, nil
}

// SearchByAuthor searches for books by author in Elasticsearch
func (r *bookRepository) SearchByAuthor(ctx context.Context, author string) ([]models.Book, error) {
	// 1. สร้าง Query Map
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"author": author, // เปลี่ยนจาก title เป็น author
			},
		},
	}

	// แปลง Query เป็น JSON
	queryJSON, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("error marshaling query: %w", err)
	}

	// 2. สร้าง Request
	req := esapi.SearchRequest{
		Index: []string{"books"},
		Body:  bytes.NewReader(queryJSON),
	}

	// 3. ยิง Request (แนะนำให้ใช้ Client จาก struct r แทน Global variable)
	// เปลี่ยน database.ESClient เป็น r.esClient ถ้าคุณ inject client เข้ามาใน repository
	res, err := req.Do(ctx, database.ESClient)
	if err != nil {
		return nil, fmt.Errorf("error executing search request: %w", err)
	}
	defer res.Body.Close()

	// 4. เช็ค Error จากฝั่ง Elasticsearch (แก้จุดที่เคย Bug)
	if res.IsError() {
		return nil, fmt.Errorf("elasticsearch returned error: %s", res.String())
	}

	// 5. เตรียม Struct มารับผลลัพธ์ (Fast Decoding)
	// วิธีนี้เร็วกว่าการใช้ map[string]interface{} มากๆ
	type ESResponse struct {
		Hits struct {
			Hits []struct {
				Source models.Book `json:"_source"` // Mapping ตรงเข้า Model เลย
			} `json:"hits"`
		} `json:"hits"`
	}

	var response ESResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error parsing response body: %w", err)
	}

	// 6. ดึงข้อมูลออกมาใส่ Slice
	var books []models.Book
	for _, hit := range response.Hits.Hits {
		books = append(books, hit.Source)
	}

	return books, nil
}
