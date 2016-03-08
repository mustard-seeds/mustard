package container

import (
	"testing"
)

func TestTrie(t *testing.T) {
	trie := NewTrie()
	trie.Insert([]byte("01234"))
	trie.Insert([]byte("01"))
	trie.Insert([]byte("012553"))
	trie.Delete([]byte("01"))
	if !trie.IsPrefix([]byte("01234")) {
		t.Errorf("is prefix error")
	}
	if !trie.IsPrefix([]byte("012345")) {
		t.Errorf("is prefix error")
	}
	if trie.IsPrefix([]byte("0123")) {
		t.Errorf("is prefix error")
	}
}
