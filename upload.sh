#!/bin/bash

MINIO_ACCESS_KEY="CNFd1AxCR1fPoXkCTLXg"
MINIO_SECRET_KEY="9mZAiJXLKaP5ualhOx93gBaBdvgq31gwmznhmsy6"
MINIO_ENDPOINT="s3.myinfra.lol"
BUCKET_NAME="uploads"

if [ $# -ne 1 ]; then
    echo "Usage: $0 <file_path>"
    exit 1
fi

FILE_PATH="$1"
FILE_EXTENSION="${FILE_PATH##*.}"

FILENAME=$(openssl rand -hex 3 | tr -dc 'a-zA-Z0-9' | head -c 5)
NEW_FILE_NAME="${FILENAME}.${FILE_EXTENSION}"
S3_KEY="$NEW_FILE_NAME"

DATE=$(date -u +"%a, %d %b %Y %H:%M:%S GMT")
CONTENT_TYPE="application/octet-stream"
SIGNATURE=$(echo -en "PUT\n\n$CONTENT_TYPE\n$DATE\n/$BUCKET_NAME/$S3_KEY" | openssl sha1 -hmac "$MINIO_SECRET_KEY" -binary | base64)

response=$(curl -s -o /dev/null -w "%{http_code}" -X PUT -T "$FILE_PATH" \
    -H "Host: $MINIO_ENDPOINT" \
    -H "Date: $DATE" \
    -H "Content-Type: $CONTENT_TYPE" \
    -H "Authorization: AWS $MINIO_ACCESS_KEY:$SIGNATURE" \
"http://$MINIO_ENDPOINT/$BUCKET_NAME/$S3_KEY")

if [[ $response -ge 200 && $response -lt 300 ]]; then
    S3_URL="https://s3.tritan.gg/$BUCKET_NAME/$S3_KEY"
    echo "File uploaded to: $S3_URL"
else
    echo "Failed to upload u fucking taint. Code: $response"
    exit 1
fi
