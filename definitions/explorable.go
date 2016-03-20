package defs

// Generic game state. Used as node in algorithm searchs
type Explorable interface {

	// Apply the movement
	Move(m Command) (err error)

	// Undoes the movement
	UndoMove(m Command) (err error)

	// Makes the playable game to put its internal parts to
	// reflect this state.
	SetState(GameState) (err error)

	// Returns a copy of current state
	State() (s GameState)

	// For a fixed state, there is a concrete set of movements that can be done. But,
	// for the sake of algorithm performance, we can minimize that set by giving more
	// information, for example, to avoid any undo action, or any action that makes a
	// closed trajectory, etc.
	ValidMovements() []Command
}
