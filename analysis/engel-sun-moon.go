package analysis

import "fmt"

import "github.com/edgarweto/puzzlopia/puzzle-solvers/finder"
import "github.com/edgarweto/puzzlopia/puzzle-solvers/games/engel"

func AnalyzeEngelSunMoon() {

	fmt.Println("Analyze Engel's SUN-MOON:")

	var colorWheels = &engel.EngelGame{}

	// Define the piece ids for each wheel
	colorWheels.Define([12]int{7, 6, 13, 14, 15, 16, 17, 18, 19, 20, 21, 8}, [12]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12})

	colorWheels.DefineIntersectionPositions([3]int{11, 0, 1}, [3]int{7, 6, 5})

	// Identify similar pieces
	colorWheels.SetAlike([][]int{

		// Right wheel is full colour, so there are two groups (rectangles and trias)
		[]int{1, 3, 5, 7, 9, 11},
		[]int{2, 4, 6, 8, 10, 12},

		// Left wheel with remaining pieces:
		[]int{13, 15, 17, 19, 21},
		[]int{14, 16, 18, 20},
	})

	const (

		// Max depth reached by finder/solver
		//MAX_DEPTH = 14
		MAX_DEPTH = 30

		// Max number of states to be processed. If 0, then ignored.
		// Can be combined with MAX_DEPTH: if either of these two values is exceeded, the algorithm stops.
		MAX_STATES = 100000

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
