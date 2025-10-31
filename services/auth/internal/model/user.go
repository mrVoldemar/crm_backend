package model

type UserInfo struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}

type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
