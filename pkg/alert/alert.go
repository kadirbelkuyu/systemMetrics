package alert

type EmailAlert struct {
	Dialer *gomail.Dialer
	From   string
	To     string
}

func NewEmailAlert(host string, port int, username, password, from, to string) *EmailAlert {
	dialer := gomail.NewDialer(host, port, username, password)
	return &EmailAlert{
		Dialer: dialer,
		From:   from,
		To:     to,
	}
}

func (e *EmailAlert) SendAlert(subject, body string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", e.From)
	msg.SetHeader("To", e.To)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/plain", body)

	return e.Dialer.DialAndSend(msg)
}
