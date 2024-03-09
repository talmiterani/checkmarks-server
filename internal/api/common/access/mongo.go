package access

import (
	"checkmarks/internal/config"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	Posts    *mongo.Collection
	Comments *mongo.Collection
	Users    *mongo.Collection
}

func initMongoConnection(c config.MongoConfig) (*Mongo, error) {
	//client options
	clientOptions := options.Client().ApplyURI(c.ConnectionString)

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		return nil, err
	}

	fmt.Println("MongoDB connection success")

	mongoInstance := &Mongo{
		Posts:    client.Database(c.DbName).Collection(c.Collections.Posts),
		Comments: client.Database(c.DbName).Collection(c.Collections.Comments),
		Users:    client.Database(c.DbName).Collection(c.Collections.Users),
	}

	fmt.Println("Collections instance is ready")

	return mongoInstance, nil
}
