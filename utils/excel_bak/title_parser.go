package excel_bak

import (
	"fmt"
	"github.com/pkg/errors"
	"reflect"
	"sort"
)

func reflectType(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Slice || t.Kind() == reflect.Ptr || t.Kind() == reflect.Interface || t.Kind() == reflect.Map {
		return reflectType(t.Elem())
	}
	return t
}

func parseTitle(sheet2 sheet) (row, error) {
	t := reflectType(reflect.TypeOf(sheet2))
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("sheet only support struct slice")
	}

	return parseTitleByType(sheet2.SheetName(), t, nil)
}

func parseTitleByType(sheetName string, t reflect.Type, parent *cell) (row, error) {
	var err error
	fieldLength := t.NumField()
	title := make(row, 0, fieldLength)

	for i := 0; i < fieldLength; i++ {
		field := t.Field(i)
		fieldRealType := reflectType(field.Type)

		_cell, _cellError := newTitleCellByStructField(sheetName, field, parent)
		if _cellError != nil {
			return nil, errors.Wrapf(_cellError, "字段[%s]表头创建失败", field.Name)
		}
		child := make(row, 0)

		switch field.Type.Kind() {
		case reflect.Slice:
			return nil, fmt.Errorf("字段[%s]是Slice类型,请使用正则展开多列模式以便生成表头", _cell.Name)
		case reflect.Map:
			if _cell.ExpandedRegexp == nil {
				return nil, fmt.Errorf("字段[%s]当前是正则展开多列模式,但未设置正则标签", _cell.Name)
			}
			if fieldRealType.Kind() != reflect.Struct {
				break
			}
			fallthrough
		case reflect.Struct:
			child, err = parseTitleByType(sheetName, fieldRealType, _cell)
			if err != nil {
				return nil, err
			}
		}
		_cell.Child = child
		title = append(title, _cell)
	}
	sort.Sort(title)

	return title, nil
}
