package service

import (
	"context"

	"github.com/begenov/effective-mobile-task/internal/model"
	"github.com/begenov/effective-mobile-task/internal/repository"
)

type Service struct {
	userRepository repository.UserRepository
}

func New(userRepository repository.UserRepository) *Service {
	return &Service{
		userRepository: userRepository,
	}
}

func (s *Service) CreateUser(ctx context.Context, user model.User) error {

	if user.Nationality == nil && (user.Nationality.ID <= 0 || len(user.Nationality.Name) == 0) {
		return model.ErrBadRequest
	}

	err := s.userRepository.Create(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdateUser(ctx context.Context, user model.User) error {
	if user.ID <= 0 {
		return model.ErrBadRequest
	}

	err := s.userRepository.Update(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteUser(ctx context.Context, userID int64) error {
	if userID <= 0 {
		return model.ErrBadRequest
	}

	err := s.userRepository.Delete(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetUsers(ctx context.Context, userID, limit, offset *int64, gender *int, nationality *string) ([]model.User, error) {

	users, err := s.userRepository.GetUsers(ctx, userID, limit, offset, gender, nationality)
	if err != nil {
		return nil, err
	}

	return users, nil
}
