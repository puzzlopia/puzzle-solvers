package games

import "fmt"
import "github.com/edgarweto/puzzlopia/puzzle-solvers/definitions"
import "github.com/edgarweto/puzzlopia/puzzle-solvers/grids"

var staticSBPStateCount_ int = 0

// Sliding Blocks Puzzle Game State
type SBPState struct {
	uid_              int
	grid              grids.Matrix2d
	movChain_         []defs.Command
	movChainCount_    int
	pieceToValue_     *defs.PieceToValue
	equivToObjective_ bool
	//originState_      *grids.Matrix2d

	// BFS
	waiting_       bool
	depth_         int
	equivalencies_ []struct {
		state_ defs.GameState
		path_  []defs.Command
		mov_   defs.Command
	}
	markedDebug_ bool

	// Graph structure
	prevState_  defs.GameState
	prevMov_    defs.Command
	nextStates_ []defs.GameState
	nextMovs_   []defs.Command
}

// Initialize the game with the starting state matrix
func (g *SBPState) Init(m *grids.Matrix2d) {
	staticSBPStateCount_++
	g.uid_ = staticSBPStateCount_
	g.grid.Copy(m)
	g.prevState_ = nil
}

func (g *SBPState) CopyGrid(s SBPState) {
	g.grid.Copy(&s.grid)
}

func (g *SBPState) UpdatePiecePositions(piecesById map[int]*grids.GridPiece2) {
	g.grid.UpdatePiecePositions(piecesById)
}

func (g *SBPState) MarkToDebug() {
	g.markedDebug_ = true
}

func (g *SBPState) MarkedToDebug() bool {
	return g.markedDebug_
}

// // Initialize the game with the starting state matrix
// func (g *SBPState) Print() {
// 	fmt.Println("\n----STATE----")
// 	fmt.Printf("PATH:")

// 	for i, m := range g.movChain_ {
// 		fmt.Printf("\n	[%d] ", i)
// 		m.Print()
// 	}

// 	fmt.Println("\nGRID:", g.grid)
// 	fmt.Println("\n")
// }

// Initialize the game with the starting state matrix
func (g *SBPState) TinyPrint() {
	fmt.Printf("\n [id:%d] depth:%d, GRID: %v\n", g.uid_, g.depth_, g.grid)
	fmt.Printf(" PATH<%d>:", g.movChainCount_)

	// Json format
	fmt.Print("[")
	for idx, m := range g.movChain_ {
		if idx > 0 {
			fmt.Print(",")
		}
		m.Print()
	}
	fmt.Print("]")

}

func (g *SBPState) TinyGoPrint() {
	fmt.Printf("\n [id:%d] depth:%d, GRID: ", g.uid_, g.depth_)

	rows := g.grid.Rows()
	cols := g.grid.Cols()
	for r := 0; r < rows; r++ {
		fmt.Printf("\n []int{")
		for c := 0; c < cols; c++ {
			if c == 0 {
				fmt.Printf("%d", g.grid.At(r, c))
			} else {
				fmt.Printf(", %d", g.grid.At(r, c))
			}
		}
		fmt.Printf("},")
	}
	fmt.Printf("\n")

	fmt.Printf(" PATH<%d>:", g.movChainCount_)

	for _, m := range g.movChain_ {
		fmt.Print(" ")
		m.Print()
	}

}

/**
 * Implement defs.GameState interface
 */
func (g *SBPState) Uid() int {
	return g.uid_
}

// Create a new defs.GameState, cloned from current
func (g *SBPState) Clone() defs.GameState {
	var c SBPState

	staticSBPStateCount_++
	c.uid_ = staticSBPStateCount_
	c.grid.Copy(&g.grid)
	c.prevState_ = nil
	return &c
}

//
func (g *SBPState) Equal(c defs.GameState) bool {
	if g.pieceToValue_ == nil {
		g.pieceToValue_ = defs.GetPieceToValueMap()
	}
	//fmt.Println("[SBPState::Equal]", g.pieceToValue_)

	rows := g.grid.Rows()
	cols := g.grid.Cols()

	s2, ok := c.(*SBPState)
	if ok {
		for i := 0; i < rows; i++ {
			for j := 0; j < cols; j++ {
				if g.pieceToValue_.At(g.grid[i][j]) != g.pieceToValue_.At(s2.grid[i][j]) {
					//fmt.Println("Different states:")
					//fmt.Printf("	At cell (%d,%d) values are %d and %d", i, j, g.pieceToValue_[g.grid[i][j]], g.pieceToValue_[s2.grid[i][j]])
					return false
				}
			}
		}
	}
	return true
}

