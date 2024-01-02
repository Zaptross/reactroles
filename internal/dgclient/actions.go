package dgclient

import (
	"reflect"
)

type ActionNames struct {
	Add           string
	Remove        string
	Update        string
	CreateChannel string
	LinkChannel   string
	Help          string
}

const (
	ActionConfigure = "configure"
)

var Actions = ActionNames{
	Add:           "add",
	Update:        "update",
	Remove:        "remove",
	CreateChannel: "create-channel",
	LinkChannel:   "link-channel",
	Help:          "help",
}

func (a ActionNames) All() []string {
	fields := reflect.VisibleFields(reflect.TypeOf(a))
	out := []string{}

	for _, field := range fields {
		out = append(out, reflect.ValueOf(a).FieldByName(field.Name).String())
	}

	return out
}
