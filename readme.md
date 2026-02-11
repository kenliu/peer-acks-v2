# Peer Acks v2

Peer Acks V2 is a port of the original Cockroach Labs peer-acks (https://github.com/andreimatei/peer-ack) to Go, adding Slack integration. It's an internal Slack integration and web app used at Cockroach Labs for employees to show appreciation to each other in the form of "peer acks".

This application is powered by Google Cloud Functions.

## Architecture

The application is deployed as serverless functions on Google Cloud Functions:

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

- Slack endpoints are authenticated using Slack's signing secret
- Database credentials are managed through environment variables

## Local Development

To run and test functions locally:

1. Install the Functions Framework:
   ```bash
   go install github.com/GoogleCloudPlatform/functions-framework-go/cmd/functions-framework@latest
   ```

2. Set up environment variables:
   ```bash
   export DATASOURCE=postgresql://user:password@host:port/dbname
   export SLACK_ACKS_CHANNELID=your_slack_channel_id
   export SLACK_SIGNING_SECRET=your_slack_signing_secret
   export SLACK_OAUTH_TOKEN=your_slack_oauth_token
   ```

3. Run a function locally:
   ```bash
   functions-framework --target=SlackEvents --port=8080
   # or
   functions-framework --target=SlackSlashCommand --port=8080
   ```

4. For local testing with Slack, use ngrok to create a public URL:
   ```bash
   ngrok http 8080
   ```
   Then configure your Slack app endpoints:
   - For Slack Events API: `https://your-id.ngrok.io` (for SlackEvents function)
   - For Slack Slash Commands: `https://your-id.ngrok.io` (for SlackSlashCommand function)

## Testing

Run the test suite:
```bash
./test.sh
# or
go test ./...
```