package modules

import (
	"gorm.io/gorm"
)

// Chat represents a chat session for a user
type Chat struct {
	gorm.Model
	UserID uint   `gorm:"not null" json:"user_id"`
	Title  string `gorm:"type:varchar(255);default:'New Chat'" json:"title"`
}
