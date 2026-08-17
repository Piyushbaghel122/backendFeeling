package modules

import (
	"gorm.io/gorm"
)

// Message represents the chat message model in the database
type Message struct {
	gorm.Model
	ChatID  uint   `gorm:"not null" json:"chat_id"`
	Content string `gorm:"type:text;not null" json:"content"`
	Role    string `gorm:"type:enum('user','ai');not null" json:"role"`
}