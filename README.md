# Lambda Templated Emails with Amazon SES

This Lambda will convert an SNS event into an SES SendTemplatedEmail API
request using the parameters in your SES template configured through
environment variables.

It will transfer the payload directly without converting or manipulating
data, and the rest is up to your email template to be handled.

## Configuration

This is obviously an AWS lambda application, so configure it accordingly
with the following environment variables set, and make sure to set the
SES:SendTemplatedEmail IAM permissions.

You can test it once configured with the example SNS payload before
associating it to an alarm with an SNS topic/subscription.

| Name                    | Purpose                                                      |
|-------------------------|--------------------------------------------------------------|
| SES_TEMPLATE            | The SES template to use for this SendTemplatedEmail request. |
| EMAIL_FROM_ADDRESS      | The email address in the `from` field.                       |
| EMAIL_TO_ADDRESSES      | The email address in the `to` field.                         |
| EMAIL_CC_ADDRESSES      | The optional email address in the `cc` field.                |
| EMAIL_REPLY_TO_ADRESSES | The email address in the `reply-to` field.                   |

## Licence

MIT