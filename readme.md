port of https://github.com/andreimatei/peer-ack to go

# Development

## Secrets configuration in local dev environment



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