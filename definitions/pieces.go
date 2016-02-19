package defs

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
var s_ptv_ *PieceToValue = nil

func GetPieceToValueMap() *PieceToValue {
	if s_ptv_ == nil {
		s_ptv_ = &PieceToValue{}
		s_ptv_.idToValue_ = make(map[int]int)
		s_ptv_.idToValue_[0] = 0
	}
	return s_ptv_
}
