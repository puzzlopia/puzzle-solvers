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
}
