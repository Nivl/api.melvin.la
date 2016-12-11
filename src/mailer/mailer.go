package mailer

import (
	"strings"

	sendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// Mailer is an object used to send email
type Mailer struct {
	APIKey      string
	DefaultFrom string
	DefaultTo   string
}

// SendStackTrace emails the current stacktrace to the default FROM
func (m *Mailer) SendStackTrace(trace []byte, endpoint string, message string, id string) error {
	msg := NewMessage("289fb147-a8db-4d56-9cf7-07912537fc9b")
	stacktrace := string(trace[:])

	msg.Body = strings.Replace(stacktrace, "\n", "<br>", -1)
	msg.SetVar("endpoint", endpoint)
	msg.SetVar("message", message)
	msg.SetVar("requestID", id)
	return m.Send(msg)
}

// Send is used to send an email
func (m *Mailer) Send(msg *Message) error {
	from := mail.NewEmail("No Reply", msg.From)
	if msg.From == "" {
		from = mail.NewEmail("No Reply", m.DefaultFrom)
	}

	to := mail.NewEmail(msg.To, msg.To)
	if msg.From == "" {
		to = mail.NewEmail(m.DefaultTo, m.DefaultTo)
	}

	content := mail.NewContent("text/html", msg.Body)
	email := mail.NewV3MailInit(from, msg.Subject, to, content)
	email.SetTemplateID(msg.TemplateID)

	for k, v := range msg.Vars {
		email.Personalizations[0].SetSubstitution(k, v)
	}

	request := sendgrid.GetRequest(m.APIKey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	_, err := sendgrid.API(request)

	return err
}

// NewMailer creates and returns a new mailer
func NewMailer(apiKey string, defaultFrom string, defaultTo string) *Mailer {
	return &Mailer{
		APIKey:      apiKey,
		DefaultFrom: defaultFrom,
		DefaultTo:   defaultTo,
	}
}
