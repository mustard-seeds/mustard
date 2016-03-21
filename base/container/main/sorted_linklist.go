package main

import (
	"fmt"
	. "mustard/base/container"
)

type IElement struct {
	Name string
	Age  int
}

func (i *IElement) Compare(e *Element) bool {
	return i.Age < e.Value.(*IElement).Age
}

func printlist(l *SortedLinkList) {
	for e := l.List.Front(); e != nil; e = e.Next() {
		fmt.Printf("=%v=", e.Value.(*Element).Value.(*IElement))
	}
	fmt.Println()
}
func main() {
	l := NewSortedLinkList()
	l.Insert(&IElement{"xx", 3})
	l.Insert(&IElement{"xx", 2})
	l.Insert(&IElement{"xx", 1})
	l.Insert(&IElement{"xx", 20})
	l.Insert(&IElement{"xx", 21})
	l.Insert(&IElement{"xx", 1})
	e := l.Insert(&IElement{"xx", 4})
	printlist(l)
	e.Value.(*IElement).Age = 233
	l.Update(e)
	printlist(l)

}
