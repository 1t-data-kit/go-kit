package excel_bak

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"
	"io"
)

/*
file是多个sheet的集合,且确定那个sheet是最终的active sheet
*/
type file interface {
	GetSheetList() []interface{}
	GetActiveSheetIndex() int
}

func Export(file2 file) (*excelize.File, error) {
	xlsx := excelize.NewFile()
	i := 0
	for _, s := range file2.GetSheetList() {
		if s == nil {
			continue
		}
		_sheet, ok := s.(sheet)
		if !ok {
			return nil, fmt.Errorf("sheet[idx=%d] can not match to sheet interface", i)
		}
		if err := exportSheet(xlsx, _sheet, i == 0); err != nil {
			return nil, errors.Wrapf(err, "sheet[%s] Export fail", _sheet.SheetName())
		}
		i++
	}
	active := file2.GetActiveSheetIndex()
	if active <= i {
		xlsx.SetActiveSheet(active)
	}
	return xlsx, nil
}

func Import(file2 file, reader io.Reader, opt ...excelize.Options) error {

	xlsx, err := excelize.OpenReader(reader, opt...)
	if err != nil {
		return errors.Wrap(err, "excel open fail")
	}
	for i, s := range file2.GetSheetList() {
		_sheet, ok := s.(sheet)
		if !ok {
			return fmt.Errorf("sheet[idx=%d] can not match to sheet interface", i)
		}
		if err = importSheet(xlsx, _sheet); err != nil {
			return errors.Wrapf(err, "sheet[%s] Import fail", _sheet.SheetName())
		}
	}
	return nil
}
