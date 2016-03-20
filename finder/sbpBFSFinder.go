package finder

import "fmt"
import "time"
import "github.com/fatih/color"

import "github.com/edgarweto/puzzlopia/puzzle-solvers/utils"
import "github.com/edgarweto/puzzlopia/puzzle-solvers/definitions"

import "github.com/edgarweto/puzzlopia/puzzle-solvers/games"

import "github.com/edgarweto/puzzlopia/puzzle-solvers/grids" //TO REMOVE DEPENDENCY

// Sliding blocks puzzle finder
type SbpBfsFinder struct {

	// Params
	limits_       FinderLimits
	silent_       bool
	hardOptimals_ bool

	// Game settings
	game_      defs.Playable
	initState_ defs.SeqGameState

	// If we are searching for a concrete state
	search_     *games.SBPState
	foundState_ *defs.SeqGameState

	// If we are searching for a concrete state
	extremals_     []defs.SeqGameState
	extremalDist_  int
	findExtremals_ bool

	// Stats
	countStates_  utils.ScalarStatistic
	nodesDegree_  utils.ScalarStatistic
	frontierSize_ utils.RangeStatistic

	// Algorithm state
	visitedStates_ map[int][]defs.SeqGameState
	frontier_      utils.Queue
	nextFrontier_  []defs.SeqGameState
	endStatus_     string
	duration_      time.Duration

	debug_         bool
	debugPath_     [][]int
	debugTemp_     bool //Activates debug for current state
	debugSteps_    int
	maxDebugSteps_ int
	fmtHeaders_    *color.Color

	//outWarning_ *color.Color
	//out_ *color.Color
	outDbg1_ *color.Color //more important
	outDbg2_ *color.Color //less important
	outDbg3_ *color.Color //different
}

func (f *SbpBfsFinder) SetDebug(b bool) {
	f.debug_ = b
}
func (f *SbpBfsFinder) SetHardOptimal(hardOptimal bool) {
	f.hardOptimals_ = hardOptimal
}
func (f *SbpBfsFinder) SetLimits(maxDepth int, maxStates int) {
	f.limits_.SetLimits(maxDepth, maxStates)
}
func (f *SbpBfsFinder) SilentMode(b bool) {
	f.silent_ = b
}

// We want to know the minimum path to this state
func (f *SbpBfsFinder) Detect(m *grids.Matrix2d) (err error) {

	f.search_ = &games.SBPState{}
	f.search_.Init(m)

	return nil
}

// Useful to run the algorithm silently until we reach a concrete path. Then, debug mode is
// activated for some amount of steps.
func (f *SbpBfsFinder) DebugPath(path [][]int) {
	f.debugPath_ = path
}

// Returns if found, and collapsed length of solution
func (f *SbpBfsFinder) GetResult() (found bool, cr int, dur time.Duration) {
	if f.foundState_ == nil {
		return false, 0, 0
	}
	return true, (*f.foundState_).CollapsedPathLen(), f.duration_
}

// Prints statistics and results
func (f *SbpBfsFinder) Resume() {

	f.fmtHeaders_.Println("\n - Condition: ", f.endStatus_)

	f.fmtHeaders_.Println("\n[STATS]")
	f.countStates_.Resume(f.outDbg2_)
	f.nodesDegree_.ResumeAv(f.outDbg2_)
	f.frontierSize_.ResumeRange(f.outDbg2_)

	if f.search_ != nil {
		f.fmtHeaders_.Println("\n\n[SOLUTION]\n")

		search := color.New(color.FgYellow, color.Bold)

		if f.foundState_ != nil {
			search.Println("Found! Path len: ", (*f.foundState_).CollapsedPathLen())

			(*f.foundState_).TinyPrint()
		} else {
			search.Println("Not found.")
		}
	}
	fmt.Println("\n\n")
}

/**
 * @summary Makes a depth-search and returns the most 'distant' states.
 *
 * @param {defs.Playable} g The sequential game
 * @param {int} maxSteps Max depth of the search
 * @param {int} maxExtremals Max number of extremal states returned
 */
func (f *SbpBfsFinder) SolvePuzzle(g defs.Playable) {

	if !f.silent_ {
		fmt.Println("Puzzle Finder v.1.0")
	}
	f.fmtHeaders_ = color.New(color.FgCyan, color.Bold)

	f.outDbg1_ = color.New(color.FgCyan)
	f.outDbg2_ = color.New(color.FgWhite)
	f.outDbg3_ = color.New(color.FgYellow)

	f.countStates_.Set("States")
	f.nodesDegree_.Set("Node degree")
	f.frontierSize_.Set("Frontier size")

	f.game_ = g
	f.visitedStates_ = make(map[int][]defs.SeqGameState)
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
	}

	if !f.silent_ {
		f.Resume()
	}
}

