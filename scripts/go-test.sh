#!/usr/bin/env bash

# For this issue https://github.com/docker/compose/issues/4076 (Probably)
sleep 1

PACKAGE=${1:-...}
ARGS=${2}

scripts/wait-dynamodb-local.sh

if [ $? -ne 0 ]; then
  echo "Failed to launch dynamodb-local"
  exit 1
fi

echo "--------------"
echo "Running go test."
echo
DISABLE_ENV_DECRYPT=1 richgo test sam-book-sample/$PACKAGE $ARGS -v
RETURNCD=$?
if [ $RETURNCD -ne 0 ]; then
  echo "--------------"
  echo "go test FAILED."
  echo "--------------"
  echo "RETURNCD: $RETURNCD"
  exit $RETURNCD
else
  echo "go test OK."
fi

echo "--------------"
echo "TEST DONE."
echo "--------------"
