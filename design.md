# routes
## /acks

renders json `{acks:[<array of acks>]}`

## /myacks

renders a form to submit new acks

## /report

renders a list of acks for a given time frame 

# pages

inforad.html --> /radiator

# Verify
* fix/test "my acks" query
* build basic report page (all acks in last 7 days)

# TODO

# KNOWN ISSUES


# PRE-DEPLOYMENT
* How to get app configuration into version control?


# FUTURE ENHANCEMENTS
* add option for posting multiple acks at once
* filter empty acks
* allow people to pick names from drop down list
* understand how tests work with go
* add all acks (beyond 7 days) to report page
* extract middleware for authentication and build logged out page
* show user email in app
* sentry integration https://github.com/gin-contrib/sentry
* add usage metrics
* ack deletion
* inline editing of existing acks
* set up cloud build trigger
* Figure out how to connect to current prod environment

# Slack integration
