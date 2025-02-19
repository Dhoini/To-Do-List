package notes

import (
	"time"
)

type Note struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"default:Untitled;size:100" json:"title"`
	Content   string    `gorm:"type:text;size:10000" json:"content"`
	Status    string    `gorm:"default:'open';size:20" json:"status"`
	UserID    string    `gorm:"not null" json:"user_id"` // Внешний ключ
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
