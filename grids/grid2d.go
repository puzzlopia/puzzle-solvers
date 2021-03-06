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

	for _, cell := range p.cells_ {
		row := p.position_[0] + cell[0] + dRow
		col := p.position_[1] + cell[1] + dCol

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

// Returns true if the piece can be moved dRow rows and dCol cols (one at a time)
func (g *Matrix2d) ValidMove(mov GridMov2) bool {

	pieceId := mov.PieceId()
	dRow, dCol := mov.Translation()

	rows := g.Rows()
	cols := g.Cols()

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			x := g.At(r, c)

			if x == pieceId {
				row := r + dRow
				col := c + dCol

				if row < 0 || row >= rows || col < 0 || col >= cols {
					return false
				}
				c := g.At(row, col)

				if c != 0 && c != pieceId {
					return false
				}
			}
		}
	}

	return true
}

// Sets to 0 all piece cells into the matrix
func (g *Matrix2d) ClearPiece(p *GridPiece2) {

	rows := g.Rows()
	cols := g.Cols()
	pieceId := p.Id()

	for _, cell := range p.cells_ {
		row := p.position_[0] + cell[0]
		col := p.position_[1] + cell[1]

		if row < 0 || row >= rows || col < 0 || col >= cols {
			fmt.Printf("\n\n		*ERROR: <grid2d::ClearPiece> Invalid cell OUTSIDE position (%d,%d)\n", row, col)
			panic("<grid2d::ClearPiece> Invalid matrix cell")
		}

		if pieceId != g.At(row, col) {
			fmt.Printf("\n\n		*ERROR: <grid2d::ClearPiece> Invalid matrix cell VALUE %d at (%d,%d) while clearing piece %d\n", g.At(row, col), row, col, pieceId)
			panic("<grid2d::ClearPiece> Invalid matrix cell")
		}

		g.SetAt(row, col, 0)
	}
}

// Sets to piece.Id() all piece cells into the matrix
func (g *Matrix2d) PlacePiece(p *GridPiece2) {

	rows := g.Rows()
	cols := g.Cols()
	pieceId := p.Id()

	for _, cell := range p.cells_ {
		row := p.position_[0] + cell[0]
		col := p.position_[1] + cell[1]

		if row < 0 || row >= rows || col < 0 || col >= cols {
			fmt.Printf("\n\n		*ERROR: <grid2d::PlacePiece> Invalid cell OUTSIDE position (%d,%d)", row, col)
			panic("<grid2d::PlacePiece> Invalid matrix cell")
		}

		if 0 != g.At(row, col) {
			fmt.Printf("\n\n		*ERROR: <grid2d::PlacePiece> Invalid matrix cell value %d at (%d,%d)", g.At(row, col), row, col)
			panic("<grid2d::PlacePiece> Invalid matrix cell")
		}

		g.SetAt(row, col, pieceId)
	}
}

// Moves each detected piece cell using the movement object.
// TODO: optimize this using only the piece's cell structure to avoid reading the whole grid.
func (g *Matrix2d) ApplyRawTranslation(pieceId int, mov GridMov2) {
	rows := g.Rows()
	cols := g.Cols()

	dRow, dCol := mov.Translation()

	// First, calc piece cells:
	cells := [][]int{}
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			x := g.At(r, c)

			if x == pieceId {
				cells = append(cells, []int{r, c})

				g.SetAt(r, c, 0) //clear old place
			}
		}
	}

	// for r := 0; r < rows; r++ {
	// 	for c := 0; c < cols; c++ {
	for _, cell := range cells {
		r := cell[0]
		c := cell[1]

		newRow := r + dRow
		newCol := c + dCol

		if newRow < 0 || newRow >= rows || newCol < 0 || newCol >= cols {
			fmt.Printf("\n\n		*ERROR::ApplyRawTranslation Cannot move (%v) piece %d to OUTSIDE cell (%d,%d)\n", mov, pieceId, newRow, newCol)
			panic("[grids::ApplyRawTranslation] invalid movement!")
		}

		// And rewrite
		y := g.At(newRow, newCol)

		if y != 0 {
			fmt.Printf("\n\n		*ERROR::ApplyRawTranslation Cannot move (%v) piece %d to cell (%d,%d) occupied by %d\n", mov, pieceId, newRow, newCol, y)
			panic("[grids::ApplyRawTranslation] invalid movement!")
		}

		g.SetAt(newRow, newCol, pieceId) //set new place

	}
}

// Reads the matrix and positionates each piece.
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
					p.SetPositionTL(r, c)
					countPieces++

					if countPieces == numPieces {
						return
					}
				}
			}
		}
	}
}

// Reads the matrix and positionates each piece.
func (g *Matrix2d) Identical(m Matrix2d) bool {

	rows := g.Rows()
	cols := g.Cols()

	if rows == m.Rows() && cols == m.Cols() {
		for r := 0; r < rows; r++ {
			for c := 0; c < cols; c++ {

				if g.At(r, c) != m.At(r, c) {
					return false
				}
			}
		}
	}

	return true
}
