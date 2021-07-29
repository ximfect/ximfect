package tool

import (
	"io"
	"os"
	"time"
	"ximfect/environ"

	"github.com/mattn/go-colorable"
)

// LevelNames holds the names of the different levels of Log entries.
var LevelNames = []string{
	"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}

// LevelColors holds ANSI color codes for the different Entry levels.
//
// Debug = Gray
// Info = Reset
// Warn = Yellow
// Error = Red
// Fatal = White (on Red)
var LevelColors = []string{
	"\033[0;90m", "\033[0m", "\033[0;33m", "\033[0;31m", "\033[41m"}

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
	start   time.Time
	file    io.WriteCloser
	fileOk  bool
	lvl     int
	stdout  io.Writer
	history []Entry
}

// NewMasterLog is self-explainatory, lvl represents the Entry filter:
// if an Entry's level is equal to or higher than the MasterLog's level, the
// Entry will be printed.
func NewMasterLog(lvl int) *MasterLog {
	ts := time.Now()
	ml := &MasterLog{
		history: []Entry{},
		start:   ts,
		lvl:     lvl,
		stdout:  colorable.NewColorableStdout(),
		fileOk:  false}
	start := ml.Sub("Master")
	fn := environ.DataPath("logs", "latest.log")
	start.Debug("Log file is " + fn)
	file, err := os.Create(fn)
	if err != nil {
		start.Warn("could not create log file: " + err.Error())
	} else {
		ml.file = file
		ml.fileOk = true
		for _, e := range ml.history {
			ml.file.Write([]byte(e.Format() + "\n"))
		}
	}
	return ml
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
	m.history = append(m.history, e)
	if e.lvl >= m.lvl {
		clr := LevelColors[e.lvl]
		full := clr + e.Format() + "\n\033[0m"
		m.stdout.Write([]byte(full))
	}
	if m.fileOk {
		m.file.Write([]byte(e.Format() + "\n"))
	}
}

// Cleanup cleans the log up and ensures it is unusable
func (m *MasterLog) Cleanup() {
	l := m.Sub("Master")
	l.Debug("Cleaning up.")

	// close the file, it's not needed anymore
	m.file.Close()
	m.fileOk = false

	// rename the latest log file to the correct name
	l.Debug("Renaming log file")
	err := os.Rename(
		environ.DataPath("logs", "latest.log"),
		environ.DataPath("logs", m.start.Format("2006-01-02_15-04")+".log"))
	if err != nil {
		l.Warn("log file could not be renamed: " + err.Error())
	}

	// ensure the thing isn't usable anymore
	m.lvl = 99
	// i can't be bothered to code in a proper check for this so an empty
	// array will have to do
	m.history = []Entry{}
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
