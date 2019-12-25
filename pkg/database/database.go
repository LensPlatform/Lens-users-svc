package database

import (
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"

	"github.com/LensPlatform/Lens-users-svc/pkg/helper"
	"github.com/LensPlatform/Lens-users-svc/pkg/model/proto/group"
	baseteam "github.com/LensPlatform/Lens-users-svc/pkg/model/proto/team"
	baseuser "github.com/LensPlatform/Lens-users-svc/pkg/model/proto/user"
	"github.com/LensPlatform/Lens-users-svc/pkg/tables"
)

var (
	GetUserByIdQuery       = "SELECT * FROM users_table WHERE id =$1"
	GetUserByUsernameQuery = "SELECT * FROM users_table WHERE username =$1"
	GetUserByEmailQuery    = "SELECT * FROM users_table WHERE email =$1"
)


type DBHandler interface {
	// Create
	CreateUser(user baseuser.User) error
	CreateTeam(founder baseuser.User, team baseteam.Team) error
	CreateGroup(owner baseuser.User, group group.Group) error

	// GET
	GetUserById(id string) (model.Team, error)
	GetUserByEmail(email string) (model.Team, error)
	GetUserByUsername(username string) (model.Team, error)
	GetPassword(id string) (string, error)
	GetAllUsers() ([]model.User, error)
	GetAllUsersFromSearchQuery(search map[string]interface{}) ([]model.User, error)
	GetUserBasedOnParam(param string, query string) (model.User, error)

	GetTeamByID(id string) (model.Team, error)
	GetTeamByName(name string) (model.Team, error)
	GetTeamsByType(teamType string) ([]model.Team, error)
	GetTeamsByIndustry(industry string) ([]model.Team, error)
	GetAllTeams() ([]model.Team, error)
	GetAllTeamsFromSearchQuery(search map[string]interface{}) ([]model.Team, error)
	GetTeamBasedOnParam(param string, query string) (model.Team, error)

	GetGroupBasedOnParam(param string, query string) (model.Group, error)
	GetGroupById(id string) (model.Group, error)
	GetGroupByName(name string) (model.Group, error)

	// Update
	UpdateUser(param map[string]interface{}, id string) (model.Team, error)

	UpdateTeamName(name string, teamId string) (model.Team, error)
	UpdateTeamType(teamType string, teamId string) (model.Team, error)
	UpdateTeamOverview(overView string, teamId string) (model.Team, error)
	updateTeamIndustry(industry string, teamId string) (model.Team, error)
	AddTeamMemberToTeam(teamMember model.TeamMember, teamId string) (model.Team, error)
	RemoveTeamMemberFromTeam(teamMember model.TeamMember, teamId string) (model.Team, error)
	AddAdvisorToTeam(advisorMember model.TeamMember, teamId string) (model.Team, error)
	RemoveAdvisorFromTeam(advisorMember model.TeamMember, teamId string) (model.Team, error)
	AddFounderToTeam(advisorMember model.TeamMember, teamId string) (model.Team, error)
	RemoveFounderFromTeam(advisorMember model.TeamMember, teamId string) (model.Team, error)

	AddMemberToGroup(memberId string, groupId string) (model.Team, error)
	RemoveMemberFromGroup(memberId string, groupId string) (model.Team, error)

	// Delete
	DeleteUserById(id string) (bool, error)
	DeleteUserByUsername(id string) (bool, error)
	DeleteUserByEmail(id string) (bool, error)
	DeleteUserBasedOnParam(param string, query string) (bool, error)

	DeleteTeamById(teamId string) (bool, error)
	DeleteTeamByName(teamName string) (bool, error)
	DeleteTeamTeamByEmail(teamEmail string) (bool, error)
	DeleteTeamBasedOnParam(param string, query string) (bool, error)

	DeleteGroupById(teamId string) (bool, error)

	// Existence
	DoesUserExist(searchParam string, query string) (bool, error)
	DoesTeamExist(searchParam string, query string) (bool, error)
}

type Database struct {
	connection *gorm.DB
}

func NewDatabase(db *gorm.DB) *Database {
	return &Database{connection: db}
}

func (db Database) CreateUser(user model.User) error {
	if user.Username == "" {
		errMsg := fmt.Sprintf("Invalid Argument provided. "+
			"the following param is null User Id : %s", user.Username)
		return errors.New(errMsg)
	}

	table, e := user.ConvertToTableRow()
	if e != nil {
		return e
	}

	e = db.connection.Table("users_table").Create(&table).Error

	if e != nil {
		return e
	}
	return nil
}

