package excel_bak

import (
	"fmt"
	"reflect"
)

func reflectValue(value reflect.Value) reflect.Value {
	if value.Kind() == reflect.Ptr || value.Kind() == reflect.Interface {
		return reflectValue(value.Elem())
	}
	return value
}

func parseData(sheet2 sheet, title row) (rowList, error) {
	_rowList := make(rowList, 0)
	sheetValue := reflectValue(reflect.ValueOf(sheet2))
	switch sheetValue.Kind() {
	case reflect.Struct:
		_row, err := parseDataByValue(title, sheetValue)
		if err != nil {
			return nil, err
		}
		if _row != nil {
			_rowList = append(_rowList, _row)
		}
	case reflect.Slice:
		for i := 0; i < sheetValue.Len(); i++ {
			_row, err := parseDataByValue(title, sheetValue.Index(i))
			if err != nil {
				return nil, err
			}
			if _row != nil {
				_rowList = append(_rowList, _row)
			}
		}
	default:
		return nil, fmt.Errorf("数据结构仅支持struct,slice")
	}

	return _rowList, nil
}

func parseDataByValue(title row, value reflect.Value) (row, error) {
	var err error
	value = reflectValue(value)
	_row := make(row, 0)

	for _, _cell := range title {
		_v := reflectValue(value.FieldByName(_cell.Name))
		_dataCell := _cell.DataCell(_v)
		child := make(row, 0)

		switch _v.Kind() {
		case reflect.Map:
			_child := make(row, 0)
			for _, _key := range _v.MapKeys() {
				_t := fmt.Sprintf("%s", _key.Interface())
				if err = _cell.ExpandTitle(_t); err != nil {
					return nil, err
				}

				_d := _v.MapIndex(_key)

				if len(_cell.Child) == 0 { //普通map[string]非struct结构的导出
					expandChild[_t] = row{_cell.DataCell(_d)}
					continue
				}
				_child, err = parseDataByValue(_cell.Child, _d)
				if err != nil {
					return nil, err
				}
				expandChild[_t] = _child
			}
		case reflect.Struct:
			child, err = parseDataByValue(_cell.Child, _v)
			if err != nil {
				return nil, err
			}
		}
		_dataCell.Child = child
		_row = append(_row, _dataCell)
	}
	return _row, nil
}
