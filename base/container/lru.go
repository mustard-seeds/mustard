package container

import (
	"container/list"
	"sync"
)

type Element struct {
	Value interface{}
}

type innerElement struct {
	key  string
	value *Element
}

type LRU struct {
	list *list.List
	index map[string]*list.Element
	Capacity int
	// for statistic
	totalReq int
	hitReq   int
	sync.RWMutex
}
func (lru *LRU)Get(key string) (*Element,bool) {
	lru.Lock()
	defer lru.Unlock()
	lru.totalReq += 1
	v,exist := lru.index[key]
	if !exist {
		return nil,exist
	}
	lru.hitReq += 1
	lru.list.MoveToFront(v)
	return v.Value.(*innerElement).value,true
}

func (lru *LRU)Set(key string, v interface{}) {
	lru.Lock()
	defer lru.Unlock()
	if lru.full() {
		last := lru.list.Back()
		delete(lru.index, last.Value.(*innerElement).key)
		lru.list.Remove(last)
	}
	lru.list.PushFront(&innerElement{
		key:key,
		value:&Element{
			Value:v,
		},
	})
	lru.index[key] = lru.list.Front()
}

func (lru *LRU)Size() int {
	return lru.list.Len()
}

func (lru *LRU)full() bool {
	return lru.Capacity == lru.Size()
}

func (lru *LRU)HitRatio() float32 {
	if lru.totalReq == 0 {
		return 0.0
	}
	return float32(lru.hitReq)/float32(lru.totalReq)
}


var _instance *LRU = nil
var _init_ctx sync.Once

func NewLRU(capacity int) *LRU {
	_init_ctx.Do(func(){
		_instance = &LRU{
			list:list.New(),
			index: make(map[string]*list.Element),
			Capacity:capacity,
		}
	})
	return _instance
}
