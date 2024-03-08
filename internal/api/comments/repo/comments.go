package repo

import (
	"awesomeProject/internal/api/common/access"
	"awesomeProject/internal/api/posts/models"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
)

type CommentsDb struct {
	*access.DbConnections
}

func New(sdc *access.DbConnections) CommentsRepo {
	return &CommentsDb{sdc}
}

func (p *CommentsDb) GetComments(ctx context.Context, postId string) ([]primitive.M, error) {

	postID, err := primitive.ObjectIDFromHex(postId)
	if err != nil {
		log.Fatal(err)
	}

	// Find the post by ID
	var post models.Post
	err = p.Mongo.Collection.FindOne(ctx, bson.M{"_id": postID}).Decode(&post)
	if err != nil {
		log.Fatal(err)
	}

	// Print all comments for the post
	for _, comment := range post.Comments {
		fmt.Printf("Author: %s\nContent: %s\nCreation Date: %s\n\n", comment.Author, comment.Content, comment.CreationDate)
	}
	return nil, nil
}
