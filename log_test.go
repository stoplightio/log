package log

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	. "gopkg.in/check.v1"
)

func TestLog(t *testing.T) { TestingT(t) }

type LogSuite struct{}

var _ = Suite(&LogSuite{})

func (s *LogSuite) SetUpTest(c *C) {
	// reset global loggers chain before every test
	loggers = []Logger{}
}

func (s *LogSuite) TestInit(c *C) {
	Init(newTestLogger("log1"), newTestLogger("log2"))
	c.Assert(loggers[0].String(), Equals, "testLogger(log1)")
	c.Assert(loggers[1].String(), Equals, "testLogger(log2)")
}

func (s *LogSuite) TestInitWithConfig(c *C) {
	InitWithConfig(LogConfig{ConsoleLoggerName, "info"}, LogConfig{SysLoggerName, "info"})
	c.Assert(loggers[0].String(), Equals, "consoleLogger(INFO)")
	c.Assert(loggers[1].String(), Equals, "sysLogger(INFO)")
}

func (s *LogSuite) TestNewLogger(c *C) {
	l, err := NewLogger(LogConfig{ConsoleLoggerName, "info"})
	c.Assert(err, IsNil)
	c.Assert(l.String(), Equals, "consoleLogger(INFO)")

	l, err = NewLogger(LogConfig{SysLoggerName, "warn"})
	c.Assert(err, IsNil)
	c.Assert(l.String(), Equals, "sysLogger(WARN)")

	l, err = NewLogger(LogConfig{UDPLoggerName, "error"})
	c.Assert(err, IsNil)
	c.Assert(l.String(), Equals, "udpLogger(ERROR)")

	l, err = NewLogger(LogConfig{"SuperDuperLogger", "info"})
	c.Assert(err, NotNil)
	c.Assert(l, IsNil)
}

func (s *LogSuite) TestInfo(c *C) {
	logger1 := newTestLogger("log1")
	logger2 := newTestLogger("log2")
	Init(logger1, logger2)

	Info("hello %s", "world")
	c.Assert(logger1.b.String(), Equals, "INFO hello world\n")
	c.Assert(logger2.b.String(), Equals, "INFO hello world\n")
}

func (s *LogSuite) TestWarn(c *C) {
	logger1 := newTestLogger("log1")
	logger2 := newTestLogger("log2")
	Init(logger1, logger2)

	Warning("hello %s", "world")
	c.Assert(logger1.b.String(), Equals, "WARN hello world\n")
	c.Assert(logger2.b.String(), Equals, "WARN hello world\n")
}

func (s *LogSuite) TestError(c *C) {
	logger1 := newTestLogger("log1")
	logger2 := newTestLogger("log2")
	Init(logger1, logger2)

	Error("hello %s", "world")
	c.Assert(logger1.b.String(), Equals, "ERROR hello world\n")
	c.Assert(logger2.b.String(), Equals, "ERROR hello world\n")
}

// testLogger helps in tests.
type testLogger struct {
	id string
	b  *bytes.Buffer
}

func newTestLogger(id string) *testLogger {
	return &testLogger{id, &bytes.Buffer{}}
}

func (l *testLogger) Info(format string, args ...interface{}) {
}

func (l *testLogger) Warning(format string, args ...interface{}) {
}

func (l *testLogger) Error(format string, args ...interface{}) {
}

func (l *testLogger) Writer(sev Severity) io.Writer {
	return l.b
}

func (l *testLogger) FormatMessage(sev Severity, caller *callerInfo, format string, args ...interface{}) string {
	return fmt.Sprintf("%s %s\n", sev, fmt.Sprintf(format, args...))
}

func (l testLogger) String() string {
	return fmt.Sprintf("testLogger(%s)", l.id)
}
