package models

import (
	"awesomeProject/internal/api/comments/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Post struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Author       string             `bson:"author"`
	Content      string             `bson:"content"`
	CreationDate time.Time          `bson:"creation_date"`
	Comments     []models.Comment   `bson:"comments"`
	Title        string             `bson:"title"`
}
