//go:build itest

// Package test contains integration tests for the package internal/api. Each
// Go file except this one and main_test.go corresponds to an endpoint.
package test

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// db is the DynamoDB client used in integration tests.
var db *dynamodb.Client

// JWTKey is the key used to sign/validate JWTs in integration tests.
var JWTKey = []byte("itest-jwt-key-0123456789qwerty")

// JWTs used in integration tests.
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
	T1InviteeToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc0FkbWluIjpmYW" +
		"xzZSwidGVhbUlEIjoiYWZlYWRjNGEtNjhiMC00YzMzLTllODMtNDY0OGQyMGZmMjZhIi" +
		"widXNlcm5hbWUiOiJ0ZWFtMUludml0ZWUifQ.yy28Kb1tCwMCkdHPSzgwH3kqgQONXAR" +
		"tP6LpFGlnDTY"
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
