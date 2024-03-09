package models

import (
	"awesomeProject/internal/api/comments/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
	"time"
)

type Post struct {
	ID       *primitive.ObjectID `json:"id,omitempty"  bson:"_id,omitempty"`
	Author   string              `json:"author,omitempty" bson:"author"`
	Content  string              `json:"content,omitempty" bson:"content"`
	Comments []models.Comment    `json:"comments,omitempty" bson:"comments"`
	Title    string              `json:"title,omitempty" bson:"title"`
	Updated  *time.Time          `json:"updated,omitempty" bson:"updated"`
}

func (p *Post) Validate() string {

	if p.Author == "" {
		return "missing author"
	}

	if p.Title == "" {
		return "missing title"
	}

	if p.Content == "" {
		return "missing content"
	}
	return ""
}

func (p *Post) Prepare() {
	p.Title = strings.TrimSpace(p.Title)
	p.Content = strings.TrimSpace(p.Content)
}
