package dtos

type RegisterRequest struct {
	Name                 string `form:"name" json:"name" validate:"required"`
	Email                string `form:"email" json:"email" validate:"required,email"`
	Password             string `form:"password" json:"password" validate:"required,min=6"`
	ConfirmationPassword string `form:"confirmation_password" json:"confirmation_password" validate:"required,eqfield=Password"`
	PhoneNumber          string `form:"phone_number" json:"phone_number" validate:"required"`
	Image                string `form:"photo_url" json:"photo_url" validate:"omitempty"`
}

type UserResponse struct {
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Image       string `json:"image"`
}
