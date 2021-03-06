package models

import (
	"time"
)

//Audition data model
type Audition struct {
	ID        int `gorm:"UNIQUE;NOT NULL;PRIMARY_KEY;AUTO_INCREMENT"`
	Club      Club
	User      User
	Status    string
	Priority  uint16
	UpdatedAt time.Time
	CreatedAt time.Time
}
