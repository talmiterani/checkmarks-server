package repo

import (
	"checkmarks/internal/api/common/access"
	"checkmarks/pkg/users"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UsersDb struct {
	*access.DbConnections
}

func New(sdc *access.DbConnections) UsersRepo {
	return &UsersDb{sdc}
}

func (u *UsersDb) Signup(ctx context.Context, user *users.User) (*primitive.ObjectID, error) {

	inserted, err := u.Mongo.Users.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	insertedID, ok := inserted.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, errors.New("failed to convert insertedID to objectID")
	}

	fmt.Println("added new user: ", inserted)

	return &insertedID, nil
}

func (u *UsersDb) CheckUniqueUsername(ctx context.Context, username string) (bool, error) {
	count, err := u.Mongo.Users.CountDocuments(ctx, bson.M{"username": username})

	if err != nil {
		return false, err
	}

	if count > 0 {
		return false, nil
	}
	return true, nil
}

func (u *UsersDb) Get(ctx context.Context, req *users.User) (*users.User, error) {
	user := users.User{}
	err := u.Mongo.Users.FindOne(ctx, bson.M{"username": req.Username}).Decode(&user)
	return &user, err
}
