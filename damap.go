// Package damap provides a map data structure on double-array trie tree and functions to search values from trie tree.
package damap

const endKey rune = 0

// DaMap is a map data structure based on double-array trie tree.
type DaMap struct {
	Base  []int
	Check []int
	Value []interface{}
}

// New creates a new empty DaMap.
func New() *DaMap {
	d := &DaMap{}
	d.allocMinSize(1)

	return d
}

// Write is used to insert value with key.
func (d *DaMap) Write(key string, value interface{}) {
	s := []rune(key)
	s = append(s, endKey)

	i := 0
	for _, k := range s {
		i = d.insert(i, k)
	}

	d.Value[i] = value
}

func (d *DaMap) insert(i int, k rune) int {
	base := 1
	if d.Base[i] != 0 {
		base = d.Base[i]
	}

	n := base + int(k) - 1
	d.allocMinSize(n + 1)

	if d.Check[n] == 0 && d.Base[n] == 0 {
		d.Base[i] = base
		d.Check[n] = i + 1
		return n
	}

	if d.Check[n] == i+1 {
		return n
	}

	var conflicts []struct {
		index int
		base  int
		check int
		value interface{}
	}

	for index, check := range d.Check {
		if check == i+1 {
			conflicts = append(
				conflicts,
				struct {
					index int
					base  int
					check int
					value interface{}
				}{
					index,
					d.Base[index],
					check,
					d.Value[index],
				},
			)
		}
	}

	s := []rune{k}
	for _, conflict := range conflicts {
		s = append(s, rune(conflict.index+1-d.Base[i]))
	}

	for {
		base++
		ok := true
		for _, k := range s {
			next := base + int(k) - 1
			d.allocMinSize(next + 1)

			ok = ok && d.Check[next] == 0 && d.Base[next] == 0
		}

		if ok {
			break
		}
	}

	for _, conflict := range conflicts {
		next := base + conflict.index + 1 - d.Base[i] - 1
		d.Base[next] = conflict.base
		d.Value[next] = conflict.value
		d.Check[next] = i + 1

		for index, check := range d.Check {
			if check == conflict.index+1 {
				d.Check[index] = next + 1
			}
		}

		d.Base[conflict.index] = 0
		d.Check[conflict.index] = 0
		d.Value[conflict.index] = nil
	}

	d.Base[i] = base

	n = base + int(k) - 1
	d.Check[n] = i + 1

	return n
}

func (d *DaMap) allocMinSize(size int) {
	if len(d.Base) < size {
		if cap(d.Base) < size {
			d.Base = append(d.Base, make([]int, size-len(d.Base))...)
		} else {
			d.Base = d.Base[:size]
		}
	}

	if len(d.Check) < size {
		if cap(d.Check) < size {
			d.Check = append(d.Check, make([]int, size-len(d.Check))...)
		} else {
			d.Check = d.Check[:size]
		}
	}

	if len(d.Value) < size {
		if cap(d.Value) < size {
			d.Value = append(d.Value, make([]interface{}, size-len(d.Value))...)
		} else {
			d.Value = d.Value[:size]
		}
	}
}

// Read is used to get value from key.
func (d *DaMap) Read(key string) interface{} {
	if len(key) == 0 {
		return false
	}

	s := []rune(key)
	s = append(s, endKey)

	i := 0
	for _, k := range s {
		n := d.Base[i] + int(k) - 1
		if d.Check[n] != i+1 {
			return nil
		}
		i = n
	}

	return d.Value[i]
}

// ExactMatchSearch provides exact-match-search function for trie tree.
func (d *DaMap) ExactMatchSearch(key string) bool {
	if len(key) == 0 {
		return false
	}

	s := []rune(key)
	s = append(s, endKey)

	i := 0
	for _, k := range s {
		n := d.Base[i] + int(k) - 1
		if d.Check[n] != i+1 {
			return false
		}
		i = n
	}

	if d.Base[i] != 0 {
		return false
	}

	return true
}

// CommonPrefixSearchResult is a value that return from CommonPrefixSearch function.
type CommonPrefixSearchResult []struct {
	Pos   int
	Key   string
	Value interface{}
}

// CommonPrefixSearch provides common prefix search functions for trie tree.
func (d *DaMap) CommonPrefixSearch(s string) CommonPrefixSearchResult {
	result := CommonPrefixSearchResult{}

	if len(s) == 0 {
		return result
	}

	for pos := 0; pos < len(s); pos++ {
		i := 0
		for j, k := range s[pos:] {
			n := d.Base[i] + int(k) - 1

			if len(d.Base) < n+1 {
				break
			}

			if d.Check[n] != i+1 {
				break
			}

			leaf := d.Base[n] + int(endKey) - 1
			if d.Base[leaf] == 0 && d.Check[leaf] == n+1 {
				result = append(
					result,
					struct {
						Pos   int
						Key   string
						Value interface{}
					}{
						Pos:   pos,
						Key:   string(s[pos : pos+j+1]),
						Value: d.Value[leaf],
					},
				)
			}

			i = n
		}
	}

	return result
}
