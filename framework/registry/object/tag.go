package object

import (
	"strings"
)

const (
	tagMark = "ror"
)

type tag struct {
	name       string
	source     string
	allowEmpty bool
}

func NewTag(name, mark string) *tag {
	_tag := &tag{
		name:   name,
		source: mark,
	}
	if mark == "" {
		return _tag
	}

	parts := strings.Split(mark, ",")
	for i := 0; i < len(parts); i++ {
		part := strings.TrimSpace(parts[i])
		switch part {
		case "null":
			_tag.allowEmpty = true
		default:
			if i == 0 {
				_tag.name = part
			}
		}
	}

	return _tag
}

func (_tag *tag) Skip() bool {
	return len(_tag.source) == 0 || _tag.name == "-"
}

func (_tag *tag) AllowEmpty() bool {
	return _tag.allowEmpty
}

func (_tag *tag) Name() string {
	return _tag.name
}
