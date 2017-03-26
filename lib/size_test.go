package lib_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ysmood/umi/lib"
)

const (
	sizeMap    = 8
	sizeString = 16
	sizeSlice  = 24
)

func TestString(t *testing.T) {
	assert.Equal(t, uintptr(22), lib.Size("我们"))
}

func TestFloat(t *testing.T) {
	assert.Equal(t, uintptr(8), lib.Size(float64(100)))
}

func TestByte(t *testing.T) {
	assert.Equal(t, uintptr(30), lib.Size([]byte("我们")))
}

func TestInt32(t *testing.T) {
	assert.Equal(t, uintptr(4), lib.Size(int32(123)))
}

func TestInt64(t *testing.T) {
	assert.Equal(t, uintptr(8), lib.Size(123))
}

func TestStruct(t *testing.T) {
	item := struct {
		A int32
		B string
	}{
		A: 11,
		B: "test",
	}

	assert.Equal(t, uintptr(24), lib.Size(item))
}

func TestPointer(t *testing.T) {
	item := uint64(10)
	p := &item

	assert.Equal(t, 8+8, int(lib.Size(p)))
}

func TestComplexStruct(t *testing.T) {
	data := &struct {
		A interface{}
	}{}

	assert.Equal(t, 8, int(lib.Size(data)))
}

func TestStructWithPrivateProp(t *testing.T) {
	defer func() {
		r := recover()

		if r == nil {
			panic("should panic")
		}
	}()

	// the size of private fields can't be auto calculated
	item := struct {
		a int32
	}{}

	lib.Size(item)
}

func TestMap(t *testing.T) {
	item := map[string]int32{
		"a": 10,
		"b": 20,
	}

	assert.Equal(t, uintptr(50), lib.Size(item))
}

func TestArr(t *testing.T) {
	item := []int32{
		10,
		20,
	}

	assert.Equal(t, uintptr(32), lib.Size(item))
}

func TestNestedStruct(t *testing.T) {
	item := struct {
		A struct {
			A string
			C int32
			D struct {
				E []byte
			}
		}
		F map[string]string
	}{
		F: map[string]string{
			"ok": "我们",
		},
	}

	assert.Equal(t, uintptr(92), lib.Size(item))
}
