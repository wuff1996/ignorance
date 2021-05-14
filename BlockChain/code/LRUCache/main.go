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
	lru.Put(3, 66)
	fmt.Println(lru.Get(3))
	fmt.Println(lru.Get(3))
	fmt.Println(lru.Get(3))

	fmt.Println(lru.header.before.before.before.value)
	fmt.Println(lru.kvMap)
	h := lru.header
	for i := 0; i < len(lru.kvMap); i++ {
		fmt.Printf("value: %8d", h.key)
		h = h.after
	}

}

//this is a node that has its cap, map, header (which used to update the least recently used data and locate the tail of the node)
type LRUCache struct {
	capacity int
	kvMap    map[int]*Link
	//Link
	header *Link
}

//this is a DoubleLink that is easy to change a node
type Link struct {
	after  *Link
	before *Link
	value  int
	key    int
}

//get a LRUCache object
func Constructor(capacity int) LRUCache {
	l := &LRUCache{
		capacity: capacity,
		kvMap:    make(map[int]*Link),
	}
	return *l
}

func (this *LRUCache) updateHead(head *Link) {
	this.header = head
}

func (l *Link) suicide() {
	if l == nil {
		return
	}
	head, tail := l.before, l.after
	head.after, tail.before = tail, head
	l.after, l.before = nil, nil
}

func (l *Link) addBefore(n *Link) {
	if l == nil || n == nil || l == n {
		return
	}
	head := n.before
	head.after = l
	l.before, l.after = head, n
	n.before = l
}

func (l *Link) addToHead(head *Link) *Link {
	if l == head {
		return l
	}
	l.suicide()
	l.addBefore(head)
	return l
}

//To find a specific key
func (this *LRUCache) Get(key int) int {
	if link, ok := this.kvMap[key]; ok {
		header := link.addToHead(this.header)
		this.updateHead(header)
		return link.value
	}
	return -1
}

//To add a Cache to the LRUCache
func (this *LRUCache) Put(key int, value int) {
	if this.header == nil {
		link := &Link{value: value, key: key}
		link.before, link.after = link, link
		this.updateHead(link)
		this.kvMap[key] = link
		return
	}
	if link, ok := this.kvMap[key]; ok {
		link.value = value
		header := link.addToHead(this.header)
		this.updateHead(header)
		return
	}
	if this.capacity > len(this.kvMap) {
		link := &Link{value: value, key: key}
		link.addBefore(this.header)
		this.updateHead(link)
		this.kvMap[key] = link
		return
	}

	old := this.kvMap[this.header.before.key]
	delete(this.kvMap, this.header.before.key)

	//reuse this old link
	old.key,old.value=key,value
	if this.header==this.header.before{
		this.kvMap[key] = old
		return
	}
	old.suicide()
	old.addBefore(this.header)
	this.updateHead(old)
	this.kvMap[key] = old
	return

}
