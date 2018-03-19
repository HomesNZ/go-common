package email

import (
	"github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

// CharSet is a default AWS character set
const CharSet = "UTF-8"

type Mailer struct {
	Session *session.Session
}

// NewClient creates a new AWS session from conf and returns a Mailer object
func NewClient(conf *aws.Config) (*Mailer, error) {
	s, err := session.NewSession(conf)
	return &Mailer{
		Session: s,
	}, err
}

// Email represents a very basic email structure
type Email struct {
	To          []*string
	CCAddresses []*string
	From        string
	Subject     string
	Body        string
}

// Send sends a simple email via AWS SES
func (m *Mailer) Send(e Email) error {
	svc := ses.New(m.Session)

	input := &ses.SendEmailInput{
		Source: aws.String(e.From),
		Destination: &ses.Destination{
			// CcAddresses: e.CCAddresses,
			ToAddresses: e.To,
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(e.Body),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CharSet),
				Data:    aws.String(e.Subject),
			},
		},
	}
	_, err := svc.SendEmail(input)
	return err
}

// MustSend sends a simple email via AWS SES
// MustSend will log any errors that occur, but they will not be returned
func (m *Mailer) MustSend(e Email) {
	svc := ses.New(m.Session)

	input := &ses.SendEmailInput{
		Source: aws.String(e.From),
		Destination: &ses.Destination{
			// CcAddresses: e.CCAddresses,
			ToAddresses: e.To,
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(e.Body),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CharSet),
				Data:    aws.String(e.Subject),
			},
		},
	}
	_, err := svc.SendEmail(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				logrus.WithError(aerr).WithField("Error Code", ses.ErrCodeMessageRejected).Error()
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				logrus.WithError(aerr).WithField("Error Code", ses.ErrCodeMailFromDomainNotVerifiedException).Error()
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				logrus.WithError(aerr).WithField("Error Code", ses.ErrCodeConfigurationSetDoesNotExistException).Error()
			default:
				logrus.WithError(aerr).Error()
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			logrus.WithError(err).Error()
		}
	}
	return
}
