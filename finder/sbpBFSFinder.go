package finder

import "fmt"
import "time"
import "github.com/fatih/color"

import "github.com/edgarweto/puzzlopia/solvers/utils"
import "github.com/edgarweto/puzzlopia/solvers/games"

import "github.com/edgarweto/puzzlopia/solvers/grids" //TO REMOVE DEPENDENCY

// Sliding blocks puzzle finder
type SbpBfsFinder struct {

	// Params
	limits_       FinderLimits
	silent_       bool
	hardOptimals_ bool

	// Game settings
	game_      games.Playable
	initState_ games.GameState

	// If we are searching for a concrete state
	search_     *games.SBPState
	foundState_ *games.GameState

	// Stats
	countStates_  utils.ScalarStatistic
	nodesDegree_  utils.ScalarStatistic
	frontierSize_ utils.RangeStatistic

	// Algorithm state
	visitedStates_ map[int][]games.GameState
	frontier_      utils.Queue
	nextFrontier_  []games.GameState
	endStatus_     string
	duration_      time.Duration

	debug_      bool
	debugPath_  [][]int
	stepDebug_  int
	fmtHeaders_ *color.Color

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
}

/**
 * @summary Makes a depth-search and returns the most 'distant' states.
 *
 * @param {games.Playable} g The sequential game
 * @param {int} maxSteps Max depth of the search
 * @param {int} maxExtremals Max number of extremal states returned
 */
func (f *SbpBfsFinder) SolvePuzzle(g games.Playable, extremals *games.GameStates) {

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
	f.visitedStates_ = make(map[int][]games.GameState)
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

		var reversePath []games.GameMov
		curState.BuildPathReversed(&reversePath)

		f.game_.SetState(curState)
		validMovs := f.game_.ValidMovementsBFS()

		f.nodesDegree_.Add(len(validMovs))

		if f.debug_ {
			f.outDbg2_.Printf("	 Valid movs:%v", len(validMovs))
		}

		for _, mov := range validMovs {

			f.game_.Move(mov)
			newState := f.game_.State()
			newState.SetPrevState(curState, mov)

			f.processState(newState, reversePath, mov)

			f.game_.UndoMove(mov)
		}

		curState = f.popFrontier()
	}
	if curState == nil {
		f.endStatus_ = "All states explored, no more states in queue"
	}
}

func (f *SbpBfsFinder) processState(s games.GameState, reversePath []games.GameMov, mov games.GameMov) {
	h := s.ToHash()

	chain := f.reversePathOn(reversePath, mov)
	s.SetMovChain(chain)

	// Stop!
	// If this state comes from a state with an equivalency, check if that path is shorter:
	if s.PrevState() != nil {
		oldLen := s.CollapsedPathLen()
		if s.PrevState().ApplyEquivalencyContinuity(s, mov) {
			newLen := s.CollapsedPathLen()
			if newLen < oldLen {
				chain = s.PathChain()
				if f.debug_ {
					f.outDbg1_.Printf("\n	 - Equivalency applied: oldLen: %d, newLen: %d", oldLen, newLen)
				}
			}
		}
	}

	if f.debug_ {
		f.outDbg1_.Printf("\n	 - Process state. Hash: %d, MOV: %v", h, mov)
	}

	if f.visitedStates_[h] == nil {
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
						st.AddEquivPath(s, chain, mov)
					} else {
						if st.SamePieceMovedNext(mov) {
							f.addToFrontier(s)
						}
					}
				}
				return
			}
		}

		// Ok, hash collision. Add this new state!
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

func (f *SbpBfsFinder) addToFrontier(s games.GameState) {
	if f.debug_ {
		f.outDbg2_.Printf("\n	  Add to frontier: state [%d], from move %v", s.Uid(), s.PrevMov())
	}
	f.frontier_.PushBack(s)
	f.frontierSize_.Add(f.frontier_.Size())
	s.SetWaiting(true)
}

func (f *SbpBfsFinder) popFrontier() games.GameState {

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

func (f *SbpBfsFinder) updateObjective(s games.GameState) {

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

func (f *SbpBfsFinder) reversePathOn(path []games.GameMov, mov games.GameMov) []games.GameMov {
	result := make([]games.GameMov, len(path)+1)

	l := len(path) - 1
	for i := l; i >= 0; i-- {
		result[l-i] = path[i]
	}
	result[len(path)] = mov
	return result
}