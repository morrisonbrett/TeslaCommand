package VehicleLib

import (
	"fmt"
	"log"

	twilio "github.com/carlosdp/twiliogo"
)

// SendText ...
func SendText(logger *log.Logger, accountSID string, authToken string, senderPhoneNumber string, recipientPhoneNumber string, body string) error {
	client := twilio.NewClient(accountSID, authToken)

	message, err := twilio.NewMessage(client, senderPhoneNumber, recipientPhoneNumber, twilio.Body(body))
	if err != nil {
		return fmt.Errorf("sendText error: %s", err)
	}

	logger.Printf("SendText status %v\n", message.Status)

	return nil
}
