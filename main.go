package main

import (
	"fmt"
	"os"
)

func main() {
	options := &Options{
		MaxFillPercent: 0.025,
		MinFillPercent: 0.0125,
		pageSize:       os.Getpagesize(),
	}

	dal, _ := newDal("mainTest2.adb", options)

	c := newEmptyCollection([]byte("Student1"), dal.rootNode)
	c.dal = dal

	_ = c.Put([]byte("fname"), []byte("Sujal"))
	_ = c.Put([]byte("lname"), []byte("Amati"))
	_ = c.Put([]byte("Age"), []byte("20"))
	_ = c.Put([]byte("Gender"), []byte("Male"))
	_ = c.Put([]byte("Key7"), []byte("Value7"))
	_ = c.Put([]byte("Key8"), []byte("Value8"))
	_ = c.Put([]byte("Key9"), []byte("Value9"))
	_ = c.Put([]byte("Key10"), []byte("Value10"))
	_ = c.Put([]byte("Key11"), []byte("Value11"))
	_ = c.Put([]byte("Key12"), []byte("Value12"))
	_ = c.Put([]byte("Key13"), []byte("Value13"))
	_ = c.Remove([]byte("Age"))
	i, _ := c.Find([]byte("Age"))
	fmt.Printf("item is: %+v\n", i)

	i,_= c.Find([]byte("Gender"))
	fmt.Printf("Key is: %s , Value is: %s \n",i.key,i.value)
	i,_= c.Find([]byte("Caste"))
	fmt.Printf("Key is: %s , Value is: %s \n",i.key,i.value)

	dal.rootNode = c.rootPgNum
	// Close the db
	dal.writeMeta()
	dal.writeFreeList()
	_ = dal.close()
}
