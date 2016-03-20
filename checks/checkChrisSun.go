package checks

import "fmt"
import "github.com/edgarweto/puzzlopia/puzzle-solvers/finder"
import "github.com/edgarweto/puzzlopia/puzzle-solvers/games"
import "github.com/edgarweto/puzzlopia/puzzle-solvers/grids"

// Result: 90
// Should be: 90
// Currently finds a solution of 67 moves, but all pages refer to optimal solution of 200. May the puzzle config be wrong? Or it is the puzzle objective?
func CheckChrisSun() {

	// Define the game
	var myPuzzle = &games.SBGame{}

	// SuperCompo
	myPuzzle.Define(&grids.Matrix2d{

		[]int{1, 1, 0, 0},
		[]int{1, 1, 6, 10},
		[]int{2, 8, 6, 5},
		[]int{2, 3, 3, 5},
		[]int{9, 4, 4, 7},
	})
	myPuzzle.AutoAlikePieces()

	// Check the puzzle is well created, and let it build its internals
	myPuzzle.Build()

	// Params for finder/solver
	const (

		// Max depth reached by finder/solver
		MAX_DEPTH = 200

		// Max number of states to be processed. If 0, then ignored.
		// Can be combined with MAX_DEPTH: if either of these two values is exceeded, the algorithm stops.
		//MAX_STATES = 4999999
		MAX_STATES = 250000

		// (Experimental),Used internally to force the algorithm to revisit some states
		// Actually, disabling it makes Pennant to be solved with non-optimal path.
		HARD_OPTIMAL = true

		// Used for tests: enables/disables console output
		SILENT_MODE = false

		// Enables/disables debug options (console output, etc.)
		DEBUG = false
	)

	// FINDER ---------------------
	var sbpFinder finder.SbpBfsFinder

	sbpFinder.SilentMode(SILENT_MODE)
	sbpFinder.SetDebug(DEBUG)
	sbpFinder.SetLimits(MAX_DEPTH, MAX_STATES)
	sbpFinder.SetHardOptimal(HARD_OPTIMAL)

	// Objective
	sbpFinder.Detect(&grids.Matrix2d{
		[]int{0, 0, 0, 0},
		[]int{0, 0, 0, 0},
		[]int{0, 0, 0, 0},
		[]int{0, 1, 1, 0},
		[]int{0, 1, 1, 0},
	})

	sbpFinder.SolvePuzzle(myPuzzle)

	found, solutionLen, _ := sbpFinder.GetResult()

	if !found {
		fmt.Println("Chris-Sun not solved!")
	}
	if solutionLen != 90 {
		fmt.Printf("Chris-Sun solution not optimal: found len = %d\n\n", solutionLen)
	}
}
