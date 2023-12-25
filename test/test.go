//go:build itest

// Package test contains integration tests for the package internal/api. Each
// Go file except this one and main_test.go corresponds to an endpoint.
package test

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/kxplxn/goteam/pkg/cookie"
)

var (
	// JWTKey is the JWT key used for signing and validating JWTs during
	// integration testing.
	JWTKey = []byte("itest-jwt-key-0123456789qwerty")

	dur = 1 * time.Hour

	AuthEncoder   = cookie.NewAuthEncoder(JWTKey, dur)
	AuthDecoder   = cookie.NewAuthDecoder(JWTKey)
	StateEncoder  = cookie.NewStateEncoder(JWTKey, dur)
	StateDecoder  = cookie.NewStateDecoder(JWTKey)
	InviteDecoder = cookie.NewInviteDecoder(JWTKey)
)

const (
	T1AdminToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJib2FyZElEcyI6WyI" +
		"5MTUzNjY2NC05NzQ5LTRkYmItYTQ3MC02ZTUyYWEzNTNhZTQiLCJmZGI4MjYzNy1mNm" +
		"E1LTRkNTUtOWRjMy05ZjYwMDYxZTYzMmYiLCIxNTU5YTMzYy01NGM1LTQyYzgtOGU1Z" +
		"i1mZTA5NmY3NzYwZmEiXSwiaXNBZG1pbiI6dHJ1ZSwidGVhbUlEIjoiYWZlYWRjNGEt" +
		"NjhiMC00YzMzLTllODMtNDY0OGQyMGZmMjZhIiwidXNlcm5hbWUiOiJ0ZWFtMUFkbWl" +
		"uIn0.bOJnHy1J6PkbZpDCfKN3FdlCO3uXwJYxgJTKI2srp6E"
	T1MemberToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJib2FyZElEcyI6WyI" +
		"5MTUzNjY2NC05NzQ5LTRkYmItYTQ3MC02ZTUyYWEzNTNhZTQiLCJmZGI4MjYzNy1mNmE" +
		"1LTRkNTUtOWRjMy05ZjYwMDYxZTYzMmYiLCIxNTU5YTMzYy01NGM1LTQyYzgtOGU1Zi1" +
		"mZTA5NmY3NzYwZmEiXSwiaXNBZG1pbiI6ZmFsc2UsInRlYW1JRCI6ImFmZWFkYzRhLTY" +
		"4YjAtNGMzMy05ZTgzLTQ2NDhkMjBmZjI2YSIsInVzZXJuYW1lIjoidGVhbTFNZW1iZXI" +
		"ifQ.lMskCZoProRSWxKsYzE5K9E4BCKKbTLnMLkwlwuXS_I"
	T2AdminToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJib2FyZElEcyI6W10s" +
		"ImlzQWRtaW4iOnRydWUsInRlYW1JRCI6IjY2Y2EwZGRmLTVmNjItNDcxMy1iY2M5LTM2" +
		"Y2IwOTU0ZWI3YiIsInVzZXJuYW1lIjoidGVhbTJBZG1pbiJ9.Y4Ah4bQHfFg9yVLf70Z" +
		"kWc3kKCDSOBoLwBB9dXW8RT4"
	T3AdminToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc0FkbWluIjp0cnVl" +
		"LCJ0ZWFtSUQiOiI3NGM4MGFlNS02NGYzLTQyOTgtYThmZi00OGY4ZjkyMGM3ZDQiLCJ1" +
		"c2VybmFtZSI6InRlYW0zQWRtaW4ifQ.eqPoE2WmFwzNgCatB9IUzyMmSRn0_t-VjIA2d" +
		"WVN3vU"
	T4AdminToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc0FkbWluIjp0cnVl" +
		"LCJ0ZWFtSUQiOiIzYzNlYzRlYS1hODUwLTRmYzUtYWFiMC0yNGU5ZTcyMjNiYmMiLCJ1" +
		"c2VybmFtZSI6InRlYW00QWRtaW4ifQ.pmbrD7hCLsP5m_ePZHkEK-JbEQfPGbY1EOR24" +
		"C2PsUA"
	T4MemberToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc0FkbWluIjpmYWx" +
		"zZSwidGVhbUlEIjoiM2MzZWM0ZWEtYTg1MC00ZmM1LWFhYjAtMjRlOWU3MjIzYmJjIiw" +
		"idXNlcm5hbWUiOiJ0ZWFtNE1lbWJlciJ9.UNjSqhfTpB_IQ68Le_ApwAKlh4lBoG7gDt" +
		"N02CFKdLw"

	EmptyStateToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJib2FyZHMiOltdf" +
		"Q.glA6vOsGSCUo4w2tsiAqyngpLelGOLA0cguBXnx-ans"
	T1StateToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJib2FyZHMiOlt7ImNv" +
		"bHVtbnMiOlt7InRhc2tzIjpbeyJpZCI6ImM2ODRhNmEwLTQwNGQtNDZmYS05ZmE1LTE0" +
		"OTdmOTg3NDU2NyIsIm9yZGVyIjoxfV19LHsidGFza3MiOltdfSx7InRhc2tzIjpbXX0s" +
		"eyJ0YXNrcyI6W119XSwiaWQiOiI5MTUzNjY2NC05NzQ5LTRkYmItYTQ3MC02ZTUyYWEz" +
		"NTNhZTQifSx7ImNvbHVtbnMiOlt7InRhc2tzIjpbeyJpZCI6IjAxYTMxNjhkLTZkMmEt" +
		"NDZmYi1hZWQ5LTcwYzI2YTRkNzFlOSIsIm9yZGVyIjoxfV19LHsidGFza3MiOltdfSx7" +
		"InRhc2tzIjpbeyJpZCI6IjlkZDljOTgyLThkMWMtNDlhYy1hNDEyLTNiMDFiYTc0YjYz" +
		"NCIsIm9yZGVyIjoxfV19LHsidGFza3MiOltdfV0sImlkIjoiZmRiODI2MzctZjZhNS00" +
		"ZDU1LTlkYzMtOWY2MDA2MWU2MzJmIn0seyJjb2x1bW5zIjpbeyJ0YXNrcyI6W119LHsi" +
		"dGFza3MiOlt7ImlkIjoiOGZiMDQwYTItOTEwYy00N2FmLWE0YWItOWRlZTQ5ZjE2ZDFk" +
		"Iiwib3JkZXIiOjF9LHsiaWQiOiJhMmU1YjU1Zi0wMWNjLTRlYWMtODg4Mi1kNzZhY2I5" +
		"NGE1YjkiLCJvcmRlciI6Mn0seyJpZCI6ImUwMDIxYTU2LTZhMWUtNDAwNy1iNzczLTM5" +
		"NWQzOTkxZmI3ZSIsIm9yZGVyIjozfV19LHsidGFza3MiOltdfSx7InRhc2tzIjpbXX1d" +
		"LCJpZCI6IjE1NTlhMzNjLTU0YzUtNDJjOC04ZTVmLWZlMDk2Zjc3NjBmYSJ9XX0.m_T4" +
		"7kdeojqex8EpW9F_L-h_6wuSh9ridCm80doNtpc"
	T3StateToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJib2FyZHMiOlt7ImNv" +
		"bHVtbnMiOlt7InRhc2tzIjpbeyJpZCI6ImMxNDY0ODZkLTcyNjAtNGQzZC05ZGE1LTI1" +
		"NDVhNTEwOWNhMSIsIm9yZGVyIjoxfV19LHsidGFza3MiOlt7ImlkIjoiMzc5YTk0YWMt" +
		"M2FmNC00Y2EwLTg0NjktNWI0MTU2N2UxYmYxIiwib3JkZXIiOjF9XX0seyJ0YXNrcyI6" +
		"W3siaWQiOiJiNTliY2ZmMy05ODI5LTQ2MzAtYTIxZi04Mzk3N2RmYzQ2NjUiLCJvcmRl" +
		"ciI6MX1dfSx7InRhc2tzIjpbeyJpZCI6IjhmZDRkMmEzLTYyNDctNGRjYy1iYzZhLTUw" +
		"NzdkOGU1N2JlMSIsIm9yZGVyIjoxfV19XSwiaWQiOiJmMGM1ZDUyMS1jY2I1LTQ3Y2Mt" +
		"YmE0MC0zMTNkZGI5MDExNjUifV19.ut1Ri0Y2bRwQwEe71KmSM_1_4ML4guJbInfsneX" +
		"UNgQ"
	T4StateToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJib2FyZHMiOlt7ImNv" +
		"bHVtbnMiOlt7InRhc2tzIjpbeyJpZCI6IjVjY2Q3NTBkLTM3ODMtNDgzMi04OTFkLTAy" +
		"NWYyNGE0OTQ0ZiIsIm9yZGVyIjowfSx7ImlkIjoiNTVlMjc1ZTQtZGU4MC00MjQxLWI3" +
		"M2ItODhlNzg0ZDU1MjJiIiwib3JkZXIiOjF9XX1dLCJpZCI6ImNhNDdmYmVjLTI2OWUt" +
		"NGVmNC1hNzRhLWJjZmJjZDU5OWZkNSJ9XX0.0m01PbRPDDBgC-dnZjqQeFdb5_leJtjA" +
		"RjpWG9Px3vU"
)

