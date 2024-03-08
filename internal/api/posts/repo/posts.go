package repo

import (
	"awesomeProject/internal/api/common/access"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostsDb struct {
	*access.DbConnections
}

func New(sdc *access.DbConnections) PostsRepo {
	return &PostsDb{sdc}
}

func (p *PostsDb) GetMovies(ctx context.Context) ([]primitive.M, error) {

	cur, err := p.Mongo.Collection.Find(ctx, bson.D{{}})

	if err != nil {
		return nil, err
	}

	var movies []primitive.M

	defer cur.Close(context.Background())

	for cur.Next(ctx) {
		var movie primitive.M

		err = cur.Decode(&movie)

		if err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}

	return movies, err
}

//func insertOneItem(movie model.Netflix) {
//	inserted, err := collection.InsertOne(context.Background(), movie)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	fmt.Println("Inserted new item: ", inserted)
//}
//
//func updateOneMovie(movieId string) {
//	id, err := primitive.ObjectIDFromHex(movieId)
//
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	filter := bson.M{"_id": id}
//	update := bson.M{"$set": bson.M{"watched": true}}
//
//	res, err := collection.UpdateOne(context.Background(), filter, update)
//
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println("Updated item: ", res)
//}
//
//func deleteOneMovie(movieId string) {
//	id, err := primitive.ObjectIDFromHex(movieId)
//
//	if err != nil {
//		log.Fatal(err)
//	}
//	filter := bson.M{"_id": id}
//
//	deleteCnt, err := collection.DeleteOne(context.Background(), filter)
//
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println("Updated item count : ", deleteCnt)
//}
//
//func deleteAllMovies() {
//	res, err := collection.DeleteMany(context.Background(), bson.D{{}})
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println("Updated item count : ", res)
//}
