package main

import "fmt"


type Collection struct{
	name []byte		// name of the Collection
	rootPgNum pgnum	// page num where the root of B-Tree is stored

	dal *dal 		// type embedding dal into collection
}

func newEmptyCollection(name []byte,root pgnum) *Collection{
	return &Collection{
		name: name,
		rootPgNum: root,
	}
}

func (c *Collection) Find(key []byte) (*Item,error){
	rootNode,err:=c.dal.getNode(c.rootPgNum)
	if err!=nil{
		return nil,err
	}
	_,containingNode,index,_,_,err:=rootNode.searchNode(key,true)
	if err!=nil{
		return nil,err
	}
	if index == -1{
		return nil,nil
	}
	return containingNode.items[index],nil
}

func (c *Collection) Put(key []byte , value []byte) error {
	item:=newItem(key,value)
	var root *Node
	var err error
	if c.rootPgNum == 0{
		root,err=c.dal.writeNode(c.dal.newNode([]*Item{item},[]pgnum{}))
		if err!=nil{
			return err
		}
		c.rootPgNum=root.pageNum
		return nil
	}


	root,err=c.dal.getNode(c.rootPgNum)


	if err!=nil{
		return err
	}
	
	wasFound,Node,index,parentIndices,parents,err:=root.searchNode(key,true)

	fmt.Println(Node.pageNum,index)
	if err!=nil{
		return err
	}

	// if key already exists in the tree, then replace the value with the new value
	if wasFound{
		Node.items[index]=item
	}else{ // key doesnt exist in the tree
		Node.insertInNode(item,index)
	}

	fmt.Println(Node)
	Node.writeNode(Node)
	
	// Rebalance the tree, from bottom to top
	for i:=len(parents)-2; i>=0; i--{
		pnode:=parents[i]
		node:=parents[i+1]
		insertIndex:=parentIndices[i]
		
		if node.isOverPopulated(){
			fmt.Printf("node %d is overpopulated",node.pageNum)
			pnode.split(node,insertIndex)
		}
	}
	
	// Balancing the root
	root = parents[0]
	if root.isOverPopulated(){
		fmt.Println("root is overpopulated")
		fmt.Printf("max page no. %d",c.dal.maxPage)
		newNode:=c.dal.newNode([]*Item{},[]pgnum{root.pageNum})
		fmt.Printf("max page no. %d",c.dal.maxPage)
		fmt.Println(root.pageNum)
		fmt.Println(newNode.pageNum)
		newNode.split(root,0)
		fmt.Println(newNode)

		newNode, err = c.dal.writeNode(newNode)
		if err != nil {
			return err
		}
		c.rootPgNum=newNode.pageNum
	}

	return nil

}