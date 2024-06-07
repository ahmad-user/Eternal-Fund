package dto

type AuthResponDto struct {
	Token string `json:"token"`
}

type AuthReqDto struct {
	Email     string `json:"email"`
	Passwords string `json:"passwords"`
}
