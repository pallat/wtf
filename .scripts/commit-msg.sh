#!/bin/sh
COMMIT_MSG=`head -n1 $1`
echo "checking commit msg"
ISSUE_FORMAT_REGEX="^([A-Z]{2,}-[0-9]{1,}[ ]{1,})|^(CHORE)|^(ADR)(.*)"
ERROR_MSG="Aborting commit. Your commit message is not following JIRA issue message format, For an example, JIRA-1234 any commit message"
if ! [[ "$COMMIT_MSG" =~ $ISSUE_FORMAT_REGEX ]]; then
    echo "$ERROR_MSG"
    exit 1
fi
