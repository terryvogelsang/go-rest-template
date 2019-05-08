package models

import (
	json "encoding/json"
	utils "mycnc-rest-api/utils"
	"time"

	gormlib "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	uuid "github.com/satori/go.uuid"
)

const (
	BoatsTable = "boats"
)

// GORMInterface : GORM Communication interface
type GORMInterface interface {
	CloseConnection() error
	CreateUser(userCreateRequestBody *UserCreateRequestBody) (*User, error)
	CreateRegatta(regattaCreateRequestBody *RegattaCreateRequestBody) (*Regatta, error)
	ReadUserFromEmail(email string) (*User, error)
	ReadUserFromID(id string) (*User, error)
	ReadUserBoatsInfos(user *User) ([]*Boat, error)
	ReadRegattaFromID(id string) (*Regatta, error)
	UpdateUserInfos(user *User, userUpdateRequestBody *UserUpdateRequestBody) error
	UpdateUserBoats(user *User, boatCreateRequestBody *BoatCreateRequestBody) error
	UpdateUserPassword(newHashedPassword string) error
	UpdateRegattaBoatChrono(regatta *Regatta, boatID string) error
	DeleteUser(user *User) error
	IsRecordNotFoundError(err error) bool
}

// GORM : GORM communication interface
type GORM struct {
	Database *gormlib.DB
}

// NewGORM : Return a new GORM abstraction struct
func NewGORM(connectionURL string) *GORM {

	// Initialize the GORM connection to a MariaDB instance running on your local machine
	db, err := gormlib.Open("mysql", connectionURL)

	if err != nil {
		utils.PanicOnError(err, "Failed to connect to MariaDB")
	}

	// db.DropTableIfExists(&User{})

	// Setup
	db.LogMode(true)
	db = db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8 auto_increment=1").Set("gorm:auto_preload", true)

	// Migrate DB Schemas
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Boat{})
	db.AutoMigrate(&BoatClass{})
	db.AutoMigrate(&Regatta{})
	db.AutoMigrate(&Ranking{})
	db.AutoMigrate(&Rank{})
	db.AutoMigrate(&ChronoEntry{})

	// Default boat class here
	dartClass := &BoatClass{
		Name:                 "Dart",
		YardstickCoefficient: 50,
	}

	db.Create(&dartClass)

	// Return new MongoDB abstraction struct
	return &GORM{
		Database: db,
	}
}

// CloseConnection : Close GORM Connection
func (gorm *GORM) CloseConnection() error {

	return gorm.Database.Close()
}

