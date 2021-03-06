package tool

// Context represents an Action's execution context
type Context struct {
	Tool *Tool
	Args ArgumentList
	Log  *Log
}

// ActionFunc represents an Action's function
type ActionFunc func(*Context) error

// Action represents a Tool's action
type Action struct {
	Func  ActionFunc
	Desc  string
	Usage ArgumentList
}

// ActionMap is a string->Action map
type ActionMap map[string]*Action
