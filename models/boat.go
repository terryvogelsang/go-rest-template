package models

import (
	gormlib "github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// Boat : Represents a boat
// A boat belongs to a Class
type Boat struct {
	ID         string    `json:"id" gorm:"primary_key;unique;not null;"`
	OwnerID    string    `json:"-" gorm:"not null;"`
	Name       string    `json:"name,omitempty" gorm:"not null;"`
	ClassID    string    `json:"-" gorm:"not null;"`
	Class      BoatClass `json:"class,omitempty" gorm:"not null;foreignkey:ClassID;"`
	SailNumber string    `json:"sailNumber,omitempty" gorm:"not null;"`
	IsCurrent  bool      `json:"isCurrent,omitempty" gorm:"not null;"`
}

// BoatClass : Represents boat type and yardstick
type BoatClass struct {
	ID                   string `json:"id" gorm:"primary_key;unique;not null;"`
	Name                 string `json:"name,omitempty" gorm:"not null;"`
	YardstickCoefficient int    `json:"yardstickCoefficient,omitempty" gorm:"not null;"`
}

//BeforeCreate : Run before DB Insertion
func (boatClass *BoatClass) BeforeCreate(scope *gormlib.Scope) error {

	// Set UserID
	boatClass.ID = uuid.NewV4().String()

	return nil
}

// BoatCreateRequestBody : Update boat class request ID
type BoatCreateRequestBody struct {
	Name       string `json:"name,omitempty" gorm:"not null;"`
	ClassID    string `json:"classID,omitempty" gorm:"not null;"`
	SailNumber string `json:"sailNumber,omitempty" gorm:"not null;"`
}
