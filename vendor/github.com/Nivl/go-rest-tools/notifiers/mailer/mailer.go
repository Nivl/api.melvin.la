package mailer

// Mailer is an object used to send email
type Mailer interface {
	// // SendStackTrace emails the current stacktrace to the default FROM
	SendStackTrace(trace []byte, endpoint string, message string, id string) error

	// Send is used to send an email
	Send(msg *Message) error
}
