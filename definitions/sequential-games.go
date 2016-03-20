package defs

// Generic game state. Used as node in algorithm searchs
type SeqGameState interface {

	//Unique identifier of the state instance
	Uid() int

	// Creates a copy of the state
	Clone() SeqGameState

	// Compares two states, returns true if they are equal
	Equal(s SeqGameState) bool

	// Compares two states, returns true if the param state is a substate. B is substate of A if all pieces in B are
	// placed at the same position of corresponding pieces in A, for example. It will depend of puzzle kind.
	EqualSub(s SeqGameState) bool

	// Flags this state as equivalent to objective
	MarkAsObjective()
	IsObjective() bool

	// Generates an integer based on the state. Useful to insert in hashes.
	ToHash() int

	// These functions have been being added while developing the BFS algorithm
	// Should refactor.
	SetMovChain([]Command, *SeqGameState)
	CollapsedPathLen() int
	RealPathLen() int
	CopyMovChainFrom(SeqGameState)
	SetPrevState(SeqGameState, Command)
	UpdateFromPrevState()
	CheckPathAndState()
	PrevState() SeqGameState
	PrevMov() Command
	UpdateFromStart(originState *SeqGameState)

	CopyMovChainAndAdd([]Command, Command, *SeqGameState)

	PathChain() []Command
	BuildPathReversed(path *[]Command)

	// BFS Algorithm
	SetWaiting(bool)
	Waiting() bool
	SetDepth(int)
	Depth() int
	MarkToDebug()
	MarkedToDebug() bool

	// Adds an equivalent node-path. Finally if this node is part of a solution, we can check all the descendant paths to origin
	// and select the shortest.
	AddEquivPath(SeqGameState, []Command, Command)
	ValidMovement(m Command) bool
	ApplyEquivalencyContinuity(SeqGameState, Command, SeqGameState) bool

	TinyPrint()
	TinyGoPrint()
}

// A slice of game states
type SeqGameStates []SeqGameState

// Basic game interface for building it
type GameDef interface {

	// Lets us to define the structure of the game
	Define(matrix SeqGameState) (err error)

	// The game builds its internals
	Build() (err error)
}

// Playable: //útil para buscar máximas diferencias
type Playable interface {

	// Apply the movement
	Move(m Command) (err error)

	// Undoes the movement
	UndoMove(m Command) (err error)

	// Makes the playable game to put its internal parts to
	// reflect this state.
	SetState(SeqGameState) (err error)

	// Returns a copy of current state
	State() (s SeqGameState)

	//PiecesById() map[int]*grids.GridPiece2

	// For a fixed state, there is a concrete set of movements that can be done. But,
	// for the sake of algorithm performance, we can minimize that set by giving more
	// information, for example, to avoid any undo action, or any action that makes a
	// closed trajectory, etc.
	ValidMovementsBFS(pieceTrajectory []Command) []Command
}
