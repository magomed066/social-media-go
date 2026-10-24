package users

type User struct {
	ID        int64  `json:"id"`
	Firstname string `json:"first_name"`
	Lastname  string `json:"last_name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	CreatedAt string `json:"created_at"`
	IsActive  bool   `json:"is_active"`
}

type Author struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

type UserWithToken struct {
	*User
	Token string `json:"-"`
}
