package lib_test

import (
	"testing"

	"github.com/ysmood/umi/lib"

	"github.com/stretchr/testify/assert"
)

func TestDB(t *testing.T) {
	type ss struct {
		A string
	}
	item := ss{
		A: "test",
	}

	db, _ := lib.Open("test.db")

	db.Set("t", item)
	var newItem ss
	db.Get("t", &newItem)

	assert.Equal(t, item.A, newItem.A)
}
