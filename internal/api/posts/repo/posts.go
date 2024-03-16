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

func (p *PostsDb) Search(ctx context.Context, req *commonModels.SearchReq) ([]models.SearchPosts, int, error) {

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
			"$lookup": bson.D{
				{"from", "comments"},
				{"localField", "_id"},
				{"foreignField", "_postId"},
				{"as", "comments"},
			},
		},
		bson.M{
			"$lookup": bson.M{
				"from": "users",
				"let":  bson.M{"userId": "$_userId"},
				"pipeline": bson.A{
					bson.M{"$match": bson.M{"$expr": bson.M{"$eq": bson.A{"$_id", "$$userId"}}}},
				},
				"as": "user",
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
				"username":     bson.M{"$arrayElemAt": bson.A{"$user.username", 0}}, // Extract username as string
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
	}
	if req.UserId != nil {
		matchStage := bson.M{
			"$match": bson.M{
				"_userId": req.UserId,
			},
		}
		pipeline = append(bson.A{matchStage}, pipeline...)
	}

	cur, err := p.Mongo.Posts.Aggregate(ctx, pipeline)

	if err != nil {
		return nil, 0, err
	}
	defer cur.Close(ctx)

	var posts []models.SearchPosts
	for cur.Next(ctx) {
		var post models.SearchPosts
		err = cur.Decode(&post)
		if err != nil {
			return nil, 0, err
		}
		posts = append(posts, post)
	}

	total, err := p.getTotalPostsCount(ctx, req.UserId, req.Query)

	return posts, total, nil
}

func (p *PostsDb) getTotalPostsCount(ctx context.Context, userId *primitive.ObjectID, query string) (int, error) {
	pipeline := bson.A{}

	if userId != nil {
		matchStage := bson.M{
			"$match": bson.M{"_userId": userId},
		}
		pipeline = append(pipeline, matchStage)
	}

	if query != "" {
		queryMatchStage := bson.M{
			"$match": bson.M{
				"$or": bson.A{
					bson.M{"title": primitive.Regex{Pattern: query, Options: "i"}},
					bson.M{"content": primitive.Regex{Pattern: query, Options: "i"}},
				},
			},
		}
		pipeline = append(pipeline, queryMatchStage)
	}

	countStage := bson.M{
		"$count": "total",
	}
	pipeline = append(pipeline, countStage)

	cur, err := p.Mongo.Posts.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cur.Close(ctx)

	var countResult struct {
		Total int `bson:"total"`
	}
	if cur.Next(ctx) {
		err := cur.Decode(&countResult)
		if err != nil {
			return 0, err
		}
	}

	return countResult.Total, nil
}
func (p *PostsDb) Get(ctx context.Context, postId string) (*bson.M, error) {

	postID, err := primitive.ObjectIDFromHex(postId)
	if err != nil {

		if err.Error() == access.ObjectIDFromHexInvalidErr || strings.Contains(err.Error(), access.ObjectIDFromHexInvalidByte) {
			return nil, nil
		}
		return nil, err
	}

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
			"$lookup": bson.M{
				"from":         "users",
				"localField":   "comments._userId",
				"foreignField": "_id",
				"as":           "commentUser",
			},
		},
		bson.M{
			"$sort": bson.M{"comments.updated": -1},
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "users",
				"localField":   "_userId",
				"foreignField": "_id",
				"as":           "postUser",
			},
		},
		bson.M{
			"$group": bson.M{
				"_id":     "$_id",
				"author":  bson.M{"$first": "$author"},
				"content": bson.M{"$first": "$content"},
				"updated": bson.M{"$first": "$updated"},
				"comments": bson.M{"$push": bson.M{
					"_id":      "$comments._id",
					"content":  "$comments.content",
					"updated":  "$comments.updated",
					"username": bson.M{"$arrayElemAt": bson.A{"$commentUser.username", 0}}, // Include comment author's username
					"userId":   bson.M{"$arrayElemAt": bson.A{"$commentUser._id", 0}},      // Include comment author's username
				}},
				"title":    bson.M{"$first": "$title"},
				"username": bson.M{"$first": "$postUser.username"}, // Include post author's username
			},
		},
	}

	cur, err := p.DbConnections.Mongo.Posts.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

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
