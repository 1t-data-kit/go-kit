package base

import (
	"fmt"
	"strings"
)

const (
	defaultCapacity = 5
)

type Errors struct {
	values []error
}

func NewErrors(opt ...Option) *Errors {
	errors := &Errors{}
	errors.values = make([]error, 0, errors.capacity(opt...))
	return errors
}

func (errors *Errors) sep(opt ...Option) string {
	sep := ";"
	if seps := Options(opt).Filter(func(item Option) bool {
		if _, ok := item.Value().(string); ok {
			return true
		}
		return false
	}); len(seps) > 0 {
		sep = seps[len(seps)-1].Value().(string)
	}
	return sep
}

func (errors *Errors) capacity(opt ...Option) int {
	_capacity := defaultCapacity
	if capacities := Options(opt).Filter(func(item Option) bool {
		if _, ok := item.Value().(int); ok {
			return true
		}
		return false
	}); len(capacities) > 0 {
		_capacity = capacities[len(capacities)-1].Value().(int)
	}
	return _capacity
}

func (errors *Errors) Error(opt ...Option) error {
	length := len(errors.values)
	if length == 0 {
		return nil
	}

	_errors := make([]string, 0, length)
	for _, err := range errors.values {
		_errors = append(_errors, err.Error())
	}

	return fmt.Errorf(strings.Join(_errors, errors.sep(opt...)))
}

func (errors *Errors) Append(err ...error) {
	errors.values = append(errors.values, err...)
}
