package comments

import (
	"context"
	commentsDTO "social-media-app/internal/shared/dto/comments"
)

type service struct {
	repository repository
}

func NewService(repository repository) service {
	return service{repository: repository}
}

func (s *service) GetByPostID(ctx context.Context, postID int64) ([]Comment, error) {
	return s.repository.GetByPostID(ctx, postID)
}

func (s *service) Create(
	ctx context.Context,
	postID, userID int64,
	payload commentsDTO.CreateCommentPayload,
) (Comment, error) {
	comment := Comment{
		Content: payload.Content,
		UserID:  userID,
		PostID:  postID,
	}

	err := s.repository.Create(ctx, &comment)
	if err != nil {
		return Comment{}, err
	}

	return comment, nil
}
