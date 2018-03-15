package email

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

// CharSet is a default AWS character set
const CharSet = "UTF-8"

var awsSession *session.Session

// Init creates a new session for SES requests
func Init() error {
	var err error
	awsSession, err = session.NewSession(&aws.Config{
		Region:      aws.String("us-west-2"),
		Credentials: credentials.NewEnvCredentials(),
	})
	return err
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
func (e *Email) Send() error {
	svc := ses.New(awsSession)

	input := &ses.SendEmailInput{
		Source: aws.String(e.From),
		Destination: &ses.Destination{
			// CcAddresses: e.CCAddresses,
			ToAddresses: e.To,
		},
		Message: &ses.Message{
			Body: &ses.Body{
				// TODO add html handling
				// Html: &ses.Content{
				// 	Charset: nil,
				// 	Data:    nil,
				// },
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
				fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
	}
	return nil
}
