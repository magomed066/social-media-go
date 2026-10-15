package posts

import (
	"context"
	postsDTO "social-media-app/internal/dto/posts"
	"social-media-app/internal/features/comments"
)

type CommentsService interface {
	GetByPostID(ctx context.Context, postID int64) ([]comments.Comment, error)
}

type service struct {
	repository repository
}

func NewService(repository repository) service {
	return service{
		repository: repository,
	}
}

func (s *service) Create(ctx context.Context, payload postsDTO.CreatePostPayload) (Post, error) {
	post := Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		UserID:  3,
	}

	err := s.repository.Create(ctx, &post)

	if err != nil {
		return Post{}, err
	}

	return post, nil
}

func (s *service) GetByID(ctx context.Context, id int64) (*Post, error) {
	post, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (s *service) Delete(ctx context.Context, id int64) error {
	return s.repository.Delete(ctx, id)
}

func (s *service) Update(ctx context.Context, post *Post) error {
	return s.repository.Update(ctx, post)
}

func (s *service) GetUserFeed(ctx context.Context, userID int64) ([]PostWithMetaData, error) {
	return s.repository.GetUserFeed(ctx, userID)
}
