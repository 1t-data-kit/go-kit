package excel_bak

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
)

const (
	defaultSheetName = "Sheet1"
)

/*
sheet是一个数据对象(struct)的slice
*/

type sheet interface {
	SheetName() string
}

func exportSheet(file2 *excelize.File, sheet2 sheet, mustReusedDefaultSheet bool) error {
	sheetName := sheet2.SheetName()
	if mustReusedDefaultSheet {
		file2.SetSheetName(defaultSheetName, sheetName)
	} else {
		file2.NewSheet(sheetName)
	}

	title, err := parseTitle(sheet2)
	if err != nil {
		return errors.Wrap(err, "parseTitle fail")
	}
	data, err := parseData(sheet2, title)
	if err != nil {
		return errors.Wrap(err, "parseData fail")
	}
	title.EndPosition()
	logrus.Debug(title)
	return nil
	if err = append(rowList{title}, data...).export(file2); err != nil {
		return errors.Wrap(err, "export fail")
	}
	return nil
}

func importSheet(file2 *excelize.File, sheet2 sheet) error {
	return nil
}
