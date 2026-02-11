# Cloud Functions

## SlackEvents
Handles Slack events and challenges
- POST requests from Slack event subscriptions
- Responds to URL verification challenges

## SlackSlashCommand  
Processes Slack slash commands
- POST requests from `/ack` slash command
- Validates request signatures
- Creates acks and posts to Slack channel

# Verify
* fix: sender_email should be not null in DB (need to migrate production DB)
* sort out schema/user creation/grants in local and mso cluster
* fix/test "my acks" query
* build basic report page (all acks in last 7 days)

# TODO
* Consider adding health check endpoint for monitoring

# PRE-DEPLOYMENT
* Set up environment variables in Cloud Functions
* Set up DB credentials as environment variables
* Configure Slack app endpoints to point to deployed Cloud Functions


## Export and reload the current prod schema
* Figure out how to connect to current prod environment

# Future enhancements
* understand how tests work with golang
* find a way to query all acks
* show user email in app
* sentry integration https://github.com/gin-contrib/sentry
* add usage metrics
* filter empty acks
* make sure acks are escaped properly
* ack deletion
* inline editing of existing ack
