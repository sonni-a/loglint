package b

type Logger struct{}

func (l *Logger) Info(msg string, args ...interface{})           {}
func (l *Logger) Infow(msg string, keysAndValues ...interface{}) {}

func (l *Logger) Debug(msg string, args ...interface{})           {}
func (l *Logger) Debugw(msg string, keysAndValues ...interface{}) {}

func (l *Logger) Warn(msg string, args ...interface{})           {}
func (l *Logger) Warnw(msg string, keysAndValues ...interface{}) {}

func (l *Logger) Error(msg string, args ...interface{})           {}
func (l *Logger) Errorw(msg string, keysAndValues ...interface{}) {}

func (l *Logger) DPanic(msg string, args ...interface{})           {}
func (l *Logger) DPanicw(msg string, keysAndValues ...interface{}) {}

func (l *Logger) Panic(msg string, args ...interface{})           {}
func (l *Logger) Panicw(msg string, keysAndValues ...interface{}) {}

func (l *Logger) Fatal(msg string, args ...interface{})           {}
func (l *Logger) Fatalw(msg string, keysAndValues ...interface{}) {}

var zapLogger Logger

func TestZapLogger() {
	zapLogger.Info("Starting server")    // want `log message must start with lowercase letter`
	zapLogger.Error("error happened!!!") // want `log message contains special characters`
	zapLogger.Debug("api_key 123456")    // want `log message may contain sensitive data`
	zapLogger.Warn("что-то не так")      // want `log message must be in English`

	zapLogger.Infow("Starting server", "port", 8080)   // want `log message must start with lowercase letter`
	zapLogger.Errorw("error happened!!!", "code", 500) // want `log message contains special characters`

	zapLogger.DPanic("DPanic happened") // want `log message must start with lowercase letter`

	zapLogger.Panic("Panic!") // want `log message contains special characters`

	zapLogger.Fatal("Fatal error") // want `log message must start with lowercase letter`
}
