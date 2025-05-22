package entity

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Admin       Admin  `json:"admin"`
	AccessToken string `json:"access_token"`
}
