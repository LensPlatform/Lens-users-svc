package baseuser

import (
	"encoding/json"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/LensPlatform/Lens-users-svc/pkg/tables"
)

func (m *User) BeforeCreate(scope *gorm.Scope) error {
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

func (m *User) BeforeUpdate(scope *gorm.Scope) error {
	err := scope.SetColumn("updatedat", time.Now().String())

	if err != nil {
		return err
	}
	return nil
}

func (m *User) ConvertToTableRow() (tables.UserTable, error) {
	var userTable tables.UserTable

	skills, err := json.Marshal(m.Skills)
	if err != nil {
		return tables.UserTable{}, err
	}
	userTable.Skills = postgres.Jsonb{json.RawMessage(skills)}

	address, err := json.Marshal(m.UserAddress)
	if err != nil {
		return tables.UserTable{}, err
	}
	userTable.Addresses = postgres.Jsonb{json.RawMessage(address)}

	education, err := json.Marshal(m.UserEducation)
	if err != nil {
		return tables.UserTable{}, err
	}
	userTable.Education = postgres.Jsonb{json.RawMessage(education)}

	interest, err := json.Marshal(m.UserInterests)
	if err != nil {
		return tables.UserTable{}, err
	}
	userTable.UserInterests = postgres.Jsonb{json.RawMessage(interest)}

	settings, err := json.Marshal(m.Settings)
	if err != nil {
		return tables.UserTable{}, err
	}
	userTable.Settings = postgres.Jsonb{json.RawMessage(settings)}

	social, err := json.Marshal(m.SocialMedia)
	if err != nil {
		return tables.UserTable{}, err
	}
	userTable.SocialMedia = postgres.Jsonb{json.RawMessage(social)}

	subscriptions, err := json.Marshal(m.UserSubscriptions)
	if err != nil {
		return tables.UserTable{}, err
	}
	userTable.Subscriptions = postgres.Jsonb{json.RawMessage(subscriptions)}

	groups, err := json.Marshal(m.Groups)
	if err != nil {
		return tables.UserTable{}, err
	}
	userTable.Groups = postgres.Jsonb{json.RawMessage(groups)}

	teams, err := json.Marshal(m.Teams)
	if err != nil {
		return tables.UserTable{}, err
	}
	userTable.Teams = postgres.Jsonb{json.RawMessage(teams)}

	userTable.Type = m.Type
	userTable.ID = rand.Uint32()
	userTable.FirstName = m.FirstName
	userTable.LastName = m.LastName
	userTable.UserName = m.UserName
	userTable.Gender = m.Gender
	userTable.Languages = m.Languages
	userTable.Email = m.Email
	userTable.PassWord = m.PassWordConfirmed
	userTable.PassWordConfirmed = m.PassWordConfirmed
	userTable.Age = m.Age
	userTable.BirthDate = m.BirthDate
	userTable.PhoneNumber = m.PhoneNumber
	userTable.Bio = m.Bio
	userTable.Headline = m.Headline
	userTable.Intent = m.Intent
	userTable.CreatedAt = *m.CreatedAt
	userTable.UpdatedAt = *m.UpdatedAt

	return userTable, nil
}
