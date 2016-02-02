package games

import "fmt"
import "github.com/edgarweto/puzzlopia/solvers/grids"

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
	movChain_         []GameMov
	movChainCount_    int
	pieceToValue_     *PieceToValue
	equivToObjective_ bool

	// BFS
	waiting_       bool
	depth_         int
	equivalencies_ []struct {
		state_ GameState
		path_  []GameMov
		mov_   GameMov
	}

	// Graph structure to be able to perform a Dijkstra-like algorithm
	prevState_  GameState
	prevMov_    GameMov
	nextStates_ []GameState
	nextMovs_   []GameMov
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

// Initialize the game with the starting state matrix
func (g *SBPState) Print() {
	fmt.Println("\n----STATE----")
	fmt.Printf("PATH:")

	for i, m := range g.movChain_ {
		fmt.Printf("\n	[%d] ", i)
		m.Print()
	}

	fmt.Println("\nGRID:", g.grid)
	fmt.Println("\n")
}

// Initialize the game with the starting state matrix
func (g *SBPState) TinyPrint() {
	fmt.Printf("\n [id:%d] depth:%d, GRID: %v\n", g.uid_, g.depth_, g.grid)
	fmt.Printf(" PATH<%d>:", g.movChainCount_)

	for _, m := range g.movChain_ {
		//fmt.Printf("\n	[%d] ", i)
		fmt.Print(" ")
		m.Print()
	}

	//fmt.Println("\n")
}

/**
 * Implement GameState interface
 */
func (g *SBPState) Uid() int {
	return g.uid_
}

// Create a new GameState, cloned from current
func (g *SBPState) Clone() GameState {
	var c SBPState

	staticSBPStateCount_++
	c.uid_ = staticSBPStateCount_
	c.grid.Copy(&g.grid)
	c.prevState_ = nil
	return &c
}

