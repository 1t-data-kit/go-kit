package excel_bak

type cellTag string

const (
	cellTagName   = cellTag("name")
	cellTagIndex  = cellTag("index")
	cellTagExpand = cellTag("expand")
	cellTagStyle  = cellTag("style")
)

var cellTagDict = map[cellTag]string{
	cellTagName:   "列名",
	cellTagIndex:  "序号",
	cellTagExpand: "展开列",
	cellTagStyle:  "样式",
}

func (tag cellTag) IsValid() bool {
	if _, exists := cellTagDict[tag]; exists {
		return true
	}
	return false
}
