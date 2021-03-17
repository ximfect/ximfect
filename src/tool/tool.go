package tool

import (
	"errors"
	"fmt"
)

// Tool represents an action-based CLI
type Tool struct {
	actions   ActionMap
	Name      string
	Desc      string
	Version   string
	MasterLog *MasterLog
	ToolLog   *Log
}

// NewTool is self-explainatory
func NewTool(name, desc, version string) *Tool {
	actions := make(ActionMap)
	master := NewMasterLog(1)
	toolLog := master.Sub("Tool")
	return &Tool{actions, name, desc, version, master, toolLog}
}

// PrintTitle prints the name and version of this Tool
func (t *Tool) PrintTitle() {
	fmt.Println(t.Name, Version, "build", Build)
}

// AddAction adds an action to this Tool
func (t *Tool) AddAction(name string, action *Action) {
	t.ToolLog.Debug("Adding action: " + name)
	_, exists := t.actions[name]
	if exists {
		t.ToolLog.Warn("Attempt to add existing action: " + name)
	} else {
		t.actions[name] = action
		t.ToolLog.Debug("Successfully added.")
	}
}

// RemAction removes an action from this Tool
func (t *Tool) RemAction(name string) {
	t.ToolLog.Debug("Removing action: " + name)
	_, exists := t.actions[name]
	if !exists {
		t.ToolLog.Warn("Attempt to remove inexsistent action: " + name)
	} else {
		delete(t.actions, name)
		t.ToolLog.Debug("Successfully removed.")
	}
}

// GetAction returns an action from this Tool
func (t *Tool) GetAction(name string) (*Action, bool) {
	// simple lookup into actions map
	action, exists := t.actions[name]
	// if the simple lookup didn't work, try to find by alias
	if !exists {
		// the action we finally found
		var found *Action
		// did we find anything?
		isFound := false
		// for every action
		for _, a := range t.actions {
			// for every alias of that action
			for _, n := range a.Aliases {
				// if the alias is the name we're looking for
				if n == name {
					// we found something
					isFound = true
					// and it's this action
					found = a
					break
				}
			}
		}
		// if we found something, return it and true
		if isFound {
			return found, true
		}
		// warn message
		t.ToolLog.Warn("Attempt to get inexsistent action: " + name)
		// return nil and false
		return nil, false
	}
	// return the result of the simple lookup and true
	return action, true
}

// GetActionList returns a list of actions in this Tool
func (t *Tool) GetActionList() []string {
	var out = []string{}
	for k := range t.actions {
		out = append(out, k)
	}
	return out
}

// RunAction runs an action in this Tool
func (t *Tool) RunAction(name string, args ArgumentList) error {
	// debug message
	t.ToolLog.Debug("Running action: " + name)
	// Get the action
	action, exists := t.GetAction(name)
	// if we couldn't get it
	if !exists {
		// error out
		t.ToolLog.Error("unknown action: " + name)
		return errors.New("unknown action: " + name)
	}
	// create a log for the action
	log := t.MasterLog.Sub("Action[" + name + "]")
	// create a context for the action
	ctx := &Context{t, args, log}
	// run the action and error out if necessary
	if err := action.Func(ctx); err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}
