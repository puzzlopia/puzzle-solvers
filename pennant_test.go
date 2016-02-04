package main

import "github.com/edgarweto/puzzlopia/puzzle-solvers/finder"
import "github.com/edgarweto/puzzlopia/puzzle-solvers/games"
import "github.com/edgarweto/puzzlopia/puzzle-solvers/grids"

import "testing"

func TestEquivalence(t *testing.T) {

	var myPuzzle = &games.SBGame{}

	// Pennant
	myPuzzle.Define(&grids.Matrix2d{
		[]int{2, 2, 1, 1},
		[]int{2, 2, 3, 3},
		[]int{5, 4, 0, 0},
		[]int{6, 7, 8, 8},
		[]int{6, 7, 9, 9},
	})

	myPuzzle.AutoAlikePieces()
	//myPuzzle.SetNotAlikePiece(5)

	// Check the puzzle is well created, and let it build its internals
	myPuzzle.Build()

	const (

		// Max depth reached by finder/solver
		MAX_DEPTH = 83

		// Max number of states to be processed. If 0, then ignored.
		// Can be combined with MAX_DEPTH: if either of these two values is exceeded, the algorithm stops.
		MAX_STATES = 0

		// (Experimental) Used internally to force the algorithm to revisit some states
		// Actually, disabling it makes Pennant to be solved with non-optimal path.
		HARD_OPTIMAL = true

		// Used for tests: enables/disables console output
		SILENT_MODE = true

		// Enables/disables debug options (console output, etc.)
		DEBUG = false
	)

	// FINDER ---------------------
	var sbpFinder finder.SbpBfsFinder

	sbpFinder.SilentMode(SILENT_MODE)
	sbpFinder.SetDebug(DEBUG)
	sbpFinder.SetLimits(MAX_DEPTH, MAX_STATES)
	sbpFinder.SetHardOptimal(HARD_OPTIMAL)

	// Pennant
	sbpFinder.Detect(&grids.Matrix2d{
		[]int{0, 0, 0, 0},
		[]int{0, 0, 0, 0},
		[]int{0, 0, 0, 0},
		[]int{2, 2, 0, 0},
		[]int{2, 2, 0, 0},
	})

	sbpFinder.SolvePuzzle(myPuzzle)

	found, solutionLen, duration := sbpFinder.GetResult()

	if !found {
		t.Errorf("Pennant not solved!")
	}
	if solutionLen != 59 {
		t.Errorf("Pennant solution not optimal: found len = %d", solutionLen)
	}
	if duration.Seconds() > 0.1 {
		t.Errorf("Should solve in less than 0.1 seconds. Current: %v", duration)
	}
}
