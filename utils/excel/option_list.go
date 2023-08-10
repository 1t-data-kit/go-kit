package excel

import "github.com/xuri/excelize/v2"

type OptionList []Option

func (list OptionList) SheetList() []sheet {
	sheetList := make([]sheet, 0)
	for _, item := range list {
		if item.typ != OptionTypeSheet {
			continue
		}
		if _sheet, ok := item.data.(sheet); ok {
			sheetList = append(sheetList, _sheet)
		}
	}

	return sheetList
}

func (list OptionList) SheetIndex() int {
	var index int
	for _, item := range list {
		if item.typ != OptionTypeSheetIndex {
			continue
		}
		if _index, ok := item.data.(int); ok {
			index = _index
		}
	}

	return index
}

func (list OptionList) ExcelOptionList() []excelize.Options {
	optionList := make([]excelize.Options, 0)
	for _, item := range list {
		if item.typ != OptionTypeExcelOption {
			continue
		}
		if _option, ok := item.data.(excelize.Options); ok {
			optionList = append(optionList, _option)
		}
	}

	return optionList
}
