package umi_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/ysmood/umi"
)

func TestGC(t *testing.T) {
	c := umi.New(&umi.Options{
		TTL:    time.Millisecond * 20,
		GCSpan: time.Millisecond * 10,
	})

	c.Set("ok", "test")

	time.Sleep(time.Millisecond * 50)

	_, has := c.Get("ok")

	assert.Equal(t, false, has)
}
