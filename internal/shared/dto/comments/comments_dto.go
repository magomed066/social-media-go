package commentsDTO

type CreateCommentPayload struct {
	Content string `json:"content" validate:"required,max=300" maxLength:"300" example:"Thanks for sharing this post!"`
}
