#!/bin/bash

aws s3api list-buckets

aws s3api create-bucket \
  --bucket ${SWAGGER_BUCKET_NAME} \
  --create-bucket-configuration LocationConstraint=ap-northeast-1
aws s3api put-bucket-lifecycle-configuration \
   --bucket ${SWAGGER_BUCKET_NAME} \
   --lifecycle-configuration file://scripts/lifecycle.json

echo "created ${SWAGGER_BUCKET_NAME}"