// Implements the BFS algorithm. Uses a priority queue to save pending nodes to be visited; explores one node at a time.
func (f *SbpBfsFinder) exploreTree() {

	statesCount := 1

	curState := f.popFrontier()
	for curState != nil {

		statesCount++
		if f.limits_.maxStates_ > 0 && statesCount >= f.limits_.maxStates_ {
			f.endStatus_ = "Max states reached."
			break
		} else if curState.Depth() >= f.limits_.maxDepth_ {
			f.endStatus_ = "Max depth reached."
			break
		}

		if !f.silent_ {
			update := f.limits_.maxStates_ / 10
			if f.limits_.maxStates_ > 0 && statesCount%update == 0 {
				pct := (100 * statesCount / f.limits_.maxStates_)
				fmt.Printf("\n%d%%", pct)
			}
		}
		var reversePath []defs.Command
		curState.BuildPathReversed(&reversePath)

		// Debug detect path:
		if f.debugTemp_ && f.debugPath_ != nil {
			f.debug_ = false
		}

		if f.debugPath_ != nil {
			//curPath := grids.GridPath2{reversePath}
			curPath := grids.PathFromSlice(reversePath)
			if curPath.IsEquivalent(f.debugPath_, true) {
				// Ok, activate temporal debug:
				f.debug_ = true
				f.debugSteps_ = 0
				f.maxDebugSteps_ = 30
			}
		}
		if curState.MarkedToDebug() {
			f.debug_ = true
		}
		if f.debug_ {
			f.debugSteps_++
		}
		if f.debugTemp_ && f.debugPath_ != nil && f.debugSteps_ >= f.maxDebugSteps_ {
			f.debug_ = false
		}

		var pieceTrajectory grids.GridPath2
		pieceTrajectory.BuildFromReversePath(reversePath)

		f.game_.SetState(curState)
		validMovs := f.game_.ValidMovementsBFS(pieceTrajectory.Path())

		f.nodesDegree_.Add(len(validMovs))

		if f.debug_ {
			f.outDbg2_.Printf("	 Valid movs:%v", len(validMovs))
		}

		frontierSize := f.frontier_.Size()

		for _, mov := range validMovs {

			f.game_.Move(mov)
			newState := f.game_.State()
			newState.SetPrevState(curState, mov)

			f.processState(newState, reversePath, mov)

			f.game_.UndoMove(mov)
		}

		// Detect extremal states:
		if f.findExtremals_ && f.frontier_.Size() == frontierSize {
			f.addExtremalState(curState)
		}

		// Let's visite next pending state
		curState = f.popFrontier()
	}
	if curState == nil {
		f.endStatus_ = "All states explored, no more states in queue"
	}
}

