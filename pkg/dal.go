package ArachneDB

import (
	"errors"
	"fmt"
	"os"
)

type pgnum uint64

type dal struct{
	file *os.File
	pageSize int
	masterCollection *Collection	//root node of master collection
	minFillPercent float32
	maxFillPercent float32
	*freelist		//type embedding freelist in dal
	*meta			//type embedding meta in dal
}

type Options struct {
	pageSize int

	MinFillPercent float32
	MaxFillPercent float32
}

var DefaultOptions = &Options{
	MinFillPercent: 0.5,
	MaxFillPercent: 0.95,
}


type page struct{
	num pgnum
	data []byte
}

func (d *dal) allocateEmptyPage() *page {
	return &page{
		data: make([]byte, d.pageSize),			//creates an empty slice of bytes of size = Database Page
	}
}

func newDal(path string, options *Options) (*dal,error){
	dal:=&dal{
		pageSize: options.pageSize,				//Size of each Database Page
		meta: newEmptyMeta(),					//initialize empty Meta struct

		masterCollection: newEmptyCollection([]byte("Master"),0),
		minFillPercent: options.MinFillPercent,
		maxFillPercent: options.MaxFillPercent,

	}
	
	//check the status of the file

	if _,err:=os.Stat(path); err==nil{			
		//file exists
		file,err:=os.OpenFile(path, os.O_CREATE|os.O_RDWR,0666)
		if err!=nil{
			err=dal.close()
			return nil,err
		}
		dal.file=file

		err=dal.readMeta()
		dal.masterCollection.rootPgNum=dal.rootNode
		// dal.masterCollection.dal=dal
		if err!=nil{
			return nil,err
		}

		err=dal.readFreeList()

		if err!=nil{
			return nil,err
		}


	} else if errors.Is(err, os.ErrNotExist) {	

		//file does not exist

		//create file with path 'path' , open with read and write only permissions, create with file permissions - rw-rw-rw-
		file,err := os.OpenFile(path, os.O_CREATE|os.O_RDWR,0666)
		
		if err!=nil{
			return nil,err;
		}
		
		dal.file=file
		dal.freelist=newFreelist()
		dal.freelistPage=dal.getNextPage()

		_,err=dal.writeFreeList()

		if err!=nil{
			return nil,err
		}

		collectionsNode, err := dal.writeNode(NewNodeForSerialization([]*Item{}, []pgnum{}))
		if err != nil {
			return nil, err
		}
		dal.rootNode = collectionsNode.pageNum
		dal.masterCollection.rootPgNum=dal.rootNode
		// dal.masterCollection.dal=dal

		_,_=dal.writeFreeList()
		_,_=dal.writeMeta()
		
	}else{
		return nil,err
	}

	return dal,nil
}

func (d *dal) close() error {
	if (d.file!=nil){
		err := d.file.Close()
		if err!=nil{
			return fmt.Errorf("could not close file: %s",err)	//if the file could not be closed,displays error
		}
		d.file=nil
	}
	return nil
}

func (d *dal) readPage(num pgnum) (*page,error){
	p:=d.allocateEmptyPage()		//create an empty page to copy the contents of the original database page

	offset:=int(num)*d.pageSize		//calculate offset = pageNum*pageSize

	_,err:=d.file.ReadAt(p.data,int64(offset))	//read from file at offset for length of page(p)

	if err!=nil{
		return nil,err
	}

	return p,err
}

func (d *dal) writePage(p *page) error{

	offset:=int(p.num)*d.pageSize		//calculate offset = pageNum*pageSize

	_,err:=d.file.WriteAt(p.data,int64(offset))		//write p.data to the file at offset for length of page(p)

	return err
}


func (d *dal) writeMeta() (*page,error){
	p:=d.allocateEmptyPage()
	p.num=metaPageNum
	d.meta.serialize(p.data)

	err:=d.writePage(p)
	if err!=nil{
		return nil,err
	}
	return p,nil
}

func (d *dal) writeFreeList() (*page, error){
	p:=d.allocateEmptyPage()
	p.num=d.freelistPage

	d.freelist.serialize(p.data)

	err:=d.writePage(p)

	d.freelistPage=p.num

	return p,err
}

func (d *dal) readMeta() error{
	p,err:=d.readPage(metaPageNum)
	if err!=nil{
		return err
	}
	d.meta.deserialize(p.data)
	return nil
}

func (d *dal) readFreeList() error{
	p,err:=d.readPage(d.freelistPage)
	
	if err!=nil{
		return err
	}
	d.freelist=newFreelist()
	d.freelist.deserialize(p.data)
	return nil
}

func (d *dal) getNode(pageNum pgnum) (*Node,error){
	p,err:=d.readPage(pageNum)
	if err!=nil{
		return nil,err
	}
	node:=newEmptyNode()
	node.deserialize(p.data)
	node.pageNum=pageNum
	// node.dal = d
	return node,nil
}

func (d *dal) writeNode(n *Node) (*Node,error){
	p:=d.allocateEmptyPage()

	if n.pageNum == 0{
		p.num = d.getNextPage()
		n.pageNum = p.num
	}else{
		p.num = n.pageNum
	}

	n.serialize(p.data)
	err:=d.writePage(p)

	if err!=nil{
		return nil,err
	}
	return n,nil
}

// func (d *dal) newNode(items []*Item,childNodes []pgnum)(*Node){
// 	node:=newEmptyNode()
// 	node.items=items
// 	node.childNodes=childNodes
// 	// node.dal=d
// 	node.pageNum=d.getNextPage()
// 	return node
// }

func (d *dal) maxThreshold() float32 {
	return d.maxFillPercent * float32(d.pageSize)
}

func (d *dal) isOverPopulated(node *Node) bool {
	return float32(node.nodeSize()) > d.maxThreshold()
}

func (d *dal) minThreshold() float32 {
	return d.minFillPercent * float32(d.pageSize)
}

func (d *dal) isUnderPopulated(node *Node) bool {
	return float32(node.nodeSize()) < d.minThreshold()
}

// returns the index where the split must occur , 
// it also returns >0 when the node can spare an element and still be more than minimum size
// else returns -1 -> indicating it cannot spare an element
func (d *dal) getSplitIndex(node *Node) int{
	size:=0
	size+=nodeHeaderSize
	for i:=0; i<len(node.items); i++{
		size+=node.itemSize(i)
		// if we have a big enough page size (more than minimum), and didn't reach the last node, which means we can
		// spare an element
		if float32(size) > d.minThreshold() && i<len(node.items){
			return i + 1
		}
	}
	return -1
}

func (d *dal) deleteNode(pageNum pgnum){
	d.releasePage(pageNum)
}
