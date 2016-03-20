package engel

import "fmt"

//import "github.com/edgarweto/puzzlopia/puzzle-solvers/definitions"

// ENGEL COMMAND
type EngelCommand struct {
	wheelId_  int
	rotation_ int
}

// Implements command interface
func (c *EngelCommand) PieceId() int {
	return c.wheelId_
}
func (c *EngelCommand) Rotation() int {
	return c.rotation_
}
func (c *EngelCommand) Inverted() interface{} {
	x := EngelCommand{c.wheelId_, -c.rotation_}
	return &x
}

func (c *EngelCommand) IsInverse(m interface{}) bool {
	mov, ok := m.(*EngelCommand)
	if ok {
		if c.wheelId_ == mov.wheelId_ && (c.rotation_+mov.rotation_)%6 == 0 {
			return true
		}
	} else {
		panic("[EngelCommand::IsInverse] arg is not an EngelCommand")
	}
	return false
}

func (c *EngelCommand) Print() {
	wheelName := 'L'
	if c.wheelId_ == 1 {
		wheelName = 'R'
	}
	fmt.Printf("[%s %d]", wheelName, c.rotation_)
}
