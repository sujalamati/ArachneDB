package main

import (
	"errors"
	"fmt"
	"os"
)

type pgnum uint64

type dal struct{
	file *os.File
	pageSize int

	*freelist		//type embedding freelist in dal
	*meta			//type embedding meta in dal
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

func newDal(path string) (*dal,error){
	dal:=&dal{
		pageSize: os.Getpagesize(),				//Size of each Database Page
		meta: newEmptyMeta(),					//initialize empty Meta struct
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

		_,err=dal.writeMeta()

		
	}else{
		return nil,err
	}

	return dal,nil
}

func (d *dal) close() error {
	if (d.file!=nil){
		err := d.file.Close()
		if err!=nil{
			return fmt.Errorf("Could not close file: %s",err)	//if the file could not be closed,displays error
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
	return node,nil
}

