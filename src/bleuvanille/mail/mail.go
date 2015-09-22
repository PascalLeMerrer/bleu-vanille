package mail

import (
	"bleuvanille/config"
	"net/smtp"
)

// Send sends an email
// address is the email address of the recipient
// TODO use a package outside of the standard lib,
// that allows sending HTML emails
func Send(address string, message []byte) error {
	auth := smtp.PlainAuth(
		"",
		config.SMTPUser,
		config.SMTPPassword,
		config.SMTPServer)

	return smtp.SendMail(
		config.SMTPServer+":"+config.SMTPPort,
		auth,
		config.NoReplyAddress,
		[]string{address},
		message)
}
