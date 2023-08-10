package excel

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
	"io"
	"reflect"
)

func Export(optionList ...Option) (*excelize.File, error) {
	list := OptionList(optionList)

	xlsx := excelize.NewFile()

	for i, _sheet := range list.SheetList() {
		if _sheet == nil {
			continue
		}
		_sheetName := _sheet.SheetName()
		if i == 0 {
			xlsx.SetSheetName(defaultSheetName, _sheetName)
		} else {
			xlsx.NewSheet(_sheetName)
		}
		if err := exportSheet(xlsx, _sheet); err != nil {
			return nil, err
		}
	}
	xlsx.SetActiveSheet(list.SheetIndex())

	return xlsx, nil
}

func Import(reader io.Reader, optionList ...Option) error {
	list := OptionList(optionList)

	xlsx, err := excelize.OpenReader(reader, list.ExcelOptionList()...)
	if err != nil {
		return errors.Wrap(err, "excel open fail")
	}

	for _, _sheet := range list.SheetList() {
		logrus.Debug(xlsx, _sheet)
	}

	return nil
}

func reflectType(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Slice || t.Kind() == reflect.Ptr || t.Kind() == reflect.Interface || t.Kind() == reflect.Map {
		return reflectType(t.Elem())
	}
	return t
}

func reflectValue(value reflect.Value) reflect.Value {
	if value.Kind() == reflect.Ptr || value.Kind() == reflect.Interface {
		return reflectValue(value.Elem())
	}
	return value
}

func axis(x int, y int) (string, error) {
	_col, err := excelize.ColumnNumberToName(x)
	if err != nil {
		return "", err
	}
	return excelize.JoinCellName(_col, y)
}
