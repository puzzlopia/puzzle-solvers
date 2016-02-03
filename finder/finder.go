package finder

import "github.com/edgarweto/puzzlopia/puzzle-solvers/games"

// Basic finder interface
type Finder interface {
	SetLimits(maxDepth int, maxStates int)
	SetDebug(bool)
	SilentMode(bool)

	// Searches for the farthest states from current state
	SolvePuzzle(g games.Playable, extremals games.GameStates)
}

// Finder limits: depth and number of different states
type FinderLimits struct {
	maxDepth_  int
	maxStates_ int
}

func (f *FinderLimits) SetLimits(maxDepth int, maxStates int) {
	f.maxDepth_ = maxDepth
	f.maxStates_ = maxStates
}
