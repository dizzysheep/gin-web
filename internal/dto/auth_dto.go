package dto

type AuthReqDTO struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthDTO struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}
