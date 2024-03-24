package bits

import "bytes"

/*
*

	This class is used for traversing the succinctly encoded trie.
*/
type FrozenTrieNode struct {
	trie       *FrozenTrie
	index      uint
	letter     byte
	final      bool
	firstChild uint
	childCount uint
}

/*
*

	Returns the number of children.
*/
func (f *FrozenTrieNode) GetChildCount() uint {
	return f.childCount
}

/*
*

	Returns the FrozenTrieNode for the given child.

	@param index The 0-based index of the child of this node. For example, if
	the node has 5 children, and you wanted the 0th one, pass in 0.
*/
func (f *FrozenTrieNode) GetChild(index uint) FrozenTrieNode {
	return f.trie.GetNodeByIndex(f.firstChild + index)
}

/*
*

	The FrozenTrie is used for looking up words in the encoded trie.

	@param data A string representing the encoded trie.

	@param directoryData A string representing the RankDirectory. The global L1
	and L2 constants are used to determine the L1Size and L2size.

	@param nodeCount The number of nodes in the trie.
*/
type FrozenTrie struct {
	data        BitString
	directory   RankDirectory
	letterStart uint
}

func (f *FrozenTrie) Init(data, directoryData string, nodeCount uint) {
	f.data.Init(data)
	f.directory.Init(directoryData, data, nodeCount*2+1, L1, L2)

	// The position of the first bit of the data in 0th node. In non-root
	// nodes, this would contain 6-bit letters.
	f.letterStart = nodeCount*2 + 1
}

/*
*

	Retrieve the FrozenTrieNode of the trie, given its index in level-order.
	This is a private function that you don't have to use.
*/
func (f *FrozenTrie) GetNodeByIndex(index uint) FrozenTrieNode {
	// retrieve the (dataBits)-bit letter.
	final := (f.data.Get(f.letterStart+index*dataBits, 1) == 1)
	letter := uint8(f.data.Get(f.letterStart+index*dataBits+1, (dataBits - 1)))
	firstChild := f.directory.Select(0, index+1) - index

	// Since the nodes are in level order, this nodes children must go up
	// until the next node's children start.
	childOfNextNode := f.directory.Select(0, index+2) - index - 1

	return FrozenTrieNode{
		trie:       f,
		index:      index,
		letter:     letter,
		final:      final,
		firstChild: firstChild,
		childCount: (childOfNextNode - firstChild),
	}
}

/*
*

	Retrieve the root node. You can use this node to obtain all of the other
	nodes in the trie.
*/
func (f *FrozenTrie) GetRoot() FrozenTrieNode {
	return f.GetNodeByIndex(0)
}

/*
*

	Look-up a word in the trie. Returns true if and only if the word exists
	in the trie.
*/
func (f *FrozenTrie) Lookup(word string) bool {
	node := f.GetRoot()
	wordBytes := []byte(word)
	for _, i := range wordBytes {
		var child FrozenTrieNode
		var j uint = 0
		for ; j < node.GetChildCount(); j++ {
			child = node.GetChild(j)
			if child.letter == i {
				break
			}
		}

		if j == node.GetChildCount() {
			return false
		}
		node = child
	}

	return node.final
}

func (f *FrozenTrie) LookupIndex(word string) (index uint, found bool) {
	node := f.GetRoot()
	wordBytes := []byte(word)
	for _, i := range wordBytes {
		var child FrozenTrieNode
		var j uint = 0
		for ; j < node.GetChildCount(); j++ {
			child = node.GetChild(j)
			if child.letter == i {
				break
			}
		}

		if j == node.GetChildCount() {
			return 0, false
		}
		node = child
	}

	return node.index, node.final
}

/*
* Apply a function to each node, traversing the trie in level order.
 */
func (t *FrozenTrie) Apply(fn func(FrozenTrieNode)) {
	var level []FrozenTrieNode
	level = append(level, t.GetRoot())
	for len(level) > 0 {
		node := level[0]
		level = level[1:]
		for i := uint(0); i < node.GetChildCount(); i++ {
			level = append(level, node.GetChild(i))
		}
		fn(node)
	}
}

/*
* Apply a function to each node, traversing the trie in pre-order.
 */
func (f *FrozenTrie) ApplyPreOrder(fn func(*FrozenTrieNode)) {
	root := f.GetRoot()
	root.ApplyPreOrder(fn)
}

/*
* Apply a function to each node, traversing the trie in pre-order.
 */
func (f *FrozenTrieNode) ApplyPreOrder(fn func(*FrozenTrieNode)) {
	fn(f)
	for i := uint(0); i < f.childCount; i++ {
		child := f.GetChild(i)
		child.ApplyPreOrder(fn)
	}
}

func (t *FrozenTrie) GetLastLexographicKey() string {
	var result bytes.Buffer
	node := t.GetRoot()

	for {
		childCount := node.GetChildCount()
		if childCount == 0 {
			return result.String()
		}
		node = node.GetChild(childCount - 1)
		result.Write([]byte{node.letter})
	}
}
