package repo

import (
	"checkmarks/pkg/users"
	"context"
)

type UsersRepo interface {
	Add(ctx context.Context, user *users.User) error
}
