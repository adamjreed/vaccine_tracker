# vaccine_tracker

This is a basic app that scrapes the Dekalb County, GA vaccine signup page and send a text message to a list of phone 
numbers if the page is not showing its "no appointments available" messaging.

## Example Config

Place the following in AWS Parameter Store under `/[stage]/vaccine_tracker/config`:

For your dev environment, you can use Twilio test credentials to simulate sending a text:
```
{
  "stage": "dev",
  "check": {
    "url": "https://wikipedia.org",
    "phrase": "The Free Encyclopedia"
  },
  "sms": {
    "from": "+15005550006",
    "to": ["+15558675321"]
  },
  "twilio": {
    "account_sid": "[your twilio test sid]",
    "auth_token": "[your twilio test auth token]"
  }
}
```

For your production environment, you should store your Twilio auth token securely using AWS Secrets Manager and pass
the secret in your config json:
```
{
  "stage": "prod",
  "check": {
    "url": "https://wikipedia.org",
    "phrase": "The Free Encyclopedia"
  },
  "sms": {
    "from": "[your twilio sms number]",
    "to": ["[a real phone number here]"]
  },
  "twilio": {
    "account_sid": "[your twilio sid]",
    "auth_token_secret": "prod/vaccine_tracker/twilio/auth_token"
  }
}
```

## Deploying

Dev environment:
`make deploy`

Prod environment:
`make deploy-prod`

## Manual Triggering

The prod environment will automatically run on a cron every five minutes. If you want to manually trigger either 
environment you can run the following AWS CLI command:
```
aws lambda invoke --function-name vaccine-tracker-[stage]-check /dev/null
```