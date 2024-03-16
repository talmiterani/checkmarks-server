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

func (s *Service) signup(ctx context.Context, user *users.User) error {

	newId, err := s.repo.Signup(ctx, user)

	if err != nil {
		return err
	}

	user.Id = newId
	return nil
}

func (s *Service) checkUniqueUsername(ctx context.Context, username string) (bool, error) {
	return s.repo.CheckUniqueUsername(ctx, username)
}

func (s *Service) get(ctx context.Context, req *users.User) (*users.User, error) {
	return s.repo.Get(ctx, req)
}
