package users

type User struct {
	ID        int64  `json:"id"`
	Firstname string `json:"firstName"`
	Lastname  string `json:"lastName"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	CreatedAt string `json:"createdAt"`
	IsActive  bool   `json:"isActive"`
}

type Author struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

type UserWithToken struct {
	*User
	Token string `json:"-"`
}

type UserWithAuthToken struct {
	User        User   `json:"user"`
	AccessToken string `json:"accessToken"`
}
