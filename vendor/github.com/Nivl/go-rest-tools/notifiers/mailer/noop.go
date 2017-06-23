package mailer

// Noop is an object used to send email through sendgrid
type Noop struct {
}

// SendStackTrace emails the current stacktrace to the default FROM
func (m *Noop) SendStackTrace(trace []byte, endpoint string, message string, id string) error {
	return nil
}

// Send is used to send an email
func (m *Noop) Send(msg *Message) error {
	return nil
}
