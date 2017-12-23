package notification

import (
	"fmt"

	"github.com/clickyab/services/notification/internal/mail"
	"github.com/sirupsen/logrus"
)

type (
	// NotifyType is the type of a notification
	NotifyType int
)

const (

	// MailType is the sms notification platform
	MailType NotifyType = iota
	// SMSType is the sms notification platform
	SMSType
)

// Duet is the name contact form for better result in some platform
// with name support, like email
type Duet struct {
	Name    string
	Contact string
}

// Send sends a notification by its notification type
func Send(platform NotifyType, subject string, msg string, From, To Duet) {
	switch platform {
	case MailType:
		mail.Send(
			subject,
			msg,
			mail.NewEmailNameAddress(From.Contact, From.Name),
			mail.NewEmailNameAddress(To.Contact, To.Name),
		)

	default:
		logrus.WithError(fmt.Errorf("not supported type %d", platform)).Panic("invalid notification type")
	}
}
