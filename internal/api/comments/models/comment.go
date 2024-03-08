package models

import "time"

type Comment struct {
	Author       string    `bson:"author"`
	Content      string    `bson:"content"`
	CreationDate time.Time `bson:"creation_date"`
}
