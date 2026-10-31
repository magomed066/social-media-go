package usersDTO

type FollowerUserPayload struct {
	// ID of the follower, supplied explicitly rather than taken from the access token.
	UserID int64 `json:"user_id" example:"2"`
} //	@name	FollowerUserPayload

type FollowPayload struct {
	// User that is being followed
	UserID     int64
	FollowerID int64
}

type RegisterUserPayload struct {
	Username  string `json:"username" validate:"required,max=50" maxLength:"50" example:"jane_doe"`
	Firstname string `json:"first_name" validate:"required,max=100" maxLength:"100" example:"Jane"`
	Lastname  string `json:"last_name" validate:"required,max=255" maxLength:"255" example:"Doe"`
	Email     string `json:"email" validate:"email,required" format:"email" example:"jane@example.com"`
	Password  string `json:"password" validate:"required,min=3,max=72" minLength:"3" maxLength:"72" example:"example-password"`
} //	@name	RegisterUserPayload

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required" example:"jane@example.com"`
	Password string `json:"password" validate:"required" example:"example-password"`
} //	@name	LoginUserPayload
