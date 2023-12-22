package main

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"log"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"
)

//go:embed templates
var emailTemplateFS embed.FS

func (app *Config) SendMail(from, to, subject, tmpl string, data interface{}) {
	log.Println("link in sendMail: ", data)

	templateToRender := fmt.Sprintf("templates/%s.html.tmpl", tmpl)

	t, err := template.New("email-html").ParseFS(emailTemplateFS, templateToRender)
	if err != nil {
		log.Println(err)
		return
	}
	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", data); err != nil {
		log.Println(err)
		return
	}
	formattedMessage := tpl.String()

	templateToRender = fmt.Sprintf("templates/%s.plain.tmpl", tmpl)

	t, err = template.New("email-plain").ParseFS(emailTemplateFS, templateToRender)
	if err != nil {
		log.Println(err)
		return
	}

	if err = t.ExecuteTemplate(&tpl, "body", data); err != nil {
		log.Println(err)
		return
	}

	plainMessage := tpl.String()
	log.Println(plainMessage)
	//send mail
	server := mail.NewSMTPClient()
	server.Host = app.EnvVars.ServerHost
	server.Port = app.EnvVars.ServerPort
	server.Username = app.EnvVars.ServerUsername
	server.Password = app.EnvVars.ServerPassword
	server.Encryption = mail.EncryptionSSL
	server.KeepAlive = false
	server.ConnectTimeout = 50 * time.Second
	server.SendTimeout = 50 * time.Second
	log.Println("before email server connect ", app.EnvVars.ServerHost)
	smtpClient, err := server.Connect()
	if err != nil {
		log.Println("before email server connect error")
		log.Println(err)
		return
	}
	log.Println("before newMSG")
	email := mail.NewMSG()
	email.SetFrom(from).
		AddTo(to).
		SetSubject(subject)

	email.SetBody(mail.TextHTML, formattedMessage)
	email.AddAlternative(mail.TextPlain, plainMessage)
	log.Println("before send email")
	err = email.Send(smtpClient)
	if err != nil {
		log.Println(err)
		return
	}

	//app.infoLog.Println("send mail")

}
