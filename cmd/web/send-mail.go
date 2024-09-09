package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/andreadebortoli2/GO-bnb/internal/models"
	mail "github.com/xhit/go-simple-mail/v2"
)

func listenForMail() {
	// anonymous function executed in background in a go routine
	go func() {
		for {
			msg := <-app.MailChan
			sendMsg(msg)
		}
	}()
}

func sendMsg(m models.MailData) {
	// set mail server with all data(some data are absent cause of the dummy mail server, needed in production, have to be found in mail server datas like username,password, encryption,...)
	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	// set the client
	client, err := server.Connect()
	if err != nil {
		errorLog.Println(err)
	}

	// create the mail message
	email := mail.NewMSG()
	// set from address, to address and subject from the MailData struct
	email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)
	// set the message body
	if m.Template == "" {
		email.SetBody(mail.TextHTML, string(m.Content))
	} else {
		data, err := os.ReadFile(fmt.Sprintf("./email-templates/%s", m.Template))
		if err != nil {
			app.ErrorLog.Println(err)
		}
		mailTemplate := string(data)
		msgToSend := strings.Replace(mailTemplate, "[%body%]", m.Content, 1)
		email.SetBody(mail.TextHTML, msgToSend)
	}

	// send the email
	err = email.Send(client)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Email sent!")
	}

}
