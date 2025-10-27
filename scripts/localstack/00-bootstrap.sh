#!/usr/bin/env bash
set -euo pipefail

awslocal s3 mb s3://financial-control-receipts || true

DLQ_URL=$(awslocal sqs get-queue-url --queue-name financial-transactions-dlq 2>/dev/null | jq -r '.QueueUrl' || true)
if [ -z "$DLQ_URL" ]; then
  DLQ_URL=$(awslocal sqs create-queue --queue-name financial-transactions-dlq --attributes VisibilityTimeout=30 | jq -r '.QueueUrl')
fi
DLQ_ARN=$(awslocal sqs get-queue-attributes --queue-url "$DLQ_URL" --attribute-names QueueArn | jq -r '.Attributes.QueueArn')

QUEUE_URL=$(awslocal sqs get-queue-url --queue-name financial-transactions-queue 2>/dev/null | jq -r '.QueueUrl' || true)
if [ -z "$QUEUE_URL" ]; then
  QUEUE_URL=$(awslocal sqs create-queue --queue-name financial-transactions-queue --attributes VisibilityTimeout=60 | jq -r '.QueueUrl')
fi
TMP_ATTR_FILE=$(mktemp)
trap 'rm -f "$TMP_ATTR_FILE"' EXIT
cat >"$TMP_ATTR_FILE" <<JSON
{
  "RedrivePolicy": "{\"deadLetterTargetArn\":\"$DLQ_ARN\",\"maxReceiveCount\":\"5\"}"
}
JSON
awslocal sqs set-queue-attributes \
  --queue-url "$QUEUE_URL" \
  --attributes file://"$TMP_ATTR_FILE"
rm -f "$TMP_ATTR_FILE"
trap - EXIT

# Cognito bootstrap (best-effort: Community edition may not support Cognito)
if awslocal cognito-idp list-user-pools --max-results 1 >/dev/null 2>&1; then
  USER_POOL_ID=$(awslocal cognito-idp list-user-pools --max-results 10 | jq -r '.UserPools[] | select(.Name=="financial-control-local") | .Id' || true)
  if [ -z "$USER_POOL_ID" ]; then
    USER_POOL_ID=$(awslocal cognito-idp create-user-pool --pool-name financial-control-local | jq -r '.UserPool.Id')
  fi

  CLIENT_ID=$(awslocal cognito-idp list-user-pool-clients --user-pool-id "$USER_POOL_ID" --max-results 10 | jq -r '.UserPoolClients[] | select(.ClientName=="financial-control-client") | .ClientId' || true)
  if [ -z "$CLIENT_ID" ]; then
    awslocal cognito-idp create-user-pool-client \
      --user-pool-id "$USER_POOL_ID" \
      --client-name financial-control-client \
      --generate-secret || true
  fi
else
  echo "Cognito IDP indisponível no LocalStack (edição gratuita); bootstrap ignorado." >&2
fi
