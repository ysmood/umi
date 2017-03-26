package lib

import (
	"bytes"
	"encoding/gob"

	"github.com/syndtr/goleveldb/leveldb"
)

// DB ...
type DB struct {
	level *leveldb.DB
}

// Open ...
func Open(path string) (*DB, error) {
	level, err := leveldb.OpenFile(path, nil)

	if err != nil {
		return nil, err
	}

	return &DB{
		level: level,
	}, nil
}

// Set ...
func (db *DB) Set(key string, val interface{}) {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	enc.Encode(val)
	db.level.Put([]byte(key), buf.Bytes(), nil)
}

// Get ...
func (db *DB) Get(key string, val interface{}) {
	data, _ := db.level.Get([]byte(key), nil)

	buf := bytes.Buffer{}
	buf.Write(data)

	dec := gob.NewDecoder(&buf)
	dec.Decode(val)
}