func (db Database) CreateTeam(founder model.User, team model.Team) error {
	if founder.ID == 0 || team.ID == 0 {
		errMsg := fmt.Sprintf("Invalid Argument provided. One of "+
			"the following params are null Founder Id : %d, Team ID : %d", founder.ID, team.ID)
		return errors.New(errMsg)
	}

	var teamMember model.TeamMember
	var founders []model.TeamMember

	founderName := fmt.Sprintf("%s %s", founder.Firstname, founder.Lastname)
	teamMember = model.TeamMember{Id: founder.ID, Name: founderName, Title: "founder"}

	founders = append(team.Founders, teamMember)
	team.Founders = founders
	team.NumberOfEmployees = len(team.TeamMembers) + len(team.Advisors) + len(team.Founders)

	e := db.connection.Create(&team).Error

	if e != nil {
		return e
	}
	return nil
}

func (db Database) CreateGroup(owner model.User, group model.Group) error {
	if owner.ID == 0 || group.ID == 0 {
		errMsg := fmt.Sprintf("Invalid Argument provided. "+
			"one of the following arguments are null User Id : %d, Group Id : %d", owner.ID, group.ID)
		return errors.New(errMsg)
	}

	group.Owner = owner.ID
	group.NumGroupMembers = 1

	e := db.connection.Create(&group).Error
	if e != nil {
		return e
	}
	return nil
}

func (db Database) GetUserById(id string) (model.User, error) {
	return db.GetUserBasedOnParam(id, GetUserByIdQuery)
}

func (db Database) GetUserByUsername(username string) (model.User, error) {
	return db.GetUserBasedOnParam(username, GetUserByUsernameQuery)
}

func (db Database) GetUserByEmail(email string) (model.User, error) {
	return db.GetUserBasedOnParam(email, GetUserByEmailQuery)
}

func (db Database) GetAllUsers() ([]model.User, error) {
	var users []model.User
	e := db.connection.Find(&users).Error
	if e != nil {
		return nil, e
	}
	return users, nil
}
func (db Database) GetAllUsersFromSearchQuery(query map[string]interface{}) ([]model.User, error) {
	var users []model.User
	// ex. db.Where(map[string]interface{}{"name": "jinzhu", "age": 20}).Find(&users)
	e := db.connection.Where(query).Find(&users).Error
	if e != nil {
		return nil, e
	}
	return users, nil
}

func (db Database) GetUserBasedOnParam(param string, query string) (model.User, error) {
	if param == "" || query == "" {
		errMsg := fmt.Sprintf("Invalid Argument provided. One of "+
			"the following params are null Search Param : %s, Query : %s", param, query)
		return model.User{}, errors.New(errMsg)
	}

	var table tables.UserTable
	rows, e := db.connection.Table("users_table").Raw(query, param).Rows()
	if e != nil {
		return model.User{}, e
	}

	defer rows.Close()
	for rows.Next() {
		_ = db.connection.ScanRows(rows, &table)
	}

	user, e := table.ConvertFromRowToUser()
	if e != nil {
		return model.User{}, e
	}

	return user, e
}

func (db Database) GetTeamByID(id string) (model.Team, error) {
	if id == "" {
		errMsg := fmt.Sprintf("Invalid Argument provided. "+
			"The following argument is null Id : %s", id)
		return model.Team{}, errors.New(errMsg)
	}
	query := "id = ?"
	return db.GetTeamBasedOnParam(id, query)
}

func (db Database) GetTeamByName(name string) (model.Team, error) {
	if name == "" {
		errMsg := fmt.Sprintf("Invalid Argument provided. "+
			"The following argument is null name : %s", name)
		return model.Team{}, errors.New(errMsg)
	}
	query := "name = ?"
	return db.GetTeamBasedOnParam(name, query)
}

func (db Database) GetTeamsByType(teamType string) (model.Team, error) {
	if teamType == "" {
		errMsg := fmt.Sprintf("Invalid Argument provided. "+
			"The following argument is null Team Type : %s", teamType)
		return model.Team{}, errors.New(errMsg)
	}
	query := "type = ?"
	return db.GetTeamBasedOnParam(teamType, query)
}

func (db Database) GetTeamsByIndustry(industry string) ([]model.Team, error) {
	if industry == "" {
		errMsg := fmt.Sprintf("Invalid Argument provided. "+
			"The following argument is null Industry Type : %s", industry)
		return nil, errors.New(errMsg)
	}
	query := "industry = ?"
	var teams []model.Team
	e := db.connection.Where(query, industry).Find(&teams).Error

	if e != nil {
		return nil, e
	}
	return teams, nil
}

