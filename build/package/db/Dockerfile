FROM amazon/dynamodb-local:2.2.1
EXPOSE 8000
WORKDIR /home/dynamodblocal
CMD ["-jar", "DynamoDBLocal.jar", "-sharedDb", "-dbPath", "./data"]