func (g *SBPState) EqualSub(c defs.GameState) bool {
	if g.pieceToValue_ == nil {
		g.pieceToValue_ = defs.GetPieceToValueMap()
	}

	rows := g.grid.Rows()
	cols := g.grid.Cols()

	anyPieceFound := false
	s2, ok := c.(*SBPState)
	if ok {
		for i := 0; i < rows; i++ {
			for j := 0; j < cols; j++ {

				// Only check pieces from 'c' state
				x := g.pieceToValue_.At(s2.grid[i][j])
				if x > 0 {
					anyPieceFound = true
					if g.pieceToValue_.At(g.grid[i][j]) != x {
						return false
					}
				}
			}
		}
	}
	return anyPieceFound
}

func (g *SBPState) MarkAsObjective() {
	g.equivToObjective_ = true
}

func (g *SBPState) IsObjective() bool {
	return g.equivToObjective_
}

// Generate an integer from the current position of pieces
func (g *SBPState) ToHash() int {
	rows := g.grid.Rows()
	cols := g.grid.Cols()
	if g.pieceToValue_ == nil {
		g.pieceToValue_ = defs.GetPieceToValueMap()
	}

	hash := 0

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			hash += (i + 1) * (rows*cols + j) * g.pieceToValue_.At(g.grid[i][j])
		}
	}
	return hash
}

func (s *SBPState) PrevState() defs.GameState {
	return s.prevState_
}

// Sets prev status and the mov that brought to this state
func (s *SBPState) SetPrevState(prev defs.GameState, mov defs.Command) {
	var ss *SBPState = nil

	if prev == nil {
		s.depth_ = 0
	} else {
		ss = prev.(*SBPState)
		s.depth_ = ss.depth_ + 1
	}

	s.prevState_ = prev
	s.prevMov_ = mov

	// // Used in BFS
	// if ss != nil {
	// 	ss.nextMovs_ = append(ss.nextMovs_, mov)
	// }

}

// After doing a reparent, we have to update
func (s *SBPState) UpdateFromPrevState() {

	prev := s.prevState_.(*SBPState)
	s.grid.Copy(&prev.grid)

	// And apply last mov
	s.applyMov(s.prevMov_)
}

// Checks only from prev
func (s *SBPState) CheckPathAndState() {
	var tempState SBPState

	prev := s.prevState_.(*SBPState)
	tempState.grid.Copy(&prev.grid)

	// And apply last mov
	tempState.applyMov(s.prevMov_)

	// Now should be identical, not only equivalent!
	if !tempState.grid.Identical(s.grid) {
		fmt.Println("\n *** CheckPathAndState FAILED *** ")

		fmt.Println("\n SHOULD BE:", tempState.grid)
		fmt.Println("\n CURRENT:", s.grid)

		panic("STOP")
	}

}

func (s *SBPState) PrevMov() defs.Command {
	return s.prevMov_
}

// // Adds a state as a next state, and the corresponding mov
// func (s *SBPState) AddNextState(u defs.GameState, mov defs.Command) {
// 	s.nextStates_ = append(s.nextStates_, u)
// 	s.nextMovs_ = append(s.nextMovs_, mov)
// }

// //Returns true if, among next movements, there is a movement of the same piece as the argument's mov.
// func (s *SBPState) SamePieceMovedNext(mov defs.Command) bool {
// 	pieceId := mov.PieceId()
// 	for _, m := range s.nextMovs_ {
// 		if m.PieceId() == pieceId && !m.IsInverse(mov) {
// 			return true
// 		}
// 	}
// 	return false
// }

func (s *SBPState) SetMovChain(movs []defs.Command, updateStateFromPath *defs.GameState) {
	s.movChain_ = movs
	s.updateChainLen()

	if updateStateFromPath != nil {
		s.UpdateFromStart(updateStateFromPath)
	}
}

