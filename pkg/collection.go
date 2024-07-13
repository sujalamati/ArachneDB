package ArachneDB

type Collection struct{
	name []byte		// name of the Collection
	rootPgNum pgnum	// page num where the root of B-Tree is stored

	// dal *dal 		// type embedding dal into collection

	tx *tx			//associated transaction
}

func newEmptyCollection(name []byte,root pgnum) *Collection{
	return &Collection{
		name: name,
		rootPgNum: root,
	}
}

func (c *Collection) Find(key []byte) (*Item,error){
	rootNode,err:=c.tx.getNode(c.rootPgNum)
	if err!=nil{
		return nil,err
	}
	_,containingNode,index,_,_,err:=rootNode.searchNode(key,false)
	if index == -1{
		return nil,nil
	}
	if err!=nil{
		return nil,err
	}
	return containingNode.items[index],nil
}

func (c *Collection) Put(key []byte , value []byte) error {
	
	if !c.tx.write {
		return writeInsideReadTxErr
	}

	item:=newItem(key,value)
	var root *Node
	var err error
	if c.rootPgNum == 0{
		root = c.tx.writeNode(c.tx.newNode([]*Item{item},[]pgnum{}))
		
		c.rootPgNum=root.pageNum
		return nil
	}


	root,err=c.tx.getNode(c.rootPgNum)
	if err!=nil{
		return err
	}
	
	wasFound,Node,index,parentIndices,parents,err:=root.searchNode(key,true)

	if err!=nil{
		return err
	}

	// if key already exists in the tree, then replace the value with the new value
	if wasFound{
		Node.items[index]=item
	}else{ // key doesnt exist in the tree
		Node.insertInNode(item,index)
	}

	Node.writeNode(Node)
	
	// Rebalance the tree, from bottom to top
	for i:=len(parents)-2; i>=0; i--{
		pnode:=parents[i]
		node:=parents[i+1]
		insertIndex:=parentIndices[i]
		
		if node.isOverPopulated(){
			pnode.split(node,insertIndex)
		}
	}
	
	// Balancing the root
	root = parents[0]
	if root.isOverPopulated(){
		newNode:=c.tx.newNode([]*Item{},[]pgnum{root.pageNum})
		newNode.split(root,0)

		newNode = c.tx.writeNode(newNode)
		
		if string(c.name)== "Master"{
			c.tx.db.rootNode=newNode.pageNum
		}else{
			c.rootPgNum=newNode.pageNum
			c.tx.db.masterCollection.Put(c.name,[]byte{byte(c.rootPgNum)})
		}
	}

	return nil

}

func (c *Collection) Remove(key []byte) error{
	if !c.tx.write {
		return writeInsideReadTxErr
	}
	
	root,err:=c.tx.getNode(c.rootPgNum)

	if err!=nil{
		return err
	}
	// locate the key in the b-tree
	Found,node,index,parentIndices,parents,err:=root.searchNode(key,false)

	if err!=nil{
		return err
	}

	if !Found{
		return nil
	}

	if node.isLeaf(){
		node.deleteFromLeaf(index)
	}else{
		affectedNodes,affectedIndices,err:=node.deleteFromInternal(index)
		if err!=nil{
			return err
		}
		parentIndices=append(parentIndices, affectedIndices...)
		parents=append(parents,affectedNodes...)
	}
	
	for i:=len(parents)-2; i>=0; i--{
		pnode:=parents[i]
		node:=parents[i+1]
		
		if node.isUnderPopulated(){
			err:=pnode.rebalance(node,parentIndices[i])
			if err!=nil{
				return err
			}
		}
	}
	rootNode:=parents[0]
	if len(rootNode.items)==0 && len(rootNode.childNodes)>0{
		c.rootPgNum=parents[1].pageNum
		rootNode.deleteNode(rootNode.pageNum)
	}
	return nil
}