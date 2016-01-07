package main

import (
	"log"
	"net/smtp"
)

func main() {
	// Set up authentication information.
	auth := smtp.PlainAuth("", "systeme@ra.fr", "system2000", "smtp.free.fr")

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	to := []string{"cbonnet@sra.fr"}
	msg := []byte("To: cbonnet@sra.fr\r\n" +
		"Subject: discount Gophers!\r\n" +
		"\r\n" +
		"This is the email body.\r\n")
	err := smtp.SendMail("smtp.free.fr:25", auth, "systeme@sra.fr", to, msg)
	if err != nil {
		log.Fatal(err)
	}
}
