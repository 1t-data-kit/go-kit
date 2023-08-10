package excel_bak

import (
	"encoding/json"
	"github.com/xuri/excelize/v2"
)

type position struct {
	X int
	Y int
}

func (p *position) Equal(_p *position) bool {
	return p.X == _p.X && p.Y == _p.Y
}

func (p *position) Axis() (string, error) {
	col, err := excelize.ColumnNumberToName(p.X)
	if err != nil {
		return "", err
	}
	return excelize.JoinCellName(col, p.Y)
}

func (p *position) String() string {
	data, _ := json.Marshal(p)
	return string(data)
}

func (p *position) Copy() *position {
	return &position{
		X: p.X,
		Y: p.Y,
	}
}
