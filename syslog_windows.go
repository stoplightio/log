// This file stubs out the sysLogger for windows, since
// syslog is not implemented for windows yet.
package log

import "io"

type sysLogger struct {
}

func NewSysLogger(conf Config) (Logger, error) {
	return &sysLogger{}, nil
}

func (l *sysLogger) SetSeverity(sev Severity) {
}

func (l *sysLogger) GetSeverity() Severity {
	//LOG_INFO
	return 6
}

func (l *sysLogger) Writer(sev Severity) io.Writer {
	return nil
}

func (l *sysLogger) FormatMessage(sev Severity, caller *CallerInfo, format string, args ...interface{}) string {
	return ""
}
