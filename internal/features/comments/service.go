package comments

import "context"

type service struct {
	repository repository
}

func NewService(repository repository) service {
	return service{repository: repository}
}

func (s *service) GetByPostID(ctx context.Context, postID int64) ([]Comment, error) {
	return s.repository.GetByPostID(ctx, postID)
}
