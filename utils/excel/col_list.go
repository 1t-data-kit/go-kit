package excel

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"
	"reflect"
	"sort"
)

type colListSortByTitle []*col

func (list colListSortByTitle) Len() int {
	return len(list)
}

func (list colListSortByTitle) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func (list colListSortByTitle) Less(i, j int) bool {
	return list[i].TagList.Title() < list[j].TagList.Title()
}

type colList []*col

func (list colList) Len() int {
	return len(list)
}

func (list colList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func (list colList) Less(i, j int) bool {
	return list[i].TagList.Index() < list[j].TagList.Index()
}

func (list colList) reflectTypeExists(typ reflect.Type) bool {
	for _, item := range list {
		if reflectType(item.Field.Type).Name() == typ.Name() {
			return true
		}
	}

	return false
}

func (list colList) fetchByExpand() colList {
	expandList := make(colList, 0)
	for _, item := range list {
		if l := colList(item.Children).fetchByExpand(); len(l) > 0 || len(item.ExpandChildren) > 0 {
			expandList = append(expandList, item)
		}
	}

	return expandList
}

func (list colList) titleExists(title string) bool {
	for _, item := range list {
		if item.TagList.Title() == title {
			return true
		}
	}

	return false
}

func (list colList) Level() int {
	return 0
}

func (list colList) fill(_sheet sheet, fillChildren bool) error {
	handle := list.fillDataByReflectValue
	handleLabel := "填充数据"
	if fillChildren {
		handle = list.fillChildrenByReflectValue
		handleLabel = "填充子表头"
	}
	data := reflectValue(reflect.ValueOf(_sheet))

	switch data.Kind() {
	case reflect.Struct:
		if err := handle(data); err != nil {
			return errors.Wrapf(err, "[正则展开多列模式]%s失败", handleLabel)
		}
	case reflect.Slice:
		for i := 0; i < data.Len(); i++ {
			item := reflectValue(data.Index(i))
			if item.Kind() != reflect.Struct {
				return fmt.Errorf("[正则展开多列模式]第[%d]行%s失败: 数据仅支持struct结构", i+1, handleLabel)
			}
			if err := handle(item); err != nil {
				return errors.Wrapf(err, "[正则展开多列模式]第[%d]行%s失败", i+1, handleLabel)
			}
		}
	default:
		return fmt.Errorf("[正则展开多列模式]%s失败: 表数据结构仅支持单行struct数据或者[]*struct多行数据", handleLabel)
	}

	return nil
}

func (list colList) fillChildrenByReflectValue(data reflect.Value) error {
	data = reflectValue(data)

	for _, _col := range list {
		if _col == nil {
			continue
		}
		_colValue := reflectValue(data.FieldByName(_col.Field.Name))
		if !_colValue.IsValid() {
			continue
		}
		if len(_col.ExpandChildren) == 0 {
			if err := colList(_col.Children).fillChildrenByReflectValue(_colValue); err != nil {
				return errors.Wrapf(err, "字段[%s]填充子表头错误", _col.Field.Name)
			}
			continue
		}

		if _colValue.Kind() != reflect.Map {
			return fmt.Errorf("字段[%s]填充子表头错误: 当前是正则展开多列模式,仅支持map数据,当前[%v]不支持，", _col.Field.Name, _colValue.Kind())
		}

		for _, key := range _colValue.MapKeys() {
			title := fmt.Sprintf("%v", key.Interface())
			if !colList(_col.Children).titleExists(title) {
				child := _col.GenerateExpandChild(title)
				if len(child.Children) > 0 {
					if err := colList(child.Children).fetchByExpand().fillChildrenByReflectValue(_colValue.MapIndex(key)); err != nil {
						return err
					}
				}
				_col.Children = append(_col.Children, child)
			}
		}
		sort.Sort(colListSortByTitle(_col.Children))
	}

	return nil
}

func (list colList) fillDataByReflectValue(data reflect.Value) error {
	data = reflectValue(data)

	for _, _col := range list {
		if _col == nil {
			continue
		}

		var _colValue reflect.Value
		if !data.IsValid() {
			_colValue = data
		} else {
			if _col.IsExpanded {
				_colValue = reflectValue(data.MapIndex(reflect.ValueOf(_col.TagList.Title())))
			} else {
				_colValue = reflectValue(data.FieldByName(_col.Field.Name))
			}
		}

		//没有子列的直接填充数据
		if len(_col.Children) == 0 {
			_col.fillData(_colValue)
			continue
		}

		if err := colList(_col.Children).fillDataByReflectValue(_colValue); err != nil {
			return errors.Wrapf(err, "字段[%s]填充子列数据错误", _col.Field.Name)
		}
	}
	return nil
}

func (list colList) export(file *excelize.File, sheetName string, offset int) (int, error) {
	var err error
	for _, _col := range list {
		offset, err = _col.export(file, sheetName, offset)
		if err != nil {
			return 0, errors.Wrapf(err, "列[%s]生成excel错误", _col.TagList.Title())
		}
	}

	return offset, nil
}

func newColList(_sheet sheet) (colList, error) {
	typ := reflectType(reflect.TypeOf(_sheet))
	if typ.Kind() != reflect.Struct {
		return nil, fmt.Errorf("sheet only support struct or struct slice")
	}

	_colList, err := newColListByReflectType(typ)
	if err != nil {
		return nil, err
	}

	expandList := _colList.fetchByExpand()
	if len(expandList) > 0 {
		if err = expandList.fill(_sheet, true); err != nil {
			return nil, err
		}
	}

	return _colList, nil
}

func newColListWithData(_sheet sheet) (colList, error) {
	_colList, err := newColList(_sheet)
	if err != nil {
		return nil, err
	}

	if err = _colList.fill(_sheet, false); err != nil {
		return nil, err
	}

	return _colList, nil
}

func newColListByReflectType(typ reflect.Type, parentList ...*col) (colList, error) {
	fieldLength := typ.NumField()
	_colList := make(colList, 0, fieldLength)

	for i := 0; i < fieldLength; i++ {
		field := typ.Field(i)
		fieldRealType := reflectType(field.Type)

		_col := newColByReflectField(field)
		title := _col.TagList.Title()

		if colList(parentList).reflectTypeExists(fieldRealType) {
			return nil, fmt.Errorf("字段[%s]类型[%s]因递归结构导致无限循环无法支持", title, fieldRealType.Name())
		}

		switch field.Type.Kind() {
		case reflect.Slice:
			return nil, fmt.Errorf("字段[%s]是Slice类型,请使用正则展开多列模式以便生成表头", title)
		case reflect.Map:
			if _col.TagList.Expand() == nil {
				return nil, fmt.Errorf("字段[%s]当前是正则展开多列模式,但未设置正则标签", title)
			}
			if fieldRealType.Kind() != reflect.Struct {
				break
			}
			fallthrough
		case reflect.Struct:
			children, err := newColListByReflectType(fieldRealType, append(parentList, _col)...)
			if err != nil {
				return nil, err
			}
			if _col.TagList.Expand() != nil {
				_col.ExpandChildren = children
			} else {
				_col.Children = children
			}
		}
		_colList = append(_colList, _col)
	}
	sort.Sort(_colList)

	return _colList, nil
}
