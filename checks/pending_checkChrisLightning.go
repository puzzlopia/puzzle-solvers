package checks

import "fmt"
import "github.com/edgarweto/puzzlopia/puzzle-solvers/finder"
import "github.com/edgarweto/puzzlopia/puzzle-solvers/games"
import "github.com/edgarweto/puzzlopia/puzzle-solvers/grids"

// Result: 105
// Should be: 104
func CheckChrisLightning() {

	// Define the game
	var myPuzzle = &games.SBGame{}

	// SuperCompo
	myPuzzle.Define(&grids.Matrix2d{

		[]int{0, 1, 1, 0},
		[]int{2, 1, 1, 10},
		[]int{2, 9, 3, 3},
		[]int{4, 4, 7, 6},
		[]int{8, 5, 5, 6},
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
		fmt.Println("Chris-Lightning not solved!")
	}
	if solutionLen != 104 {
		fmt.Printf("Chris-Lightning solution not optimal: found len = %d\n\n", solutionLen)
	}
}