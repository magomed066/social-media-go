package postsDTO

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,max=200"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags"`
}

type UpdatePostPayload struct {
	Title   *string `json:"title" validate:"omitempty,max=200"`
	Content *string `json:"content" validate:"omitempty,max=1000"`
}
