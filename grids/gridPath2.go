package grids

//import "fmt"

// Algorithms and structures to deal with paths on grids.
type GridPath2 struct {
	path_ []*GridMov2
}

// Builds a path of the last moved piece.
func (p *GridPath2) BuildFromReversePath(rp []interface{}) {
	//revPath, ok := rp.([]interface{})

	pieceId := 0
	//if ok {
	for _, item := range rp {
		mov, ok := item.(*GridMov2)
		if !ok {
			panic("[GridPath2::BuildFromReversePath] path must be built with grids.GridMov2 objects!")
		}

		if pieceId == 0 || pieceId == mov.PieceId() {
			p.path_ = append(p.path_, mov)
		} else {
			break
		}
	}
	// } else {
	// 	panic("[GridPath2::BuildFromReversePath] arg path not valid!")
	// }
}
