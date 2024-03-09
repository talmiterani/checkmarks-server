package repo

import (
	"checkmarks/internal/api/common/access"
	"checkmarks/pkg/users"
	"context"
)

type UsersDb struct {
	*access.DbConnections
}

func New(sdc *access.DbConnections) UsersRepo {
	return &UsersDb{sdc}
}

func (u *UsersDb) Add(ctx context.Context, user *users.User) error {
	_, err := u.Mongo.Users.InsertOne(ctx, user)
	return err
}
