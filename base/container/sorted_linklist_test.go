package container

import (
	"testing"
)

type IElement struct {
	Name string
	Age int
}
func (i *IElement)Compare(e *Element) bool {
	return i.Age < e.Value.(*IElement).Age
}

func TestSortList(t *testing.T) {
	l := NewSortedLinkList()
	l.Insert(&IElement{"xx",3})
	l.Insert(&IElement{"xx",2})
	l.Insert(&IElement{"xx",1})
	l.Insert(&IElement{"xx",20})
	l.Insert(&IElement{"xx",21})
	l.Insert(&IElement{"xx",1})
	e := l.Insert(&IElement{"xx",4})
	if l.Front().Value.(*IElement).Age != 21 {
	}
	e.Value.(*IElement).Age = 233
	l.Update(e)
	if l.Front().Value.(*IElement).Age != 233 {
		t.Errorf("sort update error.")
	}
}
