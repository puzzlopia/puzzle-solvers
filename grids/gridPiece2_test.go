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
			t.Errorf("GridPiece2::Equivalence 1 failed", p.cells_, q.cells_)
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
		t.Errorf("GridPiece2::Equivalence 2 failed", p.cells_, q.cells_)
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
		t.Errorf("GridPiece2::Equivalence 3 failed", p.cells_, q.cells_)
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

// Tests if 'GeneratePieces' and 'DetectAlikePieces' works as expected with 'strange' pieces.
func TestPositioningPieces(t *testing.T) {

	matrix := &Matrix2d{
		[]int{0, 2, 1, 3, 0},
		[]int{2, 2, 3, 3, 0},
		[]int{2, 3, 3, 1, 0},
		[]int{0, 0, 4, 0, 0},
		[]int{4, 4, 4, 4, 0},
		[]int{0, 0, 0, 0, 0},
	}

	var pieces []*GridPiece2
	var alikes [][]int
	matrix.GeneratePieces(&pieces, alikes)

	if len(pieces) != 4 {
		t.Errorf("[Matrix2d::GeneratePieces] Not creating pieces as expected!")
	} else {
		piecesById := make(map[int]*GridPiece2)
		for _, p := range pieces {
			piecesById[p.Id()] = p
		}

		if !piecesById[1].IsAt(0, 2) || !piecesById[2].IsAt(0, 0) || !piecesById[3].IsAt(0, 1) || !piecesById[4].IsAt(3, 0) {
			t.Errorf("[Matrix2d::GeneratePieces] Generated pieces not well positionated!")
		}

		// Now, clear all pieces to check new positions:
		for _, p := range pieces {
			matrix.ClearPiece(p)
		}

		// For each piece: move to right, move down, move left and move up
		mRight := GridMov2{0, 0, 1}
		for _, p := range pieces {
			p.Move(&mRight)
		}
		if !piecesById[1].IsAt(0, 3) || !piecesById[2].IsAt(0, 1) || !piecesById[3].IsAt(0, 2) || !piecesById[4].IsAt(3, 1) {
			t.Errorf("[Matrix2d::GeneratePieces] Move right not working!")
		}

		mDown := GridMov2{0, 1, 0}
		for _, p := range pieces {
			p.Move(&mDown)
		}
		if !piecesById[1].IsAt(1, 3) || !piecesById[2].IsAt(1, 1) || !piecesById[3].IsAt(1, 2) || !piecesById[4].IsAt(4, 1) {
			t.Errorf("[Matrix2d::GeneratePieces] Move down not working!")
		}

		mLeft := GridMov2{0, 0, -1}
		for _, p := range pieces {
			p.Move(&mLeft)
		}
		if !piecesById[1].IsAt(1, 2) || !piecesById[2].IsAt(1, 0) || !piecesById[3].IsAt(1, 1) || !piecesById[4].IsAt(4, 0) {
			t.Errorf("[Matrix2d::GeneratePieces] Move left not working!")
		}

		mUp := GridMov2{0, -1, 0}
		for _, p := range pieces {
			p.Move(&mUp)
		}
		if !piecesById[1].IsAt(0, 2) || !piecesById[2].IsAt(0, 0) || !piecesById[3].IsAt(0, 1) || !piecesById[4].IsAt(3, 0) {
			t.Errorf("[Matrix2d::GeneratePieces] Move up not working!")
		}

		// Finally, move down and right and place all:
		for _, p := range pieces {
			p.Move(&mDown)
			p.Move(&mRight)
			matrix.PlacePiece(p)
		}

		// And detect where the pieces are:
		matrix.UpdatePiecePositions(piecesById)

		// And check they are well positioned:
		if !piecesById[1].IsAt(1, 3) || !piecesById[2].IsAt(1, 1) || !piecesById[3].IsAt(1, 2) || !piecesById[4].IsAt(4, 1) {
			t.Errorf("[Matrix2d::UpdatePiecePositions] not working!")
		}
	}
}

// func printPieces(pieces []*GridPiece2) {
// 	for _, p := range pieces {
// 		fmt.Printf("\nPiece[%d]: val=%d pos=(%d,%d), tl=(%d,%d)", p.Id(), p.Value(), p.position_[0], p.position_[1], p.tl_[0], p.tl_[1])
// 	}
// }
