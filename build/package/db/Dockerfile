FROM amazon/dynamodb-local:latest
EXPOSE 8000
WORKDIR /home/dynamodblocal
CMD ["-jar", "DynamoDBLocal.jar", "-sharedDb", "-dbPath", "./data"]

