package entities

import "database/sql/driver"

type TaskStatus string

const (
	Undefined TaskStatus = ""
	Todo                 = "todo"
	Doing                = "doing"
	Done                 = "done"
)

func (s *TaskStatus) Scan(value interface{}) error {
	*s = TaskStatus(value.([]byte))
	return nil
}

func (s TaskStatus) Value() (driver.Value, error) {
	return string(s), nil
}
