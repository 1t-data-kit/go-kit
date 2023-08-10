package excel

import (
	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"
)

const (
	defaultSheetName = "Sheet1"
)

type sheet interface {
	SheetName() string
}

func exportSheet(file *excelize.File, _sheet sheet) error {
	_colList, err := newColListWithData(_sheet)
	if err != nil {
		return errors.Wrapf(err, "exportSheet[%s] error: convert fail", _sheet.SheetName())
	}

	if _, err = _colList.export(file, _sheet.SheetName(), 1); err != nil {
		return errors.Wrapf(err, "exportSheet[%s] error: export fail", _sheet.SheetName())
	}

	return nil
}
