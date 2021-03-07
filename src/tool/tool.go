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
	fmt.Println(t.Name, t.Version)
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
	action, exists := t.actions[name]
	if !exists {
		t.ToolLog.Warn("Attempt to get inexsistent action: " + name)
		return nil, false
	}
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
	t.ToolLog.Debug("Running action: " + name)
	action, exists := t.actions[name]
	if !exists {
		t.ToolLog.Error("unknown action: " + name)
		return errors.New("unknown action: " + name)
	}
	log := t.MasterLog.Sub("Action[" + name + "]")
	ctx := &Context{t, args, log}
	err := action.Func(ctx)
	if err != nil {
		log.Error(err.Error())
	}
	return err
}
