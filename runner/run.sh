#!/bin/sh

set -e

NMAP_PATH="${NMAP_PATH:-$(which nmap)}"

if [ -z "$NMAP_PATH" ]
then
	echo "error: nmap not found" >&2
	exit 1
fi

# shellcheck disable=SC3028
RUNNER_ID="${RUNNER_ID:-$HOSTNAME}"

if [ -z "$RUNNER_ID" ]
then
	echo "error: neither RUNNER_ID nor HOSTNAME are set" >&2
	exit 1
fi

rm -f "$RUNNER_ID.xml"
"$NMAP_PATH" -oX "$RUNNER_ID.xml" "$@"

SCANALYZER_ADDRESS="${SCANALYZER_ADDRESS:-scanalyzer:4280}"

curl -s "http://$SCANALYZER_ADDRESS/v1/scan" -H 'Content-Type: application/xml' -d "@$RUNNER_ID.xml"

mkdir -p "./archive"
mv "$RUNNER_ID.xml" "archive/$(date "+%Y%m%d-%H%M%S").$RUNNER_ID.xml"
