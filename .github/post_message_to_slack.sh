#!/bin/bash

set -eu

read -r COLOR STATUS_MESSAGE < <(
case "${NEEDS_PREVIOUS_JOB_RESULT}" in
  "success") echo "#74c7b8" "passed";;
  "failure") echo "#ef4f4f" "failed";;
  "cancelled") echo "#f4d160" "was cancelled";;
  "skipped") echo "#dddddd" "was skipped";;
  *) exit 1
esac
)
readonly COLOR
readonly STATUS_MESSAGE

readonly COMMIT_HASH=$(echo "${GITHUB_SHA}" | cut -c 1-7)

readonly LINK=$(
if [[ "${GITHUB_REF}" =~ ^"refs/heads/" ]]; then
  readonly BRANCH=$(echo "${GITHUB_REF}" | cut -c 12-)
  echo "branch: <https://github.com/${GITHUB_REPOSITORY}/tree/${BRANCH}|${BRANCH}>"
else
  readonly TAG=$(echo "${GITHUB_REF}" | cut -c 11-)
  echo "tag: <https://github.com/${GITHUB_REPOSITORY}/releases/tag/${TAG}|${TAG}>"
fi
)

# Message Example:
# Workflow: CI (#284) of 42milez/ProtocolStack (branch: 2021-06-07-003) passed.
# - update README.md (dce47b7) by 42milez
readonly TEXT=$(
A=${GITHUB_WORKFLOW} B=${GITHUB_REPOSITORY}          C=${GITHUB_RUN_ID} D=${GITHUB_RUN_NUMBER} E=${LINK}
F=${STATUS_MESSAGE}  G=${GITHUB_HEAD_COMMIT_MESSAGE} H=${GITHUB_SHA}    I=${COMMIT_HASH}       J=${GITHUB_ACTOR}
cat <<EOF
Workflow: ${A} (<https://github.com/${B}/actions/runs/${C}|#${D}>) of <https://github.com/${B}|${B}> (${E}) ${F}.
- ${G} (<https://github.com/${B}/commit/${H}|${I}>) by <https://github.com/${B}/pulse|${J}>
EOF
)

readonly DATA=$(cat <<EOF
{
  "attachments": [
    {
      "blocks": [
        {
          "type": "section",
          "text": {
            "type": "mrkdwn",
            "text": "${TEXT}"
          }
        }
      ],
      "color": "${COLOR}"
    },
  ],
  "blocks": [],
  "channel": "${SLACK_CHANNEL}",
  "text": "",
  "username": "${SLACK_USERNAME}"
}
EOF
)

curl -s -X POST \
     -H "Content-type: application/json; charset=utf-8" \
     -H "Authorization: Bearer ${SLACK_BOT_USER_OAUTH_TOKEN}" \
     -d "${DATA}" \
     https://slack.com/api/chat.postMessage \
  > /dev/null
