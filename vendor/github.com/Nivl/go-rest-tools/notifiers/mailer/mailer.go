package mailer

// Mailer is an object used to send email
type Mailer interface {
	// // SendStackTrace emails the current stacktrace to the default FROM
	SendStackTrace(trace []byte, message string, context map[string]string) error

	// Send is used to send an email
	Send(msg *Message) error
}
