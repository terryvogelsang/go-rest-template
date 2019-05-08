package models

import (
	"time"

	gormlib "github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// Numéro (Ordonné)
// Date
// Heure publiée de départ
// Heure réelle de départ
// Lieu
// Responsable Bouée (Membre ou externe)
// Comité départ (Membre ou externe )
// Aide 1 (Membre ou externe )
// Aide 2 (Membre ou externe)
// Responsable local (Membre Comité)
// Bateaux inscrits
// Bateaux participants (Possibilité Guest (Infos + remarques)
// Classement

// Regatta : Represents a regatta
type Regatta struct {
	ID                     string           `json:"id" gorm:"primary_key;unique;not null;"`
	Identifier             string           `json:"identifier" gorm:"not null;"`
	ScheduledStartDateTime time.Time        `json:"scheduledStartDateTime" gorm:"not null;"`
	RealStartDateTime      time.Time        `json:"realStartDateTime" gorm:"not null;"`
	Location               string           `json:"location" gorm:"not null;"`
	BuoyResponsibleID      string           `json:"-" gorm:"not null;"`                                                       // Member or External User
	BuoyResponsible        User             `json:"buoyResponsible,omitempty" gorm:"not null;foreignkey:BuoyResponsibleID"`   // Member or External User
	StartComitteeManID     string           `json:"-" gorm:"not null;"`                                                       // Member or External User
	StartComitteeMan       User             `json:"startComitteeMan,omitempty" gorm:"not null;foreignkey:StartComitteeManID"` // Member or External User
	FirstAssistantID       string           `json:"-" gorm:"not null;"`                                                       // Member or External User
	FirstAssistant         User             `json:"firstAssistant,omitempty" gorm:"not null;foreignkey:FirstAssistantID"`     // Member or External User
	SecondAssistantID      string           `json:"-" gorm:"not null;"`                                                       // Member or External User
	SecondAssistant        User             `json:"secondAssistant,omitempty" gorm:"not null;foreignkey:SecondAssistantID"`   // Member or External User
	LocalResponsibleID     string           `json:"-" gorm:"not null;"`                                                       // Member or External User
	LocalResponsible       User             `json:"localResponsible,omitempty" gorm:"not null;foreignkey:LocalResponsibleID"` // Comittee Member
	RegisteredBoats        []Boat           `json:"registeredBoats,omitempty" gorm:"not null;many2many:regatta_registeredboats"`
	ParticipatingBoats     []Boat           `json:"participatingBoats,omitempty" gorm:"not null;many2many:regatta_participatingboats"` // Possible Guests
	RankingID              string           `json:"-" gorm:"not null;"`
	Ranking                Ranking          `json:"ranking,omitempty" gorm:"not null;foreignkey:RankingID"`
	LapChronosEntries      []LapChronoEntry `json:"lapChronoEntries,omitempty" gorm:"not null;foreignkey:RegattaID"`
}

// Ranking : Regatta Ranking struct
type Ranking struct {
	ID       string `json:"id" gorm:"primary_key;unique;not null;"`
	Type     string `json:"type" gorm:"not null;"`
	IsPublic bool   `json:"isPublic" gorm:"not null;"`
	Ranks    []Rank `json:"ranks" gorm:"not null;foreignkey:RankingID;"`
}

// Rank : Represent a boat rank in a ranking
type Rank struct {
	BoatID     string `json:"boatID" gorm:"not null;primary_key;"`
	RankingID  string `json:"rankingID" gorm:"not null;primary_key;"`
	RankNumber int    `json:"rankNumber" gorm:"not null;"`
}

// LapChronoEntry : Struct representing a timestamp on which a boat finishes a regatta lap
type LapChronoEntry struct {
	BoatID    string    `json:"boatID" gorm:"not null;primary_key;"`
	RegattaID string    `json:"regattaID" gorm:"not null;primary_key;"`
	Timestamp time.Time `json:"timestamp" gorm:"not null;"`
}

//BeforeCreate : Run before DB Insertion
func (Regatta *Regatta) BeforeCreate(scope *gormlib.Scope) error {

	// Set ID
	Regatta.ID = uuid.NewV4().String()

	return nil
}

// RegattaCreateRequestBody : Update boat class request ID
type RegattaCreateRequestBody struct {
	Identifier             string    `json:"identifier"`
	ScheduledStartDateTime time.Time `json:"scheduledStartDateTime"`
	RealStartDateTime      time.Time `json:"realStartDateTime"`
	Location               string    `json:"location"`
	BuoyResponsibleID      string    `json:"buoyResponsibleID,omitempty"`  // Member or External User
	StartComitteeManID     string    `json:"startComitteeManID,omitempty"` // Committee Member or External User
	FirstAssistantID       string    `json:"firstAssistantID,omitempty"`   // Member or External User
	SecondAssistantID      string    `json:"secondAssistantID,omitempty"`  // Member or External User
	LocalResponsibleID     string    `json:"localResponsibleID,omitempty"` // Comittee Member
}
