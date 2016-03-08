package container

import (
	"container/list"
	"sync"
)

type SortedLinkList struct {
	List *list.List
	sync.RWMutex
}

func (sl *SortedLinkList)Insert(v CompareElement) *Element {
	sl.Lock()
	defer sl.Unlock()
	var location *list.Element = nil
	for e := sl.List.Front(); e != nil; e = e.Next() {
		if false == v.Compare(e.Value.(*Element)) {
			location = e
			break
		}
	}
	if location == nil {
		return sl.List.PushBack(&Element{Value:v}).Value.(*Element)
	}
	return sl.List.InsertBefore(&Element{Value:v},location).Value.(*Element)
}
func (sl *SortedLinkList)Update(ele *Element) {
	sl.Lock()
	defer sl.Unlock()
	var self *list.Element = nil
	var location *list.Element = nil
	for e := sl.List.Front(); e != nil; e = e.Next() {
		if e.Value.(*Element) == ele {
			self = e
		}
		if location == nil && false == ele.Value.(CompareElement).Compare(e.Value.(*Element)) {
			location = e
		}
		if location != nil && self != nil {
			break
		}
	}
	if self != nil {
		if location == nil {
			sl.List.MoveToBack(self)
		} else {
			sl.List.MoveBefore(self,location)
		}
	}
}
func (sl *SortedLinkList)Front() *Element {
	sl.RLock()
	defer sl.RUnlock()
	return sl.List.Front().Value.(*Element)
}
func (sl *SortedLinkList)Back() *Element {
	sl.RLock()
	defer sl.RUnlock()
	return sl.List.Back().Value.(*Element)
}
func (sl *SortedLinkList)Length() int {
	return sl.List.Len()
}

func NewSortedLinkList() *SortedLinkList {
	return &SortedLinkList{
		List:list.New(),
	}
}