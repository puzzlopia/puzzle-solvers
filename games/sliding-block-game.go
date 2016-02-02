package games

//import "fmt"
import "github.com/edgarweto/puzzlopia/solvers/grids"

// Sliding Block Game type
type SBGame struct {

	// Sliding Block game state
	state_ SBPState

	// Set of pieces
	pieces []*grids.GridPiece2

	// Map of pieces by their id.
	piecesById map[int]*grids.GridPiece2

	// Groups of pieces that should be considered alike
	alikePieces_ [][]int

	// If true, then alike pieces are automatically calculated
	autoAlikePieces_ bool

	// Set of pieces we want to maintain independent, not alike to other pieces
	notAutoalikePieces_ []int
}

// Implements GameDef interface
func (g *SBGame) Define(m *grids.Matrix2d) (err error) {
	g.state_.Init(m)

	return nil
}

// Implements GameDef interface
func (g *SBGame) SetState(s GameState) (err error) {
	b, ok := s.(*SBPState)
	if ok {
		g.state_.CopyGrid(*b)
		g.state_.UpdatePiecePositions(g.piecesById)
		g.state_.SetPrevState(s.PrevState(), s.PrevMov())
	} else {
		panic("[SBGame::SetState] GameState not of type SBPState!")
	}

	return nil
}

// Alike pieces are pieces with the same shape. Two states with two of
// these pieces at interchanged positions are considered the same state.
// This function is used to mark which pieces are equivalent or alike.
func (g *SBGame) AlikePieces(alikePieces [][]int) {
	g.alikePieces_ = alikePieces
}

// Marks all pieces with same shape as equivalent
func (g *SBGame) AutoAlikePieces() {
	g.autoAlikePieces_ = true
}

func (g *SBGame) SetNotAlikePiece(pieceId int) {
	g.notAutoalikePieces_ = append(g.notAutoalikePieces_, pieceId)
}

func (g *SBGame) Build() (err error) {

	// Create the pieces
	g.state_.BuildPieces(&g.pieces, g.alikePieces_, g.notAutoalikePieces_)

	// If autoAlikePieces, then detect equivalent pieces
	if g.autoAlikePieces_ {
		g.state_.DetectAlikePieces(g.pieces, g.notAutoalikePieces_)
	}

	g.piecesById = make(map[int]*grids.GridPiece2)
	for _, p := range g.pieces {
		g.piecesById[p.Id()] = p
	}

	return nil
}

// Playable interface
func (g *SBGame) Move(mov GameMov) (err error) {

	piece := g.piecesById[mov.PieceId()]

	// 1. Clear the piece from state's grid
	g.state_.ClearPiece(piece)

	// 2. Move the piece, updating its position
	piece.Move(mov)

	// 3. Return piece to the state's grid
	g.state_.PlacePiece(piece)

	return nil
}

// Reverts the movement
func (g *SBGame) UndoMove(mov GameMov) (err error) {

	piece := g.piecesById[mov.PieceId()]

	// 1. Clear the piece from state's grid
	g.state_.ClearPiece(piece)

	// 2. Move the piece, updating its position
	m := mov.Inverted()
	mInv := m.(GameMov)
	piece.Move(mInv)

	// 3. Return piece to the state's grid
	g.state_.PlacePiece(piece)

	return nil
}

// Return a copy of the state
func (g *SBGame) State() (s GameState) {
	return g.state_.Clone()
}

// Return a list of valid movements that can be done from this state
func (g *SBGame) ValidMovementsBFS() []GameMov {
	return g.state_.ValidMovementsBFS(g.pieces)
}