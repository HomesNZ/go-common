package email

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

// CharSet is a default AWS character set
const CharSet = "UTF-8"

// NewSession initialises an AWS session
func NewSession(sesRegion string) (*session.Session, error) {
	return session.NewSession(
		&aws.Config{
			Region:      aws.String(sesRegion),
			Credentials: credentials.NewEnvCredentials(),
		})
}

// Email represents a very basic email structure
type Email struct {
	To          []*string
	CCAddresses []*string
	From        string
	Subject     string
	Body        string
}

// Send sends a simple email via a smtp gateway using TLS
func (e *Email) Send(sess *session.Session) error {
	svc := ses.New(sess)

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
