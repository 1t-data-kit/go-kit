package excel_bak

import (
	"encoding/json"
	"github.com/xuri/excelize/v2"
)

type row []*cell

func (r row) Len() int {
	return len(r)
}

func (r row) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r row) Less(i, j int) bool {
	return r[i].Index < r[j].Index
}

func (r row) EndPosition() {
	r.endPositionWithY(r.MaxRowLength())
}

func (r row) endPositionWithY(y int) {
	for _, _cell := range r {
		_cell.endPositionWithY(y)
	}
}

func (r row) MaxRowLength() int {
	max := 0
	for _, _cell := range r {
		if y := _cell.MaxRowLength(); y > max {
			max = y
		}
	}

	return max
}

func (r row) MaxColLength() int {
	l := 0

	for _, _cell := range r {
		l += _cell.MaxColLength()
	}
	if l == 0 {
		return 0
	}

	return l
}

func (r row) export(file2 *excelize.File, offset int) (int, error) {
	var err error
	for _, _cell := range r {
		if offset, err = _cell.export(file2, offset); err != nil {
			return 0, err
		}
	}
	return offset, nil
}

func (r row) String() string {
	data, _ := json.Marshal(r)
	return string(data)
}

type rowList []row

func (list rowList) export(file2 *excelize.File) error {
	for _, _cellList := range list {
		if _, err := _cellList.export(file2, 0); err != nil {
			return err
		}
	}

	return nil
}

func (list rowList) String() string {
	data, _ := json.Marshal(list)
	return string(data)
}
