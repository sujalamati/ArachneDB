package main

import (
	"bytes"
	"encoding/binary"
)

type Item struct {
	key []byte
	value []byte
}

type Node struct{
	*dal				// type embedding dal in node

	pageNum pgnum		//each node is stored in one page
	items []*Item		
	childNodes []pgnum
}

func newEmptyNode() *Node{
	return &Node{}
}

func newItem(key []byte,value []byte) *Item{
	return &Item{
		key: key,
		value: value,
	}
}

func (n *Node) isLeaf() bool{
	return len(n.childNodes) == 0
}

func (n *Node) serialize(buf []byte){
	// using the concept of slotted pages(like PostgreSQL) to store the node in each page 

	leftPos:=0
	rightPos:=len(buf)-1

	// Page structure is
	// | Page Header | Child Pointer 1 | Offset 1 | Child Pointer 2 | Offset 2 |  ..... 
	// | Offset K | Child Pointer K+1 | ...
	// ...
	// ..........| Key K | Value K | Key K-1 | Value K-1 | ......| Key 2 | Value 2 | Key 1 | Value 1 |

	// Page Header includes isLeaf , K-V count

	isLeaf:=n.isLeaf()

	var bitSetVar uint64 = 0
	if isLeaf{
		bitSetVar=1			//set bitSetVar = 1 if it is a leaf node
	}
	// Write the page header (isLeaf and No of K-V pairs)
	buf[leftPos]=byte(bitSetVar)
	leftPos+=1

	// write the no of K-V pairs
	binary.LittleEndian.PutUint16(buf[leftPos:],uint16(len(n.items)))
	leftPos+=2

	for i:=0; i < len(n.items); i++{
		item := n.items[i]
		if !isLeaf{

			//write the ith child node
			binary.LittleEndian.PutUint64(buf[leftPos:],uint64(n.childNodes[i]))
			leftPos+=pageNumSize
		}

		
		vlen:=len(item.value)
		klen:=len(item.key)

		//write the value
		rightPos=rightPos-vlen
		copy(buf[rightPos:],item.value)

		//write the len of value
		rightPos-=1
		buf[rightPos]=byte(vlen)

		//write the key
		rightPos = rightPos - klen
		copy(buf[rightPos:],item.key)

		//write the len of key
		rightPos-=1
		buf[rightPos]=byte(klen)

		//writing the offset where the Key-Value pair is stored
		offset:=rightPos
		binary.LittleEndian.PutUint16(buf[leftPos:],uint16(offset))
		leftPos+=2
	}

	if !isLeaf{
		//write the last child node
		lastChildNode:=n.childNodes[len(n.childNodes)-1]

		binary.LittleEndian.PutUint64(buf[leftPos:],uint64(lastChildNode))
		leftPos+=pageNumSize
	}
}

func (n *Node) deserialize(buf []byte) {
	leftPos:=0
	
	//Read page header

	//read isLeaf
	isLeaf:=uint16(buf[leftPos])
	leftPos+=1

	// read the no of K-V pairs
	noOfRec:=binary.LittleEndian.Uint16(buf[leftPos:])
	leftPos+=2
	
	for i:=0 ; i<int(noOfRec); i++{
		if isLeaf == 0{
			// read the ith child node
			childNode:=binary.LittleEndian.Uint64(buf[leftPos:])
			// append ith child node
			n.childNodes = append(n.childNodes, pgnum(childNode))
			leftPos+=pageNumSize
		}
		// read offset
		offset:=binary.LittleEndian.Uint16(buf[leftPos:])
		leftPos+=2

		//read len of key
		klen:=uint16(buf[offset])
		offset=offset+1

		// read key
		key:=buf[offset:offset+klen]
		offset+=klen

		// read len of value
		vlen:=uint16(buf[offset])
		offset+=1

		// read value
		value:=buf[offset:offset+vlen]
		offset+=vlen

		item:=newItem(key,value)
		n.items = append(n.items, item)
	}
	if isLeaf == 0 {
		lastChildNode := binary.LittleEndian.Uint64(buf[leftPos:])
		n.childNodes = append(n.childNodes, pgnum(lastChildNode))
	}
}

func (n *Node) searchNode(key []byte) (*Node,int,error) {
	node,index,err:=searchNodeRec(n,key)

	if err!=nil{
		return nil,-1,err
	}
	
	return node,index,nil
}

func searchNodeRec(node *Node , key []byte) (*Node,int,error){
	
	// search for key in the node
	found,index:=node.searchInNode(key)
	if found{
		return node,index,nil
	}

	if node.isLeaf(){
		// if the node is leaf and the key is not found => key does not exist
		return nil,-1,nil
	}

	// fetch the child node where the key is present , from disk into memory
	child,err:=node.getNode(node.childNodes[index])

	if err!=nil{
		return nil,-1,err
	}

	// search in the child node
	return searchNodeRec(child,key)

}

func (n *Node) searchInNode(key []byte) (bool,int) {
	// NOTE: we are storing keys in lexicographical order

	for i,item:= range n.items{
		flag:=bytes.Compare(item.key,key)
		// flag = 0 => item.key and key matches
		// flag = 1 => item.key is lexicographically greater than key

		if flag == 0{
			// match found
			return true,i
		}
		if flag > 0{
			// item.key is lexicographically greater than key
			return false,i
		}
	}
	// desiredKey is greater than all the keys in the node => it is in the last child node 
	return false,len(n.items)
}