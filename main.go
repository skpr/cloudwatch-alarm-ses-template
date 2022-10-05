package main

import (
	"context"
	"fmt"
	"os"
	"strings"

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
	var (
		emailTemplate = os.Getenv("SES_TEMPLATE")
		fromAddress   = os.Getenv("EMAIL_FROM_ADDRESS")
		toAddresses   = os.Getenv("EMAIL_TO_ADDRESSES")
		ccAddresses   = os.Getenv("EMAIL_CC_ADDRESSES")
		replyToEmail  = os.Getenv("EMAIL_REPLY_TO_ADRESSES")
	)

	input := &ses.SendTemplatedEmailInput{
		Destination: &types.Destination{
			ToAddresses: strings.Split(toAddresses, ","),
		},
		Source:           aws.String(fromAddress),
		Template:         aws.String(emailTemplate),
		TemplateData:     aws.String(snsEvent.Records[0].SNS.Message),
		ReplyToAddresses: []string{replyToEmail},
	}

	// CC should be optional.
	if ccAddresses != "" {
		input.Destination.CcAddresses = strings.Split(ccAddresses, ",")
	}

	// Send Email.
	resp, err := sesclient.SendTemplatedEmail(ctx, input)
	if err != nil {
		return fmt.Errorf("could not send SES templated email to: %w\n", err)
	}

	fmt.Printf("Email notification %s successfully sent to %d addresses\n", *resp.MessageId, len(input.Destination.ToAddresses))

	return nil
}
