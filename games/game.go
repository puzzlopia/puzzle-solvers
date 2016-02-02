package games

// Generic game state. Used as node in algorithm searchs
type GameState interface {

	//Unique identifier of the state instance
	Uid() int

	// Creates a copy of the state
	Clone() GameState

	// Compares two states, returns true if they are equal
	Equal(s GameState) bool

	// Compares two states, returns true if the param state is a substate. B is substate of A if all pieces in B are
	// placed at the same position of corresponding pieces in A, for example. It will depend of puzzle kind.
	EqualSub(s GameState) bool

	// Flags this state as equivalent to objective
	MarkAsObjective()
	IsObjective() bool

	// Generates an integer based on the state. Useful to insert in hashes.
	ToHash() int

	SetMovChain([]GameMov)

	CollapsedPathLen() int
	RealPathLen() int
	CopyMovChainFrom(GameState)
	//PropagateUpdate()
	AddNextState(GameState, GameMov)
	SetPrevState(GameState, GameMov)
	PrevState() GameState
	PrevMov() GameMov
	SamePieceMovedNext(GameMov) bool

	TinyPrint()

	copyMovChainAndAdd([]GameMov, GameMov)

	PathChain() []GameMov
	BuildPathReversed(path *[]GameMov)

	// BFS Algorithm
	SetWaiting(bool)
	Waiting() bool
	SetDepth(int)
	Depth() int

	// Adds an equivalent node-path. Finally if this node is part of a solution, we can check all the descendant paths to origin
	// and select the shortest.
	AddEquivPath(GameState, []GameMov, GameMov)
	ApplyEquivalencyContinuity(GameState, GameMov) bool
}

// A slice of game states
type GameStates []GameState

/**
 * @summary A game command or movement
 */
type GameMov interface {
	PieceId() int
	Inverted() interface{}
	IsInverse(otherMov interface{}) bool
	Print()
}
type SequenceMov []GameMov

type MovStack struct {
	stack_ []GameMov
}

func (s *MovStack) Path() []GameMov {
	return s.stack_
}

func (s *MovStack) Push(m GameMov) {
	s.stack_ = append(s.stack_, m)
}
func (s *MovStack) Pop() (GameMov, bool) {
	l := len(s.stack_)
	if l > 0 {
		m := s.stack_[l-1]
		s.stack_ = s.stack_[:l-1]
		return m, true
	}
	return nil, false
}
func (s *MovStack) Last() GameMov {
	if len(s.stack_) > 0 {
		return s.stack_[len(s.stack_)-1]
	}
	return nil
}

// Returns number of movements.
func (s *MovStack) MovMetric() int {
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
func (s *MovStack) Clone() []GameMov {
	c := make([]GameMov, len(s.stack_))
	copy(c, s.stack_)
	return c
}

func (s *MovStack) LastPieceInvertedPath() []GameMov {
	var path []GameMov

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

func (s *MovStack) Reset() {
	s.stack_ = s.stack_[:0]
}

// Basic game interface for building it
type GameDef interface {

	// Lets us to define the structure of the game
	Define(matrix GameState) (err error)

	// The game builds its internals
	Build() (err error)
}

// Playable: //útil para buscar máximas diferencias
type Playable interface {

	// Apply the movement
	Move(m GameMov) (err error)

	// Undoes the movement
	UndoMove(m GameMov) (err error)

	// Makes the playable game to put its internal parts to
	// reflect this state.
	SetState(GameState) (err error)

	// Returns a copy of current state
	State() (s GameState)

	// Returns current valid movements
	//ValidMovements(seq *[]GameMov, lastMov GameMov, curPieceTrajectory []GameMov)

	// For a fixed state, there is a concrete set of movements that can be done. But,
	// for the sake of algorithm performance, we can minimize that set by giving more
	// information, for example, to avoid any undo action, or any action that makes a
	// closed trajectory, etc.
	ValidMovementsBFS() []GameMov
}
