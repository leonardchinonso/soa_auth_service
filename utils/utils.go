package utils

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/gomail.v2"
)

// ShouldBePresentString checks that a string in a field is required
func ShouldBePresentString(value string, prefix string, errs *[]error) {
	if value == "" {
		*errs = append(*errs, fmt.Errorf(fmt.Sprintf("%s is required", prefix)))
	}
}

// SendMailAsPlainText sends a mail to a recipient as plain text
func SendMailAsPlainText(from, to, subject, message, clientname, password string) error {
	m := gomail.NewMessage()

	// set the headers for the email to send
	m.SetHeader("From", from)
	m.SetHeader("to", to)
	m.SetHeader("Subject", subject)

	// set the body of the email as plain text
	m.SetBody("text/plain", message)

	// init settings for SMTP server
	d := gomail.NewDialer("smtp.gmail.com", 587, clientname, password)

	// TODO: should set to false on production, or use env variables
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Error sending email to recipient. Error: %v\n", err)
		return err
	}

	return nil
}

// SendSimpleMailSMTP sends a mail to a recipient as plain text with SMTP
func SendSimpleMailSMTP(from, to, subject, body, clientname, password, smtpHost, smtpPort string) error {
	log.Printf("FROM: %v TO: %v", from, to)
	smtpAuth := smtp.PlainAuth("", clientname, password, smtpHost)

	msg := fmt.Sprintf("Subject: %s\n%s", subject, body)

	err := smtp.SendMail(fmt.Sprintf("%s:%s", smtpHost, smtpPort), smtpAuth, from, []string{to}, []byte(msg))
	if err != nil {
		log.Printf("Error sending email from: %v to %v. Error: %v\n", from, to, err)
		return err
	}

	return nil
}

// CurrentPrimitiveTime returns the current time in primitive.Timestamp
func CurrentPrimitiveTime() primitive.Timestamp {
	return primitive.Timestamp{
		T: uint32(time.Now().Unix()),
	}
}
