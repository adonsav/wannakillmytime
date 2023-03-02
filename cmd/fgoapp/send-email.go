package main

import (
	"fmt"
	"github.com/adonsav/fgoapp/internal/models"
	mail "github.com/xhit/go-simple-mail/v2"
	"os"
	"strings"
	"time"
)

// listenForEmail creates a goroutine that listens for emails in the respective channels and
// then forwards them properly.
func listenForEmail() {
	go func() {
		for {
			message := <-appConfig.EmailChan
			sendEmail(message)
		}
	}()
}

// sendEmail builds and sends an email.
func sendEmail(ed models.EmailData) {
	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	client, err := server.Connect()
	if err != nil {
		appConfig.ErrorLog.Println(err)
	}

	email := mail.NewMSG()
	email.SetFrom(ed.From).AddTo(ed.To).SetSubject(ed.Subject)
	if ed.Template == "" {
		email.SetBody(mail.TextHTML, ed.Content)
	} else {
		data, err := os.ReadFile(fmt.Sprintf("internal/templates/emailtemplates/%s", ed.Template))
		if err != nil {
			appConfig.ErrorLog.Println(err)
		}
		emailTemplate := string(data)
		msgToSend := strings.Replace(emailTemplate, "[%body%]", ed.Content, 1)
		email.SetBody(mail.TextHTML, msgToSend)
	}
	err = email.Send(client)
	if err != nil {
		appConfig.ErrorLog.Println(err)
	} else {
		appConfig.InfoLog.Println("Email sent")
	}
}
