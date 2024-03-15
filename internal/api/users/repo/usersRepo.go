package repo

import (
	"checkmarks/pkg/users"
	"context"
)

type UsersRepo interface {
	Signup(ctx context.Context, user *users.User) error
	CheckUniqueUsername(ctx context.Context, username string) (bool, error)
	Get(ctx context.Context, req *users.User) (*users.User, error)
}
