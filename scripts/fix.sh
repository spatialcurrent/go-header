#!/bin/bash
set -euo pipefail
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $DIR/..
echo "******************"
echo "Fixing Headers"
goheader fix --fix-year 2019 --exit-code-on-changes 1 --verbose
