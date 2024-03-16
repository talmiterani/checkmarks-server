package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
	"time"
)

type Comment struct {
	Content string              `json:"content" bson:"content,omitempty"`
	Id      *primitive.ObjectID `json:"id,omitempty"  bson:"_id,omitempty"`
	PostId  *primitive.ObjectID `json:"post_id,omitempty"  bson:"_postId,omitempty"`
	UserId  *primitive.ObjectID `json:"user_id,omitempty"  bson:"_userId,omitempty"`
	Updated *time.Time          `json:"updated,omitempty" bson:"updated"`
}

func (c *Comment) Validate(validatePostId, validateId bool) string {

	if c.Content == "" {
		return "missing content"
	}

	if validatePostId && c.PostId == nil {
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
