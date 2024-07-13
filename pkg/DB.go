// this file contains all wrapper methods (wrapping dal) that will be called by user.
package ArachneDB

import (
	"os"
	"sync"
)

type DB struct{
	*dal
	rwlock sync.RWMutex
}

func Open(path string, options *Options) (*DB, error){
	var err error
	options.pageSize = os.Getpagesize()

	dal,err := newDal(path,options)
	if err!=nil{
		return nil,err
	}

	DB:=&DB{
		dal,
		sync.RWMutex{},
	}
	return DB,nil
}

func (db *DB) Close() error{
	return db.close()
}

func (db *DB) WriteTx() *tx{
	db.rwlock.Lock()
	tx := newTx(db,true)
	db.dal.masterCollection.tx = tx
	return tx
}

func (db *DB) ReadTx() *tx{
	db.rwlock.RLock()
	tx := newTx(db,false)
	db.dal.masterCollection.tx = tx
	return newTx(db,false)
}