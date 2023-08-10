package excel

import (
	"github.com/xuri/excelize/v2"
	"regexp"
	"strconv"
	"strings"
)

const (
	tagNameMark = "excel"
)

type tag string

func (_tag tag) String() string {
	return string(_tag)
}

func (_tag tag) Type() tagType {
	tagString := _tag.String()
	sepIndex := strings.Index(tagString, ":")
	if sepIndex == -1 {
		return tagTypeTitle
	}
	return tagType(tagString[0:sepIndex])
}

func (_tag tag) Value() string {
	tagString := _tag.String()
	sepIndex := strings.Index(tagString, ":")
	if sepIndex == -1 {
		return tagString
	}
	return tagString[sepIndex+1:]
}

func (_tag tag) Title() string {
	if _tag.Type() != tagTypeTitle {
		return ""
	}
	return _tag.Value()
}

func (_tag tag) Index() int {
	if _tag.Type() != tagTypeIndex {
		return 0
	}
	index, _ := strconv.ParseInt(_tag.Value(), 10, 64)
	return int(index)
}

func (_tag tag) Expand() *regexp.Regexp {
	if _tag.Type() != tagTypeExpand {
		return nil
	}
	return regexp.MustCompile(_tag.Value())
}

func (_tag tag) Style() *excelize.Style {
	if _tag.Type() != tagTypeStyle {
		return nil
	}
	return &excelize.Style{}
}

type tagList []tag

func newTagList(_tag string) tagList {
	list := make(tagList, 0)
	for _, item := range strings.Split(strings.TrimSpace(_tag), " ") {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		list = append(list, tag(item))
	}

	return list
}

func (list tagList) Title() string {
	for _, item := range list {
		if item.Type() != tagTypeTitle {
			continue
		}
		return item.Title()
	}
	return ""
}

func (list tagList) Index() int {
	for _, item := range list {
		if item.Type() != tagTypeIndex {
			continue
		}
		return item.Index()
	}
	return 0
}

func (list tagList) Expand() *regexp.Regexp {
	for _, item := range list {
		if item.Type() != tagTypeExpand {
			continue
		}
		return item.Expand()
	}
	return nil
}

func (list tagList) Style() *excelize.Style {
	for _, item := range list {
		if item.Type() != tagTypeStyle {
			continue
		}
		return item.Style()
	}
	return nil
}
