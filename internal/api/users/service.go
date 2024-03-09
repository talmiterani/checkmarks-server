package users

import (
	"checkmarks/internal/api/common/access"
	"checkmarks/internal/api/users/repo"
	"checkmarks/pkg/users"
	"context"
)

type Service struct {
	*access.DbConnections
	repo repo.UsersRepo
}

func NewService(sdc *access.DbConnections) *Service {
	return &Service{
		sdc,
		repo.New(sdc),
	}
}

func (s *Service) add(ctx context.Context, user *users.User) error {
	return s.repo.Add(ctx, user)
}
