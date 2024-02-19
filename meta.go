package main

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

	binary.LittleEndian.PutUint64(buf[pos:],uint64(m.rootNode))
	pos+=pageNumSize

	binary.LittleEndian.PutUint64(buf[pos:],uint64(m.freelistPage))
	pos+=pageNumSize
}

func (m *meta) deserialize(buf []byte) {
	pos:=0

	m.rootNode = pgnum(binary.LittleEndian.Uint64(buf[pos:]))
	pos+=pageNumSize
	
	m.freelistPage=pgnum(binary.LittleEndian.Uint64(buf[pos:]))
	pos+=pageNumSize
}