package defs

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

	SetMovChain([]Command, *GameState)

	CollapsedPathLen() int
	RealPathLen() int
	CopyMovChainFrom(GameState)
	//PropagateUpdate()
	AddNextState(GameState, Command)
	SetPrevState(GameState, Command)
	UpdateFromPrevState()
	CheckPathAndState()
	PrevState() GameState
	PrevMov() Command
	SamePieceMovedNext(Command) bool

	TinyPrint()
	TinyGoPrint()

	CopyMovChainAndAdd([]Command, Command, *GameState)

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
	AddEquivPath(GameState, []Command, Command)
	ValidMovement(m Command) bool
	ApplyEquivalencyContinuity(GameState, Command, GameState) bool
}

// A slice of game states
type GameStates []GameState

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
	Move(m Command) (err error)

	// Undoes the movement
	UndoMove(m Command) (err error)

	// Makes the playable game to put its internal parts to
	// reflect this state.
	SetState(GameState) (err error)

	// Returns a copy of current state
	State() (s GameState)

	//PiecesById() map[int]*grids.GridPiece2

	// For a fixed state, there is a concrete set of movements that can be done. But,
	// for the sake of algorithm performance, we can minimize that set by giving more
	// information, for example, to avoid any undo action, or any action that makes a
	// closed trajectory, etc.
	ValidMovementsBFS(pieceTrajectory []Command) []Command
}
