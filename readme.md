# Peer Acks v2

A port of the original [Cockroach Labs peer-acks](https://github.com/andreimatei/peer-ack) to Go, now powered by Google Cloud Functions.

## About

Peer Acks V2 is an internal Slack integration used at Cockroach Labs for employees to show appreciation to each other in the form of "peer acks".

## Architecture

The application consists of two serverless functions:

- `SlackEvents`: Handles Slack events and challenges
- `SlackSlashCommand`: Processes Slack slash commands

## Prerequisites

- Google Cloud SDK
- Go 1.21 or later
- CockroachDB database
- Slack app configuration

## Environment Variables

The following environment variables are required:

```
DATASOURCE=postgresql://user@crdb-host:26257/peer_acks?sslmode=verify-full&sslcert=client.crt&sslkey=client.key&sslrootcert=ca.crt
SLACK_ACKS_CHANNELID=your_slack_channel_id
SLACK_SIGNING_SECRET=your_slack_signing_secret
```

## Development

### Local Setup

1. Set up a database:
   ```bash
   # Start up local CRDB cluster and create database
   cockroach start-single-node --insecure
   cockroach sql --insecure -e "CREATE DATABASE peer_acks;"
   ```

2. Set required environment variables:
   ```bash
   export DATASOURCE="postgresql://root@localhost:26257/peer_acks?sslmode=disable"
   export SLACK_SIGNING_SECRET="your-slack-secret"
   export SLACK_ACKS_CHANNELID="your-channel-id"
   ```

3. Run a function locally:
   ```bash
   FUNCTION_TARGET=SlackEvents go run main.go
   ```

4. (Optional) To expose local server to Slack:
   ```bash
   ngrok http 8080
   ```
   Then update your Slack app's event subscription URL to the ngrok URL.

## Deployment

### Deploy to Google Cloud Functions

1. Ensure you have the Google Cloud SDK installed and configured:
   ```bash
   gcloud auth login
   gcloud config set project YOUR_PROJECT_ID
   ```

2. Deploy using Cloud Build:
   ```bash
   gcloud builds submit --config cloudbuild-functions.yaml \
     --substitutions=_DATASOURCE="your-db-connection-string",\
     _SLACK_SIGNING_SECRET="your-slack-secret",\
     _SLACK_ACKS_CHANNELID="your-channel-id"
   ```

3. After deployment, update your Slack app configuration with the new function URLs from the Google Cloud Console.

## Security

- Slack endpoints are authenticated using Slack's signing secret
- Database credentials are managed through environment variables
- Functions are deployed with HTTPS endpoints