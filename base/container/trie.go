package container


type TrieNode struct {
	Child map[byte]*TrieNode
	IsWord bool
}
func (node *TrieNode)DumpChild() (children []byte) {
	for k,_ := range node.Child {
		children = append(children, k)
	}
	return
}

type Trie struct {
	Root *TrieNode
}

func (t *Trie)Insert(word []byte) {
	node := t.Root
	for _,letter := range word {
		_,exist := node.Child[letter]
		if !exist {
			node.Child[letter] = &TrieNode{
				Child:make(map[byte]*TrieNode),
				IsWord:false,
			}
		}
		node = node.Child[letter]
	}
	node.IsWord = true
}
func (t *Trie)Delete(word []byte) bool {
	node := t.Root
	deleteLetter := []byte{}
	deleteNode := []*TrieNode{}
	for _,letter := range word {
		deleteLetter = append(deleteLetter, letter)
		deleteNode = append(deleteNode, node)
		child,exist := node.Child[letter]
		if !exist {
			return false
		}
		node = child
	}
	if !node.IsWord {
		return false
	}
	if len(node.Child) != 0 {
		node.IsWord = false
	} else {
		for i := len(deleteNode) - 1; i >= 0; i-- {
			node,letter := deleteNode[i],deleteLetter[i]
			delete(node.Child, letter)
			if len(node.Child) == 0 || node.IsWord {
				break
			}
		}
	}
	return true
}
func (t *Trie)IsPrefix(word []byte) bool {
	node := t.Root
	for _,letter := range word {
		if node.IsWord {
			return true
		}
		child,exist := node.Child[letter]
		if !exist {
			return false
		}
		node = child
	}
	return node.IsWord
}
func NewTrie() *Trie {
	return &Trie{
		Root:&TrieNode{
			Child:make(map[byte]*TrieNode),
			IsWord:false,
		},
	}
}