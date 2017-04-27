package umi

import "time"

// Options ...
type Options struct {
	// unit byte, default 100M
	MaxMemSize uint64

	// default 1sec, GC every 1sec
	GCSpan time.Duration

	// default 10, each GC round check 10 items from the tail
	GCSize int

	// default 1min
	TTL time.Duration

	// default 100, promote after 100 hits,
	// if it is less than zero, will promote on each hit
	PromoteRate int32
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
	if opts.PromoteRate == 0 {
		opts.PromoteRate = 100
	}

	return opts
}
