package gameuser

import "meatfloss/common"

// MakeBox ...
type MakeBox struct {
	UserID  int               // user id
	Lattice []*common.Lattice // 格子的数组
}

// NewMakeBox ...
func NewMakeBox(userID int) (lattice *MakeBox) {
	lattice = &MakeBox{}
	lattice.UserID = userID
	lattice.Lattice = make([]*common.Lattice, 0)
	return
}
