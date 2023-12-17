package email

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/smtp"
	"sender/config"
)

type kafkaMessage struct {
	MessageType string            `json:"messagetype"`
	Data        map[string]string `json:"data"`
}

func SendEmail(cfg *config.Config, value []byte) error {

	var mes kafkaMessage
	var err error
	err = json.Unmarshal(value, &mes)
	if err != nil {
		return err
	}

	switch mes.MessageType {
	case "Registartion":
		err = SendRegistrationEmail(cfg, mes.Data)
	default:
		err = errors.New("no message type for email")
	}

	return err
}

func SendRegistrationEmail(cfg *config.Config, data map[string]string) error {

	tmpl, err := template.ParseFiles("registration.html")
	if err != nil {
		return err
	}

	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, data)
	if err != nil {
		return err
	}

	msg := BuildMessage(cfg.Email.Sender, data["email"], "Регистрация успешна!", tpl.String())
	auth := smtp.PlainAuth("", cfg.Email.Sender, cfg.Email.Pass, cfg.Email.SmtpAuthAddress)
	err = smtp.SendMail(cfg.Email.SmtpServerAddress, auth, cfg.Email.Sender, []string{data["email"]}, []byte(msg))

	return err

}

func BuildMessage(Sender, To, Subject, Body string) string {

	msg := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	msg += fmt.Sprintf("From: %s\r\n", Sender)
	msg += fmt.Sprintf("To: %s\r\n", To)
	msg += fmt.Sprintf("Subject: %s\r\n", Subject)
	msg += fmt.Sprintf("\r\n%s\r\n", Body)

	return msg
}
