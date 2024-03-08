package access

import (
	"awesomeProject/internal/config"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	Collection *mongo.Collection
}

func initMongoConnection(c config.MongoConfig) (*Mongo, error) {
	//client options
	clientOptions := options.Client().ApplyURI(c.ConnectionString)

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		return nil, err
	}

	fmt.Println("MongoDB connection success")

	collection := client.Database(c.DbName).Collection(c.ColName)

	fmt.Println("Collection instance is ready")

	return &Mongo{collection}, nil
}
