#!/bin/bash

OPEN_API_FILE_PATH=$1
PROJECT_NAME=$2
VERSION=$3

jq --arg title "$PROJECT_NAME" \
   --arg version "$VERSION" \
   '.info.title = $title | .info.version = $version' \
   "$OPEN_API_FILE_PATH" > tmp.json && mv tmp.json "$OPEN_API_FILE_PATH"