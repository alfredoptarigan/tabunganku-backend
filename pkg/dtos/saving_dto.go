package dtos

import "time"

type SavingRequest struct {
	Name           string  `form:"name" validate:"required,min=3,max=50"`
	TargetAmount   float64 `form:"target_amount" validate:"required,gt=0"`
	CurrencyCode   string  `form:"currency_code" validate:"required,len=3"`
	FillingPlan    string  `form:"filling_plan" validate:"required,oneof=daily weekly monthly"`
	FillingNominal float64 `form:"filling_nominal" validate:"required,gt=0"`
	Image          string  `form:"image" validate:"required"`
	UserUUID       string  `json:"user_uuid"`
}

type SavingResponse struct {
	UUID           string       `json:"uuid"`
	User           UserResponse `json:"user"`
	Name           string       `json:"name"`
	TargetAmount   float64      `json:"target_amount"`
	CurrencyCode   string       `json:"currency_code"`
	CurrencyFlag   string       `json:"currency_flag"`
	Image          string       `json:"image"`
	FillingPlan    string       `json:"filling_plan"`
	FillingNominal float64      `json:"filling_nominal"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}