func (db Database) GetAllTeams() ([]model.Team, error) {
	var teams []model.Team
	e := db.connection.Find(&teams).Error

	if e != nil {
		return nil, e
	}
	return teams, nil
}

func (db Database) GetAllTeamsFromSearchQuery(search map[string]interface{}) ([]model.Team, error) {
	var teams []model.Team
	e := db.connection.Where(search).Find(&teams).Error

	if e != nil {
		return nil, e
	}
	return teams, nil
}

func (db Database) GetTeamBasedOnParam(param string, query string) (model.Team, error) {
	if param == "" || query == "" {
		errMsg := fmt.Sprintf("Invalid Argument provided. One of "+
			"the following params are null Search Param : %s, Query : %s", param, query)
		return model.Team{}, errors.New(errMsg)
	}

	var team model.Team
	e := db.connection.First(&team, query, param).Error

	if e != nil {
		return model.Team{}, e
	}
	return team, nil
}

func (db Database) GetGroupById(id string) (model.Group, error) {
	return db.GetGroupBasedOnParam(id, "id=?")
}

func (db Database) GetGroupByName(name string) (model.Group, error) {
	return db.GetGroupBasedOnParam(name, "name=?")
}

func (db Database) GetGroupBasedOnParam(param string, query string) (model.Group, error) {
	if param == "" || query == "" {
		errMsg := fmt.Sprintf("Invalid Argument provided. One of "+
			"the following params are null Search Param : %s, Query : %s", param, query)
		return model.Group{}, errors.New(errMsg)
	}

	var group model.Group
	e := db.connection.First(&group, query, param).Error

	if e != nil {
		return model.Group{}, e
	}
	return group, nil
}

func (db Database) DoesUserExist(searchParam string, query string) (bool, error) {
	// check if user exists
	var user model.User
	user, err := db.GetUserBasedOnParam(searchParam, query)

	if err != nil {
		return false, err
	}

	if user.ID != 0 {
		return true, nil
	}

	return false, nil
}

func (db Database) DoesTeamExist(searchParam string, query string) (bool, error) {
	var team model.Team
	team, err := db.GetTeamBasedOnParam(searchParam, query)

	if err != nil {
		return false, err
	}

	if team.ID != 0 {
		return true, nil
	}

	return false, nil
}

