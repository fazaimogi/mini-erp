package models

import "time"

type User struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Name         string    `json:"name" gorm:"column:name;not null"`
	Email        string    `json:"email" gorm:"column:email;unique;not null"`
	PasswordHash string    `json:"-" gorm:"column:password_hash;not null"`
	RoleID       uint      `json:"role_id" gorm:"column:role_id;not null"`
	Role         string    `json:"role" gorm:"-"`
	Status       string    `json:"status" gorm:"column:status;not null"`
	CreatedAt    time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"column:updated_at"`
}
