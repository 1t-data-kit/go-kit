package excel_bak

import (
	"github.com/1t-data-kit/go-kit/utils/excel"
	"github.com/xuri/excelize/v2"
	"regexp"
	"strconv"
	"strings"
)

type cellAttribute struct {
	tag  cellTag
	data string
}

func newCellAttribute(tag cellTag, data string) *cellAttribute {
	return &cellAttribute{
		tag:  tag,
		data: data,
	}
}

func (attribute *cellAttribute) AliasName() string {
	if attribute.tag != cellTagName {
		return ""
	}

	return attribute.data
}

func (attribute *cellAttribute) Index() int {
	if attribute.tag != cellTagIndex {
		return 0
	}
	index, _ := strconv.ParseInt(attribute.data, 10, 64)
	return int(index)
}

func (attribute *cellAttribute) Style() *excelize.Style {
	return nil
}

func (attribute *cellAttribute) Position() *excel.position {
	return nil
}

func (attribute *cellAttribute) Expand() *regexp.Regexp {
	return regexp.MustCompile(attribute.data)
}

type cellAttributeList []*cellAttribute

func parseCellAttributeByTag(source string) cellAttributeList {
	attributeList := make(cellAttributeList, 0)
	source = strings.TrimSpace(source)
	for _, g := range strings.Split(source, " ") {
		g = strings.TrimSpace(g)
		sepIndex := strings.Index(g, ":")
		tag := cellTag(g)
		data := g
		if sepIndex > -1 {
			tag = cellTag(g[0:sepIndex])
			data = g[sepIndex+1:]
		}
		attributeList = append(attributeList, newCellAttribute(tag, data))
	}

	return attributeList
}
