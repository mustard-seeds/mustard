package container

type Element struct {
	Value interface{}
}

type CompareElement interface {
	Compare(e *Element) bool
}
