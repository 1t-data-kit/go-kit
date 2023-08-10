package excel

import "github.com/xuri/excelize/v2"

type OptionType int32

const (
	OptionTypeSheet = OptionType(iota + 1)
	OptionTypeSheetIndex
	OptionTypeExcelOption
)

type Option struct {
	typ  OptionType
	data interface{}
}

func NewOption(typ OptionType, data interface{}) Option {
	return Option{
		typ:  typ,
		data: data,
	}
}

func WithSheetOption(_sheet sheet) Option {
	return NewOption(OptionTypeSheet, _sheet)
}

func WithSheetIndexOption(index int) Option {
	return NewOption(OptionTypeSheetIndex, index)
}

func WithExcelOption(option excelize.Options) Option {
	return NewOption(OptionTypeExcelOption, option)
}
