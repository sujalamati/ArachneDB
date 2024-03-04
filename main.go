package main

import (
	"fmt"
	"os"
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
	
	// d,_:=newDal("./mainTest")
	// rootNode,_ := d.getNode(d.rootNode)
	// rootNode.dal=d
	// wasFound,containingNode,index,parentIndices,parents,err:= rootNode.searchNode([]byte("Key1"),false)
	// if err!=nil{
	// 	fmt.Printf(err.Error())
	// }
	// if wasFound{
	// 	fmt.Println(parentIndices)
	// 	fmt.Println(parents)
	// 	res := containingNode.items[index]
	// 	fmt.Printf("key is: %s, value is: %s", res.key, res.value)
	// }else{
	// 	fmt.Printf("key not found")
	// }
	
	

	options:=&Options{
		MaxFillPercent: 0.025,
		MinFillPercent: 0.0125,
		pageSize: os.Getpagesize(),
	}

	

	dal,_:=newDal("mainTest2.adb",options)
	
	c:=newEmptyCollection([]byte("Student1"),dal.rootNode)
	c.dal = dal

	// _ = c.Put([]byte("Key1"), []byte("Value1"))
	// _ = c.Put([]byte("Key2"), []byte("Value2"))
	// _ = c.Put([]byte("Key3"), []byte("Value3"))
	// _ = c.Put([]byte("Key4"), []byte("Value4"))
	// _ = c.Put([]byte("Key5"), []byte("Value5"))
	// _ = c.Put([]byte("Key6"), []byte("Value6"))
	// _ = c.Put([]byte("Key7"), []byte("Value7"))
	// _ = c.Put([]byte("Key8"), []byte("Value8"))
	// _ = c.Put([]byte("Key9"), []byte("Value9"))
	// _ = c.Put([]byte("Key10"), []byte("Value10"))
	// _ = c.Put([]byte("Key11"), []byte("Value11"))
	// _ = c.Put([]byte("Key12"), []byte("Value12"))
	// _ = c.Put([]byte("Key13"), []byte("Value13"))



	fmt.Println(dal.rootNode)
	for i:=1; i<=13; i++{
		var search string = fmt.Sprint("Key",i)
		fmt.Println(search)
		i,_:= c.Find([]byte(search))
		fmt.Printf("Key is: %s , Value is: %s \n",i.key,i.value)
	}
	

	dal.rootNode=c.rootPgNum
	// Close the db
	dal.writeMeta()
	dal.writeFreeList()
	_ = dal.close()
}