// Checks whether the state is new or have been previously processed.
// If it isn't new, then compares the path length and decides if it is worth re-visiting it.
func (f *SbpBfsFinder) processState(s defs.SeqGameState, reversePath []defs.Command, mov defs.Command) {
	if f.debugTemp_ {
		s.MarkToDebug()
	}

	h := s.ToHash()
	if f.debug_ {
		f.outDbg1_.Printf("\n	 - Process state. Hash: %d, MOV: %v", h, mov)
	}

	chain := f.reversePathOn(reversePath, mov)
	s.SetMovChain(chain, nil)

	// Stop!
	// If this state comes from a state with an equivalency, check if that path is shorter:
	// If shorter, then recover that path
	if s.PrevState() != nil {
		oldLen := s.CollapsedPathLen()
		if s.PrevState().ApplyEquivalencyContinuity(s, mov, f.initState_) {
			newLen := s.CollapsedPathLen()
			if newLen < oldLen {
				chain = s.PathChain()
				if f.debug_ {
					f.outDbg1_.Printf("\n	 - Equivalency applied: oldLen: %d, newLen: %d", oldLen, newLen)
				}
			}
		}
	}

	if f.visitedStates_[h] == nil {

		// Easy, it is a new state
		f.countStates_.Incr()
		f.visitedStates_[h] = append(f.visitedStates_[h], s)

		// Let's explore its childs later
		f.addToFrontier(s)

		// If we are searching for a state, check if found:
		if f.search_ != nil {
			if s.EqualSub(f.search_) {
				f.updateObjective(s)
			}
		}

	} else {

		// We still need to compare it with other states to see if it is a new state or not.
		for _, st := range f.visitedStates_[h] {
			if st.Equal(s) {

				// We arrived to the same state from two separate chain of movements.
				// Update to the shortest path:
				lenOld := st.CollapsedPathLen()
				lenNew := s.CollapsedPathLen()

				if lenNew < lenOld {
					// Update!
					st.CopyMovChainFrom(s)

					if st.IsObjective() || s.IsObjective() {
						f.updateObjective(st)
					}

					// Discard s
					// But if st is not in the frontier, add again!
					if f.hardOptimals_ {
						if f.debug_ {
							f.outDbg1_.Printf("\n	 - Reintroduce to frontier! oldLen: %d, newLen: %d", lenOld, lenNew)
						}
						f.addToFrontier(s)
					}
				} else if lenNew == lenOld {
					if f.debug_ {
						f.outDbg1_.Printf(" (Same as [%d;CR:%d], cur: [%d;CR:%d])", st.Uid(), lenOld, s.Uid(), lenNew)
					}

					if st.Waiting() {
						if f.debug_ {
							f.outDbg1_.Printf(" Add equivalency.")
						}
						// Ok st is waiting and still not processed

						// We can add the equivalency if mov is a valid movement on directly on state st.
						// They are equivalent, but maybe, due to alike pieces, the mov command cannot be performed!
						//if st.ValidMovement(mov) {
						st.AddEquivPath(s, chain, mov) //And s is not processed by now...
						//}
					} else {
						// if st.SamePieceMovedNext(mov) {
						// 	f.addToFrontier(s)
						// }
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

		// If we are searching for a state, check if found:
		if f.search_ != nil {
			if s.EqualSub(f.search_) {
				f.updateObjective(s)
			}
		}

	}
}

// Push back to the priority queue that state.
func (f *SbpBfsFinder) addToFrontier(s defs.SeqGameState) {
	if f.debug_ {
		f.outDbg2_.Printf("\n	  Add to frontier: state [%d], from move %v", s.Uid(), s.PrevMov())
	}
	f.frontier_.PushBack(s)
	f.frontierSize_.Add(f.frontier_.Size())
	s.SetWaiting(true)
}

// Pop from start of queue (highest priority)
func (f *SbpBfsFinder) popFrontier() defs.SeqGameState {

	x := f.frontier_.PopFront()
	if x != nil {
		s := x.(*games.SBPState)
		if s != nil {
			(*s).SetWaiting(false)

			if f.debug_ {
				f.outDbg2_.Printf("\n\n### Explore state [%d]", (*s).Uid())
				(*s).TinyPrint()
			}
		}
		return s
	}
	return nil
}

// Called every time we found the objective state (solution of puzzle). We check whether the new
// solution is better or not.
func (f *SbpBfsFinder) updateObjective(s defs.SeqGameState) {

	s.MarkAsObjective()

	curPathLen := 0
	if f.foundState_ != nil {
		curPathLen = (*f.foundState_).CollapsedPathLen()
	}

	if f.foundState_ == nil || s.CollapsedPathLen() < curPathLen {
		f.foundState_ = &s

		if f.debug_ {
			f.outDbg1_.Println("\n\n --- Objective found ---")
			s.TinyPrint()
		}
	}
}

func (f *SbpBfsFinder) reversePathOn(path []defs.Command, mov defs.Command) []defs.Command {
	result := make([]defs.Command, len(path)+1)

	l := len(path) - 1
	for i := l; i >= 0; i-- {
		result[l-i] = path[i]
	}
	result[len(path)] = mov
	return result
}

/**
 * @summary Makes a depth-search and returns the most 'distant' states.
 *
 * @param {defs.Playable} g The sequential game
 * @param {int} maxSteps Max depth of the search
 * @param {int} maxExtremals Max number of extremal states returned
 */
func (f *SbpBfsFinder) FindExtremals(g defs.Playable) {

	f.findExtremals_ = true

	if !f.silent_ {
		fmt.Println("Puzzle Finder v.1.0")
	}
	f.fmtHeaders_ = color.New(color.FgCyan, color.Bold)

	f.outDbg1_ = color.New(color.FgCyan)
	f.outDbg2_ = color.New(color.FgWhite)
	f.outDbg3_ = color.New(color.FgYellow)

	f.countStates_.Set("States")
	f.nodesDegree_.Set("Node degree")
	f.frontierSize_.Set("Frontier size")

	f.game_ = g
	f.visitedStates_ = make(map[int][]defs.SeqGameState)
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
	}

	if !f.silent_ {
		f.resumeExtremals()
	}
}

// Prints statistics and results
func (f *SbpBfsFinder) resumeExtremals() {

	f.fmtHeaders_.Println("\n - Condition: ", f.endStatus_)

	f.fmtHeaders_.Println("\n[STATS]")
	f.countStates_.Resume(f.outDbg2_)
	f.nodesDegree_.ResumeAv(f.outDbg2_)
	f.frontierSize_.ResumeRange(f.outDbg2_)

	f.fmtHeaders_.Printf("\n\n[EXTREMAL STATES] Found: %d\n", len(f.extremals_))

	if len(f.extremals_) > 0 {
		f.extremals_[0].TinyGoPrint()
	}

	fmt.Println("\n\n")
}

// We are interested in extremal states, those at larger distance from the start state.
func (f *SbpBfsFinder) addExtremalState(s defs.SeqGameState) {

	if s.Depth() > f.extremalDist_ {

		// Then reset and add the state
		f.extremalDist_ = s.Depth()
		f.extremals_ = f.extremals_[:0]
		f.extremals_ = append(f.extremals_, s)
	} else if s.Depth() == f.extremalDist_ {

		f.extremals_ = append(f.extremals_, s)
	}
}
