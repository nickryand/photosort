package main

import (
	"sort"
)

type Cache struct {
	data      sort.StringSlice
	duplicate int
	failure   int
	success   int
}

func (c *Cache) IsCached(s string) (bool, int) {
	var found bool

	index := c.data.Search(s)
	if index < len(c.data) && c.data[index] == s {
		// Found it
		found = true
	} else {
		found = false
	}
	return found, index
}

func (c *Cache) Insert(s string, index int) {
	// Insert: https://code.google.com/p/go-wiki/wiki/SliceTricks
	// This code uses the append function to grow the slice by 1 element
	// and has the side effect of growing the underlying array if neccessary
	c.data = append(c.data, s)
	copy(c.data[index+1:], c.data[index:])
	c.data[index] = s
}
