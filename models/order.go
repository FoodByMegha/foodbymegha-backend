package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	UserID         uint   `json:"user_id"`
	SubscriptionID uint   `json:"subscription_id"`
	DeliveryDate   string `json:"delivery_date"` // "2026-04-07"
	Status         string `json:"status"`        // "pending", "out_for_delivery", "delivered"
	Address        string `json:"address"`
	Notes          string `json:"notes"` // "kam mirch daalna"

	// Relations
	User         User         `json:"user" gorm:"foreignKey:UserID"`
	Subscription Subscription `json:"subscription" gorm:"foreignKey:SubscriptionID"`
}
