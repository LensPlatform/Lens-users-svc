package models

import (
	"encoding/json"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
)

// User represents a single user profile
// ID should always be globally unique
type User struct {
	ID                uint32 `gorm:"primary_key"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	Type              string         `json:"user_type" validate:"required"`
	Firstname         string         `json:"first_name" validate:"required"`
	Lastname          string         `json:"last_name" validate:"required"`
	Username          string         `json:"user_name" validate:"required" gorm:"type:varchar(100);unique_index"`
	Gender            string         `json:"gender" validate:"-"`
	Languages         string         `json:"Languages" validate:"-"`
	Email             string         `json:"email" validate:"required,email"`
	PassWord          string         `json:"password" validate:"required,gte=8,lte=20"`
	PassWordConfirmed string         `json:"password_confirmed" validate:"required,gte=8,lte=20"`
	Age               int            `json:"age" validate:"gte=0,lte=120"`
	BirthDate         string         `json:"birth_date" validate:"required"`
	PhoneNumber       string         `json:"phone_number,omitempty" validate:"required"`
	Addresses         postgres.Jsonb `json:"location,omitempty" validate:"-"`
	Bio               string         `json:"bio,omitempty" validate:"required"`
	Education         Education      `json:"education,omitempty" validate:"-"`
	UserInterests     Interests      `json:"interests,omitempty" validate:"-"`
	Headline          string         `json:"headline,omitempty" validate:"max=30"`
	Subscriptions     Subscriptions  `json:"subscriptions,omitempty" validate:"-"`
	Intent            string         `json:"intent,omitempty" validate:"required"`
	Skills            Skillset       `json:"skillset,omitempty" validate:"-"`
	Teams             TeamsMemberOf  `json:"associated_teams,omitempty" validate:"-"`
	Groups            GroupsMemberOf `json: "associated_groups,omitempty" validate:"-"`
	SocialMedia       SocialMedia    `json:"social_media,omitempty" validate:"-"`
	Settings          Settings       `json:"settings,omitempty" validate:"-"`
}

type Address struct {
	City    string `json:"city" validate:"required"`
	State   string `json:"state" validate:"required"`
	Country string `json:"country" validate:"required"`
}

type TeamsMemberOf struct {
	Teams []TeamOf `json:"groups" validated:"-"`
}

type GroupsMemberOf struct {
	Groups []TeamOf `json:"groups" validated:"-"`
}

type TeamOf struct {
	Id   uint32   `json:"team_id" validate:"-"`
	Name string   `json:"team_name" validate:"-"`
	Type string   `json:"team_type" validate:"-"`
	Tags []string `json:"tags" validate:"-"`
}

type Education struct {
	JsonEmbeddable
	MostRecentInstitutionName string `json:"most_recent_institution_name" validate:"required"`
	HighestDegreeEarned       string `json:"highest_degree_earned" validate:"required"`
	Graduated                 bool   `json:"graduated" validate:"required"`
	Major                     string `json:"major" validate:"required"`
	Minor                     string `json:"minor" validate:"required"`
	YearsOfAttendance         string `json:"years_of_attendance" validate:"required"`
}

type Interests struct {
	JsonEmbeddable
	Industry []Industry `json:"industries_of_interest" validate:"omitempty"`
	Topic    []Topic    `json:"topics_of_interest" validate:"omitempty"`
}

type Topic struct {
	JsonEmbeddable
	TopicName string `json:"topic_name" validate:"required"`
	TopicType string `json:"topic_type" validate:"required"`
}

type Industry struct {
	JsonEmbeddable
	IndustryName string `json:"industry_name" validate:"required"`
}

type Subscriptions struct {
	JsonEmbeddable
	SubscriptionName string `json:"subscription_name" validate:"required"`
	Subscribe        bool   `json:"subscribe" validate:"required"`
}

type Skillset struct {
	JsonEmbeddable
	Skills []Skill `json:"skills" validate:"required"`
}

type Skill struct {
	JsonEmbeddable
	Type string `json:"skill_type" validate:"required"`
	Name string `json:"skill_name" validate:"required"`
}

func (user User) BeforeCreate(scope *gorm.Scope) error {
	id := uuid.New()
	err := scope.SetColumn("created_at", time.Now())

	if err != nil {
		return err
	}

	err = scope.SetColumn("updated_at", time.Now())

	if err != nil {
		return err
	}

	return scope.SetColumn("id", id)
}

func (user User) BeforeUpdate(scope *gorm.Scope) error {
	err := scope.SetColumn("updatedat", time.Now().String())

	if err != nil {
		return err
	}
	return nil
}

func (user User) ConvertToTableRow() (UserTable, error) {
	var userTable UserTable

	skills, err := json.Marshal(user.Skills)
	if err != nil {
		return UserTable{}, err
	}
	userTable.Skills = postgres.Jsonb{json.RawMessage(skills)}

	address, err := json.Marshal(user.Addresses)
	if err != nil {
		return UserTable{}, err
	}
	userTable.Addresses = postgres.Jsonb{json.RawMessage(address)}

	education, err := json.Marshal(user.Education)
	if err != nil {
		return UserTable{}, err
	}
	userTable.Education = postgres.Jsonb{json.RawMessage(education)}

	interest, err := json.Marshal(user.UserInterests)
	if err != nil {
		return UserTable{}, err
	}
	userTable.UserInterests = postgres.Jsonb{json.RawMessage(interest)}

	settings, err := json.Marshal(user.Settings)
	if err != nil {
		return UserTable{}, err
	}
	userTable.Settings = postgres.Jsonb{json.RawMessage(settings)}

	social, err := json.Marshal(user.SocialMedia)
	if err != nil {
		return UserTable{}, err
	}
	userTable.SocialMedia = postgres.Jsonb{json.RawMessage(social)}

	subscriptions, err := json.Marshal(user.Subscriptions)
	if err != nil {
		return UserTable{}, err
	}
	userTable.Subscriptions = postgres.Jsonb{json.RawMessage(subscriptions)}

	groups, err := json.Marshal(user.Groups)
	if err != nil {
		return UserTable{}, err
	}
	userTable.Groups = postgres.Jsonb{json.RawMessage(groups)}

	teams, err := json.Marshal(user.Teams)
	if err != nil {
		return UserTable{}, err
	}
	userTable.Teams = postgres.Jsonb{json.RawMessage(teams)}

	userTable.Type = user.Type
	userTable.ID = rand.Uint32()
	userTable.FirstName = user.Firstname
	userTable.LastName = user.Lastname
	userTable.UserName = user.Username
	userTable.Gender = user.Gender
	userTable.Languages = user.Languages
	userTable.Email = user.Email
	userTable.PassWord = user.PassWord
	userTable.PassWordConfirmed = user.PassWordConfirmed
	userTable.Age = user.Age
	userTable.BirthDate = user.BirthDate
	userTable.PhoneNumber = user.PhoneNumber
	userTable.Bio = user.Bio
	userTable.Headline = user.Headline
	userTable.Intent = user.Intent
	userTable.CreatedAt = user.CreatedAt
	userTable.UpdatedAt = user.UpdatedAt

	return userTable, nil
}
