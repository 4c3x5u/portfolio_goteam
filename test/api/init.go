//go:build itest

package api

import "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

// reqsWriteUser are the requests sent to the user table to initialise it for
// test use.
var reqsWriteUser = []types.WriteRequest{
	{PutRequest: &types.PutRequest{Item: map[string]types.AttributeValue{
		"ID": &types.AttributeValueMemberS{Value: "team1Admin"},
		"Password": &types.AttributeValueMemberB{
			Value: []byte(
				"$2a$11$kZfdRfTOjhfmel7J4WRG3eltzH9lavxp5qyrpFnzc9MIYLhZNCqTO",
			),
		},
		"IsAdmin": &types.AttributeValueMemberBOOL{
			Value: true,
		},
		"TeamID": &types.AttributeValueMemberS{
			Value: "afeadc4a-68b0-4c33-9e83-4648d20ff26a",
		},
	}}},
	{PutRequest: &types.PutRequest{Item: map[string]types.AttributeValue{
		"ID": &types.AttributeValueMemberS{Value: "team1Member"},
		"Password": &types.AttributeValueMemberB{
			Value: []byte(
				"$2a$11$kZfdRfTOjhfmel7J4WRG3eltzH9lavxp5qyrpFnzc9MIYLhZNCqTO",
			),
		},
		"IsAdmin": &types.AttributeValueMemberBOOL{
			Value: false,
		},
		"TeamID": &types.AttributeValueMemberS{
			Value: "afeadc4a-68b0-4c33-9e83-4648d20ff26a",
		},
	}}},
	{PutRequest: &types.PutRequest{Item: map[string]types.AttributeValue{
		"ID": &types.AttributeValueMemberS{Value: "team2Admin"},
		"Password": &types.AttributeValueMemberB{
			Value: []byte(
				"$2a$11$kZfdRfTOjhfmel7J4WRG3eltzH9lavxp5qyrpFnzc9MIYLhZNCqTO",
			),
		},
		"IsAdmin": &types.AttributeValueMemberBOOL{
			Value: true,
		},
		"TeamID": &types.AttributeValueMemberS{
			Value: "66ca0ddf-5f62-4713-bcc9-36cb0954eb7b",
		},
	}}},
	{PutRequest: &types.PutRequest{Item: map[string]types.AttributeValue{
		"ID": &types.AttributeValueMemberS{Value: "team2Member"},
		"Password": &types.AttributeValueMemberB{
			Value: []byte(
				"$2a$11$kZfdRfTOjhfmel7J4WRG3eltzH9lavxp5qyrpFnzc9MIYLhZNCqTO",
			),
		},
		"IsAdmin": &types.AttributeValueMemberBOOL{
			Value: false,
		},
		"TeamID": &types.AttributeValueMemberS{
			Value: "66ca0ddf-5f62-4713-bcc9-36cb0954eb7b",
		},
	}}},
	{PutRequest: &types.PutRequest{Item: map[string]types.AttributeValue{
		"ID": &types.AttributeValueMemberS{Value: "team3Admin"},
		"Password": &types.AttributeValueMemberB{
			Value: []byte(
				"$2a$11$kZfdRfTOjhfmel7J4WRG3eltzH9lavxp5qyrpFnzc9MIYLhZNCqTO",
			),
		},
		"IsAdmin": &types.AttributeValueMemberBOOL{
			Value: true,
		},
		"TeamID": &types.AttributeValueMemberS{
			Value: "74c80ae5-64f3-4298-a8ff-48f8f920c7d4",
		},
	}}},
	{PutRequest: &types.PutRequest{Item: map[string]types.AttributeValue{
		"ID": &types.AttributeValueMemberS{Value: "team4Admin"},
		"Password": &types.AttributeValueMemberB{
			Value: []byte(
				"$2a$11$kZfdRfTOjhfmel7J4WRG3eltzH9lavxp5qyrpFnzc9MIYLhZNCqTO",
			),
		},
		"IsAdmin": &types.AttributeValueMemberBOOL{
			Value: true,
		},
		"TeamID": &types.AttributeValueMemberS{
			Value: "3c3ec4ea-a850-4fc5-aab0-24e9e7223bbc",
		},
	}}},
}

