package repo

import (
	"checkmarks/internal/api/common/access"
	"checkmarks/pkg/users"
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

type UsersDb struct {
	*access.DbConnections
}

func New(sdc *access.DbConnections) UsersRepo {
	return &UsersDb{sdc}
}

func (u *UsersDb) Signup(ctx context.Context, user *users.User) error {
	_, err := u.Mongo.Users.InsertOne(ctx, user)
	return err
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
