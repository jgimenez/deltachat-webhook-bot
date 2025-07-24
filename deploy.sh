#!/bin/bash -e

# Default server
SERVER=${1:-root@gimix.tail7a312.ts.net}
IMAGE_NAME="gimix/deltachat-bot"
echo "Deploying to $SERVER..."

docker save $IMAGE_NAME | gzip | ssh $SERVER "gunzip | docker load"
ssh $SERVER "cd /opt/web/photocopies && docker compose up -d --remove-orphans"

echo "Deployment complete!" 