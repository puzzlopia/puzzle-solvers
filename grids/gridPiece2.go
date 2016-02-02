package grids

//import "fmt"

type Coords2 [2]int

// Definition of a static piece: does not contain status, like position
type GridPiece2 struct {
	pieceId_ int
	value_   int
	cells    []Coords2
	position Coords2
}

func NewGridPiece2(id int, val int) *GridPiece2 {
	value := val
	if val == 0 {
		value = id
	}
	//fmt.Printf("\nNew piece [%d] with value %d", id, value)

	return &GridPiece2{id, value, nil, Coords2{9999, 9999}}
}

// Implement GridPiece interface:
func (p *GridPiece2) Id() int {
	return p.pieceId_
}

func (p *GridPiece2) Value() int {
	return p.value_
}

func (p *GridPiece2) Len() int {
	return len(p.cells)
}

func (p *GridPiece2) Cell(i int) *Coords2 {
	return &p.cells[i]
}

func (p *GridPiece2) SetValue(v int) {
	p.value_ = v
}

// True if the pieces have the same shape
func (p *GridPiece2) Equivalent(q *GridPiece2) bool {
	//fmt.Println("		Equivalents?", p.cells, q.cells)

	if len(p.cells) != len(q.cells) {
		return false
	}

	// We can compare cells in order, because if the shape is the same, the order of the cells is also the same.
	for i, c1 := range p.cells {
		c2 := q.cells[i]

		if c1[0] != c2[0] || c1[1] != c2[1] {
			return false
		}
	}
	return true
}

func (p *GridPiece2) AddCell(c [2]int) *GridPiece2 {
	p.cells = append(p.cells, c)

	// Min cell is the most top-left cell, the 'origin' of the piece.
	if c[0] < p.position[0] {
		p.position[0] = c[0]
	}
	if c[1] < p.position[1] {
		p.position[1] = c[1]
	}

	return p
}

func (p *GridPiece2) FixPosition() {
	n := len(p.cells)
	for i := 0; i < n; i++ {
		cell := &p.cells[i]

		cell[0] -= p.position[0]
		cell[1] -= p.position[1]
	}
}

func (p *GridPiece2) SetPosition(row int, col int) {
	p.position[0] = row
	p.position[1] = col
}

func (p *GridPiece2) Move(m interface{}) {
	mov, ok := m.(*GridMov2)
	if ok {
		p.position[0] += mov.dRow
		p.position[1] += mov.dCol
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
