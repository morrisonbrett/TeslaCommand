package VehicleLib

import (
	"fmt"
	"net/smtp"
	"strconv"
)

// SendMail ...
func SendMail(mailServer string, mailServerPort int, mailServerLogin string, mailServerPassword string, fromAddress string, toAddress string, subj string, body string) error {
	// Set up authentication information.
	var auth smtp.Auth
	if len(mailServerLogin) > 0 {
		auth = smtp.PlainAuth("", mailServerLogin, mailServerPassword, mailServer)
	} else {
		auth = nil
	}

	// Connect to the server, authenticate, set the sender and recipient, and send the email in one step.
	to := []string{toAddress}
	msg := []byte("To: " + toAddress + "\r\nSubject: " + subj + "\r\n\r\n" + body + "\r\n")
	serverPort := mailServer + ":" + strconv.Itoa(mailServerPort)
	fmt.Printf("Sending mail via server %s\n", serverPort)
	err := smtp.SendMail(serverPort, auth, fromAddress, to, msg)
	if err != nil {
		return fmt.Errorf("sendMail error: %s", err)
	}

	return nil
}
