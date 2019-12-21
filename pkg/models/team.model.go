package models

import (
	"time"

	"google.golang.org/genproto/googleapis/type/date"
)

type Team struct {
	JsonEmbeddable
	ID                uint32 `json:"id" validate:"-" gorm:"primary_key"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	Name              string       `json:"name" validate:"required" sql:"name"` // team name
	Tags              []string     `json:"tags" validate:"required"`
	Email             string       `json:"email" validate:"required"`
	Type              string       `json:"type" validate:"required" sql:"type"`                   // investor or startup team
	Overview          string       `json:"overview" validate:"required" sql:"overview"`           // about the team
	IndustryOfInterest string       `json:"industry" validate:"required" sql:"industryofinterest"` // industry of interest
	FoundedDate        date.Date    `json:"founded_date" validate:"required" sql:"foundeddate"`
	Founders           []TeamMember `json:"founder" validate:"required" sql:"founders"`
	NumberOfEmployees int          `json:"number_of_employees" validate:"required" sql:"numberofemployees"` // size of team
	Headquarters      string       `json:"headquarters,omitempty" validate:"-" sql:"headquarters"`
	Interests         string       `json:"interests,omitempty" validate:"-" sql:"interests"`
	TeamMembers       []TeamMember `json:"team_members,omitempty" validate:"-" sql:"teammembers"`
	Advisors          []TeamMember `json:"advisors,omitempty" validate:"-" sql:"advisors"`
	SocialMedia       SocialMedia  `json:"social_media,omitempty" validate:"-" sql:"socialmedia"`
	Contact           Contact      `json:"contact,omitempty" validate:"-" sql:"contact"`
}

type TeamMember struct {
	JsonEmbeddable
	Id    uint32 `json:"ID" validate:"required" sql:"id"`
	Name  string `json:"name" validate:"required" sql:"name"`
	Title string `json:"status" validate:"required" sql:"title"` // co-founder
}

type Funding struct {
	JsonEmbeddable
	TotalFunding       int      `json:"total_funding" validate:"required" sql:"totalfunding"`
	LatestRound        string   `json:"latest_round" validate:"required" sql:"latestround"`
	LatestRoundDate    string   `json:"latest_round_date" validate:"required" sql:"latestrounddate"`
	LatestRoundFunding int      `json:"latest_round_funding" validate:"required" sql:"latestroundfunding"`
	Investors          []string `json:"investors" validate:"-" sql:"investors"`
	InvestorType       string   `json:"investor_type" validate:"-" sql:"investortype"`       //Accelerator
	InvestmentStage    string   `json:"investment_stage" validate:"-" sql:"investmentstage"` // Debt, Early Stage Venture, Seed
}

type Details struct {
	JsonEmbeddable
	IPOStatus   string `json:"ipo_status" validate:"-" sql:"ipostatus"`
	CompanyType string `json:"company_type" validate:"-" sql:"companytype"` // for-profit
}