// CreateUser : Store user in DB
func (gorm *GORM) CreateUser(userCreateRequestBody *UserCreateRequestBody) (*User, error) {

	var user User
	marshalled, _ := json.Marshal(userCreateRequestBody)
	json.Unmarshal(marshalled, &user)

	user.Password = userCreateRequestBody.Password

	err := gorm.Database.Create(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, gorm.Database.Save(&user).Error

}

// CreateRegatta : Store regatta in DB
func (gorm *GORM) CreateRegatta(regattaCreateRequestBody *RegattaCreateRequestBody) (*Regatta, error) {

	var regatta Regatta
	marshalled, _ := json.Marshal(regattaCreateRequestBody)
	json.Unmarshal(marshalled, &regatta)

	// Fetch users from provided IDs
	err := gorm.Database.First(&regatta.BuoyResponsible, "id = ?", regattaCreateRequestBody.BuoyResponsibleID).Error
	if err != nil {
		return nil, err
	}

	err = gorm.Database.First(&regatta.LocalResponsible, "id = ?", regattaCreateRequestBody.LocalResponsibleID).Error
	if err != nil {
		return nil, err
	}

	err = gorm.Database.First(&regatta.StartComitteeMan, "id = ?", regattaCreateRequestBody.StartComitteeManID).Error
	if err != nil {
		return nil, err
	}

	err = gorm.Database.First(&regatta.FirstAssistant, "id = ?", regattaCreateRequestBody.FirstAssistantID).Error
	if err != nil {
		return nil, err
	}

	err = gorm.Database.First(&regatta.SecondAssistant, "id = ?", regattaCreateRequestBody.SecondAssistantID).Error
	if err != nil {
		return nil, err
	}

	// Fetch committee man boat for ranking init
	var startCommitteeManUserBoat Boat

	err = gorm.Database.Find(&startCommitteeManUserBoat, "owner_id = ? AND is_current = ?", regatta.StartComitteeMan.ID, 1).Error
	if err != nil {
		return nil, err
	}

	regatta.RegisteredBoats = []Boat{}
	regatta.ParticipatingBoats = []Boat{}

	// StartCommiteeMan ranks 0 per the rules
	committeeRank := Rank{
		BoatID:     startCommitteeManUserBoat.ID,
		RankNumber: 0,
	}

	// Create regatta ranking instance
	regatta.Ranking = Ranking{
		ID:       uuid.NewV4().String(),
		Type:     "A", // Default to A type regatta (Time / Defined laps number)
		IsPublic: false,
		Ranks:    []Rank{committeeRank},
	}

	err = gorm.Database.Create(&regatta).Error

	if err != nil {
		return nil, err
	}

	return &regatta, nil

}

// ReadUserFromEmail : Read user from DB
func (gorm *GORM) ReadUserFromEmail(email string) (*User, error) {

	var user User

	return &user, gorm.Database.First(&user, "email = ?", email).Error
}

// ReadUserBoatsInfos : Read boat informations from id
func (gorm *GORM) ReadUserBoatsInfos(user *User) ([]*Boat, error) {

	var boats []*Boat

	err := gorm.Database.Model(user).Association("Boats").Find(&boats).Error

	if err != nil {
		return nil, err
	}

	return boats, nil
}

// ReadUserFromID : Read user from DB
func (gorm *GORM) ReadUserFromID(id string) (*User, error) {

	var user User

	return &user, gorm.Database.Where("id = ?", id).First(&user).Error
}

// ReadRegattaFromID : Read regatta from DB
func (gorm *GORM) ReadRegattaFromID(id string) (*Regatta, error) {

	var regatta Regatta

	err := gorm.Database.First(&regatta, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &regatta, nil
}

// UpdateUserInfos : Update user infos in DB
func (gorm *GORM) UpdateUserInfos(user *User, userUpdateRequestBody *UserUpdateRequestBody) error {

	marshalled, _ := json.Marshal(userUpdateRequestBody)
	json.Unmarshal(marshalled, user)

	return gorm.Database.Save(user).Error
}

// UpdateUserPassword : Update user password in DB
func (gorm *GORM) UpdateUserPassword(newHashedPassword string) error {

	return gorm.Database.Select("password").Updates(map[string]interface{}{"password": newHashedPassword}).Error
}

func (gorm *GORM) UpdateRegattaBoatChrono(regatta *Regatta, boatID string) error {

	var count int
	err := gorm.Database.Model(&Boat{}).Where("id = ?", boatID).Count(&count).Error

	if err != nil {
		return err
	}

	if count == 0 {
		return gormlib.ErrRecordNotFound
	}

	boatChrono := ChronoEntry{
		ID:        uuid.NewV4().String(),
		BoatID:    boatID,
		Timestamp: time.Now(),
	}

	return gorm.Database.Model(regatta).Association("ChronosEntries").Append(boatChrono).Error
}

// UpdateUserBoats : Add new boat to user in DB
func (gorm *GORM) UpdateUserBoats(user *User, boatCreateRequestBody *BoatCreateRequestBody) error {

	var boat Boat

	marshalled, _ := json.Marshal(boatCreateRequestBody)
	json.Unmarshal(marshalled, &boat)

	// Get Boat class from DB
	var boatClass BoatClass
	err := gorm.Database.First(&boatClass, "id = ?", boatCreateRequestBody.ClassID).Error

	if err != nil {
		return err
	}

	// Set isCurrent to 0 for each boat
	err = gorm.Database.Table(BoatsTable).Where("owner_id = ? AND is_current = ?", user.ID, 1).Update("is_current", 0).Error

	// Get all boats of user with changed state
	var boats []Boat
	err = gorm.Database.First(&boats, "owner_id = ?", user.ID).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return err
	}

	// Form new boat
	boat.Class = boatClass
	boat.ID = uuid.NewV4().String()
	boat.IsCurrent = true

	// Append new boat to existing boats collection
	boats = append(boats, boat)

	// Replace old boats collection
	return gorm.Database.Model(user).Association("Boats").Replace(&boats).Error

}

// DeleteUser : Delete user from DB
func (gorm *GORM) DeleteUser(user *User) error {

	return gorm.Database.Delete(&user).Error
}

// IsRecordNotFoundError : Check if error is of type RecordNotFound
func (gorm *GORM) IsRecordNotFoundError(err error) bool {
	return gormlib.IsRecordNotFoundError(err)
}
