package models

import (
	"time"

	"gorm.io/gorm"
)

type Plan struct {
	gorm.Model
	Name         string  `json:"name"`          // "Monthly Tiffin"
	Description  string  `json:"description"`   // "30 din ka tiffin"
	Price        float64 `json:"price"`         // 3999.00
	DurationDays int     `json:"duration_days"` // 30
	IsActive     bool    `json:"is_active"`     // true/false
}

type Subscription struct {
	gorm.Model
	UserID    uint      `json:"user_id"`
	PlanID    uint      `json:"plan_id"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	IsActive  bool      `json:"is_active"`

	// Relations
	User User `json:"user" gorm:"foreignKey:UserID"`
	Plan Plan `json:"plan" gorm:"foreignKey:PlanID"`
}
