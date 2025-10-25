#!/usr/bin/env bash
set -euo pipefail

awslocal s3 mb s3://finance-control-receipts || true
awslocal sqs create-queue --queue-name finance-transactions-queue || true

# Cognito bootstrap
USER_POOL_ID=$(awslocal cognito-idp list-user-pools --max-results 10 | jq -r '.UserPools[] | select(.Name=="finance-control-local") | .Id' || true)
if [ -z "$USER_POOL_ID" ]; then
  USER_POOL_ID=$(awslocal cognito-idp create-user-pool --pool-name finance-control-local | jq -r '.UserPool.Id')
fi

CLIENT_ID=$(awslocal cognito-idp list-user-pool-clients --user-pool-id "$USER_POOL_ID" --max-results 10 | jq -r '.UserPoolClients[] | select(.ClientName=="finance-control-client") | .ClientId' || true)
if [ -z "$CLIENT_ID" ]; then
  awslocal cognito-idp create-user-pool-client \
    --user-pool-id "$USER_POOL_ID" \
    --client-name finance-control-client \
    --generate-secret || true
fi
