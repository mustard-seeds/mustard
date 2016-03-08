package container

import (
	"testing"
)

type lruTestStruct struct {
	name string
	age int
}

func TestLru(t *testing.T) {
	l := NewLRU(3)
	l.Set("x",&lruTestStruct{
		name:"xx",
		age:1,
	})
	l.Set("y",&lruTestStruct{
		name:"yy",
		age:2,
	})
	l.Set("z",&lruTestStruct{
		name:"zz",
		age:4,
	})
	if l.Size() != 3 {
		t.Errorf("set value error")
	}
	_,e := l.Get("a")
	if e != false {
		t.Errorf("get value not exist.")
	}
	v,e := l.Get("x")
	if e != true {
		t.Errorf("get exist value error")
	}
	if v.Value.(*lruTestStruct).name != "xx" {
		t.Errorf("get value error")
	}
	if l.list.Front().Value.(*innerElement).value.Value.(*lruTestStruct).name != "xx" {
		t.Errorf("lru cache not sort.")
	}
	l.Set("a",&lruTestStruct{
		name:"zddz",
		age:4,
	})
	l.Set("afwe",&lruTestStruct{
		name:"zddz",
		age:4,
	})
	if len(l.index) != 3 {
		t.Errorf("lru fix size error.")
	}
	l.JustUpdateValue("a",&lruTestStruct{
		name:"newA",
		age:4,
	})
	newA,_ := l.Get("a")
	if newA.Value.(*lruTestStruct).name != "newA" {
		t.Errorf("JustUpdateValue Error.")
	}
}
