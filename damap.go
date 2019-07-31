// Package damap provides a map data structure on double-array trie tree and functions to search values from trie tree.
package damap

const endKey rune = 0

// DaMap is a map data structure based on double-array trie tree.
type DaMap struct {
	base  []int
	check []int
	value []interface{}
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

	d.value[i] = value
}

func (d *DaMap) insert(i int, k rune) int {
	base := 1
	if d.base[i] != 0 {
		base = d.base[i]
	}

	n := base + int(k) - 1
	d.allocMinSize(n + 1)

	if d.check[n] == 0 && d.base[n] == 0 {
		d.base[i] = base
		d.check[n] = i + 1
		return n
	}

	if d.check[n] == i+1 {
		return n
	}

	var conflicts []struct {
		index int
		base  int
		check int
		value interface{}
	}

	for index, check := range d.check {
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
					d.base[index],
					check,
					d.value[index],
				},
			)
		}
	}

	s := []rune{k}
	for _, conflict := range conflicts {
		s = append(s, rune(conflict.index+1-d.base[i]))
	}

	for {
		base++
		ok := true
		for _, k := range s {
			next := base + int(k) - 1
			d.allocMinSize(next + 1)

			ok = ok && d.check[next] == 0 && d.base[next] == 0
		}

		if ok {
			break
		}
	}

	for _, conflict := range conflicts {
		next := base + conflict.index + 1 - d.base[i] - 1
		d.base[next] = conflict.base
		d.value[next] = conflict.value
		d.check[next] = i + 1

		for index, check := range d.check {
			if check == conflict.index+1 {
				d.check[index] = next + 1
			}
		}

		d.base[conflict.index] = 0
		d.check[conflict.index] = 0
		d.value[conflict.index] = nil
	}

	d.base[i] = base

	n = base + int(k) - 1
	d.check[n] = i + 1

	return n
}

func (d *DaMap) allocMinSize(size int) {
	if len(d.base) < size {
		if cap(d.base) < size {
			d.base = append(d.base, make([]int, size-len(d.base))...)
		} else {
			d.base = d.base[:size]
		}
	}

	if len(d.check) < size {
		if cap(d.check) < size {
			d.check = append(d.check, make([]int, size-len(d.check))...)
		} else {
			d.check = d.check[:size]
		}
	}

	if len(d.value) < size {
		if cap(d.value) < size {
			d.value = append(d.value, make([]interface{}, size-len(d.value))...)
		} else {
			d.value = d.value[:size]
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
		n := d.base[i] + int(k) - 1
		if d.check[n] != i+1 {
			return nil
		}
		i = n
	}

	return d.value[i]
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
		n := d.base[i] + int(k) - 1
		if d.check[n] != i+1 {
			return false
		}
		i = n
	}

	if d.base[i] != 0 {
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
			n := d.base[i] + int(k) - 1

			if len(d.base) < n+1 {
				break
			}

			if d.check[n] != i+1 {
				break
			}

			leaf := d.base[n] + int(endKey) - 1
			if d.base[leaf] == 0 && d.check[leaf] == n+1 {
				result = append(
					result,
					struct {
						Pos   int
						Key   string
						Value interface{}
					}{
						Pos:   pos,
						Key:   string(s[pos : pos+j+1]),
						Value: d.value[leaf],
					},
				)
			}

			i = n
		}
	}

	return result
}
