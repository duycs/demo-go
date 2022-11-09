package entity

type Status string

const (
	Undefined Status = ""
	Todo       = "todo"
	Doing      = "doing"
	Done       = "done"
)