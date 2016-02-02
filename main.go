package main

import "fmt"

import "github.com/edgarweto/puzzlopia/solvers/finder"
import "github.com/edgarweto/puzzlopia/solvers/games"
import "github.com/edgarweto/puzzlopia/solvers/grids"

func main() {

	// Define the game
	var myPuzzle = &games.SBGame{}

	// Pennant
	myPuzzle.Define(&grids.Matrix2d{
		[]int{2, 2, 1, 1},
		[]int{2, 2, 3, 3},
		[]int{5, 4, 0, 0},
		[]int{6, 7, 8, 8},
		[]int{6, 7, 9, 9},
	})

	// myPuzzle.Define(&grids.Matrix2d{
	// 	[]int{1, 1, 1, 1, 2, 2, 3, 3},
	// 	[]int{1, 1, 1, 1, 2, 2, 3, 3},
	// 	[]int{1, 1, 1, 1, 2, 2, 3, 3},
	// 	[]int{4, 4, 5, 5, 0, 6, 6, 6},
	// 	[]int{4, 4, 5, 5, 0, 6, 6, 6},
	// 	[]int{0, 0, 0, 0, 7, 7, 7, 0},
	// 	[]int{0, 8, 8, 8, 7, 7, 7, 0},
	// 	[]int{0, 8, 8, 8, 0, 9, 9, 9},
	// 	[]int{0, 0, 0, 0, 0, 9, 9, 9},
	// })

	myPuzzle.AutoAlikePieces()
	//myPuzzle.SetNotAlikePiece(5)

	// Check the puzzle is well created, and let it build its internals
	myPuzzle.Build()

	// Params for finder/solver
	const (

		// Max depth reached by finder/solver
		MAX_DEPTH = 84

		// Max number of states to be processed. If 0, then ignored.
		// Can be combined with MAX_DEPTH: if either of these two values is exceeded, the algorithm stops.
		MAX_STATES = 0

		// (Experimental) Used internally to force the algorithm to revisit some states
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

	// Pennant
	sbpFinder.Detect(&grids.Matrix2d{
		[]int{0, 0, 0, 0},
		[]int{0, 0, 0, 0},
		[]int{0, 0, 0, 0},
		[]int{2, 2, 0, 0},
		[]int{2, 2, 0, 0},
	})

	// //Gauntlet2
	// sbpFinder.Detect(&grids.Matrix2d{
	// 	[]int{0, 0, 0, 0, 0, 0, 0, 0},
	// 	[]int{0, 0, 0, 0, 0, 0, 0, 0},
	// 	[]int{0, 0, 0, 0, 0, 0, 0, 0},
	// 	[]int{0, 0, 0, 0, 0, 0, 0, 0},
	// 	[]int{0, 0, 0, 0, 0, 0, 0, 0},
	// 	[]int{0, 0, 0, 0, 0, 0, 0, 0},
	// 	[]int{0, 0, 0, 0, 1, 1, 1, 1},
	// 	[]int{0, 0, 0, 0, 1, 1, 1, 1},
	// 	[]int{0, 0, 0, 0, 1, 1, 1, 1},
	// })

	extremals := games.GameStates{}

	sbpFinder.SolvePuzzle(myPuzzle, &extremals)

	fmt.Println("\n[Extremal states]")
	//extremals.Print()
}
