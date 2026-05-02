package utils

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendEmail(to, subject, body string) error {
	from := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")

	auth := smtp.PlainAuth("", from, password, "smtp.gmail.com")

	msg := fmt.Sprintf("From: FoodByMegha <%s>\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s", from, to, subject, body)

	err := smtp.SendMail("smtp.gmail.com:587", auth, from, []string{to}, []byte(msg))
	return err
}

func PaymentSuccessEmail(customerEmail, customerName, planName string, amount float64) {
	// Customer ko email
	customerBody := fmt.Sprintf(`
		<h2>🍱 FoodByMegha</h2>
		<p>Namaste <strong>%s</strong>!</p>
		<p>Tumhari payment successful ho gayi! 🎉</p>
		<ul>
			<li>Plan: <strong>%s</strong></li>
			<li>Amount: <strong>₹%.0f</strong></li>
		</ul>
		<p>Kal se tumhara tiffin aana shuru ho jayega! 🍛</p>
		<p>- Team FoodByMegha</p>
	`, customerName, planName, amount)

	SendEmail(customerEmail, "Payment Successful! Tiffin Pakka! 🍱", customerBody)

	// Admin ko email
	adminBody := fmt.Sprintf(`
		<h2>💰 Naya Customer!</h2>
		<p>Customer: <strong>%s</strong></p>
		<p>Plan: <strong>%s</strong></p>
		<p>Amount: <strong>₹%.0f</strong></p>
	`, customerName, planName, amount)

	SendEmail(os.Getenv("ADMIN_EMAIL"), "Naya Subscription Aaya! 💰", adminBody)
}
