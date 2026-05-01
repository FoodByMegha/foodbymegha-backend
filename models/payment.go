package models

import "gorm.io/gorm"

type Payment struct {
	gorm.Model
	UserID            uint    `json:"user_id"`
	PlanID            uint    `json:"plan_id"`
	RazorpayOrderID   string  `json:"razorpay_order_id"`   // Razorpay ka order ID
	RazorpayPaymentID string  `json:"razorpay_payment_id"` // Payment hone ke baad milega
	Amount            float64 `json:"amount"`
	Currency          string  `json:"currency"` // "INR"
	Status            string  `json:"status"`   // "created", "paid", "failed"

	// Relations
	User User `json:"user" gorm:"foreignKey:UserID"`
	Plan Plan `json:"plan" gorm:"foreignKey:PlanID"`
}
