package grids

//import "fmt"
import "testing"

func TestEquivalence(t *testing.T) {

	for j := 0; j < 10; j++ {
		p := NewGridPiece2(1, 0)
		q := NewGridPiece2(2, 0)

		for i := 0; i < 8; i++ {
			p.AddCell([2]int{(j + 1) * i, 5*i + (j + 7)})
			q.AddCell([2]int{(j+1)*i + (j+3)*(j+1), 5*i + (j + 7) + (j+3)*(j+1)})
		}
		p.FixPosition()
		q.FixPosition()

		if !p.Equivalent(q) {
			t.Errorf("GridPiece2::Equivalence 1 failed", p.cells, q.cells)
		}
	}

	// Same number of cells, different positions
	p := NewGridPiece2(1, 0)
	q := NewGridPiece2(2, 0)

	for i := 0; i < 8; i++ {
		p.AddCell([2]int{8 * i, 5*i + 3})
		q.AddCell([2]int{8 * i, 5*i + 3 + i})
	}
	p.FixPosition()
	q.FixPosition()

	if p.Equivalent(q) {
		t.Errorf("GridPiece2::Equivalence 2 failed", p.cells, q.cells)
	}

	// Different number of cells
	p = NewGridPiece2(1, 0)
	q = NewGridPiece2(2, 0)

	for i := 0; i < 8; i++ {
		p.AddCell([2]int{8 * i, 5*i + 3})

		if i%2 == 0 {
			q.AddCell([2]int{8 * i, 5*i + 3})
		}
	}
	p.FixPosition()
	q.FixPosition()

	if p.Equivalent(q) {
		t.Errorf("GridPiece2::Equivalence 3 failed", p.cells, q.cells)
	}
}

// Tests if 'DetectAlikePieces' works as expected
func TestAlikePieces(t *testing.T) {
	minValue := 5

	pieces := make([]*GridPiece2, 4)

	pieces[0] = NewGridPiece2(1, 0)
	pieces[1] = NewGridPiece2(2, 0)
	pieces[2] = NewGridPiece2(3, 0)
	pieces[3] = NewGridPiece2(4, 0)

	pieces[0].AddCell([2]int{0, 0})
	pieces[0].AddCell([2]int{0, 1})
	pieces[0].AddCell([2]int{1, 0})
	pieces[0].FixPosition()

	pieces[1].AddCell([2]int{3, 3})
	pieces[1].AddCell([2]int{3, 4})
	pieces[1].AddCell([2]int{4, 3})
	pieces[1].FixPosition()

	pieces[2].AddCell([2]int{0, 0})
	pieces[2].AddCell([2]int{0, 1})
	pieces[2].FixPosition()

	pieces[3].AddCell([2]int{6, 6})
	pieces[3].AddCell([2]int{6, 7})
	pieces[3].FixPosition()

	DetectAlikePieces(pieces, minValue)

	if pieces[0].Value() != pieces[1].Value() {
		t.Errorf("::DetectAlikePieces 1 failed, should be equivalent!")
	}
	if pieces[2].Value() != pieces[3].Value() {
		t.Errorf("::DetectAlikePieces 2 failed, should be equivalent!")
	}
	if pieces[0].Value() == pieces[2].Value() {
		t.Errorf("::DetectAlikePieces 3 failed, should be different!")
	}
}
