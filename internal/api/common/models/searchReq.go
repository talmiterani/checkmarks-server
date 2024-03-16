package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type SearchReq struct {
	PageSize int                 `json:"page_size"`
	Page     int                 `json:"page"`
	Query    string              `json:"query,omitempty"`
	UserId   *primitive.ObjectID `json:"user_id" bson:"_userId"`
}
