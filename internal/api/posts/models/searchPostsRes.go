package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type SearchPosts struct {
	Id          *primitive.ObjectID `json:"id,omitempty"  bson:"_id,omitempty"`
	Author      string              `json:"author,omitempty" bson:"author"`
	Content     string              `json:"content,omitempty" bson:"content"`
	Title       string              `json:"title,omitempty" bson:"title"`
	Updated     *time.Time          `json:"updated,omitempty" bson:"updated"`
	CommentsCnt int                 `json:"comments_cnt,omitempty" bson:"comments_cnt"`
	Username    string              `json:"username,omitempty" bson:"username"`
}
type SearchPostsRes struct {
	Posts []SearchPosts `json:"posts,omitempty"`
	Total int           `json:"total,omitempty" bson:"total"`
}
