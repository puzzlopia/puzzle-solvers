package main

import "github.com/edgarweto/puzzlopia/puzzle-solvers/finder"
import "github.com/edgarweto/puzzlopia/puzzle-solvers/games"
import "github.com/edgarweto/puzzlopia/puzzle-solvers/grids"

func doingTests() {

	// Define the game
	var myPuzzle = &games.SBGame{}

	// // BrokenPennant: too much time!!
	// myPuzzle.Define(&grids.Matrix2d{
	// 	[]int{1, 2, 5, 5},
	// 	[]int{3, 4, 6, 6},
	// 	[]int{7, 8, 0, 0},
	// 	[]int{10, 11, 12, 12},
	// 	[]int{10, 11, 9, 9},
	// })

	// // Simplicity2 (http://www.mathpuzzle.com/MAA/31-Sliding%20Block%20Puzzles/mathgames_12_13_04.html)
	// myPuzzle.Define(&grids.Matrix2d{
	// 	[]int{1, 1, 2, 3},
	// 	[]int{4, 5, 10, 3},
	// 	[]int{4, 5, 6, 6},
	// 	[]int{8, 8, 7, 7},
	// 	[]int{9, 0, 0, 7},
	// })
	// myPuzzle.AutoAlikePieces()
	// myPuzzle.SetNotAlikePiece(1)

	// Super-Century
	myPuzzle.Define(&grids.Matrix2d{

		[]int{4, 6, 1, 2},
		[]int{4, 3, 10, 10},
		[]int{5, 3, 10, 10},
		[]int{5, 8, 8, 7},
		[]int{0, 0, 9, 9},
	})
	myPuzzle.AutoAlikePieces()

	// // SuperCompo
	// myPuzzle.Define(&grids.Matrix2d{

	// 	[]int{0, 1, 1, 0},
	// 	[]int{9, 1, 1, 10},
	// 	[]int{2, 3, 4, 6},
	// 	[]int{2, 3, 4, 6},
	// 	[]int{7, 5, 5, 8},
	// })
	// myPuzzle.AutoAlikePieces()
	//myPuzzle.SetNotAlikePiece(5)

	// // Neil's puzzles
	// myPuzzle.Define(&grids.Matrix2d{
	// 	[]int{1, 1, 2, 3},
	// 	[]int{4, 5, 10, 3},
	// 	[]int{4, 5, 6, 6},
	// 	[]int{8, 8, 7, 7},
	// 	[]int{9, 0, 0, 7},
	// })
	// myPuzzle.AutoAlikePieces()
	// myPuzzle.SetNotAlikePiece(1)

	// // BUG------------------------------------------Toulouzas's puzzles
	// myPuzzle.Define(&grids.Matrix2d{
	// 	[]int{1, 0, 0, 2, 4, 3},
	// 	[]int{1, 5, 0, 2, 0, 3},
	// 	[]int{6, 7, 8, 9, 10, 11},
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
	// 	[]int{11, 11, 0, 0, 7, 7, 7, 0},
	// 	[]int{11, 11, 0, 0, 7, 7, 7, 0},
	// 	[]int{10, 10, 8, 8, 8, 9, 9, 9},
	// 	[]int{10, 10, 8, 8, 8, 9, 9, 9},
	// })
	// myPuzzle.AlikePieces([][]int{
	// 	[]int{8, 9},
	// 	[]int{1, 2, 3},
	// 	[]int{4, 5},
	// })

	//myPuzzle.AutoAlikePieces()
	//myPuzzle.SetNotAlikePiece(1)
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
		MAX_DEPTH = 350

		// Max number of states to be processed. If 0, then ignored.
		// Can be combined with MAX_DEPTH: if either of these two values is exceeded, the algorithm stops.
		//MAX_STATES = 4999999
		MAX_STATES = 1999999

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

	// // BrokenPennant
	sbpFinder.Detect(&grids.Matrix2d{
		[]int{0, 0, 0, 0},
		[]int{0, 0, 0, 0},
		[]int{0, 0, 0, 0},
		[]int{0, 10, 10, 0},
		[]int{0, 10, 10, 0},
	})

	// sbpFinder.DebugPath([][]int{

	// 	[]int{9, -1, 0},
	// 	[]int{2, -1, 0},
	// 	[]int{7, -1, 0},
	// 	[]int{5, 0, -1},
	// 	[]int{8, 0, -1},
	// 	[]int{6, 1, 0},
	// 	[]int{10, 1, 0},
	// 	[]int{1, 0, 1},
	// 	[]int{9, 0, 1},
	// 	[]int{2, -1, 0},
	// 	[]int{3, 0, -1},
	// 	[]int{10, 0, -1},
	// 	[]int{6, -1, 0},
	// 	[]int{8, 0, 1},
	// 	[]int{5, 0, 1},
	// 	[]int{7, 1, 0},
	// 	[]int{4, 0, -1},
	// 	[]int{10, 1, 0},
	// 	[]int{3, 0, 1},
	// 	[]int{2, 1, 0},
	// 	[]int{9, 0, -1},
	// 	[]int{1, 0, -1},
	// 	[]int{6, -1, 0},
	// 	[]int{6, -1, 0},
	// 	[]int{10, 0, 1},
	// 	[]int{3, 0, 1},
	// 	[]int{4, 0, 1},
	// 	[]int{7, -1, 0},
	// 	[]int{5, 0, -1},

	// 	[]int{8, 0, -1},
	// 	[]int{10, 1, 0},
	// 	[]int{4, 0, 1},
	// 	[]int{7, 0, 1},
	// 	[]int{7, -1, 0},
	// 	[]int{5, -1, 0},
	// 	[]int{8, 0, -1},
	// 	[]int{8, 0, -1},
	// 	[]int{10, 0, -1},
	// 	[]int{10, 0, -1},
	// 	[]int{4, 1, 0},
	// 	[]int{5, 0, 1},
	// 	[]int{5, 0, 1},

	// 	[]int{8, -1, 0},
	// 	[]int{8, 0, 1},
	// 	[]int{2, 1, 0},
	// 	[]int{2, 1, 0},
	// 	[]int{7, 0, -1},
	// 	[]int{7, -1, 0},
	// 	[]int{2, -1, 0},
	// 	[]int{8, -1, 0},
	// 	[]int{10, -1, 0},
	// 	[]int{4, 0, -1},
	// 	[]int{4, 0, -1},
	// 	[]int{5, 1, 0},
	// 	[]int{10, 0, 1},
	// 	[]int{10, 0, 1},

	// 	[]int{8, 1, 0},
	// 	[]int{3, 0, -1},
	// 	[]int{10, -1, 0},
	// 	[]int{5, -1, 0},
	// 	[]int{4, 0, 1},
	// 	[]int{4, 0, 1},
	// 	[]int{2, 1, 0},
	// 	[]int{3, 0, -1},
	// 	[]int{10, 0, -1},
	// 	[]int{8, 1, 0},

	// 	[]int{5, 0, -1},
	// 	[]int{6, 1, 0},
	// 	[]int{6, 1, 0},
	// 	[]int{1, 0, 1},
	// 	[]int{7, 0, 1},
	// 	[]int{7, -1, 0},
	// 	[]int{3, -1, 0},
	// 	[]int{2, -1, 0},
	// })

	sbpFinder.SolvePuzzle(myPuzzle)

	//sbpFinder.FindExtremals(myPuzzle)
}
