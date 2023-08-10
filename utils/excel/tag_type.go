package excel

type tagType string

const (
	tagTypeTitle  = tagType("title")
	tagTypeIndex  = tagType("index")
	tagTypeExpand = tagType("expand")
	tagTypeStyle  = tagType("style")
)

var tagTypeDict = map[tagType]string{
	tagTypeTitle:  "表头",
	tagTypeIndex:  "序号",
	tagTypeExpand: "展开列",
	tagTypeStyle:  "样式",
}

func (typ tagType) String() string {
	if !typ.IsValid() {
		return ""
	}
	return tagTypeDict[typ]
}

func (typ tagType) IsValid() bool {
	if _, exists := tagTypeDict[typ]; exists {
		return true
	}
	return false
}
