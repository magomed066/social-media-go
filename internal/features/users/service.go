package users

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"social-media-app/internal/features/roles"
	"social-media-app/internal/shared/consts"
	usersDTO "social-media-app/internal/shared/dto/users"

	"github.com/google/uuid"
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

func (s *service) Register(ctx context.Context, payload usersDTO.RegisterUserPayload) (*UserWithToken, error) {
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
		Role: roles.Role{
			Name: "user",
		},
	}

	plainToken := uuid.New().String()
	hashed := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hashed[:])

	userWithToken := UserWithToken{
		User:  user,
		Token: plainToken,
	}

	err = s.repository.CreateAndInvite(ctx, user, hashToken, consts.MailExpires)
	if err != nil {
		return nil, err
	}

	return &userWithToken, nil
}

func (s *service) ActivateUser(ctx context.Context, token string) error {
	return s.repository.ActivateUser(ctx, token)
}

func (s *service) Delete(ctx context.Context, userID int64) error {
	return s.repository.Delete(ctx, userID)
}

func (s *service) Login(ctx context.Context, payload usersDTO.LoginUserPayload) (User, error) {
	return s.repository.Login(ctx, payload)
}
