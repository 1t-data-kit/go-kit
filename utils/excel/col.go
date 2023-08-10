package excel

import (
	"github.com/mohae/deepcopy"
	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"
	"reflect"
)

type col struct {
	Field          reflect.StructField
	TagList        tagList
	Data           []interface{}
	Children       []*col
	ExpandChildren []*col
	IsExpanded     bool
}

func (_col *col) fillData(value reflect.Value) {
	data := _col.EmptyValue()
	if value.IsValid() {
		data = value.Interface()
	}
	_col.Data = append(_col.Data, data)
}

func (_col *col) EmptyValue() interface{} {
	switch _col.Field.Type.Kind() {
	case reflect.String:
		return ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return int64(0)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return uint64(0)
	case reflect.Float32, reflect.Float64:
		return float64(0)
	}
	return nil
}

func (_col *col) GenerateExpandChild(title string) *col {
	return &col{
		TagList:        tagList{tag(title)},
		Data:           make([]interface{}, 0),
		Children:       deepcopy.Copy(_col.ExpandChildren).([]*col),
		ExpandChildren: make([]*col, 0),
		IsExpanded:     true,
	}
}

func (_col *col) export(file *excelize.File, sheetName string, x int) (int, error) {
	childrenLength := len(_col.Children)
	if childrenLength > 0 {
		return colList(_col.Children).export(file, sheetName, x)
	}
	if _col.TagList.Expand() != nil {
		return x, nil
	}
	for i, _data := range _col.Data {
		_axis, err := axis(x, i+1)
		if err != nil {
			return 0, errors.Wrapf(err, "列[%s]第%d行生成excel坐标失败", _col.TagList.Title(), i+1)
		}
		if err = file.SetCellValue(sheetName, _axis, _data); err != nil {
			return 0, errors.Wrapf(err, "列[%s]第%d行[%s]导出数据失败", _col.TagList.Title(), i+1, _axis)
		}
	}
	return x + 1, nil
}

func newColByReflectField(field reflect.StructField) *col {
	return &col{
		Field:          field,
		TagList:        newTagList(field.Tag.Get(tagNameMark)),
		Data:           make([]interface{}, 0),
		Children:       make([]*col, 0),
		ExpandChildren: make([]*col, 0),
	}
}
