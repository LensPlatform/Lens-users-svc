package models

type Contact struct {
	JsonEmbeddable
	Email       string `json:"email" validate:"required" sql:"email"`
	PhoneNumber string `json:"phonenumber" validated:"required" sql:"phonenumber"`
}

type SocialMedia struct {
	JsonEmbeddable
	Website   string `json:"website" validate:"-" sql:"website"`
	Facebook  string `json:"facebook" validate:"-" sql:"facebook"`
	Twitter   string `json:"twitter" validate:"-" sql:"twitter"`
	LinkedIn  string `json:"linkedIn" validate:"-" sql:"linkedin"`
	Instagram string `json:"instagram" validate:"-" sql:"instagram"`
	Youtube   string `json:"youtube" validate:"-" sql:"youtube"`
}

type Experience struct {
	JsonEmbeddable
	CompanyName string `json:"company_name" validate:"required" sql:"company"`
	StartDate   string `json:"start_date" validate:"required" sql:"start_date"`
	EndDate     string `json:"end_date" validate:"required" sql:"end_date"`
	Title       string `json:"title" validate:"reqiuired" sql:"title"`
}

type Investment struct {
	JsonEmbeddable
	CompanyName string `json:"company_name" validate:"required" sql:"company"`
	Industry    string `json:"industry" validate:"required" sql:"industry"`
}
