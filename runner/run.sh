#!/bin/sh

set -e

NMAP_PATH="${NMAP_PATH:-$(which nmap)}"

if [ -z "$NMAP_PATH" ]
then
	echo "error: nmap not found" >&2
	exit 1
fi

# shellcheck disable=SC2039
RUNNER_ID="${RUNNER_ID:-$HOSTNAME}"

if [ -z "$RUNNER_ID" ]
then
	echo "error: neither RUNNER_ID nor HOSTNAME are set" >&2
	exit 1
fi

rm -f "$RUNNER_ID.xml"
"$NMAP_PATH" -oX "$RUNNER_ID.xml" "$@"

mkdir -p "./reports"
rn "$RUNNER_ID.xml" "reports/$(date "+%Y%m%d-%H%M%S").$RUNNER_ID.xml"
