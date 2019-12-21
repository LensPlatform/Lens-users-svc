package models

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	"go.uber.org/zap"
)

type UserTable struct {
	ID                uint32 `gorm:"primary_key"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	Type              string         `json:"user_type" validate:"required"`
	FirstName         string         `json:"first_name" validate:"required" gorm:"type:varchar(100);column:firstname"`
	LastName          string         `json:"last_name" validate:"required" gorm:"type:varchar(100);column:lastname"`
	UserName          string         `json:"user_name" validate:"required" gorm:"type:varchar(100);unique_index;column:username"`
	Gender            string         `json:"gender" validate:"-"`
	Languages         string         `json:"Languages" validate:"-"`
	Email             string         `json:"email" validate:"required,email"`
	PassWord          string         `json:"password" validate:"required,gte=8,lte=20"`
	PassWordConfirmed string         `json:"password_confirmed" validate:"required,gte=8,lte=20"`
	Age               int            `json:"age" validate:"gte=0,lte=120"`
	BirthDate         string         `json:"birth_date" validate:"required"`
	PhoneNumber       string         `json:"phone_number,omitempty" validate:"required"`
	Bio               string         `json:"bio,omitempty" validate:"required"`
	Headline          string         `json:"headline,omitempty" validate:"max=30"`
	Intent            string         `json:"intent,omitempty" validate:"required"`
	Addresses         postgres.Jsonb `json:"location,omitempty" validate:"-"`
	Education         postgres.Jsonb `json:"education,omitempty" validate:"-"`
	UserInterests     postgres.Jsonb `json:"interests,omitempty" validate:"-"`
	Subscriptions     postgres.Jsonb `json:"subscriptions,omitempty" validate:"-"`
	Skills            postgres.Jsonb `json:"skillset,omitempty" validate:"-"`
	Teams             postgres.Jsonb `json:"associated_teams,omitempty" validate:"-"`
	Groups            postgres.Jsonb `json: "associated_groups,omitempty" validate:"-"`
	SocialMedia       postgres.Jsonb `json:"social_media,omitempty" validate:"-"`
	Settings          postgres.Jsonb `json:"settings,omitempty" validate:"-"`
}

func (table UserTable) MigrateSchemaOrCreateTable(db *gorm.DB, logger *zap.Logger) {
	t := reflect.TypeOf(table)
	tableName := t.Name()

	if db.HasTable(&table) {
		err := db.AutoMigrate(&table).Error
		if err != nil {
			logger.Error(fmt.Sprintf("Cannot Migrate: %s Schema", tableName))
			logger.Error(err.Error())
		}
		logger.Info(fmt.Sprintf("Successfully Migrated %s Schema", tableName))
	} else {
		err := db.Table("users_table").CreateTable(&table).Error
		if err != nil {
			logger.Error(fmt.Sprintf("Cannot Create %s Table", tableName))
			logger.Error(err.Error())
		}
		logger.Info(fmt.Sprintf("Sucessfully Created %s Table", tableName))
	}
}

func (table UserTable) ConvertFromRowToUser() (User, error) {
	var user User
	err := json.Unmarshal(table.Education.RawMessage, &user.Education)
	if err != nil {
		return User{}, err
	}
	err = json.Unmarshal(table.UserInterests.RawMessage, &user.UserInterests)
	if err != nil {
		return User{}, err
	}
	err = json.Unmarshal(table.SocialMedia.RawMessage, &user.SocialMedia)
	if err != nil {
		return User{}, err
	}
	err = json.Unmarshal(table.Subscriptions.RawMessage, &user.Subscriptions)
	if err != nil {
		return User{}, err
	}
	err = json.Unmarshal(table.Groups.RawMessage, &user.Groups)
	if err != nil {
		return User{}, err
	}
	err = json.Unmarshal(table.Teams.RawMessage, &user.Teams)
	if err != nil {
		return User{}, err
	}
	err = json.Unmarshal(table.Skills.RawMessage, &user.Skills)
	if err != nil {
		return User{}, err
	}
	err = json.Unmarshal(table.Addresses.RawMessage, &user.Addresses)
	if err != nil {
		return User{}, err
	}
	err = json.Unmarshal(table.Settings.RawMessage, &user.Settings)
	if err != nil {
		return User{}, err
	}

	user.Type = table.Type
	user.ID = table.ID
	user.Firstname = table.FirstName
	user.Lastname = table.LastName
	user.Username = table.UserName
	user.Gender = table.Gender
	user.Languages = table.Languages
	user.Email = table.Email
	user.PassWord = table.PassWord
	user.PassWordConfirmed = table.PassWordConfirmed
	user.Age = table.Age
	user.BirthDate = table.BirthDate
	user.PhoneNumber = table.PhoneNumber
	user.Bio = table.Bio
	user.Headline = table.Headline
	user.Intent = table.Intent
	user.CreatedAt = table.CreatedAt
	user.UpdatedAt = table.UpdatedAt
	return user, nil
}
