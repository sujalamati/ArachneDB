package main

import (
	"fmt"
)

func main(){
	// dal,_:=newDal("arachne.adb")
	// p:=dal.allocateEmptyPage()
	// p.num=dal.getNextPage()
	// copy(p.data[:],"rohit taparia")
	// _=dal.writePage(p)
	// _,_=dal.writeFreeList()

	// _=dal.close()

	// dal,_=newDal("arachne.adb")
	// p = dal.allocateEmptyPage()
	// p.num=dal.getNextPage()
	// copy(p.data[:],"akash prasad")
	// _=dal.writePage(p)

	// _,_=dal.writeFreeList()

	//testing read

	// dal,_:=newDal("arachne.adb")
	// p,err:=dal.readPage(3)
	// if err!=nil{
	// 	fmt.Print("error while reading page")
	// }	
	// fmt.Print(p.data)
	d,_:=newDal("./mainTest")
	rootNode,_ := d.getNode(d.rootNode)
	rootNode.dal=d
	containingNode, index, _ := rootNode.searchNode([]byte("Key1"))
	res := containingNode.items[index]

	fmt.Printf("key is: %s, value is: %s", res.key, res.value)
	// Close the db
	_ = d.close()

}