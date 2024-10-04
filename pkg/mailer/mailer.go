package mailer

import (
	"gopkg.in/gomail.v2"
	"systemMetric/config"
)

// New Mail dialer
func NewMailDialer(cfg *config.Config) *gomail.Dialer {
	return gomail.NewDialer(cfg.Smtp.Host, cfg.Smtp.Port, cfg.Smtp.User, cfg.Smtp.Password)
}
