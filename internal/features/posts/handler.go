package posts

import (
	"errors"
	"net/http"
	"social-media-app/internal/json"
	postsDTO "social-media-app/internal/shared/dto/posts"
	mainErrors "social-media-app/internal/shared/errors/main"
	"social-media-app/internal/shared/pagination"
)

type Handler struct {
	service         service
	commentsService CommentsService
}

func NewHandler(service service, commentsService CommentsService) Handler {
	return Handler{
		service:         service,
		commentsService: commentsService,
	}
}

// Create godoc
//
//	@Summary		Create a post
//	@Description	Creates a post with a title, content, and optional tags. The current implementation assigns user ID 3 as the author instead of the authenticated user.
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		postsDTO.CreatePostPayload	true	"Post details"
//	@Success		201		{object}	json.envelope{data=Post}
//	@Failure		400		{object}	json.envelope	"Invalid request body or post details"
//	@Failure		401		{object}	json.envelope	"Missing, invalid, or expired access token"
//	@Failure		500		{object}	json.envelope	"Internal server error"
//	@Security		ApiKeyAuth
//	@Router			/posts/ [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var payload postsDTO.CreatePostPayload

	err := json.ReadJson(w, r, &payload)
	if err != nil {
		mainErrors.BadRequestResponse(w, r, err)
		return
	}

	if err := json.Validate.Struct(payload); err != nil {
		mainErrors.BadRequestResponse(w, r, err)
		return
	}

	post, err := h.service.Create(r.Context(), payload)
	if err != nil {
		mainErrors.InternalServerError(w, r, err)
		return
	}

	json.WriteJSON(w, http.StatusCreated, post)
}

// GetByID godoc
//
//	@Summary		Fetch a post
//	@Description	Returns a post by ID, including its comments.
//	@Tags			Posts
//	@Produce		json
//	@Param			postID	path		int	true	"Post ID"
//	@Success		200		{object}	json.envelope{data=Post}
//	@Failure		400		{object}	json.envelope	"Invalid post ID"
//	@Failure		401		{object}	json.envelope	"Missing, invalid, or expired access token"
//	@Failure		404		{object}	json.envelope	"Post not found"
//	@Failure		500		{object}	json.envelope	"Internal server error"
//	@Security		ApiKeyAuth
//	@Router			/posts/{postID}/ [get]
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	post, ok := PostFromContext(r.Context())
	if !ok {
		mainErrors.InternalServerError(
			w, r, errors.New("post missing from request context"),
		)
		return
	}

	comments, err := h.commentsService.GetByPostID(r.Context(), post.ID)
	if err != nil {
		mainErrors.InternalServerError(w, r, err)
		return
	}

	post.Comments = comments

	json.WriteJSON(w, http.StatusOK, post)
}

// Update godoc
//
//	@Summary		Update a post
//	@Description	Updates the supplied title and content fields. Omitted or null fields remain unchanged. Requires ownership of the post or a role with at least moderator privileges.
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Param			postID	path		int							true	"Post ID"
//	@Param			payload	body		postsDTO.UpdatePostPayload	true	"Post details"
//	@Success		200		{object}	json.envelope{data=Post}
//	@Failure		400		{object}	json.envelope	"Invalid post ID or request body"
//	@Failure		401		{object}	json.envelope	"Missing, invalid, or expired access token"
//	@Failure		403		{object}	json.envelope	"Insufficient permissions to modify this post"
//	@Failure		404		{object}	json.envelope	"Post not found"
//	@Failure		500		{object}	json.envelope	"Internal server error"
//	@Security		ApiKeyAuth
//	@Router			/posts/{postID}/ [patch]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	post, ok := PostFromContext(r.Context())
	if !ok {
		mainErrors.InternalServerError(
			w, r, errors.New("post missing from request context"),
		)
		return
	}

	var payload postsDTO.UpdatePostPayload
	if err := json.ReadJson(w, r, &payload); err != nil {
		mainErrors.BadRequestResponse(w, r, err)
		return
	}

	if err := json.Validate.Struct(payload); err != nil {
		mainErrors.BadRequestResponse(w, r, err)
		return
	}

	if payload.Title != nil {
		post.Title = *payload.Title
	}
	if payload.Content != nil {
		post.Content = *payload.Content
	}

	err := h.service.Update(r.Context(), post)
	if err != nil {
		if errors.Is(err, mainErrors.ErrNotFound) {
			mainErrors.NotFoundResponse(w, r, err)
			return
		}

		mainErrors.InternalServerError(w, r, err)
		return
	}

	json.WriteJSON(w, http.StatusOK, post)
}

