package bits

import (
	"bytes"
)

/*
FrozenTrieMap maps words in a trie onto indices.

	@param data A string representing the encoded trie.

	@param directoryData A string representing the RankDirectory. The global L1
	and L2 constants are used to determine the L1Size and L2size.

	@param nodeCount The number of nodes in the trie.
*/
type FrozenTrieMap struct {
	Ft                        FrozenTrie
	mapLevelOrderToInfixOrder map[uint]uint
	mapInfixOrderToLevelOrder []uint
	Words                     uint
}

func (f *FrozenTrieMap) Create(teData string, nodeCount uint) {

	f.mapLevelOrderToInfixOrder = make(map[uint]uint)
	// f.mapInfixOrderToLevelOrder = make([]uint)

	rd := CreateRankDirectory(teData, nodeCount*2+1, L1, L2)
	f.Ft.Init(teData, rd.GetData(), nodeCount)

	f.Ft.ApplyPreOrder(func(node *FrozenTrieNode) {
		if node.final {
			f.mapInfixOrderToLevelOrder = append(f.mapInfixOrderToLevelOrder, node.index)
			f.mapLevelOrderToInfixOrder[node.index] = f.Words
			f.Words++
		}
	})
}

func (f *FrozenTrieMap) LookupIndex(word string) (index uint, found bool) {
	node := f.Ft.GetRoot()
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

	return f.mapLevelOrderToInfixOrder[node.index], node.final
}

func (f *FrozenTrieMap) ReverseLookup(keyIndex uint) (word string) {
	var resultBytes []byte
	trieNodeNumber := f.mapInfixOrderToLevelOrder[keyIndex]
	for trieNodeNumber > 0 {
		node := f.Ft.GetNodeByIndex(trieNodeNumber)
		resultBytes = append([]byte{node.letter}, resultBytes...)
		parentOffset := f.Ft.directory.Select(1, trieNodeNumber+1)
		trieNodeNumber = f.Ft.directory.Rank(0, parentOffset) - 1
	}
	return string(resultBytes)
}

func (f *FrozenTrieMap) GetBuffer() []byte {
	var result bytes.Buffer
	result.WriteString(f.Ft.data.GetData())
	result.WriteString(f.Ft.directory.GetData())
	return result.Bytes()
}
