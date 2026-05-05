package cron

import (
	"log"
	"time"

	"github.com/FoodByMegha/foodbymegha-backend/config"
	"github.com/FoodByMegha/foodbymegha-backend/models"
	robfigcron "github.com/robfig/cron/v3"
)

func StartCronJobs() {
	c := robfigcron.New()
	c.AddFunc("0 0 * * *", CreateDailyOrders)
	c.Start()
	log.Println("✅ Cron job shuru ho gaya!")
}

func CreateDailyOrders() {
	today := time.Now().Format("2006-01-02")
	log.Println("🍱 Aaj ke orders ban rahe hain:", today)

	var subscriptions []models.Subscription
	config.DB.Preload("User").Where("is_active = ? AND end_date > ?", true, time.Now()).Find(&subscriptions)

	for _, sub := range subscriptions {
		var existing models.Order
		result := config.DB.Where("user_id = ? AND delivery_date = ?", sub.UserID, today).First(&existing)
		if result.Error == nil {
			log.Printf("⚠️ User %d ka aaj ka order already hai\n", sub.UserID)
			continue
		}

		var lastOrder models.Order
		addressResult := config.DB.Where("user_id = ?", sub.UserID).
			Order("created_at DESC").First(&lastOrder)

		address := "Sarojini Nagar, Lucknow"
		if addressResult.Error == nil {
			address = lastOrder.Address
		}

		order := models.Order{
			UserID:         sub.UserID,
			SubscriptionID: sub.ID,
			DeliveryDate:   today,
			Status:         "pending",
			Address:        address,
			Notes:          "",
		}

		config.DB.Create(&order)
		log.Printf("✅ User %d ka order create ho gaya!\n", sub.UserID)
	}
}
