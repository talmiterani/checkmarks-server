package repo

import (
	"checkmarks/internal/api/comments/models"
	"checkmarks/internal/api/common/access"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CommentsDb struct {
	*access.DbConnections
}

func New(sdc *access.DbConnections) CommentsRepo {
	return &CommentsDb{sdc}
}

func (c *CommentsDb) GetByPostId(ctx context.Context, postId string) ([]models.Comment, error) {

	filter := bson.M{"postId": postId}
	sortOptions := options.Find().SetSort(bson.D{{"updated", -1}})

	cur, err := c.Mongo.Comments.Find(ctx, filter, sortOptions)

	if err != nil {
		return nil, err
	}

	var comments []models.Comment

	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var comment models.Comment

		err = cur.Decode(&comment)

		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, err
}

func (c *CommentsDb) Add(ctx context.Context, comment *models.Comment) (*primitive.ObjectID, error) {

	inserted, err := c.Mongo.Comments.InsertOne(ctx, comment)
	if err != nil {
		return nil, err
	}

	insertedID, ok := inserted.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, errors.New("failed to convert InsertedID to ObjectID")
	}

	fmt.Println("added new comment: ", inserted)

	return &insertedID, nil

}

func (c *CommentsDb) Update(ctx context.Context, comment *models.Comment) error {

	filter := bson.M{"_id": comment.Id}
	update := bson.M{"$set": bson.M{
		"content": comment.Content,
		"updated": comment.Updated,
	}}

	res, err := c.Mongo.Comments.UpdateOne(ctx, filter, update)

	if err != nil {
		return err
	}

	fmt.Println("updated comment: ", res)

	return nil
}

func (c *CommentsDb) Delete(ctx context.Context, commentId string) error {

	id, err := primitive.ObjectIDFromHex(commentId)

	if err != nil {
		return err
	}

	filter := bson.M{"_id": id}

	deleteCnt, err := c.Mongo.Comments.DeleteOne(ctx, filter)

	if err != nil {
		return err
	}
	fmt.Printf("deleted post count: %v, comment id: %s", deleteCnt, commentId)

	return nil
}
