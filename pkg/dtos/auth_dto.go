package dtos

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=100"`
}

type LoginResponse struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Email        string `json:"email"`
	UserUuid     string `json:"user_uuid"`
	Name         string `json:"name"`
}
