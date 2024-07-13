package ArachneDB

import (
	"bytes"
	"encoding/binary"
)

type Item struct {
	key   []byte
	value []byte
}

type Node struct {
	// *dal // type embedding dal in node

	pageNum    pgnum //each node is stored in one page
	items      []*Item
	childNodes []pgnum

	tx *tx		//associated transaction
}

func newEmptyNode() *Node {
	return &Node{}
}

func NewNodeForSerialization(items []*Item, childNodes []pgnum) *Node {
	return &Node{
		childNodes: childNodes,
		items:      items,
	}
}

func newItem(key []byte, value []byte) *Item {
	return &Item{
		key:   key,
		value: value,
	}
}

func (i *Item) Key() []byte{
	return i.key
}

func (i *Item) Value() []byte{
	return i.value
}

func (n *Node) isLeaf() bool {
	return len(n.childNodes) == 0
}

func(n *Node) getNode(pageNum pgnum) (*Node,error) {
	return n.tx.getNode(pageNum)
}

func(n *Node) deleteNode(pageNum pgnum){
	n.tx.deleteNode(pageNum)
}

func (n *Node) serialize(buf []byte) {
	// using the concept of slotted pages(like PostgreSQL) to store the node in each page

	leftPos := 0
	rightPos := len(buf) - 1

	// Page structure is
	// | Page Header | Child Pointer 1 | Offset 1 | Child Pointer 2 | Offset 2 |  .....
	// | Offset K | Child Pointer K+1 | ...
	// ...
	// ..........| Key K | Value K | Key K-1 | Value K-1 | ......| Key 2 | Value 2 | Key 1 | Value 1 |

	// Page Header includes isLeaf , K-V count

	isLeaf := n.isLeaf()

	var bitSetVar uint64 = 0
	if isLeaf {
		bitSetVar = 1 //set bitSetVar = 1 if it is a leaf node
	}
	// Write the page header (isLeaf and No of K-V pairs)
	buf[leftPos] = byte(bitSetVar)
	leftPos += 1

	// write the no of K-V pairs
	binary.LittleEndian.PutUint16(buf[leftPos:], uint16(len(n.items)))
	leftPos += 2

	for i := 0; i < len(n.items); i++ {
		item := n.items[i]
		if !isLeaf {

			//write the ith child node
			binary.LittleEndian.PutUint64(buf[leftPos:], uint64(n.childNodes[i]))
			leftPos += pageNumSize
		}

		vlen := len(item.value)
		klen := len(item.key)

		//write the value
		rightPos = rightPos - vlen
		copy(buf[rightPos:], item.value)

		//write the len of value
		rightPos -= 1
		buf[rightPos] = byte(vlen)

		//write the key
		rightPos = rightPos - klen
		copy(buf[rightPos:], item.key)

		//write the len of key
		rightPos -= 1
		buf[rightPos] = byte(klen)

		//writing the offset where the Key-Value pair is stored
		offset := rightPos
		binary.LittleEndian.PutUint16(buf[leftPos:], uint16(offset))
		leftPos += 2
	}

	if !isLeaf {
		//write the last child node
		lastChildNode := n.childNodes[len(n.childNodes)-1]

		binary.LittleEndian.PutUint64(buf[leftPos:], uint64(lastChildNode))
		leftPos += pageNumSize
	}
}

func (n *Node) deserialize(buf []byte) {
	leftPos := 0

	//Read page header

	//read isLeaf
	isLeaf := uint16(buf[leftPos])
	leftPos += 1

	// read the no of K-V pairs
	noOfRec := binary.LittleEndian.Uint16(buf[leftPos:])
	leftPos += 2

	for i := 0; i < int(noOfRec); i++ {
		if isLeaf == 0 {
			// read the ith child node
			childNode := binary.LittleEndian.Uint64(buf[leftPos:])
			// append ith child node
			n.childNodes = append(n.childNodes, pgnum(childNode))
			leftPos += pageNumSize
		}
		// read offset
		offset := binary.LittleEndian.Uint16(buf[leftPos:])
		leftPos += 2

		//read len of key
		klen := uint16(buf[offset])
		offset = offset + 1

		// read key
		key := buf[offset : offset+klen]
		offset += klen

		// read len of value
		vlen := uint16(buf[offset])
		offset += 1

		// read value
		value := buf[offset : offset+vlen]
		offset += vlen

		item := newItem(key, value)
		n.items = append(n.items, item)
	}
	if isLeaf == 0 {
		lastChildNode := binary.LittleEndian.Uint64(buf[leftPos:])
		n.childNodes = append(n.childNodes, pgnum(lastChildNode))
	}
}

