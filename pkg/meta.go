package ArachneDB

import "encoding/binary"

type meta struct{
	rootNode pgnum
	freelistPage pgnum
}

func newEmptyMeta() *meta{
	return &meta{}
}

func (m *meta) serialize(buf []byte) {
	pos:=0

	binary.LittleEndian.PutUint16(buf[pos:],magicNumber)
	pos+=magicNumberSize

	binary.LittleEndian.PutUint64(buf[pos:],uint64(m.rootNode))
	pos+=pageNumSize

	binary.LittleEndian.PutUint64(buf[pos:],uint64(m.freelistPage))
	pos+=pageNumSize
}

func (m *meta) deserialize(buf []byte) {
	pos:=0

	magic:=binary.LittleEndian.Uint16(buf[pos:])
	pos+=magicNumberSize

	if(magic!=magicNumber){
		panic("Not an Arachne DB file")
	}

	m.rootNode = pgnum(binary.LittleEndian.Uint64(buf[pos:]))
	pos+=pageNumSize
	
	m.freelistPage=pgnum(binary.LittleEndian.Uint64(buf[pos:]))
	pos+=pageNumSize
}