package mailer

// Mailer is an object used to send email
type Mailer interface {
	SendStackTrace(trace []byte, endpoint string, message string, id string) error
	Send(msg *Message) error
}
