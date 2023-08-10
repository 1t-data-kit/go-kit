package object

import "reflect"

type chains []reflect.Value

func (list chains) Exists(value reflect.Value) bool {
	for _, chain := range list {
		if chain == value {
			return true
		}
	}
	return false
}