func (s *SBPState) updateChainLen() {
	// Two or more consecutive movs on the same piece count as only one movement!
	s.movChainCount_ = 0
	lastPieceId := 0
	for _, m := range s.movChain_ {
		if lastPieceId != m.PieceId() {
			s.movChainCount_++
			lastPieceId = m.PieceId()
		}
	}
}

func (s *SBPState) SetDepth(d int) {
	s.depth_ = d
}
func (s *SBPState) Depth() int {
	return s.depth_
}

func (s *SBPState) CollapsedPathLen() int {
	return s.movChainCount_
}

func (s *SBPState) RealPathLen() int {
	return len(s.movChain_)
}

func (s *SBPState) CopyMovChainFrom(gs defs.GameState) {
	from, ok := gs.(*SBPState)
	if ok {
		s.movChain_ = make([]defs.Command, len(from.movChain_))
		copy(s.movChain_, from.movChain_)
		s.updateChainLen()
	}
}

// Sets a new path to the state, adding a final step, and updates its real length
func (s *SBPState) CopyMovChainAndAdd(path []defs.Command, mov defs.Command, updateStateFromPath *defs.GameState) {
	s.movChain_ = make([]defs.Command, len(path))
	copy(s.movChain_, path)
	s.movChain_ = append(s.movChain_, mov)
	s.updateChainLen()

	if updateStateFromPath != nil {
		s.UpdateFromStart(updateStateFromPath)
	}
}

// Removes the state and restarts: it starts with the original state and then applies the chain of movs.
func (s *SBPState) UpdateFromStart(originState *defs.GameState) {

	//fmt.Println("\t**** UpdateFromStart ****")

	o := (*originState).(*SBPState)
	s.grid.Copy(&o.grid)

	// fmt.Println("Origin grid:", s.grid)
	// fmt.Println("Path: [")
	for _, mov := range s.movChain_ {

		// if idx > 0 {
		// 	fmt.Printf(",")
		// }
		// mov.Print()
		s.applyMov(mov)
	}
	//	fmt.Println("]")
}

// Builds a set of pieces from current state
func (s *SBPState) BuildPieces(pieces *[]*grids.GridPiece2, alikePieces [][]int, notAutoalikePieces []int) int {
	n := s.grid.GeneratePieces(pieces, alikePieces)
	s.updatePieceToValue(*pieces, notAutoalikePieces)
	return n
}

func (s *SBPState) DetectAlikePieces(pieces []*grids.GridPiece2, notAutoalikePieces []int) {
	grids.DetectAlikePieces(pieces, s.grid.Max()+1)
	s.updatePieceToValue(pieces, notAutoalikePieces)
}

func (s *SBPState) updatePieceToValue(pieces []*grids.GridPiece2, notAutoalikePieces []int) {

	// Edit not-autoalike piece values
	for _, id := range notAutoalikePieces {
		for _, q := range pieces {
			if q.Id() == id {
				q.SetValue(q.Id())
				break
			}
		}
	}

	// Assign the map
	s.pieceToValue_ = defs.GetPieceToValueMap()

	// Finally edit the map
	for _, p := range pieces {
		s.pieceToValue_.Set(p.Id(), p.Value())
	}
}

func (s *SBPState) ClearPiece(p *grids.GridPiece2) {
	s.grid.ClearPiece(p)
}

func (s *SBPState) PlacePiece(p *grids.GridPiece2) {
	s.grid.PlacePiece(p)
}

// We need to update a state from the original state.
func (s *SBPState) applyMov(mov defs.Command) (err error) {

	pieceId := mov.PieceId()

	gMov, ok := mov.(*grids.GridMov2)
	if !ok || gMov == nil {
		panic("[SBPState::applyMov] mov is not a GridMov2!")
	}

	//fmt.Println("  apply mov: ", gMov)

	s.grid.ApplyRawTranslation(pieceId, *gMov)

	return nil
}

func (s *SBPState) BuildPathReversed(path *[]defs.Command) {

	if s.prevState_ != nil {
		*path = append(*path, s.prevMov_)
		s.prevState_.BuildPathReversed(path)
	}
}

/**
 * ====================== BFS ===============================
 */
