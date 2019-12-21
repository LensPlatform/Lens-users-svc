package models

type InvestorTeam struct {
	Team
	InvestorType        string `json:"investor_type" validate:"-" sql:"investortype"`
	InvestmentStage     string `json:"investment_stage" validate:"-" sql:"investmentstage"`
	NumberOfExits       int    `json:"number_of_exits" validate:"-" sql:"numberofexits"`
	NumberOfinvestments int    `json:"number_of_investments" validate:"-" sql:"numberofinvestments"`
	NumberOfFunds       int    `json:"number_of_funds" validate:"-" sql:"numberoffunds"`
}