// searchNode searches for a key inside the tree. Once the key is found, the node and the correct index are returned
// so the key itself can be accessed in the following way node[index]. A list of the parent nodes is also returned.
// If the key isn't found, we have 2 options. If mode is false, it means we expect searchNode
// to find the key. If mode is true, then searchNode is used to locate where a new key should be
// inserted so the position is returned.

func (n *Node) searchNode(key []byte, mode bool) (bool, *Node, int, []int, []*Node, error) {
	
	parentIndices := []int{}
	parents := []*Node{}
	wasFound, node, index, parentIndices, parents, err := searchNodeRec(n, key, parentIndices, parents, mode)

	if err != nil {
		return false, nil, -1, []int{}, []*Node{}, err
	}

	return wasFound, node, index, parentIndices, parents, nil
}

func searchNodeRec(node *Node, key []byte, parentIndices []int, parents []*Node, mode bool) (bool, *Node, int, []int, []*Node, error) {

	// search for key in the node
	found, index := node.searchInNode(key)
	// append index to parentIndices
	parentIndices = append(parentIndices, index)
	// append node to parentNodes
	parents = append(parents, node)
	if found {
		return true, node, index, parentIndices, parents, nil
	}

	if node.isLeaf() {
		// if the node is leaf and the key is not found => key does not exist
		if mode {
			return false, node, index, parentIndices, parents, nil
		}
		return false, nil, -1, []int{}, []*Node{}, nil
	}

	// fetch the child node where the key is present , from disk into memory
	child, err := node.getNode(node.childNodes[index])

	if err != nil {
		return false, nil, -1, []int{}, []*Node{}, err
	}

	// search in the child node
	return searchNodeRec(child, key, parentIndices, parents, mode)

}

func (n *Node) searchInNode(key []byte) (bool, int) {
	// NOTE: we are storing keys in lexicographical order

	for i, item := range n.items {
		flag := bytes.Compare(item.key, key)
		// flag = 0 => item.key and key matches
		// flag = 1 => item.key is lexicographically greater than key

		if flag == 0 {
			// match found
			return true, i
		}
		if flag > 0 {
			// item.key is lexicographically greater than key
			return false, i
		}
	}
	// desiredKey is greater than all the keys in the node => it is in the last child node
	return false, len(n.items)
}

func (n *Node) insertInNode(item *Item, index int) {
	if index == len(n.items) {
		n.items = append(n.items, item)
		return
	}
	// move all the items starting from index to right by one place
	n.items = append(n.items[:index+1], n.items[index:]...)

	//insert item at index
	n.items[index] = item
}

func (n *Node) isOverPopulated() bool {
	return n.tx.db.isOverPopulated(n)
}

func (n *Node) isUnderPopulated() bool {
	return n.tx.db.isUnderPopulated(n)
}

func (n *Node) nodeSize() int { // returns size of each node
	size := 0
	size += nodeHeaderSize // size of Page Header

	for i := range n.items {
		size += n.itemSize(i) // size of each Item
	}
	size += pageNumSize // size of last ChildNode
	return size
}

func (n *Node) itemSize(i int) int {
	size := 0
	size += len(n.items[i].key) + 1   // size of key
	size += len(n.items[i].value) + 1 // size of value
	size += 2                         // size of offset
	size += pageNumSize               // size of childNode
	return size
}

func (n *Node) split(node *Node, parentSplitIndex int) {
	splitIndex := node.tx.db.getSplitIndex(node)

	midItem := node.items[splitIndex]

	var newNode *Node
	if node.isLeaf() {
		newNode = n.writeNode(n.tx.newNode(node.items[splitIndex+1:], []pgnum{}))
		node.items = node.items[:splitIndex]
	} else {
		newNode = n.writeNode(n.tx.newNode(node.items[splitIndex+1:], node.childNodes[splitIndex+1:]))
		node.items = node.items[:splitIndex]
		node.childNodes = node.childNodes[:splitIndex+1]
	}
	n.insertInNode(midItem, parentSplitIndex)
	if parentSplitIndex == len(n.childNodes) {
		n.childNodes = append(n.childNodes, newNode.pageNum)
	} else {
		n.childNodes = append(n.childNodes[:parentSplitIndex+1], n.childNodes[parentSplitIndex:]...)
		n.childNodes[parentSplitIndex+1] = newNode.pageNum
	}
	n.writeNode(n)
	n.writeNode(node)
}

func (n *Node) writeNode(node *Node) *Node {
	return n.tx.writeNode(node)
}

func (n *Node) writeNodes(nodes ...*Node) {
	for _, node := range nodes {
		n.writeNode(node)
	}
}

func (n *Node) deleteFromLeaf(index int) {
	n.items = append(n.items[:index], n.items[index+1:]...)
	n.writeNode(n)
}

func (n *Node) canSpareKey() bool {
	return n.tx.db.getSplitIndex(n) != -1
}

