package models

type Investor struct {
	Founder
	Investments []Investment `json:"investments" validate:"required" sql:"investments"`
}
