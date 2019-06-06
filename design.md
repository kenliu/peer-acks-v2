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
* fix: sender_email should be not null in DB (need to migrate production DB)
* sort out schema/user creation/grants in local and mso cluster
* fix/test "my acks" query
* build basic report page (all acks in last 7 days)

# TODO
* get email from IAP (or environment)

# PRE-DEPLOYMENT
* Understand how K8S deployment works again
* Set up DB credentials in k8s app
* Set up cert for CRDB in k8s app
* How to get app configuration into version control


## Export and reload the current prod schema
* Figure out how to connect to current prod environment

# Future enhancements
* understand how tests work with golang
* find a way to query all acks
* extract middleware for authentication and logged out page
* show user email in app
* sentry integration https://github.com/gin-contrib/sentry
* add usage metrics
* filter empty acks
* make sure acks are escaped properly
* ack deletion
* create liveness endpoint for k8s

* inline editing of existing ack
