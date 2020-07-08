package NotifyLib
// Most mail servers require a FROM: address. 

import (
	"fmt"
	"log"
	"net/smtp"
	"strconv"
)

// SendMail ...
func SendMail(logger *log.Logger, mailServer string, mailServerPort int, mailServerLogin string, mailServerPassword string, fromAddress string, toAddress1 string, toAddress2 string, subj string, body string) error {
	// Set up authentication information.
	var auth smtp.Auth
	if len(mailServerLogin) > 0 {
		auth = smtp.PlainAuth("", mailServerLogin, mailServerPassword, mailServer)
	}

	//set up 2 email addresses
	to := []string{toAddress1, toAddress2}
	// Connect to the server, authenticate, set the sender and recipient, and send the email in one step.
	msg := []byte("To: " + toAddress1 + "\r\nFrom: " + fromAddress + "\r\nSubject: " + subj + "\r\n" + body + "\r\n")
	serverPort := mailServer + ":" + strconv.Itoa(mailServerPort)
	logger.Printf("Sending mail via server %s\n", serverPort)
	err := smtp.SendMail(serverPort, auth, fromAddress, to, msg)
	if err != nil {
		return fmt.Errorf("sendMail error: %s", err)
	}

	return nil
}
