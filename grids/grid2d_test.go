package grids

//import "fmt"
import "testing"

func TestCopy(t *testing.T) {

	m1 := &Matrix2d{
		[]int{1, 1, 1},
		[]int{1, 1, 1},
	}

	var m2 Matrix2d

	m2.Copy(m1)

	if m1.Rows() != m2.Rows() || m1.Cols() != m2.Cols() {
		t.Errorf("Grid2::Copy 1 failed", m1, m2)
	}

	m2.SetAt(0, 0, 3)

	if m1.At(0, 0) != 1 {
		t.Errorf("Grid2::Copy 2 failed", m1, m2)
	}
}
