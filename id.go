package cellnet

import (
	"fmt"
)

type CellID int64

func (self CellID) Region() int32 {

	return int32(self >> 32)
}

func (self CellID) Valid() bool {
	return self != 0
}

func (self CellID) Index() int32 {

	return int32(self & 0x00000000ffffffff)
}

func (self CellID) String() string {

	return fmt.Sprintf("[%d.%d]", self.Region(), self.Index())
}

func NewCellID(region, index int32) CellID {

	return CellID(int64(region)<<32 | int64(index))
}
