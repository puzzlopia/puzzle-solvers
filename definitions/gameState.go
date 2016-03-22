package defs

// Generic game state. Used as node in algorithm searchs
type GameState interface {

	//Unique identifier of the state instance
	Uid() int

	// Creates a copy of the state
	Clone() GameState

	// Compares two states, returns true if they are equal
	Equal(s GameState) bool

	// Generates an integer based on the state. Useful to insert in hashes.
	ToHash() int

	Print()

	// Used in BFS algorithms
	Depth() int
	SetPrevMov(Command)
	PrevMov() Command

	// Sometimes we can optimize if we avoid re-visiting states
	AddPrevMov(Command)

	// Returns whether the state is the root state of all explorations
	Initial() bool
}
