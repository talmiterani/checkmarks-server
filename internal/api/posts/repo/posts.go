package repo

import (
	"checkmarks/internal/api/common/access"
	commonModels "checkmarks/internal/api/common/models"
	"checkmarks/internal/api/posts/models"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
)

type PostsDb struct {
	*access.DbConnections
}

func New(sdc *access.DbConnections) PostsRepo {
	return &PostsDb{sdc}
}

func (p *PostsDb) Search(ctx context.Context, req *commonModels.SearchReq) ([]models.SearchPostsRes, error) {

	pipeline := bson.A{
		bson.M{
			"$match": bson.M{
				"$or": bson.A{
					bson.M{"title": primitive.Regex{Pattern: req.Query, Options: "i"}},
					bson.M{"content": primitive.Regex{Pattern: req.Query, Options: "i"}},
				},
			},
		},
		bson.M{
			"$sort": bson.M{"updated": -1},
		},
		bson.M{
			"$skip": (req.Page - 1) * req.PageSize,
		},
		bson.M{
			"$limit": req.PageSize,
		},
		bson.M{
			"$lookup": bson.D{
				{"from", "comments"},
				{"localField", "_id"},
				{"foreignField", "_postId"},
				{"as", "comments"},
			},
		},
		bson.M{
			"$project": bson.M{
				"_id":          1,
				"author":       1,
				"title":        1,
				"content":      1,
				"updated":      1,
				"comments_cnt": bson.M{"$size": "$comments"},
			},
		},
	}

	cur, err := p.Mongo.Posts.Aggregate(ctx, pipeline)

	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var posts []models.SearchPostsRes
	for cur.Next(ctx) {
		var post models.SearchPostsRes
		err = cur.Decode(&post)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (p *PostsDb) Get(ctx context.Context, postId string) (*bson.M, error) {

	postID, err := primitive.ObjectIDFromHex(postId)
	if err != nil {

		if err.Error() == access.ObjectIDFromHexInvalidErr || strings.Contains(err.Error(), access.ObjectIDFromHexInvalidByte) {
			return nil, nil
		}
		return nil, err
	}

	// Define the pipeline with $lookup stage
	pipeline := bson.A{
		bson.M{
			"$match": bson.M{"_id": postID},
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "comments",
				"localField":   "_id",
				"foreignField": "_postId",
				"as":           "comments",
			},
		},
		bson.M{
			"$unwind": bson.M{
				"path":                       "$comments",
				"preserveNullAndEmptyArrays": true,
			},
		},
		bson.M{
			"$sort": bson.M{"comments.updated": -1},
		},
		bson.M{
			"$group": bson.M{
				"_id":      "$_id",
				"author":   bson.M{"$first": "$author"},
				"content":  bson.M{"$first": "$content"},
				"updated":  bson.M{"$first": "$updated"},
				"comments": bson.M{"$push": "$comments"},
				"title":    bson.M{"$first": "$title"},
			},
		},
	}

	// Use aggregate with the defined pipeline
	cur, err := p.DbConnections.Mongo.Posts.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	// Decode the result
	var results []bson.M
	if err := cur.All(ctx, &results); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, nil
	}

	return &results[0], nil
}

func (p *PostsDb) Add(ctx context.Context, post *models.Post) (*primitive.ObjectID, error) {

	inserted, err := p.Mongo.Posts.InsertOne(ctx, post)
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
	update := bson.M{"$set": bson.M{
		"author":  post.Author,
		"content": post.Content,
		"title":   post.Title,
		"updated": post.Updated,
	}}

	res, err := p.Mongo.Posts.UpdateOne(ctx, filter, update)

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

	deleteCnt, err := p.Mongo.Posts.DeleteOne(ctx, filter)

	if err != nil {
		return err
	}
	fmt.Printf("deleted post count: %v, post id: %s", deleteCnt, postId)

	// Filter to delete comments associated with the post
	commentsFilter := bson.M{"_postId": id}

	// Delete comments associated with the post
	deleteCommentsCnt, err := p.Mongo.Comments.DeleteMany(ctx, commentsFilter)
	if err != nil {
		return err
	}
	fmt.Printf("Deleted comments count: %v, post id: %s\n", deleteCommentsCnt, postId)

	return nil
}