// Delete godoc
//
//	@Summary		Delete a post
//	@Description	Deletes a post by ID. Requires ownership of the post or a role with at least administrator privileges.
//	@Tags			Posts
//	@Produce		json
//	@Param			postID	path	int	true	"Post ID"
//	@Success		204		"No Content"
//	@Failure		400		{object}	json.envelope	"Invalid post ID"
//	@Failure		401		{object}	json.envelope	"Missing, invalid, or expired access token"
//	@Failure		403		{object}	json.envelope	"Insufficient permissions to modify this post"
//	@Failure		404		{object}	json.envelope	"Post not found"
//	@Failure		500		{object}	json.envelope	"Internal server error"
//	@Security		ApiKeyAuth
//	@Router			/posts/{postID}/ [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	post, ok := PostFromContext(r.Context())
	if !ok {
		mainErrors.InternalServerError(
			w, r, errors.New("post missing from request context"),
		)
		return
	}

	err := h.service.Delete(r.Context(), post.ID)
	if err != nil {
		if errors.Is(err, mainErrors.ErrNotFound) {
			mainErrors.NotFoundResponse(w, r, err)
			return
		}

		mainErrors.InternalServerError(w, r, err)
		return
	}

	json.WriteJSON(w, http.StatusNoContent, nil)
}

// GetUserFeed godoc
//
//	@Summary		Get the user feed
//	@Description	Returns posts with author details and comment counts, filtered and ordered by creation time. The current implementation uses user ID 3 instead of the authenticated user's ID.
//	@Tags			Users
//	@Produce		json
//	@Param			limit	query		int			true	"Number of posts to return"													minimum(1)	maximum(20)
//	@Param			offset	query		int			true	"Number of posts to skip"													minimum(0)
//	@Param			sort	query		string		false	"Creation time sort order"													Enums(asc, desc)	default(desc)
//	@Param			tags	query		[]string	false	"Up to 5 comma-separated tags; posts must contain all supplied tags"		collectionFormat(csv)
//	@Param			search	query		string		false	"Case-insensitive text search in title or content (maximum 10 characters)"	maxLength(10)
//	@Param			from	query		string		false	"Inclusive creation time, formatted as YYYY-MM-DD HH:MM:SS"
//	@Param			to		query		string		false	"Exclusive creation time, formatted as YYYY-MM-DD HH:MM:SS; only parsed when from is supplied"
//	@Success		200		{object}	json.envelope{data=[]PostWithMetaData}
//	@Failure		400		{object}	json.envelope	"Missing or invalid pagination or filter parameters"
//	@Failure		401		{object}	json.envelope	"Missing, invalid, or expired access token"
//	@Failure		500		{object}	json.envelope	"Internal server error"
//	@Security		ApiKeyAuth
//	@Router			/users/feed [get]
func (h *Handler) GetUserFeed(w http.ResponseWriter, r *http.Request) {
	pagination := pagination.Pagination{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
	}

	pagination, err := pagination.Parse(r)
	if err != nil {
		mainErrors.BadRequestResponse(w, r, err)
		return
	}

	if err := json.Validate.Struct(pagination); err != nil {
		mainErrors.BadRequestResponse(w, r, err)
		return
	}

	userID := 3
	posts, err := h.service.GetUserFeed(r.Context(), int64(userID), pagination)
	if err != nil {
		mainErrors.InternalServerError(w, r, err)
		return
	}

	if err := json.WriteJSON(w, http.StatusOK, posts); err != nil {
		mainErrors.InternalServerError(w, r, err)
		return
	}

}
