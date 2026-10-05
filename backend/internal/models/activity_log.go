package models

import "time"

type ActivityLog struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      *uint     `json:"user_id" gorm:"column:user_id"`
	UserName    string    `json:"user_name" gorm:"column:user_name;not null"`
	Action      string    `json:"action" gorm:"column:action;not null"`
	Module      string    `json:"module" gorm:"column:module;not null"`
	Description string    `json:"description" gorm:"column:description;not null"`
	IPAddress   string    `json:"ip_address" gorm:"column:ip_address"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at"`
}
