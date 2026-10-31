package roles

import "context"

type service struct {
	repository repository
}

func NewService(repo repository) service {
	return service{
		repository: repo,
	}
}

func (s *service) GetByName(ctx context.Context, roleName string) (Role, error) {
	return s.repository.GetByName(ctx, roleName)
}
