#!/bin/bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
set -eu
cd $DIR/..
pkgs=$(go list ./... )
year=$(date '+%Y')
echo "******************"
echo "Validating Headers"
count=$(goheader dump -f jsonl -i . | jq -c "select(.copyright.year != $year)" | wc -l)
if [[ $count -gt 0 ]]; then
  echo "The following files have out of date headers"
  goheader dump -f jsonl -i . | jq -c "select(.copyright.year != $year) | .path"
  exit 1
fi
