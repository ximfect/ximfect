package tool

import (
	"errors"
	"fmt"
)

// Tool represents an action-based CLI
type Tool struct {
	Categories CategoryMap
	Name       string
	Desc       string
	Version    string
	MasterLog  *MasterLog
	ToolLog    *Log
}

// NewTool is self-explainatory
func NewTool(name, desc, version string) *Tool {
	cats := make(CategoryMap)
	master := NewMasterLog(1)
	toolLog := master.Sub("Tool")
	return &Tool{cats, name, desc, version, master, toolLog}
}

// PrintTitle prints the name and version of this Tool
func (t *Tool) PrintTitle() {
	fmt.Println(t.Name, Version, "build", Build)
}

// AddAction adds an action to this Tool
func (t *Tool) AddAction(cat string, action *Action) {
	t.ToolLog.Debug("Adding action: " + action.Name)

	_, catExists := t.Categories[cat]
	if !catExists {
		t.ToolLog.Debug("Creating category: " + cat)
		t.Categories[cat] = Category{NoDesc, []*Action{}}
	}

	categ := t.Categories[cat]
	categ.Actions = append(categ.Actions, action)
	t.Categories[cat] = categ
}

// AddManyActions adds multiple actions to one category at a time
func (t *Tool) AddManyActions(cat string, action ...*Action) {
	for _, a := range action {
		t.AddAction(cat, a)
	}
}

// RemAction removes an action from this Tool
func (t *Tool) RemAction(name string) {
	t.ToolLog.Debug("Removing action: " + name)

	for _, cat := range t.Categories {
		tmp := []*Action{}
		for _, act := range cat.Actions {
			if act.Name != name {
				tmp = append(tmp, act)
			}
		}
		cat.Actions = tmp
	}
}

// GetAction returns an action from this Tool
func (t *Tool) GetAction(name string) (*Action, bool) {
	for _, cat := range t.Categories {
		for _, act := range cat.Actions {
			if act.Name == name {
				return act, true
			}
			for _, als := range act.Alias {
				if als == name {
					return act, true
				}
			}
		}
	}

	return nil, false
}

// RunAction runs an action in this Tool
func (t *Tool) RunAction(name string, args ArgList) error {
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
	if err := action.Func(ctx, args); err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

// SetCategoryDesc changes a category's description
func (t *Tool) SetCategoryDesc(cat, desc string) {
	t.ToolLog.Debug("Updating category description: " + cat)

	categ := t.Categories[cat]
	categ.Desc = desc
	t.Categories[cat] = categ
}
