package grids

//import "fmt"

const (
	PIECE_HUGE_SIZE = 9999
)

type Coords2 [2]int

// Definition of a static piece: does not contain status, like position
type GridPiece2 struct {
	pieceId_  int
	value_    int
	cells_    []Coords2
	position_ Coords2

	// Top-left position: is the cell at top-left position, the first that we would see in
	// a row search, from top to bottom, left to right.
	tl_ Coords2
}

func NewGridPiece2(id int, val int) *GridPiece2 {
	value := val
	if val == 0 {
		value = id
	}
	//fmt.Printf("\nNew piece [%d] with value %d", id, value)

	return &GridPiece2{id, value, nil, Coords2{PIECE_HUGE_SIZE, PIECE_HUGE_SIZE}, Coords2{PIECE_HUGE_SIZE, PIECE_HUGE_SIZE}}
}

// Implement GridPiece interface:
func (p *GridPiece2) Id() int {
	return p.pieceId_
}

func (p *GridPiece2) Value() int {
	return p.value_
}

func (p *GridPiece2) IsAt(row int, col int) bool {
	return p.position_[0] == row && p.position_[1] == col
}

func (p *GridPiece2) Len() int {
	return len(p.cells_)
}

func (p *GridPiece2) Cell(i int) *Coords2 {
	return &p.cells_[i]
}

func (p *GridPiece2) SetValue(v int) {
	p.value_ = v
}

// True if the pieces have the same shape
func (p *GridPiece2) Equivalent(q *GridPiece2) bool {
	//fmt.Println("		Equivalents?", p.cells_, q.cells_)

	if len(p.cells_) != len(q.cells_) {
		return false
	}

	// We can compare cells in order, because if the shape is the same, the order of the cells is also the same.
	for i, c1 := range p.cells_ {
		c2 := q.cells_[i]

		if c1[0] != c2[0] || c1[1] != c2[1] {
			return false
		}
	}
	return true
}

func (p *GridPiece2) AddCell(c [2]int) *GridPiece2 {
	p.cells_ = append(p.cells_, c)

	// Min cell is the most top-left cell, the 'origin' of the piece.
	if c[0] < p.position_[0] {
		p.position_[0] = c[0]
	}
	if c[1] < p.position_[1] {
		p.position_[1] = c[1]
	}

	// Update the Most-top-left origin:
	if c[0] == p.tl_[0] {
		if c[1] < p.tl_[1] {
			p.tl_[1] = c[1]
		}
	} else if c[0] < p.tl_[0] {
		p.tl_[0] = c[0]
		p.tl_[1] = c[1]
	}

	return p
}

func (p *GridPiece2) FixPosition() {
	n := len(p.cells_)
	for i := 0; i < n; i++ {
		cell := &p.cells_[i]

		cell[0] -= p.position_[0]
		cell[1] -= p.position_[1]
	}
	p.tl_[0] -= p.position_[0]
	p.tl_[1] -= p.position_[1]
}

func (p *GridPiece2) SetPosition(row int, col int) {
	p.position_[0] = row
	p.position_[1] = col
}

// When we read a grid in order to see where a piece is placed,
// we detect pieces by its id. The top-left positions is the origin of
// the piece. But if the piece is like this:
//
//	00X00
//	XXX00
//	0X000
//
// then it has a 0 in its top-left boundary box. We detect the piece,
// reading by rows, top to bottom and left to right, in the place (row:0, col:2).
// But the real position of the piece (its boundary box) is (0,0). That's why we
// need this function. We will call it SetPositionTL(0,2) because is the position
// where we've detected the piece, but it correctly will position it at (0,0)
//
//
func (p *GridPiece2) SetPositionTL(row int, col int) {
	//fmt.Printf("\n	-SetPositionTL(%d,%d)", row, col)
	p.position_[0] = row - p.tl_[0]
	p.position_[1] = col - p.tl_[1]

	//fmt.Printf("\n	- pos = (%d,%d)", p.position_[0], p.position_[1])
}

func (p *GridPiece2) Move(m interface{}) {
	mov, ok := m.(*GridMov2)
	if ok {
		p.position_[0] += mov.dRow
		p.position_[1] += mov.dCol
	} else {
		panic("[GridPiece2::Move] Unknown type of move!")
	}
}

// Sets the same value for all alike pieces (those with exactly the same shape)
func DetectAlikePieces(pieces []*GridPiece2, minValue int) {

	// Save a representant for each piece
	var shapeClasses []*GridPiece2
	var values []int

	value := minValue

	for _, p := range pieces {
		found := false
		foundIdx := 0
		for i, q := range shapeClasses {
			if p.Equivalent(q) {
				//fmt.Printf("\n		- equivalent pieces: %d - %d", p.Id(), q.Id())
				found = true
				foundIdx = i
				break
			}
		}
		if found {
			p.SetValue(values[foundIdx])
		} else {
			shapeClasses = append(shapeClasses, p)
			values = append(values, value)
			p.SetValue(value)
			value++
		}
	}
}
