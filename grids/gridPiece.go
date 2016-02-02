package grids

// //A cell is usually a square, a point, a 'mino' of a polyominoe
// type Cell interface{}
// type Cells []Cell

/**
 * @summary A grid piece, interface that assumes piece structure and status.
 */
type GridPiece interface {

	// Returns the id of the piece, unique for each piece
	Id() int

	// Value used to compare states. Some times there are equivalent pieces, they have this same value
	Value() int

	// Returns the width of its bounding box
	Width() int

	// Returns the height of its bounding box
	Height() int

	// Returns the number of cells
	Len() int

	Move(interface{})

	// // True if the pieces have the same shape
	// Equivalent(*GridPiece) bool
}
