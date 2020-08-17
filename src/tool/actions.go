package tool

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// Action is a function of an Action
type Action func(*Tool, ArgumentList) error

// ActionMap is a string->Action map
type ActionMap map[string]Action

// Tool handles actions that are registered to it
type Tool struct {
	actions     ActionMap
	actionsDesc map[string]string
	name        string
	version     string
	desc        string
	silent      bool
}

// NewTool makes a blank tool
func NewTool(name, version, desc string) *Tool {
	tmp := new(Tool)
	tmp.ClearActions()
	tmp.SetInfo(name, version, desc)
	return tmp
}

// SetInfo set's the tool's name, version and description
func (t *Tool) SetInfo(name, version, desc string) {
	t.name = name
	t.version = version
	t.desc = desc
}

// GetVersion retunts the tool's version
func (t *Tool) GetVersion() string {
	return t.version
}

// ClearActions clears the action map of the Tool
func (t *Tool) ClearActions() {
	t.actions = make(ActionMap)
	t.actionsDesc = make(map[string]string)
}

// RunAction parses the arguments provided and runs any action detected
func (t *Tool) RunAction(args []string) error {
	argList := GetArgv(args)

	if len(argList.PosArgs) < 1 {
		return errors.New("not enough positional arguments")
	}

	_, t.silent = argList.NamedArgs["silent"]
	t.Welcome()

	actionName := strings.ToLower(argList.PosArgs[0].Value)
	action, hasAction := t.actions[actionName]

	if !hasAction {
		return fmt.Errorf("could not find action: %s", actionName)
	}

	actionReturn := action(t, argList)

	if actionReturn != nil {
		return fmt.Errorf("%s: %v", actionName, actionReturn)
	}

	return nil
}

// AddAction adds an action to the tool
func (t *Tool) AddAction(name string, action Action, desc string) {
	actionName := strings.ToLower(name)
	t.actions[actionName] = action
	t.actionsDesc[actionName] = desc
}

// DelAction deletes an action from the tool
func (t *Tool) DelAction(name string) {
	delete(t.actions, strings.ToLower(name))
	delete(t.actionsDesc, strings.ToLower(name))
}

// VerboseLn functions like a Println but only prints when the Tool isn't silent
func (t *Tool) VerboseLn(a ...interface{}) {
	if t.silent {
		return
	}
	fmt.Println(a...)
}

// VerboseF functions like a Printf but only prints when the Tool isn't silent
func (t *Tool) VerboseF(format string, a ...interface{}) {
	if t.silent {
		return
	}
	fmt.Printf(format, a...)
}

// ErrorExit closes the application with the given error and a non-zero exit code
func (t *Tool) ErrorExit(a ...interface{}) {
	fmt.Println(a...)
	os.Exit(1)
}

// ErrorExitF functions like ErrorExit but allows formatting
func (t *Tool) ErrorExitF(format string, a ...interface{}) {
	fmt.Printf(format, a...)
	os.Exit(1)
}

// Welcome prints a welcome message
func (t *Tool) Welcome() {
	t.VerboseF("%s v%s\n%s\n\n", t.name, t.version, t.desc)
}

// Init adds an additional help action
func (t *Tool) Init() {
	t.AddAction("help", func(tool *Tool, a ArgumentList) error {
		fmt.Println("Here are the availible actions:")
		for actionName := range t.actions {
			fmt.Println("\t", actionName, "-", t.actionsDesc[actionName])
		}
		return nil
	}, "Shows help")
}
