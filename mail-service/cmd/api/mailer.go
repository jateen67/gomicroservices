package main

import (
	"bytes"
	"html/template"
	"time"

	"github.com/vanng822/go-premailer/premailer"
	mail "github.com/xhit/go-simple-mail/v2"
)

// type used for setting up an instance of the mailer with the appropriate configuration
type Mail struct {
	Domain      string
	Host        string
	Port        int
	Username    string
	Password    string
	Encryption  string
	FromAddress string
	FromName    string
}

// describe the kind of content we will use for an individual email message
type Message struct {
	From        string
	FromName    string
	To          string
	Subject     string
	Attachments []string
	Data        any
	DataMap     map[string]any
}

// create function to send email
func (m *Mail) SendSMTPMessage(msg Message) error {
	// make sure theres a valid	'FromAddress' and 'FromName'
	// if not specified, well use the defaults
	if msg.From == "" {
		msg.From = m.FromAddress
	}

	if msg.FromName == "" {
		msg.FromName = m.FromName
	}

	// calling two templates, one for html mail and one for plain text mail, and passing data into those templates
	data := map[string]any{
		"message": msg.Data,
	}

	msg.DataMap = data

	// get two forms of the message: one in html and one in plain text

	// html version of the message
	formattedMessage, err := m.buildHTMLMessage(msg)
	if err != nil {
		return err
	}

	// plain message version of the message
	plainMessage, err := m.buildPlainTextMessage(msg)
	if err != nil {
		return err
	}

	// create a mail server
	server := mail.NewSMTPClient()
	// describe the server
	server.Host = m.Host
	server.Port = m.Port
	server.Username = m.Username
	server.Password = m.Password
	server.Encryption = m.getEncryption(m.Encryption)
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	// create client and connect to the host specificed above
	smtpClient, err := server.Connect()
	if err != nil {
		return err
	}

	// now we create an email message
	email := mail.NewMSG()
	email.SetFrom(msg.From).
		AddTo(msg.To).
		SetSubject(msg.Subject).
		SetBody(mail.TextPlain, plainMessage).
		AddAlternative(mail.TextHTML, formattedMessage)

	// check for attachments if there are any
	if len(msg.Attachments) > 0 {
		for _, x := range msg.Attachments {
			email.AddAttachment(x)
		}
	}

	// finally we send the message
	err = email.Send(smtpClient)
	if err != nil {
		return err
	}

	return nil
}

func (m *Mail) buildHTMLMessage(msg Message) (string, error) {
	// point to the specific html template
	templateToRender := "./templates/mail.html.gohtml"

	t, err := template.New("email-html").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	// declare variable that we will read into
	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	// we now know that there are no errors in our template file so now we will get our formatted message from 'tpl'

	// now in string format
	formattedMessage := tpl.String()
	// now we inline the css
	formattedMessage, err = m.inlineCSS(formattedMessage)
	if err != nil {
		return "", err
	}

	return formattedMessage, nil
}

func (m *Mail) buildPlainTextMessage(msg Message) (string, error) {
	// point to the specific plain text template
	templateToRender := "./templates/mail.plain.gohtml"

	t, err := template.New("email-plain").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	// declare variable that we will read into
	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	plainMessage := tpl.String()
	// since its plain text we dont need to inline the css like with the html version

	return plainMessage, nil
}

func (m *Mail) inlineCSS(s string) (string, error) {
	options := premailer.Options{
		RemoveClasses:     false,
		CssToAttributes:   false,
		KeepBangImportant: true,
	}

	prem, err := premailer.NewPremailerFromString(s, &options)
	if err != nil {
		return "", err
	}

	html, err := prem.Transform()
	if err != nil {
		return "", err
	}

	return html, nil
}

func (m *Mail) getEncryption(s string) mail.Encryption {
	switch s {
	case "tls":
		return mail.EncryptionSTARTTLS
	case "ssl":
		return mail.EncryptionSSLTLS
	case "none", "":
		return mail.EncryptionNone
	default:
		return mail.EncryptionSTARTTLS
	}
}
