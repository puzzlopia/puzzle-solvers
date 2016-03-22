package defs

// Command is a generic interface that represents a game movement. Puzzle games can be very different, so a command
// could be a sliding block piece move, or a number guess in sudoku, etc.
type Command interface {
	PieceId() int
	Inverted() interface{}
	IsInverse(otherMov interface{}) bool
	Print()
	Equals(otherMov interface{}) bool
}

type SequenceMov []Command

// Algorithm ad-hoc structure, a stack of commands (movements)
type CmdStack struct {
	stack_ []Command
}

func (s *CmdStack) Path() []Command {
	return s.stack_
}

// Adds a new movement to the top of the stack.
func (s *CmdStack) Push(m Command) {
	s.stack_ = append(s.stack_, m)
}

// Returns the top of the stack and removes it from stack.
func (s *CmdStack) Pop() (Command, bool) {
	l := len(s.stack_)
	if l > 0 {
		m := s.stack_[l-1]
		s.stack_ = s.stack_[:l-1]
		return m, true
	}
	return nil, false
}

// Returns last stack element (top). Does not modify the stack.
func (s *CmdStack) Last() Command {
	if len(s.stack_) > 0 {
		return s.stack_[len(s.stack_)-1]
	}
	return nil
}

// Returns number of movements/commands using 'move metric': two consecutive movements on the same piece count as 1 movement.
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

// Detects all movements of last piece command. Then returns that path.
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

// Resets the stack
func (s *CmdStack) Reset() {
	s.stack_ = s.stack_[:0]
}