func (n *Node) deleteFromInternal(index int) ([]*Node, []int, error) {
	childNode, err := n.getNode(n.childNodes[index])

	if err!=nil{
		return nil,nil,err
	}

	var affectedNodes []*Node = []*Node{childNode}
	var affectedIndices []int = []int{}

	for !childNode.isLeaf() {
		nextIndex := len(childNode.childNodes) - 1
		nextChildPage := childNode.childNodes[nextIndex]
		childNode, err = n.getNode(nextChildPage)
		if err != nil {
			return nil, nil, err
		}
		affectedNodes = append(affectedNodes, childNode)
		affectedIndices = append(affectedIndices, nextIndex)
	}
	// replace the key with its immediate predecessor
	pred := childNode.items[len(childNode.items)-1]

	n.items[index] = pred
	//delete the immediate predecessor from its original node
	childNode.items = childNode.items[:len(childNode.items)-1]
	n.writeNode(childNode)
	n.writeNode(n)
	return affectedNodes, affectedIndices, nil
}

func (n *Node) rebalance(node *Node, index int) error {

	if index != 0 {
		leftNode, err := n.getNode(n.childNodes[index-1])
		if err != nil {
			return err
		}
		if leftNode.canSpareKey() {
			// perform a right rotation
			rightRotate(leftNode, node, n, index)
			n.writeNodes(n, leftNode, node)
			return nil
		}
	}

	if index != len(n.childNodes)-1 {
		rightNode, err := n.getNode(n.childNodes[index+1])
		if err != nil {
			return err
		}
		if rightNode.canSpareKey() {
			//perform a left rotation
			leftRotate(node, rightNode, n, index)
			n.writeNodes(rightNode, n, node)
			return nil
		}
	}
	// The merge function merges a given node with its node to the left. So by default, we merge an unbalanced node
	// with its left sibling. In the case where the unbalanced node is the leftmost, we have to change the merge
	// parameters, so the unbalanced node right sibling, will be merged into the unbalanced node.
	if index == 0 {
		rightNode, err := n.getNode(n.childNodes[index+1])
		if err != nil {
			return err
		}
		merge(node, rightNode, n, index+1)
		return nil
	}
	leftNode, err := n.getNode(n.childNodes[index-1])
	if err != nil {
		return err
	}
	merge(leftNode, node, n, index)
	return nil
}

func rightRotate(lNode, rNode, pNode *Node, index int) {
	// 	           p                                    p
	//                 4                                    3
	//	      /        \           ------>         /          \
	//	   a           b (unbalanced)            a        b (unbalanced)
	//      1,2,3             5                     1,2            4,5

	leftItem := lNode.items[len(lNode.items)-1]
	parentItem := pNode.items[index-1]

	rNode.insertInNode(parentItem, 0)
	pNode.items[index-1] = leftItem
	lNode.items = lNode.items[:len(lNode.items)-1]

	if !lNode.isLeaf() {
		lChild := lNode.childNodes[len(lNode.childNodes)-1]
		lNode.childNodes = lNode.childNodes[:len(lNode.childNodes)-1]
		rNode.childNodes = append([]pgnum{lChild}, rNode.childNodes...)
	}
}

func leftRotate(lNode, rNode, pNode *Node, index int) {
	// 	           p                                     p
	//                 2                                     3
	//	      /        \           ------>         /          \
	//  a(unbalanced)       b                 a(unbalanced)        b
	//   1                3,4,5                   1,2             4,5

	rightItem := rNode.items[0]
	parentItem := pNode.items[index]

	lNode.items = append(lNode.items, parentItem)
	pNode.items[index-1] = rightItem
	rNode.items = rNode.items[1:]

	if !rNode.isLeaf() {
		rChild := rNode.childNodes[0]
		rNode.childNodes = rNode.childNodes[1:]
		lNode.childNodes = append(lNode.childNodes, rChild)
	}
}

func merge(lNode, rNode, pNode *Node, index int) {
	// 	               p                                     p
	//                3,5                                    5
	//	      /        |       \       ------>         /          \
	//       l   	   r        c                     l            c
	//     1,2         4        6,7                 1,2,3,4         6,7

	// Take the item from the parent, remove it and add it to the unbalanced node
	pItem := pNode.items[index-1]
	lNode.items = append(lNode.items, pItem)
	lNode.items = append(lNode.items, rNode.items...)
	pNode.childNodes = append(pNode.childNodes[:index], pNode.childNodes[index+1:]...)
	pNode.items = append(pNode.items[:index-1], pNode.items[index:]...)

	if !lNode.isLeaf() {
		lNode.childNodes = append(lNode.childNodes, rNode.childNodes...)
	}
	pNode.writeNodes(pNode, lNode)
	pNode.deleteNode(rNode.pageNum)

}