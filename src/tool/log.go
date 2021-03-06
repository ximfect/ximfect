package tool

import (
	"io"
	"time"

	"github.com/mattn/go-colorable"
)

// LevelNames holds the names of the different levels of Log entries.
var LevelNames = []string{
	"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}

// LevelColors holds ANSI color codes for the different Entry levels.
var LevelColors = []string{
	"\033[0m", "\033[0m", "\033[0;33m", "\033[0;31m", "\033[41m"}

// Entry represents an entry in a Log.
type Entry struct {
	ts  time.Time
	src string
	lvl int
	msg string
}

// Format returns a human-readable version of the Entry.
func (e Entry) Format() string {
	tsF := e.ts.Format("15:04:05")
	lvlF := LevelNames[e.lvl]

	return "[" + tsF + "] " + e.src + ": " + lvlF + ": " + e.msg
}

// MasterLog represents the main log used in an application, and should only be
// instantiated once.
type MasterLog struct {
	lvl    int
	stdout io.Writer
}

// NewMasterLog is self-explainatory, lvl represents the Entry filter:
// if an Entry's level is equal to or higher than the MasterLog's level, the
// Entry will be printed.
func NewMasterLog(lvl int) *MasterLog {
	return &MasterLog{lvl, colorable.NewColorableStdout()}
}

// SetLevel changes this MasterLog's filter level
func (m *MasterLog) SetLevel(lvl int) {
	m.lvl = lvl
}

// Sub returns a named Log
func (m *MasterLog) Sub(name string) *Log {
	return &Log{m, name}
}

// Emit prints an Entry, if it can pass the level filter (see NewMasterLog)
func (m *MasterLog) Emit(e Entry) {
	if e.lvl >= m.lvl {
		clr := LevelColors[e.lvl]
		full := clr + e.Format() + "\n\033[0m"
		m.stdout.Write([]byte(full))
	}
}

// Log represents a named sub-log of it's MasterLog
type Log struct {
	master *MasterLog
	name   string
}

func (l *Log) emit(ts time.Time, lvl int, msg string) {
	e := Entry{ts, l.name, lvl, msg}
	l.master.Emit(e)
}

// Sub returns a named sub-Log, with the same MasterLog as this one
func (l *Log) Sub(name string) *Log {
	return &Log{l.master, l.name + "." + name}
}

// Debug emits a debug-level Entry
func (l *Log) Debug(msg string) {
	l.emit(time.Now(), 0, msg)
}

// Info emits an info-level Entry
func (l *Log) Info(msg string) {
	l.emit(time.Now(), 1, msg)
}

// Warn emits a warn-level Entry
func (l *Log) Warn(msg string) {
	l.emit(time.Now(), 2, msg)
}

// Error emits a error-level Entry
func (l *Log) Error(msg string) {
	l.emit(time.Now(), 3, msg)
}

// Fatal emits a fatal-level Entry
func (l *Log) Fatal(msg string) {
	l.emit(time.Now(), 4, msg)
}
