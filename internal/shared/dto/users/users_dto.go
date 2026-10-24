package usersDTO

type FollowerUserPayload struct {
	UserID int64 `json:"user_id"`
}

type FollowPayload struct {
	// User that is being followed
	UserID     int64
	FollowerID int64
}

type RegisterUserPayload struct {
	Username  string `json:"username" validate:"required,max=50"`
	Firstname string `json:"first_name" validate:"required,max=100"`
	Lastname  string `json:"last_name" validate:"required,max=255"`
	Email     string `json:"email" validate:"email,required"`
	Password  string `json:"password" validate:"required,min=3,max=72"`
}