// reqsWriteTeam are the requests sent to the team table to initialise it for
// test use.
var reqsWriteTeam = []types.WriteRequest{
	{PutRequest: &types.PutRequest{Item: map[string]types.AttributeValue{
		"ID": &types.AttributeValueMemberS{
			Value: "afeadc4a-68b0-4c33-9e83-4648d20ff26a",
		},
		"Members": &types.AttributeValueMemberL{
			Value: []types.AttributeValue{
				&types.AttributeValueMemberS{Value: "team1Admin"},
				&types.AttributeValueMemberS{Value: "team1Member"},
			},
		},
		"Boards": &types.AttributeValueMemberL{
			Value: []types.AttributeValue{
				&types.AttributeValueMemberM{
					Value: map[string]types.AttributeValue{
						"ID": &types.AttributeValueMemberS{
							Value: "91536664-9749-4dbb-a470-6e52aa353ae4",
						},
						"Name": &types.AttributeValueMemberS{
							Value: "Team 1 Board 1",
						},
					},
				},
				&types.AttributeValueMemberM{
					Value: map[string]types.AttributeValue{
						"ID": &types.AttributeValueMemberS{
							Value: "fdb82637-f6a5-4d55-9dc3-9f60061e632f",
						},
						"Name": &types.AttributeValueMemberS{
							Value: "Team 1 Board 2",
						},
					},
				},
				&types.AttributeValueMemberM{
					Value: map[string]types.AttributeValue{
						"ID": &types.AttributeValueMemberS{
							Value: "1559a33c-54c5-42c8-8e5f-fe096f7760fa",
						},
						"Name": &types.AttributeValueMemberS{
							Value: "Team 1 Board 3",
						},
					},
				},
			},
		},
	}}},
	{PutRequest: &types.PutRequest{Item: map[string]types.AttributeValue{
		"ID": &types.AttributeValueMemberS{
			Value: "66ca0ddf-5f62-4713-bcc9-36cb0954eb7b",
		},
		"Members": &types.AttributeValueMemberL{
			Value: []types.AttributeValue{
				&types.AttributeValueMemberS{Value: "team2Admin"},
				&types.AttributeValueMemberS{Value: "team2Member"},
			},
		},
		"Boards": &types.AttributeValueMemberL{Value: []types.AttributeValue{}},
	}}},
	{PutRequest: &types.PutRequest{Item: map[string]types.AttributeValue{
		"ID": &types.AttributeValueMemberS{
			Value: "74c80ae5-64f3-4298-a8ff-48f8f920c7d4",
		},
		"Members": &types.AttributeValueMemberL{
			Value: []types.AttributeValue{
				&types.AttributeValueMemberS{Value: "team3Admin"},
			},
		},
		"Boards": &types.AttributeValueMemberL{
			Value: []types.AttributeValue{
				&types.AttributeValueMemberM{
					Value: map[string]types.AttributeValue{
						"ID": &types.AttributeValueMemberS{
							Value: "f0c5d521-ccb5-47cc-ba40-313ddb901165",
						},
						"Name": &types.AttributeValueMemberS{
							Value: "Team 3 Board 1",
						},
					},
				},
			},
		},
	}}},
	{PutRequest: &types.PutRequest{Item: map[string]types.AttributeValue{
		"ID": &types.AttributeValueMemberS{
			Value: "3c3ec4ea-a850-4fc5-aab0-24e9e7223bbc",
		},
		"Members": &types.AttributeValueMemberL{
			Value: []types.AttributeValue{
				&types.AttributeValueMemberS{Value: "team4Admin"},
			},
		},
		"Boards": &types.AttributeValueMemberL{Value: []types.AttributeValue{}},
	}}},
}
