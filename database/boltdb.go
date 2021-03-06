// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package database

import (
    "fmt"
    "bytes"
	"github.com/FactomProject/simplecoin"
    "github.com/boltdb/bolt"
)

// This database stores and retrieves IBlock instances.  To do that, it
// needs a list of buckets that the using function wants, so it can make sure
// all those buckets exist.  (Avoids checking and building buckets in every 
// write).  
//
// It also needs a map of a hash to a IBlock instance.  To support this, 
// every block needs to be able to give the database a Hash for its type.
// This has to match the reverse, where looking up the hash gives the 
// database the type for the hash.  This way, the database can marshal
// and unmarshal IBlocks for storage in the database.  And since the IBlocks
// can provide the hash, we don't need two maps.  Just the Hash to the
// IBlock.
type BoltDB struct {
	SCDatabase
    
    db          *bolt.DB                        // Pointer to the bolt db
    instances   map[[32]byte]simplecoin.IBlock  // Maps a hash to an instance of an IBlock
}

var _ ISCDatabase = (*BoltDB)(nil)
// We have to make accomadation for many Init functions.  But what we really
// want here is:
//
//      Init(bucketList [][]byte, instances map[[32]byte]IBlock)
//
func (d *BoltDB) Init(a ...interface{}) {
    simplecoin.Prtln("NEED TO CONFIGURE DB")
    
    bucketList := a[0].([][]byte)
    instances  := a[1].(map[[32]byte]simplecoin.IBlock)
    
    tdb, err := bolt.Open("/tmp/bolt/my.db", 0600, nil)
    d.db = tdb
    
    if err != nil {
        panic("Database was not found, and could not be created.")
    }
    
    for _,bucket := range bucketList {
        d.db.Update(func(tx *bolt.Tx) error {
            _, err := tx.CreateBucketIfNotExists(bucket)
            if err != nil {
                return fmt.Errorf("create bucket: %s", err)
            }
            return nil
        })
    }
    
    d.instances = instances
}

func (d *BoltDB) Close() {
    d.db.Close()
}

func (d *BoltDB) GetRaw(bucket []byte, key []byte) (value simplecoin.IBlock) {
    var v []byte
    d.db.View(func(tx *bolt.Tx) error {
        b := tx.Bucket(bucket)
        v = b.Get([]byte(key))
        return nil
    })
    var vv[32]byte
    copy(vv[:],v)
    var instance simplecoin.IBlock = d.instances[vv]
    if instance == nil {
        panic("This should not happen.  Object stored in the database has no IBlock instance")
    }
    
    r := instance.GetNewInstance()
    r.UnmarshalBinaryData(v)
    
    return r
}


func (d *BoltDB) PutRaw(bucket []byte, key []byte, value simplecoin.IBlock) {
    var out bytes.Buffer
    hash := value.GetDBHash()
    out.Write(hash.Bytes())
    data, err := value.MarshalBinary()
    out.Write(data)
    
    if err != nil {
        panic("This should not happen.  Failed to marshal IBlock for BoltDB")
    }
    d.db.Update(func(tx *bolt.Tx) error {
        b := tx.Bucket(bucket)
        err := b.Put(key, data)
    return err
    })
}

func (db *BoltDB) Get(bucket string, key simplecoin.IHash) (value simplecoin.IBlock) {
    return db.GetRaw([]byte(bucket), key.Bytes())
}

func (db *BoltDB) GetKey(key IDBKey) (value simplecoin.IBlock) {
    return db.GetRaw(key.GetBucket(),key.GetKey())
}

func (db *BoltDB) Put(bucket string, key simplecoin.IHash, value simplecoin.IBlock) {
    b := []byte(bucket)
    k := key.Bytes()
    db.PutRaw(b, k, value)
}

func (db *BoltDB) PutKey(key IDBKey, value simplecoin.IBlock) {
    db.PutRaw(key.GetBucket(), key.GetKey(), value)
}