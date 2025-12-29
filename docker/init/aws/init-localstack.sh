#!/bin/bash

# Create bucket
awslocal s3 mb s3://uploads

# Create SQS queue
awslocal sqs create-queue --queue-name event-queue
 
echo "LocalStack initialization complete"
