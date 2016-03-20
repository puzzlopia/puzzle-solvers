package finder

import "fmt"
import "time"
import "github.com/fatih/color"

import "github.com/edgarweto/puzzlopia/puzzle-solvers/utils"
import "github.com/edgarweto/puzzlopia/puzzle-solvers/definitions"

// Sliding blocks puzzle finder
type Analyzer struct {

	// Params
	limits_       FinderLimits
	silent_       bool
	hardOptimals_ bool

	// Game settings
	game_      defs.Explorable
	initState_ defs.GameState

	// If we are searching for a concrete state
	extremals_     []defs.GameState
	extremalDist_  int
	findExtremals_ bool

	// Stats
	countStates_  utils.ScalarStatistic
	nodesDegree_  utils.ScalarStatistic
	frontierSize_ utils.RangeStatistic
	maxDepth_     utils.RangeStatistic
	depthDistr_   utils.RangeHistogram

	// Algorithm state
	visitedStates_  map[int][]defs.GameState
	frontier_       utils.Queue
	farthestStates_ utils.Queue
	maxFarthest_    int
	nextFrontier_   []defs.GameState
	endStatus_      string
	duration_       time.Duration

	initialized_ bool
	debug_       bool
	nextDepth_   int

	fmtHeaders_ *color.Color
	outDbg1_    *color.Color //more important
	outDbg2_    *color.Color //less important
	outDbg3_    *color.Color //different
}

func (f *Analyzer) SetDebug(b bool) {
	f.debug_ = b
}

func (f *Analyzer) SetLimits(maxDepth int, maxStates int) {
	f.limits_.SetLimits(maxDepth, maxStates)
}
func (f *Analyzer) SilentMode(b bool) {
	f.silent_ = b
}

func (f *Analyzer) init() {
	f.fmtHeaders_ = color.New(color.FgCyan, color.Bold)

	f.outDbg1_ = color.New(color.FgCyan)
	f.outDbg2_ = color.New(color.FgWhite)
	f.outDbg3_ = color.New(color.FgYellow)

	f.countStates_.Set("States")
	f.nodesDegree_.Set("Node degree")
	f.frontierSize_.Set("Frontier size")
	f.maxDepth_.Set("Max depth")
	f.depthDistr_.Set("Depth states distribution")

	f.initialized_ = true
	f.nextDepth_ = 0
	f.maxFarthest_ = 12
}

// Prints statistics and results
func (f *Analyzer) Resume() {
	if !f.initialized_ {
		f.init()
	}

	f.fmtHeaders_.Println("\n[STATS]")
	f.countStates_.Resume(f.outDbg2_)
	f.nodesDegree_.ResumeAv(f.outDbg2_)
	f.frontierSize_.ResumeRange(f.outDbg2_)
	f.maxDepth_.ResumeRange(f.outDbg2_)
	f.depthDistr_.ResumeHistogram(f.outDbg2_)

	//Farthest states: print 3
	f.fmtHeaders_.Println("\n[FARTHEST STATES]")
	for i := 0; i < 3 && i < f.farthestStates_.Size(); i++ {
		x := f.farthestStates_.PopFront()
		if x != nil {
			s := x.(defs.GameState)

			f.outDbg2_.Printf("\n\n State [%d]:\n", i+1)
			s.Print()
		}
	}

	fmt.Println("\n\n")
}

// Explores all possible reachable states
func (f *Analyzer) Explore(g defs.Explorable) {
	if !f.initialized_ {
		f.init()
	}
	if !f.silent_ {
		fmt.Println("Puzzle Explorer v.1.0")
	}

	f.game_ = g
	f.visitedStates_ = make(map[int][]defs.GameState)
	f.initState_ = f.game_.State()

	h := f.initState_.ToHash()
	f.visitedStates_[h] = append(f.visitedStates_[h], f.initState_)
	f.addToFrontier(f.initState_)
	f.countStates_.Incr()

	tStart := time.Now()
	if !f.silent_ {
		f.fmtHeaders_.Println("\n[WORKING]...")
	}

	f.exploreTree()

	tEnd := time.Now()
	f.duration_ = tEnd.Sub(tStart)
	if !f.silent_ {
		f.fmtHeaders_.Println("\n[DONE] ", f.duration_)
		f.fmtHeaders_.Println("\n - End condition: ", f.endStatus_)
	}
}

