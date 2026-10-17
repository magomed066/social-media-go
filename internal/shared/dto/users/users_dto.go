package usersDTO

type FollowerUserPayload struct {
	UserID int64 `json:"user_id"`
}

type FollowPayload struct {
	// User that is being followed
	UserID     int64
	FollowerID int64
}
