package common

import "time"

type QueryPrepareBindings string

//  AllDate to string, All bool to int
func (qpb QueryPrepareBindings) PrepareBindings(bindings ...interface{}) []interface{} {
	var _bindings []interface{}
	for _, binding := range bindings {
		switch binding.(type) {
		case time.Time:
			data := binding.(time.Time).String()
			_bindings = append(_bindings, data)
		case bool:
			data := 0
			if binding == true {
				data = 1
			}
			_bindings = append(_bindings, data)
		default:
			_bindings = append(_bindings, binding)
		}
	}
	return bindings
}
