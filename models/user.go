package models

import (
	utils "vulnlabs-rest-api/utils"

	gormlib "github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

const (
	DEFAULT_ROLE = "H4x0r"
)

// User : User Account Struct
type User struct {
	ID                string         `json:"id" gorm:"primary_key;unique;not null;"`
	Email             string         `json:"email,omitempty" gorm:"unique;not null;"`
	Password          string         `json:"-" gorm:"not null;"`
	FirstName         string         `json:"firstName,omitempty" gorm:"not null;"`
	LastName          string         `json:"lastName,omitempty" gorm:"not null;"`
	PhoneNumber       string         `json:"phoneNumber,omitempty" gorm:"not null;"`
	ProfilePictureURL string         `json:"profilePictureURL,omitempty"`
	IsCommittee       bool           `json:"isCommittee" gorm:"not null;"`
	Role              ReadOnlyString `json:"role,omitempty" gorm:"not null;"`
}

type UserCreateRequestBody struct {
	Email             string `json:"email,omitempty"`
	Password          string `json:"password,omitempty"`
	FirstName         string `json:"firstName,omitempty"`
	LastName          string `json:"lastName,omitempty"`
	PhoneNumber       string `json:"phoneNumber,omitempty"`
	ProfilePictureURL string `json:"profilePictureURL,omitempty"`
}

type UserUpdateRequestBody struct {
	Email             string `json:"email,omitempty"`
	FirstName         string `json:"firstName,omitempty"`
	LastName          string `json:"lastName,omitempty"`
	PhoneNumber       string `json:"phoneNumber,omitempty"`
	ProfilePictureURL string `json:"profilePictureURL,omitempty"`
}

//BeforeCreate : Run before DB Insertion
func (user *User) BeforeCreate(scope *gormlib.Scope) error {

	// Set UserID
	user.ID = uuid.NewV4().String()

	// If user is in admin list, get role admin
	if utils.IsStringIn(user.Email, GlobalConfig.AdminUsers) {
		user.Role = "Admin"
	} else {
		user.Role = DEFAULT_ROLE
	}

	return nil
}
