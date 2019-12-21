package models

import (
	"fmt"
	"reflect"

	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	"go.uber.org/zap"
)

type GroupTable struct {
	ID              uint32         `json:"id" validate:"-" gorm:"primary_key"`
	Name            string         `json:"group_name" validate:"required" gorm:"varchar(100)"`
	Type            string         `json:"type" validate:"required" gorm:"varchar(100)"` // tech, social, ... etc
	Owner           int            `json:"group_owner" validate:"required"`
	Bio             string         `json:"group_bio" validate:"required" gorm:"varchar(300)"`
	Tags            postgres.Jsonb `json:"tags" validate:"required"`
	NumGroupMembers int            `json:"num_group_members" validate:"required"`
	GroupMembers    postgres.Jsonb `json:"group_member_ids" valudate:"-"`
}

func (table GroupTable) MigrateSchemaOrCreateTable(db *gorm.DB, logger *zap.Logger) {
	t := reflect.TypeOf(table)
	tableName := t.Name()

	if db.HasTable(&table) {
		err := db.AutoMigrate(&table).Error
		if err != nil {
			logger.Error(fmt.Sprintf("Cannot Migrate: %s Schema", tableName))
			logger.Error(err.Error())
		} else {
			logger.Info(fmt.Sprintf("Successfully Migrated %s Schema", tableName))
		}
	} else {
		err := db.Table("groups_table").CreateTable(&table).Error
		if err != nil {
			logger.Error(fmt.Sprintf("Cannot Create %s Table", tableName))
			logger.Error(err.Error())
		} else {
			logger.Info(fmt.Sprintf("Sucessfully Created %s Table", tableName))
		}
	}
}
