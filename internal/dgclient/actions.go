package dgclient

import (
	"reflect"
)

type ActionNames struct {
	Add    string
	Remove string
	Help   string
}

var Actions = ActionNames{
	Add:    "add",
	Remove: "remove",
	Help:   "help",
}

func (a ActionNames) All() []string {
	fields := reflect.VisibleFields(reflect.TypeOf(a))
	out := []string{}

	for _, field := range fields {
		out = append(out, reflect.ValueOf(a).FieldByName(field.Name).String())
	}

	return out
}
