package repo

import (
	"awesomeProject/internal/api/common/access"
	"awesomeProject/internal/api/posts/models"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostsDb struct {
	*access.DbConnections
}

func New(sdc *access.DbConnections) PostsRepo {
	return &PostsDb{sdc}
}

func (p *PostsDb) GetAll(ctx context.Context) ([]models.Post, error) {

	cur, err := p.Mongo.Collection.Find(ctx, bson.D{{}})

	if err != nil {
		return nil, err
	}

	var posts []models.Post

	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var post models.Post

		err = cur.Decode(&post)

		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, err
}

func (p *PostsDb) Add(ctx context.Context, post *models.Post) (*primitive.ObjectID, error) {

	inserted, err := p.Mongo.Collection.InsertOne(ctx, post)
	if err != nil {
		return nil, err
	}

	insertedID, ok := inserted.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, errors.New("failed to convert InsertedID to ObjectID")
	}

	fmt.Println("added new post: ", inserted)

	return &insertedID, nil
}

func (p *PostsDb) Update(ctx context.Context, post *models.Post) error {

	filter := bson.M{"_id": post.Id}
	update := bson.D{{"$set", post}}

	res, err := p.Mongo.Collection.UpdateOne(ctx, filter, update)

	if err != nil {
		return err
	}

	fmt.Println("updated post: ", res)

	return nil
}

func (p *PostsDb) Delete(ctx context.Context, postId string) error {

	id, err := primitive.ObjectIDFromHex(postId)

	if err != nil {
		return err
	}

	filter := bson.M{"_id": id}

	deleteCnt, err := p.Mongo.Collection.DeleteOne(ctx, filter)

	if err != nil {
		return err
	}
	fmt.Printf("deleted post count: %v, post id: %s", deleteCnt, postId)

	return nil
}
