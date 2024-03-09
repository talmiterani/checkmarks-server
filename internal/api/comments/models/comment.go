package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
	"time"
)

type Comment struct {
	Author  string              `bson:"author"`
	Content string              `bson:"content"`
	Id      *primitive.ObjectID `json:"id,omitempty"  bson:"_id,omitempty"`
	PostId  string              `json:"post_id,omitempty"`
	Updated *time.Time          `json:"updated,omitempty" bson:"updated"`
}

func (c *Comment) Validate(validateAuthor, validatePostId, validateId bool) string {

	if validateAuthor && c.Author == "" { //todo Author should be base on the connected user (author = username)
		return "missing author"
	}

	if c.Content == "" {
		return "missing content"
	}

	if validatePostId && c.PostId == "" {
		return "missing post id"
	}

	if validateId && c.Id == nil {
		return "missing id"
	}
	return ""
}

func (c *Comment) Prepare() {
	c.Content = strings.TrimSpace(c.Content)
}
