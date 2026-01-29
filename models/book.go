package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Book struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title       string             `bson:"title,omitempty" json:"title,omitempty"`
	Author      string             `bson:"author,omitempty" json:"author,omitempty"`
	ISBN        string             `bson:"isbn,omitempty" json:"isbn,omitempty"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	Publisher   string             `bson:"publisher,omitempty" json:"publisher,omitempty"`
	PublishDate time.Time          `bson:"publish_date,omitempty" json:"publish_date,omitempty"`
	Pages       int                `bson:"pages,omitempty" json:"pages,omitempty"`
	Language    string             `bson:"language,omitempty" json:"language,omitempty"`
	CreatedAt   time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt   time.Time          `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}
