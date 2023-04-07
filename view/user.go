package view

type User struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type LoginReq struct {
	User User `json:"user"`
}
type LoginData User
