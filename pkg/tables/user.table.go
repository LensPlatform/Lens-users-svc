package tables

import (
	"fmt"
	"reflect"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	"go.uber.org/zap"
)

type UserTable struct {
	ID                uuid.UUID       `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	CreatedAt         timestamp.Timestamp
	UpdatedAt         timestamp.Timestamp
	Type              string         `json:"user_type" validate:"required" gorm:"type:varchar(100)"`
	FirstName         string         `json:"first_name" validate:"required" gorm:"type:varchar(100);column:firstname"`
	LastName          string         `json:"last_name" validate:"required" gorm:"type:varchar(100);column:lastname"`
	UserName          string         `json:"user_name" validate:"required" gorm:"type:varchar(100);unique_index;column:username"`
	Gender            string         `json:"gender" validate:"-" gorm:"type:varchar(100)"`
	Languages         string         `json:"Languages" validate:"-" gorm:"type:text"`
	Email             string         `json:"email" validate:"required,email" gorm:"type:varchar(100)"`
	PassWord          string         `json:"password" validate:"required,gte=8,lte=20" gorm:"type:varchar(200)"`
	PassWordConfirmed string         `json:"password_confirmed" validate:"required,gte=8,lte=20" gorm:"type:varchar(200)"`
	Age               int32           `json:"age" validate:"gte=0,lte=120"`
	BirthDate         string         `json:"birth_date" validate:"required"`
	PhoneNumber       string         `json:"phone_number,omitempty" validate:"required"`
	Bio               string         `json:"bio,omitempty" validate:"required"`
	Headline          string         `json:"headline,omitempty" validate:"max=30"`
	Intent            string         `json:"intent,omitempty" validate:"required"`
	Addresses         postgres.Jsonb `json:"location,omitempty" validate:"-" gorm:"type:jsonb;not null;default '{}'"`
	Education         postgres.Jsonb `json:"education,omitempty" validate:"-" gorm:"type:jsonb;not null;default '{}'"`
	UserInterests     postgres.Jsonb `json:"interests,omitempty" validate:"-" gorm:"type:jsonb;not null;default '{}'"`
	Subscriptions     postgres.Jsonb `json:"subscriptions,omitempty" validate:"-" gorm:"type:jsonb;not null;default '{}'"`
	Skills            postgres.Jsonb `json:"skillset,omitempty" validate:"-" gorm:"type:jsonb;not null;default '{}'"`
	Teams             postgres.Jsonb `json:"associated_teams,omitempty" validate:"-" gorm:"type:jsonb;not null;default '{}'"`
	Groups            postgres.Jsonb `json:"associated_groups,omitempty" validate:"-" gorm:"type:jsonb;not null;default '{}'"`
	SocialMedia       postgres.Jsonb `json:"social_media,omitempty" validate:"-" gorm:"type:jsonb;not null;default '{}'"`
	Settings          postgres.Jsonb `json:"settings,omitempty" validate:"-" gorm:"type:jsonb;not null;default '{}'"`
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
