package users

import (
	"context"
	usersDTO "social-media-app/internal/shared/dto/users"

	"golang.org/x/crypto/bcrypt"
)

type service struct {
	repository repository
}

func NewService(repository repository) service {
	return service{
		repository: repository,
	}
}

func (s *service) GetByID(ctx context.Context, userID int64) (User, error) {
	user, err := s.repository.GetByID(ctx, userID)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (s *service) Follow(ctx context.Context, payload usersDTO.FollowPayload) error {
	return s.repository.Follow(ctx, payload)
}

func (s *service) UnFollow(ctx context.Context, payload usersDTO.FollowPayload) error {
	return s.repository.UnFollow(ctx, payload)
}

func (s *service) Register(ctx context.Context, payload usersDTO.RegisterUserPayload) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		Firstname: payload.Firstname,
		Lastname:  payload.Lastname,
		Email:     payload.Email,
		Username:  payload.Username,
		Password:  string(hash),
	}

	if err := s.repository.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
