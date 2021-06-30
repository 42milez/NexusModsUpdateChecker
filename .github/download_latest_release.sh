#!/bin/bash

readonly REPO='42milez/NexusModsUpdateChecker'
readonly ASSET_NAME='watcher-linux64.tar.gz'

readonly TAG_ID=$(curl --show-error --silent \
  -H "Authorization: token ${GITHUB_TOKEN}" \
  -X 'GET' \
  "https://api.github.com/repos/${REPO}/releases/latest" \
| jq -r ".id")

test "${TAG_ID}" = '' && exit 1

readonly DOWNLOAD_URL=$(curl --location --show-error --silent \
  -H "Accept: application/vnd.github.v3.raw" \
  -H "Authorization: token ${GITHUB_TOKEN}" \
  -X 'GET' \
  "https://api.github.com/repos/${REPO}/releases/${TAG_ID}" \
| jq -r ".assets[] | select(.name==\"${ASSET_NAME}\").url")

test "${DOWNLOAD_URL}" = '' && exit 1

curl --location --show-error --silent \
  -H "Accept: application/octet-stream" \
  -H "Authorization: token ${GITHUB_TOKEN}" \
  -X 'GET' \
  -o "${ASSET_NAME}" \
  "${DOWNLOAD_URL}"
