#!/usr/bin/env bash

aws s3 cp swagger.yml s3://${SWAGGER_BUCKET_NAME}/swagger.yml

sam package \
 --template-file ./template.yml \
 --s3-bucket "${SWAGGER_BUCKET_NAME}" \
 --output-template-file ./packaged.yml && \
sam deploy \
 --template-file ./packaged.yml \
 --stack-name "$STACK_NAME" \
 --capabilities CAPABILITY_IAM \
 --no-fail-on-empty-changeset \
 --parameter-overrides SwaggerBucketName=${SWAGGER_BUCKET_NAME}
