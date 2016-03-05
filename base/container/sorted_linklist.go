package container

import (
	"container/list"
)

type SortedLinkList struct {
	List *list.List
}

func (sl *SortedLinkList)Insert(v CompareElement) *Element {
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
func (sl *SortedLinkList)Front() *Element {
	return sl.List.Front().Value.(*Element)
}
func (sl *SortedLinkList)Back() *Element {
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