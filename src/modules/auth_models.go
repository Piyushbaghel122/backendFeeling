package modules

import (
	"time"
)

type User struct {
	UserId    uint      `gorm:"primarykey;autoIncrement" json:"user_id"`
	UserName  string    `gorm:"not null" json:"username"`
	Email     string    `gorm:"unique;not null" json:"email"`
	Password  string    `gorm:"not null" json:"password"`
	AvatarURL string    `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}