package umi

import "time"

// Options ...
type Options struct {
	// unit byte, default 100M
	MaxMemSize uint64

	// default 1sec, GC every 1sec
	// if less than zero, gc will be disabled
	GCSpan time.Duration

	// default 10, each GC round check 10 items from the tail
	// if less than zero, gc will be disabled
	GCSize int

	// default 1min
	TTL time.Duration

	OnEvicted func(*Item)
}

func defaultOptions(opts *Options) *Options {
	if opts == nil {
		opts = &Options{}
	}

	if opts.MaxMemSize == 0 {
		opts.MaxMemSize = 100 * 1024 * 1024
	}
	if opts.GCSize == 0 {
		opts.GCSize = 10
	}
	if opts.GCSpan == 0 {
		opts.GCSpan = time.Second
	}
	if opts.TTL == 0 {
		opts.TTL = time.Minute
	}

	return opts
}
