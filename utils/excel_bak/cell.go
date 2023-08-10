package excel_bak

import (
	"encoding/json"
	"fmt"
	"github.com/1t-data-kit/go-kit/utils/excel"
	"github.com/xuri/excelize/v2"
	"reflect"
	"regexp"
	"sort"
)

const (
	tagName = "excel"
)

type cell struct {
	IsTitle        bool
	SheetName      string
	Name           string
	Index          int
	Start          *excel.position
	End            *excel.position
	ExpandedRegexp *regexp.Regexp //展开为多列时的正则匹配
	Style          *excelize.Style
	Data           interface{}
	Child          row
}

func newTitleCellByStructField(sheetName string, field reflect.StructField, parent *cell) (*cell, error) {
	c := &cell{
		IsTitle:   true,
		SheetName: sheetName,
		Name:      field.Name,
		Child:     make(row, 0),
		Data:      "",
	}
	for _, attribute := range parseCellAttributeByTag(field.Tag.Get(tagName)) {
		if !attribute.tag.IsValid() && c.Data == "" {
			c.Data = attribute.data
		}
		switch attribute.tag {
		case cellTagName:
			c.Data = attribute.AliasName()
		case cellTagIndex:
			c.Index = attribute.Index()
		case cellTagStyle:
			c.Style = attribute.Style()
		case cellTagExpand:
			c.ExpandedRegexp = attribute.Expand()
		}
	}
	if c.Data == "" {
		c.Data = field.Name
	}
	if c.Index == 0 && len(field.Index) > 0 {
		c.Index = field.Index[0] + 1
	}

	x := c.Index
	y := 1
	if parent != nil {
		//x += parent.Start.X - 1
		y += parent.Start.Y
	}
	c.Start = &excel.position{
		X: x,
		Y: y,
	}

	return c, nil
}

func (c *cell) export(file2 *excelize.File, offset int) (int, error) {
	var err error
	realStartPosition := &excel.position{
		X: c.Start.X + offset,
		Y: c.Start.Y,
	}
	realEndPosition := &excel.position{
		X: c.End.X + offset,
		Y: c.End.Y,
	}
	startAxis, err := realStartPosition.Axis()
	if err != nil {
		return 0, err
	}
	endAxis, err := realEndPosition.Axis()
	if err != nil {
		return 0, err
	}

	if err = file2.SetCellValue(c.SheetName, startAxis, c.Data); err != nil {
		return 0, err
	}
	if !realStartPosition.Equal(realEndPosition) {
		file2.MergeCell(c.SheetName, startAxis, endAxis)
	}
	if len(c.Child) > 0 {
		if offset, err = c.Child.export(file2, offset); err != nil {
			return 0, err
		}
	}

	return offset + len(c.Child), nil
}

func (c *cell) String() string {
	data, _ := json.Marshal(c)
	return string(data)
}

func (c *cell) AppendChild(list ...*cell) {
	c.Child = append(c.Child, list...)
}

func (c *cell) MaxRowLength() int {
	if len(c.Child) == 0 {
		return c.Start.Y
	}
	return c.Child.MaxRowLength()
}

func (c *cell) MaxColLength() int {
	if len(c.Child) == 0 {
		return 1
	}
	return c.Child.MaxColLength() - 1
}

func (c *cell) endPositionWithY(y int) {
	childLength := len(c.Child)
	if childLength == 0 {
		c.End = &excel.position{
			X: c.Start.X,
			Y: y,
		}
		return
	}
	c.End = &excel.position{
		X: c.Start.X + c.Child.MaxColLength(),
		Y: c.Start.Y,
	}
	c.Child.endPositionWithY(y)
}

func (c *cell) ExpandTitle(title string) error {
	if c.ExpandedRegexp == nil {
		return fmt.Errorf("展开多列模式时[%s]未设置正则匹配方式", c.Name)
	}
	if !c.ExpandedRegexp.MatchString(title) {
		return fmt.Errorf("展开多列模式时[%s]正则规则匹配内容[%s]失败", c.Name, title)
	}

	modify := false
	if !sliceContainString(c.ExpandSortedTitle, title) {
		c.ExpandSortedTitle = append(c.ExpandSortedTitle, title)
		modify = true
		c.ExpandChild[title] = row{
			c.ExpandCell(title),
		}
	}
	if modify {
		sort.Strings(c.ExpandSortedTitle)
	}

	return nil
}

func (c *cell) ExpandCell(title string) *cell {
	_cell := &cell{
		SheetName:      c.SheetName,
		Name:           c.Name,
		Index:          c.Index,
		Start:          c.Start.Copy(),
		ExpandedRegexp: c.ExpandedRegexp,
		Style:          c.Style,
		Data:           title,
		Child:          make(row, 0),
		//TODO:坐标
	}
	if c.End != nil {
		_cell.End = c.End.Copy()
	}
	return _cell
}

func (c *cell) DataCell(value reflect.Value) *cell {
	_cell := &cell{
		SheetName:      c.SheetName,
		Name:           c.Name,
		Index:          c.Index,
		Start:          c.Start.Copy(),
		ExpandedRegexp: c.ExpandedRegexp,
		Style:          c.Style,
		Data:           value.Interface(),
		Child:          make(row, 0),
		//TODO:坐标
	}
	if c.End != nil {
		_cell.End = c.End.Copy()
	}
	return _cell
}
