#!/usr/bin/sh

aws dynamodb create-table --endpoint-url http://localhost:8000 --cli-input-json '{
  "TableName": "goteam-user",
  "AttributeDefinitions": [
    {
      "AttributeName": "Username",
      "AttributeType": "S"
    }
  ],
  "KeySchema": [
    {
      "AttributeName": "Username",
      "KeyType": "HASH"
    }
  ],
  "ProvisionedThroughput": {
    "ReadCapacityUnits": 1,
    "WriteCapacityUnits": 1
  }
}'

aws dynamodb create-table --endpoint-url http://localhost:8000 --cli-input-json '{
  "TableName": "goteam-team",
  "AttributeDefinitions": [
    {
      "AttributeName": "ID",
      "AttributeType": "S"
    }
  ],
  "KeySchema": [
    {
      "AttributeName": "ID",
      "KeyType": "HASH"
    }
  ],
  "ProvisionedThroughput": {
    "ReadCapacityUnits": 1,
    "WriteCapacityUnits": 1
  }
}'

aws dynamodb create-table --endpoint-url http://localhost:8000 --cli-input-json '{
  "TableName": "goteam-task",
  "AttributeDefinitions": [
    {
      "AttributeName": "ID",
      "AttributeType": "S"
    },
    {
      "AttributeName": "TeamID",
      "AttributeType": "S"
    },
    {
      "AttributeName": "BoardID",
      "AttributeType": "S"
    }
  ],
  "KeySchema": [
    {
      "AttributeName": "TeamID",
      "KeyType": "HASH"
    },
    {
      "AttributeName": "ID",
      "KeyType": "Range"
    }
  ],
  "ProvisionedThroughput": {
    "ReadCapacityUnits": 1,
    "WriteCapacityUnits": 1
  },
  "GlobalSecondaryIndexes": [
    {
      "IndexName": "BoardID-index",
      "KeySchema": [
        {
          "AttributeName": "BoardID",
          "KeyType": "HASH"
        },
        {
          "AttributeName": "ID",
          "KeyType": "RANGE"
        }
      ],
      "Projection": {
        "ProjectionType": "ALL"
      },
      "ProvisionedThroughput": {
        "ReadCapacityUnits": 1,
        "WriteCapacityUnits": 1
      }
    }
  ]
}'
