package commentsDTO

type CreateCommentPayload struct {
	Content string `json:"content" validate:"required,max=300"`
}
