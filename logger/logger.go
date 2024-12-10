// logger/logger.go
package logger

import (
    "os"
    "fmt"
    "time"
    "strings"
    "io"
    "sync"
)

type Level int

const (
    DEBUG Level = iota
    INFO
    WARN
    ERROR
)

var levelNames = map[Level]string{
    DEBUG: "DEBUG",
    INFO:  "INFO",
    WARN:  "WARN",
    ERROR: "ERROR",
}

type Logger struct {
    level     Level
    output    io.Writer
    mu        sync.Mutex
    prefix    string
    timeFormat string
}

var defaultLogger = New(os.Stdout, INFO)

func New(output io.Writer, level Level) *Logger {
    return &Logger{
        level:      level,
        output:     output,
        timeFormat: "2006-01-02 15:04:05",
    }
}

func (l *Logger) SetLevel(level Level) {
    l.mu.Lock()
    defer l.mu.Unlock()
    l.level = level
}

func (l *Logger) SetPrefix(prefix string) {
    l.mu.Lock()
    defer l.mu.Unlock()
    l.prefix = prefix
}

func (l *Logger) SetTimeFormat(format string) {
    l.mu.Lock()
    defer l.mu.Unlock()
    l.timeFormat = format
}

func (l *Logger) log(level Level, v ...interface{}) {
    if level < l.level {
        return
    }

    l.mu.Lock()
    defer l.mu.Unlock()

    // Build the log message
    timeStr := time.Now().Format(l.timeFormat)
    levelStr := levelNames[level]
    prefix := ""
    if l.prefix != "" {
        prefix = "[" + l.prefix + "] "
    }

    msg := fmt.Sprint(v...)
    lines := strings.Split(msg, "\n")
    for _, line := range lines {
        if line == "" {
            continue
        }
        fmt.Fprintf(l.output, "%s [%s] %s%s\n", timeStr, levelStr, prefix, line)
    }
}

func (l *Logger) logf(level Level, format string, v ...interface{}) {
    if level < l.level {
        return
    }
    l.log(level, fmt.Sprintf(format, v...))
}

func (l *Logger) Debug(v ...interface{})                 { l.log(DEBUG, v...) }
func (l *Logger) Info(v ...interface{})                  { l.log(INFO, v...) }
func (l *Logger) Warn(v ...interface{})                  { l.log(WARN, v...) }
func (l *Logger) Error(v ...interface{})                 { l.log(ERROR, v...) }
func (l *Logger) Debugf(format string, v ...interface{}) { l.logf(DEBUG, format, v...) }
func (l *Logger) Infof(format string, v ...interface{})  { l.logf(INFO, format, v...) }
func (l *Logger) Warnf(format string, v ...interface{})  { l.logf(WARN, format, v...) }
func (l *Logger) Errorf(format string, v ...interface{}) { l.logf(ERROR, format, v...) }

// Package-level functions that use the default logger
func Debug(v ...interface{})                 { defaultLogger.Debug(v...) }
func Info(v ...interface{})                  { defaultLogger.Info(v...) }
func Warn(v ...interface{})                  { defaultLogger.Warn(v...) }
func Error(v ...interface{})                 { defaultLogger.Error(v...) }
func Debugf(format string, v ...interface{}) { defaultLogger.Debugf(format, v...) }
func Infof(format string, v ...interface{})  { defaultLogger.Infof(format, v...) }
func Warnf(format string, v ...interface{})  { defaultLogger.Warnf(format, v...) }
func Errorf(format string, v ...interface{}) { defaultLogger.Errorf(format, v...) }

func SetLevel(level Level)      { defaultLogger.SetLevel(level) }
func SetPrefix(prefix string)   { defaultLogger.SetPrefix(prefix) }
func SetOutput(output io.Writer) { defaultLogger.output = output }
