package defs

// A command represents a movement in a game. It can be a piece move, a number guess in sudoku, etc.
type Command interface {
	PieceId() int
	Inverted() interface{}
	IsInverse(otherMov interface{}) bool
	Print()
}
type SequenceMov []Command

type CmdStack struct {
	stack_ []Command
}

func (s *CmdStack) Path() []Command {
	return s.stack_
}

func (s *CmdStack) Push(m Command) {
	s.stack_ = append(s.stack_, m)
}
func (s *CmdStack) Pop() (Command, bool) {
	l := len(s.stack_)
	if l > 0 {
		m := s.stack_[l-1]
		s.stack_ = s.stack_[:l-1]
		return m, true
	}
	return nil, false
}
func (s *CmdStack) Last() Command {
	if len(s.stack_) > 0 {
		return s.stack_[len(s.stack_)-1]
	}
	return nil
}

// Returns number of movements.
func (s *CmdStack) MovMetric() int {
	movs := 0
	lastPieceId := 0
	for _, m := range s.stack_ {
		if lastPieceId != m.PieceId() {
			movs++
			lastPieceId = m.PieceId()
		}
	}
	return movs
}

// Makes a copy of the list of movements
func (s *CmdStack) Clone() []Command {
	c := make([]Command, len(s.stack_))
	copy(c, s.stack_)
	return c
}

func (s *CmdStack) LastPieceInvertedPath() []Command {
	var path []Command

	pieceId := 0
	for i := len(s.stack_) - 1; i >= 0; i-- {

		m := s.stack_[i]

		if pieceId == 0 {
			path = append(path, m)
			pieceId = m.PieceId()
		} else if pieceId == m.PieceId() {
			path = append(path, m)
		} else {
			return path
		}
	}

	return path

}

func (s *CmdStack) Reset() {
	s.stack_ = s.stack_[:0]
}
