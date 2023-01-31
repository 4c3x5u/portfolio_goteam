package log

// Fakeloger is a test fake for Logger.
type FakeLogger struct {
	InLevel   Level
	InMessage string
}

// Log implements the Logger interface on FakeLogger. It assigns the parameters
// passed into it to their corresponding In... fields on the fake instance.
func (f *FakeLogger) Log(level Level, message string) {
	f.InLevel, f.InMessage = level, message
}
