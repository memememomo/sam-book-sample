#!/usr/bin/env bash

# For this issue https://github.com/docker/compose/issues/4076 (Probably)
sleep 1

echo "--------------"
echo "Running go vet & golint."
echo

RESULT=0
LINT_RET=$(go list ./... | grep -v /vendor/ | xargs golint -set_exit_status 2>&1 | grep -v 'should have comment' | grep -v 'failing' | grep -v '^\s*$')
if [ "$LINT_RET" != "" ];then
  echo "golint error"
  echo $LINT_RET
  RESULT=$(expr $RESULT + 1)
else
  echo "golint OK."
fi

VET_RET=$(go list ./... | grep -v /vendor/ | xargs go vet 2>&1)
if [ "$VET_RET" != "" ];then
  echo "govet error"
  echo $VET_RET
  RESULT=$(expr $RESULT + 1)
else
  echo "govet OK."
fi


echo "--------------"
echo "Running gofmt."
echo
FMT_RET=$(gofmt -s -l $(find . -type f -name '*.go' -not -path './vendor/*' -not -path './scripts/*'))
if [ "$FMT_RET" != "" ]; then
  echo "Invalid files: $FMT_RET"
  RESULT=$(expr $RESULT + 1)
else
  echo "gofmt OK."
fi

if [ $RESULT -ne 0 ]; then
  echo "**************"
  echo "go vet or golint or gofmt FAILED."
  echo "**************"
  exit 1
fi
