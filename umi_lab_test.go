package umi_test

import (
	"bytes"
	"encoding/gob"
	"testing"
	"unsafe"
)

type S struct {
	Content string
}

func noop(any interface{}) {}

func TestLab(t *testing.T) {
	// cache := umi.Cache{}

	// a := &S{}

	// cache.Add("ok", a)

	a := &S{}
	var b int32
	var cc complex128 = 4 + 5i
	s := []byte{'1', '2'}
	ss := "aa"

	res := []interface{}{
		unsafe.Sizeof(a),
		unsafe.Sizeof(b),
		unsafe.Sizeof(s),
		unsafe.Sizeof(ss),
		unsafe.Sizeof(cc),
		unsafe.Sizeof(map[string]int32{
			"a": 12,
			"b": 12,
		}),
		unsafe.Sizeof([]int32{1, 2}),
		unsafe.Sizeof(struct {
			a int8
			// b struct {
			// 	int32
			// }
		}{}),
	}

	noop(res)
}

func TestLabGob(t *testing.T) {
	s := S{
		Content: "testing",
	}

	var b bytes.Buffer

	enc := gob.NewEncoder(&b)

	enc.Encode(s)

	var bb bytes.Buffer
	bb.Write(b.Bytes())
	dec := gob.NewDecoder(&bb)
	var o S
	dec.Decode(&o)
}
