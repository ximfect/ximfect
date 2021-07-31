package tool

import "strings"

// Context represents an Action's execution context
type Context struct {
	Tool *Tool
	Args ArgList
	Log  *Log
}

// ActionFunc represents an Action's function
type ActionFunc func(*Context, ArgList) error

// Action represents a Tool's action
type Action struct {
	Name  string
	Alias []string
	Desc  string
	Usage ArgList
	Func  ActionFunc
}

// NewAction is self-explainatory
func NewAction(n string, a []string, d string, u ArgList, f ActionFunc) *Action {
	tmp := new(Action)
	tmp.Name = strings.ToLower(n)
	tmp.Alias = a
	tmp.Desc = strings.TrimSpace(d)
	tmp.Usage = u
	tmp.Func = f
	return tmp
}

// Category is a tool action category
type Category struct {
	Desc    string
	Actions []*Action
}

// CategoryMap is a string->Category map
type CategoryMap map[string]Category
