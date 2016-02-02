package grids

import "fmt"

// Planar matrix: 0 are void cells, numbers represents pieces. Each piece is represented by a unique number.
// m[2][3] => row 2, column 3
type Matrix2d [][]int

type Grid2d interface {

	// Are they equal? Have orientation, so reflections or rotations make them different!
	Equal(m *Matrix2d) (res bool)
	Clone() *Matrix2d
	Rows() (r int)
	Cols() (c int)
	Max() int
	At(row int, col int) int
	SetAt(row int, col int, value int)

	GeneratePieces(pieces *[]GridPiece2, alikePieces [][]int) int

	PieceMovements(*GridPiece2) []*GridMov2

	CanPieceMove(p *GridPiece2, dRow int, dCol int) bool
}

// Implments Grid2d interface on Matrix2d
func (g *Matrix2d) Rows() (r int) {
	if g != nil {
		return len(*g)
	}
	return 0
}

// Implments Grid2d interface on Matrix2d
func (g *Matrix2d) Cols() (r int) {
	if g != nil {
		if len(*g) > 0 {
			return len((*g)[0])
		}
	}
	return 0
}

// Implments Grid2d interface on Matrix2d
func (g *Matrix2d) At(row int, col int) int {
	return (*g)[row][col]
}

func (g *Matrix2d) SetAt(row int, col int, value int) {
	(*g)[row][col] = value
}

// Implments Grid2d interface on Matrix2d
func (g *Matrix2d) Copy(m *Matrix2d) {

	if m != nil && len(*m) > 0 {
		rows := len(*m)
		cols := len((*m)[0])

		*g = make([][]int, len(*m))

		for i := 0; i < rows; i++ {
			(*g)[i] = make([]int, cols)

			copy((*g)[i], (*m)[i])
		}
	}
}

// Implments Grid2d interface on Matrix2d
func (g *Matrix2d) Max() int {
	rows := g.Rows()
	cols := g.Cols()

	max := 0
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			v := g.At(r, c)
			if max < v {
				max = v
			}
		}
	}
	return max
}

/**
 * @summary Generates one GridPiece2 for each piece in the matrix.
 * @returns {int} Number of extracted pieces
 */
func (g *Matrix2d) GeneratePieces(pieces *[]*GridPiece2, alikePieces [][]int) int {
	rows := g.Rows()
	cols := g.Cols()

	alikeValue := g.Max() + 1
	alikeValuesMap := make(map[int]int)

	for _, equivPieceIds := range alikePieces {
		for _, id := range equivPieceIds {
			alikeValuesMap[id] = alikeValue
		}
		alikeValue++
	}

	foundPieces := make(map[int]*GridPiece2)
	numPieces := 0

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			v := g.At(r, c)
			if v != 0 {
				var p *GridPiece2 = nil

				// If the piece doesn't exist, create it
				p = foundPieces[v]
				if p == nil {
					//*pieces = append(*pieces, &GridPiece2{v, nil})
					*pieces = append(*pieces, NewGridPiece2(v, alikeValuesMap[v]))
					p = (*pieces)[numPieces]
					numPieces++

					foundPieces[v] = p
				}
				p.AddCell([2]int{r, c})
			}
		}
	}

	// Fix positions, making cells' origin the piece origin
	for _, p := range *pieces {
		p.FixPosition()
	}
	return len(*pieces)
}

// Returns array of valid movements for this piece
func (g *Matrix2d) PieceMovements(p *GridPiece2) (movs []*GridMov2) {

	if g.CanPieceMove(p, 1, 0) {
		movs = append(movs, &GridMov2{p.Id(), 1, 0})
	}
	if g.CanPieceMove(p, -1, 0) {
		movs = append(movs, &GridMov2{p.Id(), -1, 0})
	}
	if g.CanPieceMove(p, 0, 1) {
		movs = append(movs, &GridMov2{p.Id(), 0, 1})
	}
	if g.CanPieceMove(p, 0, -1) {
		movs = append(movs, &GridMov2{p.Id(), 0, -1})
	}

	return movs
}

// Returns true if the piece can be moved dRow rows and dCol cols (one at a time)
func (g *Matrix2d) CanPieceMove(p *GridPiece2, dRow int, dCol int) bool {

	rows := g.Rows()
	cols := g.Cols()
	pieceId := p.Id()

	for _, cell := range p.cells {
		row := p.position[0] + cell[0] + dRow
		col := p.position[1] + cell[1] + dCol

		if row < 0 || row >= rows || col < 0 || col >= cols {
			return false
		}
		c := g.At(row, col)

		if c != 0 && c != pieceId {
			return false
		}
	}

	return true
}

// Sets to 0 all piece cells into the matrix
func (g *Matrix2d) ClearPiece(p *GridPiece2) {

	rows := g.Rows()
	cols := g.Cols()
	pieceId := p.Id()

	for _, cell := range p.cells {
		row := p.position[0] + cell[0]
		col := p.position[1] + cell[1]

		if row < 0 || row >= rows || col < 0 || col >= cols {
			fmt.Printf("\n\n		*ERROR: <grid2d::ClearPiece> Invalid cell position (%d,%d)", row, col)
		}

		if pieceId != g.At(row, col) {
			fmt.Printf("\n\n		*ERROR: <grid2d::ClearPiece> Invalid matrix cell value %d at (%d,%d)", g.At(row, col), row, col)
		}

		g.SetAt(row, col, 0)
	}
}

// Sets to 0 all piece cells into the matrix
func (g *Matrix2d) PlacePiece(p *GridPiece2) {

	rows := g.Rows()
	cols := g.Cols()
	pieceId := p.Id()

	for _, cell := range p.cells {
		row := p.position[0] + cell[0]
		col := p.position[1] + cell[1]

		if row < 0 || row >= rows || col < 0 || col >= cols {
			fmt.Printf("\n\n		*ERROR: <grid2d::PlacePiece> Invalid cell position (%d,%d)", row, col)
		}

		if 0 != g.At(row, col) {
			fmt.Printf("\n\n		*ERROR: <grid2d::PlacePiece> Invalid matrix cell value %d at (%d,%d)", g.At(row, col), row, col)
		}

		g.SetAt(row, col, pieceId)
	}
}

// Sets to 0 all piece cells into the matrix
func (g *Matrix2d) UpdatePiecePositions(piecesById map[int]*GridPiece2) {

	rows := g.Rows()
	cols := g.Cols()

	numPieces := len(piecesById)
	foundPieces := make(map[int]bool)
	countPieces := 0

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			v := g.At(r, c)
			if v != 0 {
				_, ok := foundPieces[v]
				if !ok {
					p := piecesById[v]
					foundPieces[v] = true
					p.SetPosition(r, c)
					countPieces++

					if countPieces == numPieces {
						return
					}
				}
			}
		}
	}
}
