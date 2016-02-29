package VehicleLib

import (
	"fmt"

	twilio "github.com/carlosdp/twiliogo"
)

// SendText ...
func SendText(accountSID string, authToken string, senderPhoneNumber string, recipientPhoneNumber string, body string) error {
	client := twilio.NewClient(accountSID, authToken)

	message, err := twilio.NewMessage(client, senderPhoneNumber, recipientPhoneNumber, twilio.Body(body))
	if err != nil {
		return fmt.Errorf("sendText error: %s", err)
	}

	fmt.Printf("SendText status %v\n", message.Status)

	return nil
}
