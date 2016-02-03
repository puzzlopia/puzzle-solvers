package grids

//import "fmt"
import "github.com/edgarweto/puzzlopia/puzzle-solvers/definitions"

// Algorithms and structures to deal with paths on grids.
type GridPath2 struct {
	//path2_ []*GridMov2
	path_ []defs.Command
}

// Builds a path of the last moved piece. The result is an inverse path (starts with last movement)
func (p *GridPath2) BuildFromReversePath(rp []defs.Command) {

	pieceId := 0
	for _, mov := range rp {

		if pieceId == 0 || pieceId == mov.PieceId() {
			pieceId = mov.PieceId()
			p.path_ = append(p.path_, mov)
		} else {
			break
		}
	}
}

func (p *GridPath2) Path() []defs.Command {
	return p.path_
}

// Returns true if the 'mov' makes the trajectory to touch itself.
// Example
/**
 *	x->x->x->x
 *					 |
 *					 x
 *					 |
 *	   x<-x<-x
 *		 |  ^-touches!
 *		 x->y
 *
 * where 'y' is the new position casued by mov: it touches the trajectory on its top face!
 */
func TrajectoryTouchesWithMov(invertedPath []defs.Command, mov defs.Command) bool {

	// Short paths cannot have loops, the shortest one is 3-length
	if len(invertedPath) < 2 {
		return false
	}

	// Check we are still moving the same piece!
	if invertedPath[0].PieceId() != mov.PieceId() {
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
			// fmt.Println("\nMOV:", mov)
			// fmt.Println("\nPATH (inverted):")
			// for _, mm := range invertedPath {
			// 	//fmt.Printf("	P[%d] mov (%d, %d)\n", mm.pieceId, mm.dRow, mm.dCol)
			// 	mm.Print()
			// }

			return true
		}
	}

	return false
}
