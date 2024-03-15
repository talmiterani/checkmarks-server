package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
	"time"
)

type Post struct {
	Id      *primitive.ObjectID `json:"id,omitempty"  bson:"_id,omitempty"`
	UserId  *primitive.ObjectID `json:"user_id,omitempty"  bson:"_userId,omitempty"`
	Author  string              `json:"author,omitempty" bson:"author"`
	Content string              `json:"content,omitempty" bson:"content"`
	Title   string              `json:"title,omitempty" bson:"title"`
	Updated *time.Time          `json:"updated,omitempty" bson:"updated"`
}

func (p *Post) Validate(validateId bool) string {

	if p.Author == "" {
		return "missing author"
	}

	if p.Title == "" {
		return "missing title"
	}

	if p.Content == "" {
		return "missing content"
	}

	if validateId && p.Id == nil {
		return "missing post id"
	}
	return ""
}

func (p *Post) Prepare() {
	p.Author = strings.TrimSpace(p.Author)
	p.Title = strings.TrimSpace(p.Title)
	p.Content = strings.TrimSpace(p.Content)
}
