package engel

import "github.com/edgarweto/puzzlopia/puzzle-solvers/definitions"

// Defines a puzzle consisting of intersecting wheels (Engel designs).
type EngelGame struct {

	// Initial state of the game
	state_ EngelState

	pieces_       []int
	pieceToValue_ *defs.PieceToValue
}

// Arrays are pieces unique id's.
// All wheels are indexed from 0 to 11, 0 being the right-middle rectangle piece position.
// So odd positions are occupied by rectangle pieces, while even positions correspond to triangle pieces.
func (g *EngelGame) Define(leftWheelPieces [12]int, rightWheelPieces [12]int) {
	g.state_.Init(leftWheelPieces, rightWheelPieces)

	g.pieces_ = leftWheelPieces[:]
	g.pieces_ = append(g.pieces_, rightWheelPieces[:]...)
}

// Defines the indexs of common piece positions for the wheels.
func (g *EngelGame) DefineIntersectionPositions(leftPositions [3]int, rightPositions [3]int) {
	g.state_.DefineIntersectionPositions(leftPositions, rightPositions)
}

// Identifies groups of pieces as being interchangeable. Creates a new identifier for each group.
func (g *EngelGame) SetAlike(pieceIdGroups [][]int) {

	// Assign the map
	g.pieceToValue_ = defs.GetPieceToValueMap()

	maxId := 0
	for _, id := range g.pieces_ {
		if id > maxId {
			maxId = id
		}
	}

	// Init map with read ids:
	for _, id := range g.pieces_ {
		g.pieceToValue_.Set(id, id)
	}

	// Then, change identified pieces values by their new group value:
	for idx, group := range pieceIdGroups {
		newValue := maxId + idx + 1

		for _, id := range group {
			g.pieceToValue_.Set(id, newValue)
		}
	}
}

// Implement Explorable interface

// Apply the movement
// wheelId is 0 (left) or 1 (right)
// steps is 0 (identity), 1 (60 degree rotation clockwise), etc.
func (g *EngelGame) Move(m defs.Command) (err error) {
	g.state_.Move(m.(*EngelCommand))
	return nil
}

// Undoes the movement
func (g *EngelGame) UndoMove(mov defs.Command) (err error) {
	m := mov.Inverted()
	mInv := m.(defs.Command)

	g.state_.Move(mInv.(*EngelCommand))
	return nil
}

// Makes the playable game to put its internal parts to
// reflect this state.
func (g *EngelGame) SetState(s defs.GameState) (err error) {
	g.state_.Assign(*(s.(*EngelState)))
	return nil
}

// Returns a copy of current state
func (g *EngelGame) State() (s defs.GameState) {
	return g.state_.Clone()
}

func (g *EngelGame) ValidMovements() []defs.Command {
	var movs []defs.Command

	prevMov := g.state_.PrevMov()
	if prevMov != nil {
		engelMov := prevMov.(*EngelCommand)
		altWheel := 1 - engelMov.PieceId()

		//fmt.Printf("\n <ValidMovs> 5, wheel: %d", altWheel)
		movs = []defs.Command{
			&EngelCommand{altWheel, 1},
			&EngelCommand{altWheel, 2},
			&EngelCommand{altWheel, 3},
			&EngelCommand{altWheel, 4},
			&EngelCommand{altWheel, 5},
		}

	} else {
		// Always we can perform wheel movements since they are never blocked.
		movs = []defs.Command{
			&EngelCommand{0, 1},
			&EngelCommand{0, 2},
			&EngelCommand{0, 3},
			&EngelCommand{0, 4},
			&EngelCommand{0, 5},
			&EngelCommand{1, 1},
			&EngelCommand{1, 2},
			&EngelCommand{1, 3},
			&EngelCommand{1, 4},
			&EngelCommand{1, 5},
		}
	}
	return movs
}
