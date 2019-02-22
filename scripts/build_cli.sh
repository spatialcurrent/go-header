#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

LDFLAGS="-X main.gitBranch=$(git branch | grep \* | cut -d ' ' -f2) -X main.gitCommit=$(git rev-list -1 HEAD)"

DEST=$(realpath ${1:-$DIR/../bin})
mkdir -p $DEST

echo "******************"
echo "Building programs for go-header"
cd $DEST
for GOOS in darwin linux windows; do
  GOOS=${GOOS} GOARCH=amd64 go build -o "goheader_${GOOS}_amd64" -ldflags "$LDFLAGS" github.com/spatialcurrent/go-header/cmd/goheader
done
if [[ "$?" != 0 ]] ; then
    echo "Error building program for goheader (goheader_${GOOS}_amd64)"
    exit 1
fi
echo "Executables built at $DEST"
