package utils

import (
	"fmt"
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
	gomail "gopkg.in/mail.v2"
)

// SendMail sends top up link message
func SendMail(senderMail, destinationMail, invoiceID, paymentLink string) error {
	// Create a new message
	message := gomail.NewMessage()

	// Set email headers with multiple recipients
	message.SetHeader("From", senderMail)
	message.SetHeader("To", destinationMail)
	message.SetHeader("Subject", "Book Rent Top Up Payment")

	// Set email body
	message.SetBody("text/html", fmt.Sprintf(`
        <html>
            <body>
                <h1>Top Up Payment</h1>
                <p>The following are your top up payment link.</p>
				<p><strong>Invoice ID: </strong> %s<br>
				<strong>Payment Link: </strong><a href="%s">payment-link</a></p>
                <p>Thank you,<br>Book Rent</p>
            </body>
        </html>
		
    `, invoiceID, paymentLink,
	))

	// Set up the SMTP dialer
	mailtrapPort, err := strconv.Atoi(os.Getenv("MAILTRAP_PORT"))
	if err != nil {
		return fmt.Errorf("sending top up email: %w", err)
	}
	dialer := gomail.NewDialer(
		os.Getenv("MAILTRAP_HOST"),
		mailtrapPort,
		os.Getenv("MAILTRAP_USERNAME"),
		os.Getenv("MAILTRAP_PASSWORD"),
	)
	// Send the email
	return dialer.DialAndSend(message)
}
