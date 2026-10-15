package users

import "context"

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
