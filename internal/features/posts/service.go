package posts

import (
	"context"
	postsDTO "social-media-app/internal/dto/posts"
)

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
