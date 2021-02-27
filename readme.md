peer-acks-v2 is a port of the original Cockroach Labs peer-acks https://github.com/andreimatei/peer-ack to Go, adding Slack integration

# About

Peer Acks V2 is an internal web app and Slack integration used at Cockroach Labs for employees to show appreciation
to each other in the form of "peer acks".

TODO screenshot

# Development

## Setting up local dev environment

### start up local CRDB cluster



### start up local server

go build

ngrok http -subdomain=peer-acks 8888

https://peer-acks.ngrok.io/slack/events


## Secrets configuration in local dev environment

TODO local env setup

## Secrets configuration on the server

TODO note conn string and DB client cert

## Building and pushing to k8s

### Create Docker image
```sh
BRANCH=$(git symbolic-ref --short HEAD)-$USER
SHA=$(git rev-parse --short HEAD)-$USER
gcloud --project cockroach-dev-inf builds submit --substitutions=BRANCH_NAME=$BRANCH,SHORT_SHA=$SHA
```

### Deploying container image to k8s
Update `peer-acks-v2:36` to the sha of the generated Docker image above.

```sh
# Upsert the configuration to GKE
kubectl apply -f peer-acks-v2.yaml
# Get a friendly view of the status.
kubectl describe deployment h2hello
kubectl describe service h2hello
```