// logger/logger.go
package logger

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strings"
    "sync"
    "time"
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
    level      Level
    output     io.Writer
    mu         sync.Mutex
    prefix     string
    timeFormat string
}

var defaultLogger = New(os.Stdout, INFO)

// New creates a new logger instance
func New(output io.Writer, level Level) *Logger {
    return &Logger{
        level:      level,
        output:     output,
        timeFormat: "2006-01-02 15:04:05",
    }
}

// FileLogger extends Logger with file operations
type FileLogger struct {
    *Logger
    filePath string
    file     *os.File
}

// NewFileLogger creates a new file logger
func NewFileLogger(path string, level Level) (*FileLogger, error) {
    // Ensure directory exists
    dir := filepath.Dir(path)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return nil, fmt.Errorf("failed to create log directory: %w", err)
    }

    // Open log file
    file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        return nil, fmt.Errorf("failed to open log file: %w", err)
    }

    return &FileLogger{
        Logger:   New(file, level),
        filePath: path,
        file:     file,
    }, nil
}

// Close closes the log file
func (fl *FileLogger) Close() error {
    if fl.file != nil {
        return fl.file.Close()
    }
    return nil
}

// Rotate rotates the log file
func (fl *FileLogger) Rotate() error {
    fl.mu.Lock()
    defer fl.mu.Unlock()

    // Close existing file
    if err := fl.file.Close(); err != nil {
        return fmt.Errorf("failed to close current log file: %w", err)
    }

    // Rotate file (add timestamp)
    timestamp := time.Now().Format("20060102150405")
    rotatedPath := fmt.Sprintf("%s.%s", fl.filePath, timestamp)
    if err := os.Rename(fl.filePath, rotatedPath); err != nil {
        return fmt.Errorf("failed to rotate log file: %w", err)
    }

    // Open new file
    file, err := os.OpenFile(fl.filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        return fmt.Errorf("failed to open new log file: %w", err)
    }

    fl.file = file
    fl.output = file
    return nil
}

// SetLevel changes the logging level
func (l *Logger) SetLevel(level Level) {
    l.mu.Lock()
    defer l.mu.Unlock()
    l.level = level
}

// SetPrefix sets the logger prefix
func (l *Logger) SetPrefix(prefix string) {
    l.mu.Lock()
    defer l.mu.Unlock()
    l.prefix = prefix
}

// SetTimeFormat sets the time format string
func (l *Logger) SetTimeFormat(format string) {
    l.mu.Lock()
    defer l.mu.Unlock()
    l.timeFormat = format
}

// SetOutput changes the output writer
func (l *Logger) SetOutput(output io.Writer) {
    l.mu.Lock()
    defer l.mu.Unlock()
    l.output = output
}

// log handles the actual logging
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

// logf handles formatted logging
func (l *Logger) logf(level Level, format string, v ...interface{}) {
    if level < l.level {
        return
    }
    l.log(level, fmt.Sprintf(format, v...))
}

// Debug logs debug messages
func (l *Logger) Debug(v ...interface{})                 { l.log(DEBUG, v...) }
func (l *Logger) Info(v ...interface{})                  { l.log(INFO, v...) }
func (l *Logger) Warn(v ...interface{})                  { l.log(WARN, v...) }
func (l *Logger) Error(v ...interface{})                 { l.log(ERROR, v...) }

// Debugf logs formatted debug messages
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

// Package-level configuration functions
func SetLevel(level Level)       { defaultLogger.SetLevel(level) }
func SetPrefix(prefix string)    { defaultLogger.SetPrefix(prefix) }
func SetOutput(output io.Writer) { defaultLogger.SetOutput(output) }
func SetTimeFormat(format string) { defaultLogger.SetTimeFormat(format) }