package models

import (
	"fmt"
	"reflect"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/type/date"
)

type TeamTable struct {
	ID                 uint `json:"id" validate:"-" gorm:"primary_key"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
	Name               string         `json:"name" validate:"required" gorm:"type:varchar(100)"` // team name
	Tags               postgres.Jsonb `json:"tags" validate:"required"`
	Email              string         `json:"email" validate:"required" gorm:"type:varchar(100)"`
	Type               string         `json:"type" validate:"required" gorm:"type:varchar(100)"`     // investor or startup team
	Overview           string         `json:"overview" validate:"required" gorm:"type:varchar(300)"` // about the team
	IndustryOfInterest string         `json:"industry" validate:"required" gorm:"type:varchar(100)"` // industry of interest
	FoundedDate        date.Date      `json:"founded_date "validate:"required""`
	Founders           postgres.Jsonb `json:"founder" validate:"required"`
	NumberOfEmployees  int            `json:"number_of_employees" validate:"required"` // size of team
	Headquarters       string         `json:"headquarters,omitempty" validate:"-" gorm:"type:varchar(100)"`
	Interests          string         `json:"interests,omitempty" validate:"-" gorm:"type:varchar(400)"`
	TeamMembers        postgres.Jsonb `json:"team_members,omitempty" validate:"-"`
	Advisors           postgres.Jsonb `json:"advisors,omitempty" validate:"-"`
	SocialMedia        postgres.Jsonb `json:"social_media,omitempty" validate:"-"`
	Contact            postgres.Jsonb `json:"contact,omitempty" validate:"-"`
}

func (table TeamTable) MigrateSchemaOrCreateTable(db *gorm.DB, logger *zap.Logger) {
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
		err := db.Table("teams_table").CreateTable(&table).Error
		if err != nil {
			logger.Error(fmt.Sprintf("Cannot Create %s Table", tableName))
			logger.Error(err.Error())
		}
		logger.Info(fmt.Sprintf("Sucessfully Created %s Table", tableName))
	}
}
