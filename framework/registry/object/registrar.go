package object

import (
	"fmt"
	"github.com/1t-data-kit/go-kit/base"
	"github.com/pkg/errors"
	"reflect"
)

type Registrar struct {
	data map[string]interface{}
}

func NewRegistrar() *Registrar {
	return &Registrar{
		data: map[string]interface{}{},
	}
}

func NewRegistrarOption(registrar *Registrar) base.Option {
	return base.NewOption(registrar)
}

func (r *Registrar) Register(objects ...interface{}) error {
	objectsLen := len(objects)
	if objectsLen%2 != 0 {
		return fmt.Errorf("arguments must be key,value,key,value... ")
	}
	for i := 0; i < objectsLen/2; i++ {
		kSerial := 2 * i
		vSerial := kSerial + 1
		k, v := objects[kSerial], objects[vSerial]
		kString, ok := k.(string)
		if !ok {
			return fmt.Errorf("arguments on [%d] must be string", vSerial)
		}
		r.data[kString] = v
	}

	return nil
}

func (r *Registrar) Bind(target ...interface{}) error {
	if len(target) == 0 {
		return nil
	}
	for _, _target := range target {
		if err := r.BindTarget(_target); err != nil {
			return err
		}
	}
	return nil
}

func (r *Registrar) BindTarget(target interface{}, chain ...reflect.Value) error {
	value := reflect.ValueOf(target)
	if value.IsNil() {
		return nil
	}

	value = reflect.Indirect(value)
	if chains(chain).Exists(value) {
		return nil
	}
	if value.Kind() != reflect.Struct {
		return notSupportError
	}

	return r.bindTarget(value, chain...)
}

func (r *Registrar) bindTarget(target reflect.Value, chain ...reflect.Value) error {
	for i := 0; i < target.NumField(); i++ {
		field := target.Field(i)
		fieldType := target.Type().Field(i)

		_tag := NewTag(fieldType.Name, fieldType.Tag.Get(tagMark))
		if _tag.Skip() {
			continue
		}
		_tagName := _tag.Name()

		if fieldType.Anonymous {
			if err := r.BindTarget(field, append(chain, target)...); err != nil {
				if err == notSupportError ||
					errors.Cause(err) == notSupportError {
					continue
				} else {
					return errors.Wrapf(err, "bind error on embed field %s.%s[%s]", target.Type(), fieldType.Name, field.Type())
				}
			}
			continue
		}

		fieldValue, exists := r.data[_tagName]
		if !exists {
			if _tag.AllowEmpty() {
				continue
			} else {
				return fmt.Errorf("bind error on field %s.%s[%s]: not found %s", target.Type(), fieldType.Name, field.Type(), _tagName)
			}
		}
		if err := r.BindTarget(fieldValue, append(chain, target)...); err != nil {
			if err == notSupportError ||
				errors.Cause(err) == notSupportError {
				continue
			} else {
				return errors.Wrapf(err, "bind error on field %s.%s[%s]", target.Type(), fieldType.Name, field.Type())
			}
		}
		if !field.CanSet() {
			return fmt.Errorf("bind error on field %s.%s[%s]: can not set", target.Type(), fieldType.Name, field.Type())
		}

		fieldValueType := reflect.TypeOf(fieldValue)
		if !fieldValueType.AssignableTo(fieldType.Type) {
			return fmt.Errorf("bind error on field %s.%s[%s]: field value[%s] not assignable", target.Type(), fieldType.Name, field.Type(), fieldValueType.Name())
		}
		field.Set(reflect.ValueOf(fieldValue))
	}

	return nil
}
