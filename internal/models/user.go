package models

import "time"

type User struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null;size:200" json:"name"`
	Email     string    `gorm:"unique;not null" json:"email"`
	Password  string    `gorm:"not null;size:200" json:"password"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	Notes     []Note    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"notes"`
}