func (s *SBPState) ValidMovementsBFS(pieces []*grids.GridPiece2, pieceTrajectory []defs.Command) []defs.Command {
	var seq []defs.Command

	// Prioritize consecutive movements with the same piece
	var samePieceMovs []defs.Command
	var otherMovs []defs.Command
	pieceId := 0
	if s.prevMov_ != nil {
		pieceId = s.prevMov_.PieceId()
	}

	for _, p := range pieces {

		movs := s.grid.PieceMovements(p)

		if pieceId > 0 && pieceId == p.Id() {
			for _, m := range movs {
				if !m.IsInverse(s.prevMov_) {

					// We need to avoid, when moving one piece consecutively, trajectories that touch themselves!
					if len(pieceTrajectory) == 0 || !grids.TrajectoryTouchesWithMov(pieceTrajectory, m) {
						samePieceMovs = append(samePieceMovs, m)
					}
				}
			}
		} else {
			for _, m := range movs {
				otherMovs = append(otherMovs, m)
			}
		}
	}

	for _, m := range samePieceMovs {
		seq = append(seq, m)
	}
	for _, m := range otherMovs {
		seq = append(seq, m)
	}

	return seq
}

func (s *SBPState) SetWaiting(b bool) {
	s.waiting_ = b
}
func (s *SBPState) Waiting() bool {
	return s.waiting_
}

func (s *SBPState) AddEquivPath(a defs.GameState, path []defs.Command, m defs.Command) {

	// Duplicate path
	p := make([]defs.Command, len(path))
	copy(p, path)

	// Add to equivalences
	s.equivalencies_ = append(s.equivalencies_,
		struct {
			state_ defs.GameState
			path_  []defs.Command
			mov_   defs.Command
		}{a, p, m})
}

// Returns true if the movement m can be applied to the state
func (s *SBPState) ValidMovement(mov defs.Command) bool {
	gMov, ok := mov.(*grids.GridMov2)
	if !ok || gMov == nil {
		panic("[SBPState::ValidMovement] mov is not a GridMov2!")
	}

	return s.grid.ValidMove(*gMov)
}

// Returns true if current movement is on the same piece as an equivalent position of this state.
// It indicates that the 'other' path is shorter (the length is checked when adding the equivalency)
func (s *SBPState) ApplyEquivalencyContinuity(a defs.GameState, mov defs.Command, origin defs.GameState) bool {

	//fmt.Println("[ApplyEquivalencyContinuity]")

	// If two states are equivalent, that doesn't means all pieces are equally positionated. It only means that
	// all alike pieces are at the same position. Here we are going to update the path to a concrete state with
	// the path of an equivalence state. Then, we may end up with an inconsistency between alike pieces, since that
	// new path will bring to the same 'equivalent' state, but with undetermined alike piece switches!
	// To solve this, we need to update the state from the path.
	updateStateFromPath := &origin

	for _, x := range s.equivalencies_ {

		//fmt.Printf("\n\t Equivalency update. MOV (%v), GRID: %v", x.mov_, x.state_.(*SBPState).grid)

		m := x.mov_

		// x.mov_ is going to be the previous movement to reach state a, so we should ensure that it can be undone on
		// that state:
		//mInv := x.mov_.Inverted().(defs.Command)

		//if m.PieceId() == mov.PieceId() && !m.IsInverse(mov) && x.state_.ValidMovement(mov) && a.ValidMovement(mInv) {
		//if m.PieceId() == mov.PieceId() && x.state_.ValidMovement(mov) && a.ValidMovement(mInv) {
		if m.PieceId() == mov.PieceId() && !m.IsInverse(mov) && x.state_.ValidMovement(mov) {

			// DO NOT MODIFY S, the prev state:
			//s is the prev state of a.
			//s.SetMovChain(x.path_, updateStateFromPath)

			//a state has that path plus mov.
			a.CopyMovChainAndAdd(x.path_, mov, updateStateFromPath)

			// REPARENT:
			//a.SetPrevState(x.state_, x.mov_)
			a.SetPrevState(x.state_, mov)

			// a should have a correct state because we've updated it from the origin!
			//a.CheckPathAndState()

			//a.UpdateFromStart(updateStateFromPath)

			//a.UpdateFromPrevState() // Now we should update a's grid state!! Because its new parent may have alike pieces switched!!!
			//a.CheckPathAndState()

			return true
		}
	}
	return false
}

func (s *SBPState) PathChain() []defs.Command {
	return s.movChain_
}
