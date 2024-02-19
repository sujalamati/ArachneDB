package main

import "encoding/binary"

const(
	initialPage=0     					//uses one page (zeroth page) for metadata
)

type freelist struct{
	maxPage pgnum						//highest pageID of the database
	releasedPages []pgnum				//list of released pageID's 
}

func newFreelist() *freelist{
	return &freelist{
		maxPage: initialPage,						
		releasedPages: []pgnum{},		//initialised an empty slice of pageID's
	}
}

func (fl *freelist) getNextPage() pgnum{
	if(len(fl.releasedPages)!=0){										//check if there are any released pages
		pageID:=fl.releasedPages[len(fl.releasedPages)-1]				//fetch the last page
		fl.releasedPages=fl.releasedPages[:len(fl.releasedPages)-1]		//remove the last page from releasedPages
		return pageID													//return the last pageID
	}
	fl.maxPage+=1
	return fl.maxPage													//inc maxPage and return maxPage
}

func (fl *freelist) releasePage(pageID pgnum){
	fl.releasedPages = append(fl.releasedPages, pageID)			//append released page ID to releasedPagesList
}

func (fl *freelist) serialize(buf []byte) {

	//we are using Little Endian Byte Sequencing order
	pos:=0 				//pos is like a cursor in a file

	//after encoding the maxPage of freelist, we need to move the cursor forward by size of maxPage in bytes
	//otherwise it will overwrite 
	binary.LittleEndian.PutUint64(buf[pos:],uint64(fl.maxPage))
	pos+=pageNumSize

	binary.LittleEndian.PutUint64(buf[pos:],uint64(len(fl.releasedPages)))
	pos+=pageNumSize

	for _,page:=range fl.releasedPages{
		binary.LittleEndian.PutUint64(buf[pos:],uint64(page))
		pos+=pageNumSize
	}
}

func (fl *freelist) deserialize(buf []byte){
	pos:=0

	fl.maxPage=pgnum(binary.LittleEndian.Uint64(buf[pos:]))
	pos+=pageNumSize

	releasedPagesLen:=int(binary.LittleEndian.Uint64(buf[pos:]))
	pos+=pageNumSize

	for i:=0 ; i<releasedPagesLen; i++{
		fl.releasedPages = append(fl.releasedPages, pgnum(binary.LittleEndian.Uint64(buf[pos:])))
		pos+=pageNumSize
	}
}