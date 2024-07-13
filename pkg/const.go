package ArachneDB

import(
	"errors"
)

const(
	pageNumSize = 8 	//the size of a page number (in bytes).
	metaPageNum = 0 	//the page no. where meta data is stored
	nodeHeaderSize = 3  // size of the page header of each node in bytes
	magicNumber uint16 = 0x2801	// Magic Number 
	magicNumberSize = 2 // size of the Magic Number in bytes
)

var writeInsideReadTxErr = errors.New("can't perform a write operation inside a read transaction")
