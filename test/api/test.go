//go:build itest

// Package test contains integration tests for the package internal/api. Each
// Go file except this one and main_test.go corresponds to a endpoint.
package api

import (
	"database/sql"
	"net/http"
)

// db is the database connection pool used during integration testing.
// It is set in main_test.go/TestMain.
// TODO: remove once fully migrated to DynamoDB
var db *sql.DB

const (
	// jwtKey is the JWT key used for signing and validating JWTs during
	// integration testing.
	jwtKey = "itest-jwt-key-0123456789qwerty"

	// JWTs to be used for testing purposes
	// TODO: remove
	jwtTeam1Admin = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZWFtMUF" +
		"kbWluIn0.hdiH2HHc8QFT9VbkpfXKubtV5-mMIT__tmMmYZHMVeA"
	jwtTeam1Member = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZWFtMU" +
		"1lbWJlciJ9.uJbS6vSFZzH1Nfbbto3ega9COg9dMuo63iYHmMYJ6bc"
	jwtTeam2Admin = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZWFtMkF" +
		"kbWluIn0.vjQ93bx9-LK7SZEmhuzISf-Mcf_-A2bZ6VbLn27THPY"
	jwtTeam2Member = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZWFtMk" +
		"1lbWJlciJ9.g4FxHf1WupHGzzlvvi-8my1shFhpNuaWZKfJSV-Edxs"
	jwtTeam3Admin = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZWFtM0F" +
		"kbWluIn0.QHFI2okGYug7GNwMwwpwYyTtZkx53I-R-uNjlodCwTU"
	jwtTeam4Admin = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZWFtNEF" +
		"kbWluIn0.BxguaMUSynY33m3CB3jsV-l4ZC0bTE8_8XJJ8VFNo3o"

	tkTeam1Admin = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJib2FyZElEcyI6WyI" +
		"5MTUzNjY2NC05NzQ5LTRkYmItYTQ3MC02ZTUyYWEzNTNhZTQiLCJmZGI4MjYzNy1mNm" +
		"E1LTRkNTUtOWRjMy05ZjYwMDYxZTYzMmYiLCIxNTU5YTMzYy01NGM1LTQyYzgtOGU1Z" +
		"i1mZTA5NmY3NzYwZmEiXSwiaXNBZG1pbiI6dHJ1ZSwidGVhbUlEIjoiYWZlYWRjNGEt" +
		"NjhiMC00YzMzLTllODMtNDY0OGQyMGZmMjZhIiwidXNlcm5hbWUiOiJ0ZWFtMUFkbWl" +
		"uIn0.bOJnHy1J6PkbZpDCfKN3FdlCO3uXwJYxgJTKI2srp6E"
	tkTeam1Member = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJib2FyZElEcyI6WyI" +
		"5MTUzNjY2NC05NzQ5LTRkYmItYTQ3MC02ZTUyYWEzNTNhZTQiLCJmZGI4MjYzNy1mNmE" +
		"1LTRkNTUtOWRjMy05ZjYwMDYxZTYzMmYiLCIxNTU5YTMzYy01NGM1LTQyYzgtOGU1Zi1" +
		"mZTA5NmY3NzYwZmEiXSwiaXNBZG1pbiI6ZmFsc2UsInRlYW1JRCI6ImFmZWFkYzRhLTY" +
		"4YjAtNGMzMy05ZTgzLTQ2NDhkMjBmZjI2YSIsInVzZXJuYW1lIjoidGVhbTFNZW1iZXI" +
		"ifQ.lMskCZoProRSWxKsYzE5K9E4BCKKbTLnMLkwlwuXS_I"
	tkTeam2Admin = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJib2FyZElEcyI6W10s" +
		"ImlzQWRtaW4iOnRydWUsInRlYW1JRCI6IjY2Y2EwZGRmLTVmNjItNDcxMy1iY2M5LTM2" +
		"Y2IwOTU0ZWI3YiIsInVzZXJuYW1lIjoidGVhbTJBZG1pbiJ9.Y4Ah4bQHfFg9yVLf70Z" +
		"kWc3kKCDSOBoLwBB9dXW8RT4"

	tkStateEmpty = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJib2FyZHMiOltdfQ.g" +
		"lA6vOsGSCUo4w2tsiAqyngpLelGOLA0cguBXnx-ans"
	tkStateTeam1 = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJib2FyZHMiOlt7ImNv" +
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
)

// addCookieAuth is used in various test cases to authenticate the request
// being sent to a handler.
func addCookieAuth(token string) func(*http.Request) {
	return func(req *http.Request) {
		req.AddCookie(&http.Cookie{Name: "auth-token", Value: token})
	}
}

// addCookieState adds the given token as the state cookie value to the request.
func addCookieState(token string) func(*http.Request) {
	return func(req *http.Request) {
		req.AddCookie(&http.Cookie{Name: "state-token", Value: token})
	}
}
