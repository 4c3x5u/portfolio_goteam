package log

import (
	"bytes"
	"log"
	"testing"

	"server/assert"
)

// TestAppLogger tests the AppLogger to assert that it uses the correct
// prefixes for each log Level.
func TestAppLogger(t *testing.T) {
	sut := NewAppLogger()

	for _, c := range []struct {
		name    string
		level   Level
		message string
		wantLog string
	}{
		{
			name:    "LevelInfo",
			level:   LevelInfo,
			message: "some information",
			wantLog: "--[INFO]-- some information\n",
		},
		{
			name:    "LevelError",
			level:   LevelError,
			message: "an error occured",
			wantLog: "--[ERROR]-- an error occured\n",
		},
		{
			name:    "LevelFatal",
			level:   LevelFatal,
			message: "fatal error occured",
			wantLog: "--[FATAL]-- fatal error occured\n",
		},
		{
			name:    "InvalidLevel",
			level:   12,
			message: "an error occured",
			wantLog: "--[WARNING]-- invalid log level for message: " +
				"\n  an error occured\n",
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			var res bytes.Buffer
			log.SetOutput(&res)

			sut.Log(c.level, c.message)

			// Take only the len(wantLog) amount of characters from the end of
			// the actual log to ignore the date that is printed before it.
			resEnd := res.String()[len(res.String())-len(c.wantLog):]
			if err := assert.Equal(c.wantLog, resEnd); err != nil {
				t.Error(err)
			}
		})
	}
}
