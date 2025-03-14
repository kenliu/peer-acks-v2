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

# Peer Acks v2

A system for peer recognition and acknowledgments, now powered by Google Cloud Functions.

## Architecture

The application is now split into several serverless functions:

- `GetAcks`: Retrieves all acknowledgments
- `CreateAck`: Creates a new acknowledgment
- `GetMyAcks`: Retrieves acknowledgments for the current user
- `GetAcksReport`: Generates a report of all acknowledgments
- `SlackEvents`: Handles Slack events and challenges
- `SlackSlashCommand`: Processes Slack slash commands

## Prerequisites

- Google Cloud SDK
- Go 1.21 or later
- PostgreSQL database
- Slack app configuration

## Environment Variables

The following environment variables need to be configured in Google Cloud:

```
DATASOURCE=postgresql://user:password@host:port/dbname
SLACK_ACKS_CHANNELID=your_slack_channel_id
SLACK_SIGNING_SECRET=your_slack_signing_secret
```

## Deployment

1. Ensure you have the Google Cloud SDK installed and configured:
   ```bash
   gcloud auth login
   gcloud config set project YOUR_PROJECT_ID
   ```

2. Deploy all functions using Cloud Build:
   ```bash
   gcloud builds submit --config cloudbuild-functions.yaml
   ```

3. After deployment, you can find the function URLs in the Google Cloud Console or by running:
   ```bash
   gcloud functions describe FUNCTION_NAME --gen2 --region=us-central1
   ```

## Security

- All ack-related functions require authentication through Google Cloud IAP
- Slack endpoints are authenticated using Slack's signing secret
- Database credentials are managed through environment variables

## Development

To run functions locally for development:

1. Install the Functions Framework:
   ```bash
   go install github.com/GoogleCloudPlatform/functions-framework-go/cmd/functions-framework@latest
   ```

2. Run a function locally:
   ```bash
   functions-framework --target=FUNCTION_NAME
   ```