package dto

type UserProfileUpdate struct {
	Username   string `json:"username" validate:"required"`
	Password   string `json:"password" validate:"required"`
	Email      string `json:"email" validate:"required"`
	Role       string `json:"role" validate:"required"`
	University string `json:"university" validate:"required"`
	Interests  string `json:"interests" validate:"required"`
}
