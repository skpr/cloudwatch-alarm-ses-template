package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

func main() {
	lambda.Start(HandleEvents)
}

// HandleEvents will handle the SNS event.
func HandleEvents(ctx context.Context, snsEvent events.SNSEvent) error {

	// Authentication
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to setup client: %d", err)
	}
	sesclient := ses.NewFromConfig(cfg)

	// Get SES configuration from lambda env variables
	EmailTemplate := os.Getenv("SES_TEMPLATE")
	SourceEmail := os.Getenv("EMAIL_FROM_ADDRESS")
	DestinationToEmail := os.Getenv("EMAIL_TO_ADDRESSES")
	DestinationCcEmail := os.Getenv("EMAIL_CC_ADDRESSES")
	ReplyToEmail := os.Getenv("EMAIL_REPLY_TO_ADRESSES")

	emailDetails := &ses.SendTemplatedEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{DestinationToEmail},
		},
		Source:           aws.String(SourceEmail),
		Template:         aws.String(EmailTemplate),
		TemplateData:     aws.String(snsEvent.Records[0].SNS.Message),
		ReplyToAddresses: []string{ReplyToEmail},
	}

	// CC should be optional.
	if DestinationCcEmail != "" {
		emailDetails.Destination.CcAddresses = []string{DestinationCcEmail}
	}

	// Send Email.
	resp, err := sesclient.SendTemplatedEmail(context.TODO(), emailDetails)

	// Report errors if needed.
	if err != nil {
		return fmt.Errorf("could not send SES templated email to '%s': %s\n", SourceEmail, err.Error())
	}

	fmt.Printf("Email notification %s successfully sent to %s from %s\n", *resp.MessageId, DestinationToEmail, SourceEmail)
	return nil
}
