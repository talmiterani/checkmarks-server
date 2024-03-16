package repo

import (
	"checkmarks/pkg/users"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UsersRepo interface {
	Signup(ctx context.Context, user *users.User) (*primitive.ObjectID, error)
	CheckUniqueUsername(ctx context.Context, username string) (bool, error)
	Get(ctx context.Context, req *users.User) (*users.User, error)
}
