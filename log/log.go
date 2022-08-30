package log




var L Logger

type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
}

func SetLogger(newLogger Logger) {
	L = newLogger
}

type EmptyLogger struct {}

func (l *EmptyLogger) Debug(args ...interface{})
func (l *EmptyLogger) Debugf(format string, args ...interface{})
func (l *EmptyLogger) Error(args ...interface{})
func (l *EmptyLogger) Errorf(format string, args ...interface{})

func init(){
	L = &EmptyLogger{}
}