package games

import "fmt"
import "github.com/edgarweto/puzzlopia/puzzle-solvers/definitions"
import "github.com/edgarweto/puzzlopia/puzzle-solvers/grids"

// A static, global or common for all states, struct that maps each piece id to its value used for comparisons.
type PieceToValue struct {
	idToValue_ map[int]int
}

func (ptv *PieceToValue) At(id int) int {
	return ptv.idToValue_[id]
}
func (ptv *PieceToValue) Set(id int, value int) {
	ptv.idToValue_[id] = value
}

// Shared among all states, static map.
var ptv *PieceToValue = nil

func getPieceToValueMap() *PieceToValue {
	if ptv == nil {
		ptv = &PieceToValue{}
		ptv.idToValue_ = make(map[int]int)
		ptv.idToValue_[0] = 0
	}
	return ptv
}

var staticSBPStateCount_ int = 0

// Sliding Blocks Puzzle Game State
type SBPState struct {
	uid_              int
	grid              grids.Matrix2d
	movChain_         []defs.Command
	movChainCount_    int
	pieceToValue_     *PieceToValue
	equivToObjective_ bool

	// BFS
	waiting_       bool
	depth_         int
	equivalencies_ []struct {
		state_ defs.GameState
		path_  []defs.Command
		mov_   defs.Command
	}

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
		g.pieceToValue_ = getPieceToValueMap()
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
		g.pieceToValue_ = getPieceToValueMap()
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
		g.pieceToValue_ = getPieceToValueMap()
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

	// Used in BFS
	if ss != nil {
		ss.nextMovs_ = append(ss.nextMovs_, mov)
	}

}

func (s *SBPState) PrevMov() defs.Command {
	return s.prevMov_
}

// Adds a state as a next state, and the corresponding mov
func (s *SBPState) AddNextState(u defs.GameState, mov defs.Command) {
	s.nextStates_ = append(s.nextStates_, u)
	s.nextMovs_ = append(s.nextMovs_, mov)
}

//Returns true if, among next movements, there is a movement of the same piece as the argument's mov.
func (s *SBPState) SamePieceMovedNext(mov defs.Command) bool {
	pieceId := mov.PieceId()
	for _, m := range s.nextMovs_ {
		if m.PieceId() == pieceId && !m.IsInverse(mov) {
			return true
		}
	}
	return false
}

func (s *SBPState) SetMovChain(movs []defs.Command) {
	s.movChain_ = movs
	s.updateChainLen()
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
func (s *SBPState) CopyMovChainAndAdd(path []defs.Command, mov defs.Command) {
	s.movChain_ = make([]defs.Command, len(path))
	copy(s.movChain_, path)
	s.movChain_ = append(s.movChain_, mov)
	s.updateChainLen()
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
	s.pieceToValue_ = getPieceToValueMap()

	// Finally edit the map
	for _, p := range pieces {
		s.pieceToValue_.Set(p.Id(), p.Value())
	}
}

// func (g *SBPState) ValidMovements(pieces []*grids.GridPiece2, seq *[]defs.Command, lastMov defs.Command, curPieceTrajectory []defs.Command) {

// 	// Prioritize consecutive movements with the same piece
// 	var samePieceMovs []defs.Command
// 	var otherMovs []defs.Command
// 	pieceId := 0
// 	if lastMov != nil {
// 		pieceId = lastMov.PieceId()
// 	}

// 	var piecePath []*grids.GridMov2 = nil
// 	if curPieceTrajectory != nil && len(curPieceTrajectory) > 0 {
// 		for _, m := range curPieceTrajectory {
// 			gMov, ok := m.(*grids.GridMov2)
// 			if !ok {
// 				panic("[SbpState::ValidMovements] mov is not a []*GridMov2!")
// 			}
// 			piecePath = append(piecePath, gMov)
// 		}
// 	}

// 	for _, p := range pieces {

// 		movs := g.grid.PieceMovements(p)

// 		if pieceId > 0 && pieceId == p.Id() {
// 			for _, m := range movs {
// 				if !m.IsInverse(lastMov) {

// 					// We need to avoid, when moving one piece consecutively, trajectories that touch themselves!
// 					if piecePath == nil || !grids.TrajectoryTouchesWithMov(piecePath, m) {
// 						samePieceMovs = append(samePieceMovs, m)
// 					}
// 				}
// 			}

// 		} else {
// 			for _, m := range movs {
// 				otherMovs = append(otherMovs, m)
// 			}
// 		}

// 		// for _, m := range movs {
// 		// 	*seq = append(*seq, m)
// 		// }
// 	}

// 	for _, m := range samePieceMovs {
// 		*seq = append(*seq, m)
// 	}
// 	for _, m := range otherMovs {
// 		*seq = append(*seq, m)
// 	}
// 	//*seq = append(*seq, samePieceMovs, otherMovs)
// }

func (s *SBPState) ClearPiece(p *grids.GridPiece2) {
	s.grid.ClearPiece(p)
}

func (s *SBPState) PlacePiece(p *grids.GridPiece2) {
	s.grid.PlacePiece(p)
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

	p := make([]defs.Command, len(path))
	copy(p, path)
	s.equivalencies_ = append(s.equivalencies_,
		struct {
			state_ defs.GameState
			path_  []defs.Command
			mov_   defs.Command
		}{a, p, m})
}

// Returns true if current movement is on the same piece as an equivalent position of this state.
// It indicates that the 'other' path is shorter (the length is checked when adding the equivalency)
func (s *SBPState) ApplyEquivalencyContinuity(a defs.GameState, mov defs.Command) bool {

	for _, x := range s.equivalencies_ {
		m := x.mov_

		if m.PieceId() == mov.PieceId() && !m.IsInverse(mov) {
			//s is the prev state of a.
			s.SetMovChain(x.path_)

			//a state has that path plus mov.
			a.CopyMovChainAndAdd(x.path_, mov)

			// REPARENT:
			a.SetPrevState(x.state_, x.mov_)

			return true
		}
	}
	return false
}

func (s *SBPState) PathChain() []defs.Command {
	return s.movChain_
}
