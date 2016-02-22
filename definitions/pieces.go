package defs

// Struct that maps each piece id to its value used for comparisons.
// It holds the 'alike pieces' piece id to common value. All alike pieces will have the
// same value, so when comparing states, any switch between alike pieces will generate equivalent
// states.
// Static, global or common for all states.
type PieceToValue struct {
	idToValue_ map[int]int
}

// Value corresponding to a piece id
func (ptv *PieceToValue) At(id int) int {
	return ptv.idToValue_[id]
}

// Set a piece id-value correspondence
func (ptv *PieceToValue) Set(id int, value int) {
	ptv.idToValue_[id] = value
}

// Shared among all states, static map.
var s_ptv_ *PieceToValue = nil

// Returns the global map
func GetPieceToValueMap() *PieceToValue {
	if s_ptv_ == nil {
		s_ptv_ = &PieceToValue{}
		s_ptv_.idToValue_ = make(map[int]int)
		s_ptv_.idToValue_[0] = 0
	}
	return s_ptv_
}