//
func (g *SBPState) Equal(c GameState) bool {
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

func (g *SBPState) EqualSub(c GameState) bool {
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

// func (s *SBPState) PropagateUpdate() {

// 	// Back:
// 	parent := s.prevState_
// 	if parent != nil {
// 		//if parent.CollapsedPathLen() > s.movChainCount_+1 {
// 		//
// 		// Well, it seems that two states are at least at distance 1. But this is not true if the movement
// 		// is performed on the last moved piece. So if the path to reach state1 ends moving pieceX, then the
// 		// state2 that results in moving pieceX again has the same length!
// 		if parent.CollapsedPathLen() > s.movChainCount_ {
// 			//fmt.Print("[propagate back]")
// 			movToPrev := s.prevMov_.Inverted().(*grids.GridMov2)

// 			parent.copyMovChainAndAdd(s.movChain_, movToPrev)
// 			parent.PropagateUpdate()
// 		}
// 	}

// 	// Forward
// 	for idx, next := range s.nextStates_ {
// 		//if next.CollapsedPathLen() > s.movChainCount_+1 {
// 		if next.CollapsedPathLen() > s.movChainCount_ {
// 			//fmt.Print("[propagate forward]")
// 			next.copyMovChainAndAdd(s.movChain_, s.nextMovs_[idx])
// 			next.PropagateUpdate()
// 		}
// 	}
// }

func (s *SBPState) PrevState() GameState {
	return s.prevState_
}

// Sets prev status and the mov that brought to this state
func (s *SBPState) SetPrevState(prev GameState, mov GameMov) {
	var ss *SBPState = nil

	if prev == nil {
		s.depth_ = 0
	} else {
		ss = prev.(*SBPState)
		s.depth_ = ss.depth_ + 1
	}

	//s.AddEquivPath(prev, mov)

	s.prevState_ = prev
	s.prevMov_ = mov

	// Used in BFS
	if ss != nil {
		ss.nextMovs_ = append(ss.nextMovs_, mov)
	}

}

func (s *SBPState) PrevMov() GameMov {
	return s.prevMov_
}

// Adds a state as a next state, and the corresponding mov
func (s *SBPState) AddNextState(u GameState, mov GameMov) {
	s.nextStates_ = append(s.nextStates_, u)
	s.nextMovs_ = append(s.nextMovs_, mov)
}

//Returns true if, among next movements, there is a movement of the same piece as the argument's mov.
func (s *SBPState) SamePieceMovedNext(mov GameMov) bool {
	pieceId := mov.PieceId()
	for _, m := range s.nextMovs_ {
		if m.PieceId() == pieceId && !m.IsInverse(mov) {
			return true
		}
	}
	return false
}

func (s *SBPState) SetMovChain(movs []GameMov) {
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

// func (s *SBPState) MovChain() GameMov {
// 	return s.movChain_
// }

func (s *SBPState) CollapsedPathLen() int {
	//return len(s.movChain_)
	return s.movChainCount_
}

func (s *SBPState) RealPathLen() int {
	return len(s.movChain_)
}

func (s *SBPState) CopyMovChainFrom(gs GameState) {
	from, ok := gs.(*SBPState)
	if ok {
		s.movChain_ = make([]GameMov, len(from.movChain_))
		copy(s.movChain_, from.movChain_)
		s.updateChainLen()
	}
}

// Sets a new path to the state, adding a final step, and updates its real length
func (s *SBPState) copyMovChainAndAdd(path []GameMov, mov GameMov) {
	s.movChain_ = make([]GameMov, len(path))
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

	// Maps each piece id to the piece value, useful for states comparison
	//s.pieceToValue_ = make(map[int]int)
	s.pieceToValue_ = getPieceToValueMap()

	// Add the grid's 0 value
	//s.pieceToValue_[0] = 0
	for _, id := range notAutoalikePieces {
		for _, q := range pieces {
			if q.Id() == id {
				q.SetValue(q.Id())
				break
			}
		}
	}

	for _, p := range pieces {
		s.pieceToValue_.Set(p.Id(), p.Value())
	}

}

// func (g *SBPState) ValidMovements(pieces []*grids.GridPiece2, seq *[]GameMov, lastMov GameMov, curPieceTrajectory []GameMov) {

// 	// Prioritize consecutive movements with the same piece
// 	var samePieceMovs []GameMov
// 	var otherMovs []GameMov
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

// // Builds reverse path to reach root parent
// func (s *SBPState) BuildPath(path []GameMov) {

// 	if s.prevState_ != nil {
// 		path = append(path, s.prevMov_)
// 		s.prevState_.BuildPath(path)
// 	}
// }

func (s *SBPState) BuildPathReversed(path *[]GameMov) {

	if s.prevState_ != nil {
		*path = append(*path, s.prevMov_)
		s.prevState_.BuildPathReversed(path)
	}
}

/**
 * ====================== BFS ===============================
 */
func (s *SBPState) ValidMovementsBFS(pieces []*grids.GridPiece2) []GameMov {
	var seq []GameMov

	// Prioritize consecutive movements with the same piece
	var samePieceMovs []GameMov
	var otherMovs []GameMov
	pieceId := 0
	if s.prevMov_ != nil {
		pieceId = s.prevMov_.PieceId()
	}
	//fmt.Printf("\nPrev moved piece: %d", pieceId)

	for _, p := range pieces {

		movs := s.grid.PieceMovements(p)

		if pieceId > 0 && pieceId == p.Id() {
			for _, m := range movs {
				if !m.IsInverse(s.prevMov_) {

					// We need to avoid, when moving one piece consecutively, trajectories that touch themselves!
					//if piecePath == nil || !grids.TrajectoryTouchesWithMov(piecePath, m) {
					samePieceMovs = append(samePieceMovs, m)
					//}
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

	// for _, p := range pieces {

	// 	movs := s.grid.PieceMovements(p)

	// 	// if pieceId > 0 && pieceId == p.Id() {
	// 	// 	for _, m := range movs {
	// 	// 		if !m.IsInverse(lastMov) {

	// 	// 			// We need to avoid, when moving one piece consecutively, trajectories that touch themselves!
	// 	// 			if piecePath == nil || !grids.TrajectoryTouchesWithMov(piecePath, m) {
	// 	// 				samePieceMovs = append(samePieceMovs, m)
	// 	// 			}
	// 	// 		}
	// 	// 	}

	// 	// } else {
	// 	// 	for _, m := range movs {
	// 	// 		otherMovs = append(otherMovs, m)
	// 	// 	}
	// 	// }

	// 	for _, m := range movs {
	// 		seq = append(seq, m)
	// 	}
	// }

	return seq
}

func (s *SBPState) SetWaiting(b bool) {
	s.waiting_ = b
}
func (s *SBPState) Waiting() bool {
	return s.waiting_
}

func (s *SBPState) AddEquivPath(a GameState, path []GameMov, m GameMov) {

	p := make([]GameMov, len(path))
	copy(p, path)
	s.equivalencies_ = append(s.equivalencies_,
		struct {
			state_ GameState
			path_  []GameMov
			mov_   GameMov
		}{a, p, m})
}

// Returns true if current movement is on the same piece as an equivalent position of this state.
// It indicates that the 'other' path is shorter (the length is checked when adding the equivalency)
func (s *SBPState) ApplyEquivalencyContinuity(a GameState, mov GameMov) bool {

	for _, x := range s.equivalencies_ {
		m := x.mov_

		if m.PieceId() == mov.PieceId() && !m.IsInverse(mov) {
			//s is the prev state of a.
			s.SetMovChain(x.path_)

			//a state has that path plus mov.
			a.copyMovChainAndAdd(x.path_, mov)

			// REPARENT:
			a.SetPrevState(x.state_, x.mov_)

			return true
		}
	}
	return false
}

func (s *SBPState) PathChain() []GameMov {
	return s.movChain_
}
