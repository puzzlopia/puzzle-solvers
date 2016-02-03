package grids

import "fmt"

type GridMov2 struct {
	pieceId int
	dRow    int
	dCol    int
}

func (m *GridMov2) PieceId() int {
	return m.pieceId
}

func (m *GridMov2) SameMov(dRow int, dCol int) bool {
	return m.dRow == dRow && m.dCol == dCol
}

func (m *GridMov2) Inverted() interface{} {
	x := GridMov2{m.pieceId, 0, 0}
	x.dRow = -1 * m.dRow
	x.dCol = -1 * m.dCol

	return &x
}

func (m *GridMov2) IsInverse(x interface{}) bool {
	y, ok := x.(*GridMov2)
	if ok {
		return m.pieceId == y.pieceId && (m.dRow+y.dRow == 0) && (m.dCol+y.dCol == 0)
	} else {
		panic("[GridMov2::IsInverse] param is not a GridMov2 instance!")
	}
	return false
}

func (m *GridMov2) Print() {
	fmt.Printf("[%d]~(%d, %d)", m.pieceId, m.dRow, m.dCol)
}
