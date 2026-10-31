package postsDTO

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,max=200" maxLength:"200" example:"My first post"`
	Content string   `json:"content" validate:"required,max=1000" maxLength:"1000" example:"Sharing what I learned today."`
	Tags    []string `json:"tags" example:"go,learning"`
}

// UpdatePostPayload changes only fields with non-null values.
type UpdatePostPayload struct {
	Title   *string `json:"title" validate:"omitempty,max=200" maxLength:"200" example:"Updated title"`
	Content *string `json:"content" validate:"omitempty,max=1000" maxLength:"1000" example:"Updated post content."`
}