// AddAuthCk is used in various test cases to authenticate the request being
// sent to a handler.
func AddAuthCk(token string) func(*http.Request) {
	return func(r *http.Request) {
		r.AddCookie(&http.Cookie{Name: "auth-token", Value: token})
	}
}

// AddStateCk adds the given token as the state cookie value to the request.
func AddStateCk(token string) func(*http.Request) {
	return func(r *http.Request) {
		r.AddCookie(&http.Cookie{Name: "state-token", Value: token})
	}
}

// CreateTable creates a DynamoDB table with the given name, and given sort and
// partition keys.
func CreateTable(
	svc *dynamodb.Client,
	name *string,
	partKey string,
	sortKey string,
	secINames ...string,
) (func() error, error) {
	attrDefs := []types.AttributeDefinition{
		{AttributeName: &partKey, AttributeType: types.ScalarAttributeTypeS},
	}

	var secIs []types.GlobalSecondaryIndex
	for _, iname := range secINames {
		attrDefs = append(attrDefs, types.AttributeDefinition{
			AttributeName: &iname, AttributeType: types.ScalarAttributeTypeS,
		})

		secIs = append(secIs, types.GlobalSecondaryIndex{
			IndexName: aws.String(iname + "_index"),
			KeySchema: []types.KeySchemaElement{
				{AttributeName: &iname, KeyType: types.KeyTypeHash},
				{AttributeName: &partKey, KeyType: types.KeyTypeRange},
			},
			Projection: &types.Projection{
				ProjectionType: types.ProjectionTypeAll,
			},
			ProvisionedThroughput: &types.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(25),
				WriteCapacityUnits: aws.Int64(25),
			},
		})
	}

	keySchema := []types.KeySchemaElement{
		{AttributeName: &partKey, KeyType: types.KeyTypeHash},
	}
	if sortKey != "" {
		attrDefs = append(attrDefs, types.AttributeDefinition{
			AttributeName: &sortKey, AttributeType: types.ScalarAttributeTypeS,
		})
		keySchema = append(keySchema, types.KeySchemaElement{
			AttributeName: &sortKey, KeyType: types.KeyTypeRange,
		})
	}

	_, err := svc.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
		TableName:            name,
		AttributeDefinitions: attrDefs,
		KeySchema:            keySchema,
		BillingMode:          types.BillingModeProvisioned,
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(25),
			WriteCapacityUnits: aws.Int64(25),
		},
		GlobalSecondaryIndexes: secIs,
	})
	if err != nil {
		return TearDownNone, err
	}

	// create user table teardown function
	return func() error {
		svc.DeleteTable(context.TODO(), &dynamodb.DeleteTableInput{
			TableName: name,
		})
		return nil
	}, nil
}

// EnsureTableActive checks whether the test table is created and its status is
// "ACTIVE" every 500 milliseconds until it is true.
func EnsureTableActive(svc *dynamodb.Client, tableName string) error {
	fmt.Println("ensuring all test tables are active")
	var teamTableActive bool
	for {
		if !teamTableActive {
			resp, err := svc.DescribeTable(
				context.TODO(), &dynamodb.DescribeTableInput{
					TableName: &tableName,
				},
			)
			if err != nil {
				return err
			}
			if resp.Table.TableStatus == types.TableStatusActive {
				teamTableActive = true
			}
		} else {
			break
		}

		time.Sleep(500 * time.Millisecond)
	}
	return nil
}

// TearDownNone is returned when there is nothing to tear down.
func TearDownNone() error { return nil }
