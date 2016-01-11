package main

import "fmt"
import "gopkg.in/gomail.v2"

func main() {
	fmt.Println("DEBUT")
	m := gomail.NewMessage()
	m.SetHeader("From", "x3@sra.fr")
	m.SetHeader("To", "cbonnet@sra.fr")
	//m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/html", "Hello <b>Bob</b> and <i>Cora</i>!")
	//m.Attach("/home/Alex/lolcat.jpg")

	d := gomail.NewPlainDialer("smtp.free.fr", 25, "systeme@sra.fr", "systeme2000")

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
