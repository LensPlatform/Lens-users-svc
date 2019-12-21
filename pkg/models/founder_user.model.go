package models

type Founder struct {
	User
	Experience []Experience `json:"experience" validate:"required" sql:"experience"`
}
