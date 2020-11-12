package tool

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
	"ximfect/environ"
)

// ActionFunc represents an Action's function
type ActionFunc func(*Tool, ArgumentList) error

// Action represents a Tool's action
type Action struct {
	Func  ActionFunc
	Desc  string
	Usage string
}

//type Action func(*Tool, ArgumentList) error

// ActionMap is a string->Action map
type ActionMap map[string]Action

// Tool handles actions that are registered to it
type Tool struct {
	actions ActionMap
	name    string
	version string
	desc    string
	Log     *Logger
	LogName string
}

// NewTool makes a blank tool
func NewTool(name, version, desc string) *Tool {
	tmp := new(Tool)
	tmp.ClearActions()
	tmp.SetInfo(name, version, desc)
	t := time.Now()
	logName := fmt.Sprintf("%d-%d-%d_%d-%d-%d",
		t.Day(), t.Month(), t.Year(), t.Hour(), t.Minute(), t.Second())
	logFilePath := environ.AppdataPath("logs", logName+".log")
	tmp.Log = NewLogger(logFilePath, logName)
	return tmp
}

// SetInfo set's the tool's name, version and description
func (t *Tool) SetInfo(name, version, desc string) {
	t.name = name
	t.version = version
	t.desc = desc
}

// GetVersion returns the tool's version
func (t *Tool) GetVersion() string {
	return t.version
}

// ClearActions clears the action map of the Tool
func (t *Tool) ClearActions() {
	t.actions = make(ActionMap)
}

// RunAction parses the arguments provided and runs any action detected
func (t *Tool) RunAction(args []string) error {
	argList := GetArgv(args)

	if len(argList.PosArgs) < 1 {
		return errors.New("not enough positional arguments")
	}

	_, verbose := argList.NamedArgs["verbose"]
	t.Log.SetVerbose(verbose)
	t.Welcome()

	actionName := strings.ToLower(argList.PosArgs[0].Value)
	t.VerboseLn("Finding action:", actionName)
	action, hasAction := t.actions[actionName]
	if !hasAction {
		return fmt.Errorf("could not find action: %s", actionName)
	}

	t.VerboseLn("Running action...")
	actionReturn := action.Func(t, argList)
	if actionReturn != nil {
		return fmt.Errorf("%s: %v", actionName, actionReturn)
	}
	t.VerboseLn("Finished!")
	return nil
}

// AddAction adds an action to the tool
func (t *Tool) AddAction(name string, action Action) {
	actionName := strings.ToLower(name)
	t.actions[actionName] = action
	t.VerboseLn("Added action:", name)
}

// AddActionQuick adds an action without the need to instantiate an Action
func (t *Tool) AddActionQuick(name, desc, usage string, act ActionFunc) {
	action := Action{act, desc, usage}
	t.AddAction(name, action)
}

// DelAction deletes an action from the tool
func (t *Tool) DelAction(name string) {
	delete(t.actions, strings.ToLower(name))
}

// VerboseLn maps directly to Logger.VerboseLn
func (t *Tool) VerboseLn(a ...interface{}) {
	t.Log.VerboseLn(a...)
}

// VerboseF maps directly to Logger.VerboseF
func (t *Tool) VerboseF(format string, a ...interface{}) {
	t.Log.VerboseF(format, a...)
}

// PrintLn maps directly to Logger.PrintLn
func (t *Tool) PrintLn(a ...interface{}) {
	t.Log.PrintLn(a...)
}

// PrintF maps directly to Logger.PrintF
func (t *Tool) PrintF(format string, a ...interface{}) {
	t.Log.PrintF(format, a...)
}

// ErrorExit closes the application with the given error and a non-zero exit code
func (t *Tool) ErrorExit(a ...interface{}) {
	t.Log.PrintLn(a...)
	os.Exit(1)
}

// ErrorExitF functions like ErrorExit but allows formatting
func (t *Tool) ErrorExitF(format string, a ...interface{}) {
	t.Log.PrintF(format, a...)
	os.Exit(1)
}

// Welcome prints a welcome message
func (t *Tool) Welcome() {
	t.VerboseF("%s v%s\n%s\n\n", t.name, t.version, t.desc)
}

// Init adds an additional help action
func (t *Tool) Init() {
	t.AddActionQuick("help", "Shows help", "", func(tool *Tool, a ArgumentList) error {
		t.PrintLn("Here are the available actions:\n")
		for name, action := range t.actions {
			if len(action.Usage) == 0 {
				t.PrintLn("\t* ", name, " -- ", action.Desc)
			} else {
				t.PrintLn("\t* ", name, " ", action.Usage, " -- ", action.Desc)
			}
		}
		return nil
	})

	//keys := make([]string, 0, len(t.actions))
	//for k := range t.actions {
	//	keys = append(keys, k)
	//}
	//sort.Strings(keys)
	//
	//sorted := make(ActionMap)
	//for _, k := range keys {
	//	sorted[k] = t.actions[k]
	//}
	//
	//t.actions = sorted
}
