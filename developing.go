package main

//import "github.com/edgarweto/puzzlopia/puzzle-solvers/definitions"
import "github.com/edgarweto/puzzlopia/puzzle-solvers/finder"
import "github.com/edgarweto/puzzlopia/puzzle-solvers/games"
import "github.com/edgarweto/puzzlopia/puzzle-solvers/grids"

func doingTests() {

	// Define the game
	var myPuzzle = &games.SBGame{}

	// BUG!
	myPuzzle.Define(&grids.Matrix2d{
		// []int{1, 1, 2, 2},
		// []int{1, 1, 3, 4},
		// []int{7, 7, 3, 4},
		// []int{9, 8, 5, 0},
		// []int{8, 8, 6, 0},
		[]int{2, 2, 0, 0},
		[]int{2, 1, 0, 0},
		[]int{1, 1, 0, 0},
	})

	// // BrokenPennant: too much time!!
	// myPuzzle.Define(&grids.Matrix2d{
	// 	[]int{1, 2, 5, 5},
	// 	[]int{3, 4, 6, 6},
	// 	[]int{7, 8, 0, 0},
	// 	[]int{10, 11, 12, 12},
	// 	[]int{10, 11, 9, 9},
	// })

	// // Simple Sparse SBP
	// myPuzzle.Define(&grids.Matrix2d{
	// 	[]int{0, 0, 0},
	// 	[]int{0, 0, 0},
	// 	[]int{1, 2, 0},
	// })

	// // Pennant (http://www.puzzlopia.com/puzzles/pennant/play)
	// myPuzzle.Define(&grids.Matrix2d{
	// 	[]int{2, 2, 1, 1},
	// 	[]int{2, 2, 3, 3},
	// 	[]int{5, 4, 0, 0},
	// 	[]int{6, 7, 8, 8},
	// 	[]int{6, 7, 9, 9},
	// })

	// //Ninja II (http://www.puzzlopia.com/puzzles/ninja-ii/play)
	// // Requires about 106 independent steps!
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

	//myPuzzle.AutoAlikePieces()
	// myPuzzle.SetNotAlikePiece(1)
	// myPuzzle.SetNotAlikePiece(2)
	// myPuzzle.SetNotAlikePiece(3)
	// myPuzzle.SetNotAlikePiece(4)

	// //BrokenPennant
	// myPuzzle.AlikePieces([][]int{
	// 	[]int{1, 2, 3, 4},
	// 	[]int{7, 8},
	// 	[]int{5, 6, 9, 12},
	// 	[]int{10, 11},
	// })

	// Check the puzzle is well created, and let it build its internals
	myPuzzle.Build()

	// Params for finder/solver
	const (

		// Max depth reached by finder/solver
		MAX_DEPTH = 5

		// Max number of states to be processed. If 0, then ignored.
		// Can be combined with MAX_DEPTH: if either of these two values is exceeded, the algorithm stops.
		MAX_STATES = 99

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

	sbpFinder.Detect(&grids.Matrix2d{
		[]int{0, 0, 0, 0},
		[]int{0, 2, 2, 0},
		[]int{0, 2, 0, 0},
	})

	// // BrokenPennant
	// sbpFinder.Detect(&grids.Matrix2d{
	// 	[]int{0, 0, 0, 0},
	// 	[]int{0, 0, 0, 0},
	// 	[]int{0, 0, 0, 0},
	// 	[]int{1, 2, 0, 0},
	// 	[]int{3, 4, 0, 0},
	// })

	// // Pennant
	// sbpFinder.Detect(&grids.Matrix2d{
	// 	[]int{0, 0, 0, 0},
	// 	[]int{0, 0, 0, 0},
	// 	[]int{0, 0, 0, 0},
	// 	[]int{2, 2, 0, 0},
	// 	[]int{2, 2, 0, 0},
	// })

	// //Ninja II
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

	sbpFinder.SolvePuzzle(myPuzzle)

	//sbpFinder.FindExtremals(myPuzzle)
}
