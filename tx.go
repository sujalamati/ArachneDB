package ArachneDB

type tx struct{

	dirtyNodes map[pgnum]*Node	//stores all the dirty nodes in memory
	pagesToDelete []pgnum		//nodes that are deleted in transaction
	allocatedPageNums []pgnum	//list of all new pages allocated in transaction

	write bool	//read tx or write tx

	db *DB
}

func newTx(db *DB,write bool) *tx{
	return &tx{
		map[pgnum]*Node{},
		make([]pgnum, 0),
		make([]pgnum, 0),
		write,
		db,
	}
}

func (tx *tx) Commit() error{
	if !tx.write{
		tx.db.rwlock.RUnlock()
		return nil
	}

	for _,node := range tx.dirtyNodes{
		_,err := tx.db.writeNode(node)
		if err!=nil{
			return err
		}
	}

	for _,pageNum := range tx.pagesToDelete{
		tx.db.deleteNode(pageNum)
	}
	_,err:=tx.db.writeFreeList()
	if err!=nil {
		return err
	}
	tx.dirtyNodes = nil
	tx.allocatedPageNums = nil
	tx.pagesToDelete = nil
	tx.db.rwlock.Unlock()
	return nil
}

func (tx *tx) Rollback(){
	if !tx.write{
		tx.db.rwlock.RUnlock()
		return
	}

	for _,pageNum := range tx.allocatedPageNums{
		tx.db.freelist.releasePage(pageNum)
	}
	tx.allocatedPageNums = nil
	tx.dirtyNodes = nil
	tx.pagesToDelete = nil
	tx.db.rwlock.Unlock()
}

func (tx *tx) newNode(items []*Item,childNodes []pgnum)(*Node){
	node:=newEmptyNode()
	node.tx = tx
	node.childNodes = childNodes
	node.items = items
	node.pageNum = tx.db.getNextPage()

	tx.allocatedPageNums = append(tx.allocatedPageNums, node.pageNum)
	return node
}

func (tx *tx) writeNode(n *Node) *Node{
	tx.dirtyNodes[n.pageNum] = n
	n.tx = tx
	return n
}

func(tx *tx) getNode(pageNum pgnum) (*Node,error){
	node,ok := tx.dirtyNodes[pageNum]
	if ok{
		return node,nil
	}

	node,err:= tx.db.getNode(pageNum)
	if err!=nil{
		return nil,err
	}
	node.tx = tx
	return node,nil
}

func (tx *tx) deleteNode(pageNum pgnum){
	tx.pagesToDelete = append(tx.pagesToDelete, pageNum)
}


func (tx *tx) CreateCollection(name []byte) (*Collection,error){
	if !tx.write {
		return nil, writeInsideReadTxErr
	}

	newNode,err:=tx.db.writeNode(tx.newNode([]*Item{},[]pgnum{}))
	if err!=nil{
		return nil,err
	}
 	collection:=newEmptyCollection(name,newNode.pageNum)
	// collection.dal=d
	collection.tx = tx
	// fmt.Println(collection.rootPgNum)
 	tx.db.masterCollection.Put(name,[]byte{byte(collection.rootPgNum)})
	return collection,nil
}
 
func (tx *tx) GetCollection(name []byte) (*Collection,error){
	i,err:=tx.db.masterCollection.Find(name)
	if err!=nil{
		return nil,err
	}
	if i==nil{
		return nil,nil
	}else{
		collection:= newEmptyCollection(name,pgnum(i.value[0]))
		// collection.dal=d
		collection.tx = tx
		return collection,nil
	}
}

func (tx *tx) DeleteCollection(name []byte) (error){
	if !tx.write {
		return writeInsideReadTxErr
	}

	return tx.db.masterCollection.Remove(name)
}