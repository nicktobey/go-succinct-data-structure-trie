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
	ft    FrozenTrie
	keys  RankDirectory
	words uint
}

func (f *FrozenTrieMap) Create(trie Trie) {
	finalNodes := BitWriter{}

	teData := trie.Encode()
	rd := CreateRankDirectory(teData, trie.GetNodeCount()*2+1, L1, L2)

	f.ft.Init(teData, rd.GetData(), trie.GetNodeCount())

	f.ft.Apply(func(node FrozenTrieNode) {
		if node.final {
			finalNodes.Write(1, 1)
			f.words++
		} else {
			finalNodes.Write(0, 1)
		}
	})

	f.keys = CreateRankDirectory(finalNodes.GetData(), trie.GetNodeCount(), L1, L2)
}

func (f *FrozenTrieMap) Init(ft FrozenTrie, keys RankDirectory) {
	f.ft = ft
	f.keys = keys
}

func (f *FrozenTrieMap) LookupIndex(word string) (index uint, found bool) {
	node := f.ft.GetRoot()
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

	return f.keys.Rank(1, node.index), node.final
}

func (f *FrozenTrieMap) ReverseLookup(keyIndex uint) (word string) {
	var resultBytes []byte
	trieNodeNumber := f.keys.Select(1, keyIndex)
	for trieNodeNumber > 0 {
		node := f.ft.GetNodeByIndex(trieNodeNumber)
		resultBytes = append([]byte{node.letter}, resultBytes...)
		parentOffset := f.ft.directory.Select(1, trieNodeNumber+1)
		trieNodeNumber = f.ft.directory.Rank(0, parentOffset) - 1
	}
	return string(resultBytes)
}

func (f *FrozenTrieMap) GetBuffer() []byte {
	var result bytes.Buffer
	result.WriteString(f.ft.data.GetData())
	result.WriteString(f.ft.directory.GetData())
	return result.Bytes()
}

func (f *FrozenTrieMap) GetOffsets() []byte {
	var result bytes.Buffer
	nodeCount := (f.ft.letterStart - 1) / 2
	// TODO: This assumes fewer than 256 nodes.
	result.WriteByte(byte(nodeCount))
	result.WriteString(f.keys.data.GetData())
	result.WriteString(f.keys.directory.GetData())
	return result.Bytes()
}
