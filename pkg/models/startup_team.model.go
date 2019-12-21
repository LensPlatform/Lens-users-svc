package models

type StartupTeam struct {
	Team
	Funding        Funding `json:"funding" validate:"-" sql:"funding"`
	CompanyDetails Details `json:"company_details" validate:"-" sql:"companydetails"`
}
