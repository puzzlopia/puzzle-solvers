package analysis

import "fmt"

import "github.com/edgarweto/puzzlopia/puzzle-solvers/finder"
import "github.com/edgarweto/puzzlopia/puzzle-solvers/games/engel"

//import "github.com/edgarweto/puzzlopia/puzzle-solvers/grids"

func AnalyzeEngelColorWheels() {

	fmt.Println("Analyze Engel's COLOR WHEELS:")
	//var puzzle = &games.EngelGame{}

	var colorWheels = &engel.EngelGame{}

	colorWheels.Define([12]int{7, 6, 13, 14, 15, 16, 17, 18, 19, 20, 21, 8}, [12]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12})

	colorWheels.DefineIntersectionPositions([3]int{11, 0, 1}, [3]int{7, 6, 5})

	// Identify similar pieces
	colorWheels.SetAlike([][]int{

		// The six rectangles
		[]int{1, 5, 9, 13, 17, 21},

		// Triangles in the same group
		[]int{2, 4},
		[]int{6, 8},
		[]int{10, 12},
		[]int{14, 16},
		[]int{18, 20},
	})

	const (

		// Max depth reached by finder/solver
		MAX_DEPTH = 30

		// Max number of states to be processed. If 0, then ignored.
		// Can be combined with MAX_DEPTH: if either of these two values is exceeded, the algorithm stops.
		MAX_STATES = 60000000

		// Used for tests: enables/disables console output
		SILENT_MODE = false

		// Enables/disables debug options (console output, etc.)
		DEBUG = false
	)

	var analyzer finder.Analyzer

	analyzer.SetDebug(DEBUG)
	analyzer.SetLimits(MAX_DEPTH, MAX_STATES)

	analyzer.Explore(colorWheels)

	analyzer.Resume()
}
