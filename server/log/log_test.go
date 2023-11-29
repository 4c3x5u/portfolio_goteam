//go:build utest

package log

import (
	"bytes"
	"log"
	"testing"

	"github.com/kxplxn/goteam/server/assert"
)

func TestLog(t *testing.T) {
	sut := New()

	for _, c := range []struct {
		name    string
		logFunc func(string)
		msg     string
		wantLog string
	}{
		{
			name:    "Info",
			logFunc: sut.Info,
			msg:     "some information",
			wantLog: "--[INFO]-- some information\n",
		},
		{
			name:    "Error",
			logFunc: sut.Error,
			msg:     "an error occured",
			wantLog: "--[ERROR]-- an error occured\n",
		},
		{
			name:    "Fatal",
			logFunc: sut.Fatal,
			msg:     "fatal error occured",
			wantLog: "--[FATAL]-- fatal error occured\n",
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			var res bytes.Buffer
			log.SetOutput(&res)

			c.logFunc(c.msg)

			// Take only the len(wantLog) amount of characters from the end of
			// the actual log to ignore the date/time that is printed before it.
			resEnd := res.String()[len(res.String())-len(c.wantLog):]
			assert.Equal(t.Error, resEnd, c.wantLog)
		})
	}
}
