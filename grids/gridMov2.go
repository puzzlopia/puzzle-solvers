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

// Returns true if the 'mov' makes the trajectory to touch itself.
// Example
/**
 *	x->x->x->x
 *					 |
 *					 x
 *					 |
 *	   x<-x<-x
 *		 |  ^
 *		 x->y
 *
 * where 'y' is the new position casued by mov: it touches the trajectory on its top face!
 */
func TrajectoryTouchesWithMov(invertedPath []*GridMov2, mov interface{}) bool {

	// Short paths cannot have loops, the shortest one is 3-length
	if len(invertedPath) < 2 {
		return false
	}

	gMov, ok := mov.(*GridMov2)
	if !ok {
		panic("[grids::TrajectoryTouchesWithMov] mov is not a GridMov2!")
	}
	// Set the origin at the las mov when it is done. Then apply the inverse of the mov and then the inverse of
	// the path, starting from the end. Then, if we arrive at any adjacent of the origin at any moment, the path touches itself!
	position := [2]int{0, 0}

	// Move to the only one valid adjacent position of the origin
	movI := gMov.Inverted().(*GridMov2)
	position[0] += movI.dRow
	position[1] += movI.dCol

	//for i := len(path) - 1; i >= 0; i-- {
	for _, m := range invertedPath {
		//m := path[i]
		mInv := m.Inverted().(*GridMov2)

		position[0] += mInv.dRow
		position[1] += mInv.dCol

		if position[0] == movI.dRow && position[1] == movI.dCol {
			//Then we are on the last part of the path..again! The algorithm should have detected it before!

			// fmt.Println("**** ERROR: [grids::TrajectoryTouchesWithMov] algorithm detecting loop inside path!")
			// fmt.Println("PATH (inverted):")
			// for _, mm := range invertedPath {
			// 	fmt.Printf("	P[%d] mov (%d, %d)\n", mm.pieceId, mm.dRow, mm.dCol)
			// }
			// fmt.Println("MOV:", mov)
			// fmt.Println("POSITION:", position)

			panic("[grids::TrajectoryTouchesWithMov] algorithm detecting loop inside path!")
		}

		// Adjacent positions are (-1, 0), (1, 0), (0, 1) and (0, -1), the only ones with module 1:
		if position[0]*position[0]+position[1]*position[1] == 1 {
			// fmt.Println("**** DETECTED LOOP!: [grids::TrajectoryTouchesWithMov] algorithm detecting loop inside path!")
			// fmt.Println("PATH (inverted):")
			// for _, mm := range invertedPath {
			// 	fmt.Printf("	P[%d] mov (%d, %d)\n", mm.pieceId, mm.dRow, mm.dCol)
			// }
			// fmt.Println("MOV:", mov)
			// fmt.Println("POSITION:", position)

			return true
		}
	}

	return false
}