func (db Database) UpdateUser(param map[string]interface{}, id string) (model.User, error) {
	var user model.User
	user, err := db.GetUserById(id)
	if err != nil {
		return model.User{}, err
	}
	// db.Model(&user).Updates(map[string]interface{}{"name": "hello", "age": 18, "actived": false})
	err = db.connection.Model(&user).Updates(param).Error
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (db Database) UpdateTeamName(name string, teamId string) (model.Team, error) {
	var team model.Team
	team, err := db.GetTeamByID(teamId)

	if err != nil {
		return model.Team{}, err
	}

	err = db.connection.Model(&team).UpdateColumn("name", name).Error
	if err != nil {
		return model.Team{}, err
	}
	return team, nil
}

func (db Database) UpdateTeamType(teamType string, teamId string) (model.Team, error) {
	var team model.Team
	team, err := db.GetTeamByID(teamId)

	if err != nil {
		return model.Team{}, err
	}

	err = db.connection.Model(&team).UpdateColumn("type", teamType).Error

	if err != nil {
		return model.Team{}, err
	}
	return team, nil
}

func (db Database) UpdateTeamOverview(overView string, teamId string) (model.Team, error) {
	var team model.Team
	team, err := db.GetTeamByID(teamId)

	if err != nil {
		return model.Team{}, err
	}

	err = db.connection.Model(&team).UpdateColumn("overview", overView).Error
	if err != nil {
		return model.Team{}, err
	}
	return team, nil
}

func (db Database) UpdateTeamIndustry(industry string, teamId string) (model.Team, error) {
	var team model.Team
	team, err := db.GetTeamByID(teamId)

	if err != nil {
		return model.Team{}, err
	}

	err = db.connection.Model(&team).UpdateColumn("industry", industry).Error
	if err != nil {
		return model.Team{}, err
	}
	return team, nil
}

func (db Database) AddTeamMemberToTeam(teamMember model.TeamMember, teamId string) (model.Team, error) {
	if teamMember.Id == 0 || teamId == "" {
		errMsg := fmt.Sprintf("Invalid Argument provided. One of the following params are null Team Id : %d, Team Member"+
			" ID : %s", teamMember.Id, teamId)
		return model.Team{}, errors.New(errMsg)
	}

	var team model.Team
	var teamMembers []model.TeamMember
	team, err := db.GetTeamByID(teamId)

	if err != nil {
		return model.Team{}, err
	}

	teamMembers = append(team.TeamMembers, teamMember)
	team.TeamMembers = teamMembers
	team.NumberOfEmployees = len(team.TeamMembers) + len(team.Advisors) + len(team.Founders)

	err = db.connection.Save(&team).Error

	if err != nil {
		return model.Team{}, err
	}
	return team, nil
}

func (db Database) RemoveTeamMemberFromTeam(teamMember model.TeamMember, teamId string) (model.Team, error) {
	if teamMember.Id == 0 || teamId == "" {
		errMsg := fmt.Sprintf("Invalid Argument provided. One of the following params are null Team Id : %d, Team Member"+
			" ID : %s", teamMember.Id, teamId)
		return model.Team{}, errors.New(errMsg)
	}

	var team model.Team
	team, err := db.GetTeamByID(teamId)

	if err != nil {
		return model.Team{}, err
	}

	var index int

	for i, item := range team.TeamMembers {
		if item.Id == teamMember.Id {
			index = i
			break
		}
	}

	team.TeamMembers[index] = team.TeamMembers[len(team.TeamMembers)-1]
	team.TeamMembers = team.TeamMembers[:len(team.TeamMembers)-1]

	team.NumberOfEmployees = len(team.TeamMembers) + len(team.Advisors) + len(team.Founders)
	err = db.connection.Save(&team).Error

	if err != nil {
		return model.Team{}, err
	}
	return team, nil
}

func (db Database) AddAdvisorToTeam(advisorMember model.TeamMember, teamId string) (model.Team, error) {
	if advisorMember.Id == 0 || teamId == "" {
		errMsg := fmt.Sprintf("Invalid Argument provided. One of the following params are null Team Id : %d, Team Member"+
			" ID : %s", advisorMember.Id, teamId)
		return model.Team{}, errors.New(errMsg)
	}

	var team model.Team
	var advisorMembers []model.TeamMember
	team, err := db.GetTeamByID(teamId)

	if err != nil {
		return model.Team{}, err
	}

	advisorMembers = append(team.Advisors, advisorMember)
	team.Advisors = advisorMembers
	team.NumberOfEmployees = len(team.TeamMembers) + len(team.Advisors) + len(team.Founders)

	err = db.connection.Save(&team).Error

	if err != nil {
		return model.Team{}, err
	}
	return team, nil
}

func (db Database) RemoveAdvisorFromTeam(advisorMember model.TeamMember, teamId string) (model.Team, error) {
	if advisorMember.Id == 0 || teamId == "" {
		errMsg := fmt.Sprintf("Invalid Argument provided. One of the following params are null Team Id : %d, Team Member"+
			" ID : %s", advisorMember.Id, teamId)
		return model.Team{}, errors.New(errMsg)
	}

	var team model.Team
	team, err := db.GetTeamByID(teamId)

	if err != nil {
		return model.Team{}, err
	}

	var index int

	for i, item := range team.Advisors {
		if item.Id == advisorMember.Id {
			index = i
			break
		}
	}

	team.Advisors[index] = team.Advisors[len(team.Advisors)-1]
	team.Advisors = team.Advisors[:len(team.Advisors)-1]

	team.NumberOfEmployees = len(team.TeamMembers) + len(team.Advisors) + len(team.Founders)
	err = db.connection.Save(&team).Error

	if err != nil {
		return model.Team{}, err
	}
	return team, nil
}

func (db Database) AddFounderToTeam(founderMember model.TeamMember, teamId string) (model.Team, error) {
	if founderMember.Id == 0 || teamId == "" {
		errMsg := fmt.Sprintf("Invalid Argument provided. One of the following params are null Team Id : %d, Team Member"+
			" ID : %s", founderMember.Id, teamId)
		return model.Team{}, errors.New(errMsg)
	}

	var team model.Team
	var founderMembers []model.TeamMember
	team, err := db.GetTeamByID(teamId)

	if err != nil {
		return model.Team{}, err
	}

	founderMembers = append(team.Founders, founderMember)
	team.Advisors = founderMembers
	team.NumberOfEmployees = len(team.TeamMembers) + len(team.Advisors) + len(team.Founders)

	err = db.connection.Save(&team).Error

	if err != nil {
		return model.Team{}, err
	}
	return team, nil
}

func (db Database) RemoveFounderFromTeam(founderMember model.TeamMember, teamId string) (model.Team, error) {
	if founderMember.Id == 0 || teamId == "" {
		errMsg := fmt.Sprintf("Invalid Argument provided. One of the following params are null Team Id : %d, Team Member"+
			" ID : %s", founderMember.Id, teamId)
		return model.Team{}, errors.New(errMsg)
	}

	var team model.Team
	team, err := db.GetTeamByID(teamId)

	if err != nil {
		return model.Team{}, err
	}

	var index int

	for i, item := range team.Advisors {
		if item.Id == founderMember.Id {
			index = i
			break
		}
	}

	team.Founders[index] = team.Founders[len(team.Founders)-1]
	team.Founders = team.Founders[:len(team.Founders)-1]

	team.NumberOfEmployees = len(team.TeamMembers) + len(team.Advisors) + len(team.Founders)

	err = db.connection.Save(&team).Error

	if err != nil {
		return model.Team{}, err
	}
	return team, nil
}

func (db Database) AddMemberToGroup(memberId string, groupId string) (model.Group, error) {
	group, err := db.GetGroupById(groupId)
	if err != nil {
		return model.Group{}, err
	}

	if memberId == "" {
		return model.Group{}, helper.ErrInvalidArgumentProvided
	}

	group.GroupMembers = append(group.GroupMembers, memberId)
	group.NumGroupMembers = len(group.GroupMembers)

	err = db.connection.Save(&group).Error

	if err != nil {
		return model.Group{}, err
	}
	return group, nil

}

func (db Database) RemoveMemberFromGroup(memberId string, groupId string) (model.Group, error) {
	group, err := db.GetGroupById(groupId)
	if err != nil {
		return model.Group{}, err
	}

	if memberId == "" {
		return model.Group{}, helper.ErrInvalidArgumentProvided
	}

	var index int

	for i, item := range group.GroupMembers {
		if item == memberId {
			index = i
			break
		}
	}

	group.GroupMembers[index] = group.GroupMembers[len(group.GroupMembers)-1]
	group.GroupMembers = group.GroupMembers[:len(group.GroupMembers)-1]
	group.NumGroupMembers = len(group.GroupMembers)

	err = db.connection.Save(&group).Error

	if err != nil {
		return model.Group{}, err
	}

	return group, nil
}

func (db Database) DeleteUserById(id string) (bool, error) {
	query := "id = ?"
	return db.DeleteUserBasedOnParam(id, query)
}

func (db Database) DeleteUserByUsername(id string) (bool, error) {
	query := "username = ?"
	return db.DeleteUserBasedOnParam(id, query)
}

func (db Database) DeleteUserByEmail(id string) (bool, error) {
	query := "email = ?"
	return db.DeleteUserBasedOnParam(id, query)
}

func (db Database) DeleteUserBasedOnParam(param string, query string) (bool, error) {
	var user model.User
	user, err := db.GetUserBasedOnParam(param, query)

	if err != nil {
		return false, err
	}

	if user.ID == 0 {
		return false, nil
	}

	err = db.connection.Set("gorm:delete_option", "OPTION (OPTIMIZE FOR UNKNOWN)").Delete(&user).Error

	if err != nil {
		return false, nil
	}
	return true, nil
}

func (db Database) DeleteTeamById(teamId string) (bool, error) {
	query := "id = ?"
	return db.DeleteTeamBasedOnParam(teamId, query)
}

func (db Database) DeleteTeamByName(teamName string) (bool, error) {
	query := "name = ?"
	return db.DeleteTeamBasedOnParam(teamName, query)
}

func (db Database) DeleteTeamTeamByEmail(teamEmail string) (bool, error) {
	query := "email = ?"
	return db.DeleteTeamBasedOnParam(teamEmail, query)
}

func (db Database) DeleteTeamBasedOnParam(param string, query string) (bool, error) {
	var team model.Team
	team, err := db.GetTeamBasedOnParam(param, query)

	if err != nil {
		return false, err
	}

	if team.ID == 0 {
		return false, nil
	}

	err = db.connection.Set("gorm:delete_option", "OPTION (OPTIMIZE FOR UNKNOWN)").Delete(&team).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func (db Database) DeleteGroupById(groupId string) (bool, error) {
	group, err := db.GetGroupById(groupId)
	if err != nil {
		return false, err
	}

	if group.ID == 0 {
		return false, nil
	}

	err = db.connection.Set("gorm:delete_option", "OPTION (OPTIMIZE FOR UNKNOWN)").Delete(&group).Error
	if err != nil {
		return false, err
	}
	return true, nil
}
