package engel

import "fmt"
import "github.com/edgarweto/puzzlopia/puzzle-solvers/definitions"

var staticGameStateCount_ int = 0

type EngelWheel struct {

	// Wheels are defined with piece values, we don't need real piece pointers or piece real identifiers.
	// Piece values are piece identifiers for differentiated pieces, and other shared values for all interchangeable pieces.
	pieces_ [12]int

	intersectionPositions_ [3]int
}

func (w *EngelWheel) At(idx int) int {
	return w.pieces_[idx]
}

func (w *EngelWheel) Pieces() [12]int {
	return w.pieces_
}

func (w *EngelWheel) SetPieces(ids [12]int) {
	w.pieces_ = ids
}
func (w *EngelWheel) SetIntersection(poss [3]int) {
	w.intersectionPositions_ = poss
}

func (w *EngelWheel) Copy(v EngelWheel) {
	w.pieces_ = v.pieces_
	w.intersectionPositions_ = v.intersectionPositions_
}

func (w *EngelWheel) Rotate(r int) {
	t := w.pieces_
	for i := 0; i < 12; i++ {
		w.pieces_[(i+2*r)%12] = t[i]
	}
}

func (w *EngelWheel) Equal(v EngelWheel, pieceToVal *defs.PieceToValue) bool {
	for i := 0; i < 12; i++ {
		if pieceToVal.At(w.pieces_[i]) != pieceToVal.At(v.pieces_[i]) {
			return false
		}
	}
	return true
}

func (w *EngelWheel) UpdateIntersection(v EngelWheel) {
	w.pieces_[w.intersectionPositions_[0]] = v.pieces_[v.intersectionPositions_[0]]
	w.pieces_[w.intersectionPositions_[1]] = v.pieces_[v.intersectionPositions_[1]]
	w.pieces_[w.intersectionPositions_[2]] = v.pieces_[v.intersectionPositions_[2]]
}

// Sliding Blocks Puzzle Game State
type EngelState struct {
	// Graph structure
	uid_     int
	depth_   int
	prevMov_ defs.Command
	//prevMovs_ []defs.Command
	prevMovWheels_ [2]bool
	isInitial_     bool

	// Structure state
	wheelLeft_  EngelWheel
	wheelRight_ EngelWheel

	// Utils
	pieceToValue_ *defs.PieceToValue
}

func (s *EngelState) SetInitial() {
	s.isInitial_ = true
}
func (s *EngelState) Initial() bool {
	return s.isInitial_
}

func (s *EngelState) Init(leftIds [12]int, rightIds [12]int) {
	s.wheelLeft_.SetPieces(leftIds)
	s.wheelRight_.SetPieces(rightIds)
	s.prevMovWheels_ = [2]bool{false, false}
	s.depth_ = 0

	// Assign the map
	s.pieceToValue_ = defs.GetPieceToValueMap()
}

func (s *EngelState) DefineIntersectionPositions(leftPositions [3]int, rightPositions [3]int) {
	s.wheelLeft_.SetIntersection(leftPositions)
	s.wheelRight_.SetIntersection(rightPositions)
}

func (s *EngelState) Assign(e EngelState) {
	s.uid_ = e.uid_
	s.wheelLeft_.Copy(e.wheelLeft_)
	s.wheelRight_.Copy(e.wheelRight_)
	s.pieceToValue_ = e.pieceToValue_
	s.depth_ = e.depth_
	s.prevMov_ = e.prevMov_
	s.prevMovWheels_ = e.prevMovWheels_
}

// Interface for sequential game states:
func (g *EngelState) Uid() int {
	return g.uid_
}

func (s *EngelState) Clone() defs.GameState {
	var c EngelState

	staticGameStateCount_++
	c.uid_ = staticGameStateCount_
	c.wheelLeft_.Copy(s.wheelLeft_)
	c.wheelRight_.Copy(s.wheelRight_)
	c.pieceToValue_ = s.pieceToValue_
	c.depth_ = s.depth_ + 1
	c.prevMov_ = s.prevMov_
	c.prevMovWheels_ = s.prevMovWheels_

	return &c
}

func (s *EngelState) Equal(o defs.GameState) bool {
	if s.pieceToValue_ == nil {
		s.pieceToValue_ = defs.GetPieceToValueMap()
	}
	e := o.(*EngelState)
	if s.wheelLeft_.Equal(e.wheelLeft_, s.pieceToValue_) && s.wheelRight_.Equal(e.wheelRight_, s.pieceToValue_) {
		return true
	}
	return false
}

func (s *EngelState) ToHash() int {
	h := 1
	for idx, id := range s.wheelLeft_.Pieces() {
		h += idx * s.pieceToValue_.At(id)
	}
	for idx, id := range s.wheelRight_.Pieces() {
		h += 3 * idx * s.pieceToValue_.At(id)
	}
	return h
}

func (s *EngelState) Print() {
	fmt.Printf("[%v | %v]", s.wheelLeft_, s.wheelRight_)
}
func (s *EngelState) Depth() int {
	return s.depth_
}

func (s *EngelState) SetPrevMov(m defs.Command) {
	s.prevMov_ = m
}

func (s *EngelState) PrevMov() defs.Command {
	return s.prevMov_
}

func (s *EngelState) AddPrevMov(m defs.Command) {
	mov := m.(*EngelCommand)

	if mov.PieceId() == 0 {
		s.prevMovWheels_[0] = true
	} else {
		s.prevMovWheels_[1] = true
	}

	// for _, mov := range s.prevMovs_ {
	// 	if mov.Equals(m) {
	// 		return
	// 	}
	// }
	// s.prevMovs_ = append(s.prevMovs_, m)
}

// // Filters a command: returns true if the command would generate
// // a visited state.
// func (s *EngelState) IgnoreMovement(m EngelCommand) bool {
// 	for _, mov := range s.prevMovs_ {
// 		if mov.IsInverse(m) {
// 			return true
// 		}
// 	}
// 	return false
// }

// Returns true if the state has been reached by moving left wheel by one path, and right wheel by other path.
func (s *EngelState) CanIgnoreState() bool {
	return s.prevMovWheels_[0] && s.prevMovWheels_[1]
}

func (s *EngelState) Move(mov *EngelCommand) {
	wheelId := mov.PieceId()
	rotation := mov.Rotation()

	rotation = rotation % 6
	if rotation < 0 {
		rotation += 6
	} else if rotation == 0 {
		// Identity movement
		return
	}

	// Rotate wheel
	var rotWheel, staticWheel *EngelWheel
	if wheelId == 0 {
		rotWheel = &s.wheelLeft_
		staticWheel = &s.wheelRight_
	} else {
		rotWheel = &s.wheelRight_
		staticWheel = &s.wheelLeft_
	}

	// Cyclic movement
	rotWheel.Rotate(rotation)

	// Update non-rotated wheel common pieces
	staticWheel.UpdateIntersection(*rotWheel)
}
