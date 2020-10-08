/*
This is my test for LRU (least recently used) working process
*/

package main

import "fmt"

func main() {
	lru := Constructor(5)

	lru.Put(1, 1)
	lru.Put(2, 2)
	lru.Put(3, 3)
	lru.Put(4, 4)
	lru.Put(5, 5)
	lru.Put(6, 6)
	lru.Put(3, 7)
	lru.Put(5435, 34)
	lru.Put(9, 3)
	lru.Get(3)
	fmt.Println(lru.Get(3))
	fmt.Println(lru.header.before.before.before.key)
	fmt.Println(lru.kvMap)
	h := lru.header
	for i := 0; i < len(lru.kvMap); i++ {
		fmt.Printf("%8d", h.key)
		h = h.after
	}

}

//this is a node that has its cap, map, header (which used to update the least recently used data and locate the tail of the node)
type LRUCache struct {
	capacity int
	kvMap    map[int]int
	//Link
	header *Link
}

//this is a DoubleLink that is easy to change a node
type Link struct {
	after  *Link
	before *Link
	key    int
}

//get a LRUCache object
func Constructor(capacity int) LRUCache {
	l := &LRUCache{
		capacity: capacity,
		kvMap:    make(map[int]int),
	}
	return *l
}

//To find a specific key
func (this *LRUCache) Get(key int) int {
	if v, ok := this.kvMap[key]; ok {
		if key == this.header.key {
			this.header = this.header.after
			return v
		} else if key == this.header.before.key {
			return v
		}
		after := this.header
		for i := 0; i < len(this.kvMap); i++ {
			if after.key == key {
				after.before.after = after.after
				after.after.before = after.before
				break
			}
			after = after.after
		}
		tail := this.header.before
		entry := &Link{key: key}
		entry.before = tail
		entry.after = this.header
		tail.after = entry
		this.header.before = entry
		return v
	}
	return -1
}

//To add a Cache to the LRUCache
func (this *LRUCache) Put(key int, value int) {
	if _, ok := this.kvMap[key]; ok {
		if this.header.key == key {
			this.header = this.header.after
			this.kvMap[key] = value
			return
		}
		if this.header.before.key == key {
			this.kvMap[key] = value
			return
		}
		after := this.header
		for i := 0; i < this.capacity; i++ {
			if after.key == key {
				after.before.after = after.after
				after.after.before = after.before
				break
			}
			after = after.after
		}
		tail := this.header.before
		entry := &Link{key: key}
		tail.after = entry
		entry.after = this.header
		entry.before = tail
		this.header.before = entry
		this.kvMap[key] = value
		return
	} else {
		if len(this.kvMap) == 0 {
			this.header = &Link{key: key}
			this.header.after = this.header
			this.header.before = this.header
			this.kvMap[key] = value
			return
		}
		if this.capacity > len(this.kvMap) {
			tail := this.header.before
			entry := &Link{key: key}
			tail.after = entry
			entry.after = this.header
			entry.before = tail
			this.header.before = entry
			this.kvMap[key] = value
			return
		} else {
			delete(this.kvMap, this.header.key)
			this.kvMap[key] = value
			entry := &Link{key: key}
			tail := this.header.before
			this.header = this.header.after
			this.header.before = entry
			entry.after = this.header
			entry.before = tail
			tail.after = entry
			return
		}
	}
}
