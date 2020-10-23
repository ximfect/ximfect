package tool

import (
	"fmt"
	"os"
)

type Logger struct {
	Output     *os.File
	OutputPath string
	hasOutput  bool
	isVerbose  bool
	Name       string
}

// NewLogger creates a new Logger with a File and a name.
func NewLogger(path, name string) *Logger {
	tmp := &Logger{Name: name}
	tmp.Name = name

	opened, err := os.Open(path)
	if err != nil {
		created, err := os.Create(path)
		if err == nil {
			tmp.SetOutput(created, path)
		}
	} else {
		tmp.SetOutput(opened, path)
	}
	defer tmp.PanicHandler()
	return tmp
}

// NewLoggerNoFile creates a new Logger with a name and no File.
func NewLoggerNoFile(name string) *Logger {
	return &Logger{Name: name}
}

// SetVerbose sets the verbosity of the Logger.
// (doesn't affect file output)
func (l *Logger) SetVerbose(verbose bool) {
	l.isVerbose = verbose
}

// SetOutput sets the file output.
func (l *Logger) SetOutput(out *os.File, path string) {
	l.Output = out
	l.OutputPath = path
	l.hasOutput = true
	defer l.Output.Close()
}

func (l *Logger) writeOut(msg string) {
	if l.hasOutput {
		_, err := l.Output.Write([]byte(msg + "\n"))
		if err != nil {
			l.VerboseLn("There was a problem writing to file:", err,
				"(file output is now disabled)")
			l.hasOutput = false
		}
	}
}

func (l *Logger) PanicHandler() {
	if r := recover(); r != nil {
		l.PrintLn("A critical error has occurred and the program must exit:\n", r)
		if l.hasOutput {
			l.PrintLn("\nThe log can be found here:", l.OutputPath)
		}
		os.Exit(-1)
	}
}

// VerboseLn is a Println that only appears if isVerbose is true.
// (it is still written to the file output)
func (l *Logger) VerboseLn(a ...interface{}) {
	msg := fmt.Sprint(a...)
	if l.isVerbose {
		fmt.Println(msg)
	}
	l.writeOut(msg)
}

// VerboseF is a Printf that only appears if isVerbose is true.
// (it is still written to the file output)
func (l *Logger) VerboseF (format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	if l.isVerbose {
		fmt.Println(msg)
	}
	l.writeOut(msg)
}

// PrintLn is a Println that also write to file output.
func (l *Logger) PrintLn(a ...interface{}) {
	msg := fmt.Sprint(a...)
	fmt.Println(msg)
	l.writeOut(msg)
}

// PrintF is a Printf that also write to file output.
func (l *Logger) PrintF (format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	fmt.Println(msg)
	l.writeOut(msg)
}