func (f *Analyzer) exploreTree() {

	statesCount := 1

	curState := f.popFrontier()
	for curState != nil {

		statesCount++
		if f.limits_.maxStates_ > 0 && statesCount >= f.limits_.maxStates_ {
			f.endStatus_ = "Max states reached."
			break
		} else if curState.Depth() > f.limits_.maxDepth_ {
			f.endStatus_ = "Max depth reached."
			break
		}

		if !f.silent_ {
			if f.nextDepth_ < curState.Depth() {
				fmt.Printf("\n\n DEPTH %d", f.nextDepth_)
				fmt.Printf(" -------------------------------------")
				fmt.Printf("\n States explored: %d", statesCount-1)
				fmt.Printf("\n\n")
				f.nextDepth_ = curState.Depth()
			}

			update := f.limits_.maxStates_ / 100
			if f.limits_.maxStates_ > 0 && statesCount%update == 0 {
				pct := (100 * statesCount / f.limits_.maxStates_)
				fmt.Printf("\n  - Explored states pct: %d%%, frontier: %d", pct, f.frontier_.Size())
			}
		}

		f.game_.SetState(curState)
		validMovs := f.game_.ValidMovements()

		f.nodesDegree_.Add(len(validMovs))

		if f.debug_ {
			f.outDbg2_.Printf("	 Valid movs:%v", len(validMovs))
		}

		for _, mov := range validMovs {

			f.game_.Move(mov)
			newState := f.game_.State()

			newState.SetPrevMov(mov)

			f.processState(newState)

			f.game_.UndoMove(mov)
		}

		// Let's visite next pending state
		curState = f.popFrontier()
	}
	if curState == nil {
		f.endStatus_ = "All states explored, no more states in queue"
	}
}

func (f *Analyzer) processState(s defs.GameState) {

	h := s.ToHash()
	if f.debug_ {
		f.outDbg1_.Printf("\n	 - Process state. Hash: %d", h)
	}

	if f.visitedStates_[h] == nil {

		// Easy, it is a new state
		f.countStates_.Incr()
		f.visitedStates_[h] = append(f.visitedStates_[h], s)

		// Let's explore its childs later
		f.addToFrontier(s)
	} else {

		// Compare with potential equivalent states
		for _, st := range f.visitedStates_[h] {
			if st.Equal(s) {
				if s.Depth() < st.Depth() {
					if f.debug_ {
						f.outDbg1_.Printf("\n	 - Detected state with inferior depth!")
					}
				}
				return
			}
		}

		// Ok, new state with hash collision.
		f.countStates_.Incr()
		f.visitedStates_[h] = append(f.visitedStates_[h], s)

		// Let's explore its childs later
		f.addToFrontier(s)
	}
}

// Push back to the priority queue that state.
func (f *Analyzer) addToFrontier(s defs.GameState) {
	if f.debug_ {
		f.outDbg2_.Printf("\n	  Add to frontier: state [%d]", s.Uid())
	}
	f.frontier_.PushBack(s)
	f.frontierSize_.Add(f.frontier_.Size())
	f.maxDepth_.Add(s.Depth())

	f.depthDistr_.Add(s.Depth(), 1)

	// Farthest states are always the last added
	f.farthestStates_.PushBack(s)
	if f.farthestStates_.Size() > f.maxFarthest_ {
		f.farthestStates_.PopFront()
	}
}

// Pop from start of queue (highest priority)
func (f *Analyzer) popFrontier() defs.GameState {

	x := f.frontier_.PopFront()
	if x != nil {
		s := x.(defs.GameState)
		if f.debug_ {
			f.outDbg2_.Printf("\n\n### Explore state [%d]", s.Uid())
			s.Print()
		}
		return s
	}
	return nil
}
