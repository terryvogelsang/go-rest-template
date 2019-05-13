package models

import (
	json "encoding/json"
	utils "vulnlabs-rest-api/utils"

	gormlib "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const (
	BoatsTable = "boats"
)

// GORMInterface : GORM Communication interface
type GORMInterface interface {
	CloseConnection() error
	CreateUser(userCreateRequestBody *UserCreateRequestBody) (*User, error)
	ReadUserFromEmail(email string) (*User, error)
	ReadUserFromID(id string) (*User, error)
	UpdateUserInfos(user *User, userUpdateRequestBody *UserUpdateRequestBody) error
	UpdateUserPassword(newHashedPassword string) error
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

// ReadUserFromEmail : Read user from DB
func (gorm *GORM) ReadUserFromEmail(email string) (*User, error) {

	var user User

	return &user, gorm.Database.First(&user, "email = ?", email).Error
}

// ReadUserFromID : Read user from DB
func (gorm *GORM) ReadUserFromID(id string) (*User, error) {

	var user User

	return &user, gorm.Database.Where("id = ?", id).First(&user).Error
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

// DeleteUser : Delete user from DB
func (gorm *GORM) DeleteUser(user *User) error {

	return gorm.Database.Delete(&user).Error
}

// IsRecordNotFoundError : Check if error is of type RecordNotFound
func (gorm *GORM) IsRecordNotFoundError(err error) bool {
	return gormlib.IsRecordNotFoundError(err)
}
