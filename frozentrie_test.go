package bits

import "testing"

func lookupTestCase(t *testing.T, ft *FrozenTrie, keys RankDirectory, word string, expected bool) {
	_, found := ft.LookupIndex(word)
	if found != expected {
		t.Error(word)
	}
}

func TestLookup(t *testing.T) {
	te := Trie{}
	te.Init()
	insertNotInAlphabeticalOrder(&te)
	finalNodes := BitWriter{}

	teData, _ := te.Encode()
	rd := CreateRankDirectory(teData, te.GetNodeCount()*2+1, L1, L2)

	ft := FrozenTrie{}
	ft.Init(teData, rd.GetData(), te.GetNodeCount())

	t.Log("Frozen Trie")
	var words uint = 0
	var nodes uint = 0
	ft.Apply(func(node FrozenTrieNode) {
		t.Log(node.index, node.letter, node.final)
		if node.final {
			finalNodes.Write(1, 1)
			nodes++
			words++
		} else {
			finalNodes.Write(0, 1)
			nodes++
		}
	})

	t.Log(teData)
	keys := CreateRankDirectory(finalNodes.GetData(), nodes, L1, L2)
	t.Log(rd.data.length)        // 3 bits + 1 character per node
	t.Log(rd.directory.length)   // assume at worst this doubles
	t.Log(keys.data.length)      // 1 bit per node
	t.Log(keys.directory.length) // assume at worst this doubles
	// Max space: 8 bits + 2 characters per node

	t.Log("Words")
	for i := uint(0); i < words+1; i++ {
		t.Log(i, keys.Select(1, i))
	}
	t.Log("Ranks")
	for i := uint(0); i < keys.numBits; i++ {
		t.Log(i, keys.Rank(0, i), keys.Rank(1, i))
	}

	lookupTestCase(t, &ft, keys, "apple", true)
	lookupTestCase(t, &ft, keys, "appl", false)
	lookupTestCase(t, &ft, keys, "applea", false)
	lookupTestCase(t, &ft, keys, "orange", true)
	lookupTestCase(t, &ft, keys, "lamp", true)
	lookupTestCase(t, &ft, keys, "hello", true)
	lookupTestCase(t, &ft, keys, "jello", true)
	lookupTestCase(t, &ft, keys, "quiz", true)
	lookupTestCase(t, &ft, keys, "quize", false)
	lookupTestCase(t, &ft, keys, "alphaph", false)
	lookupTestCase(t, &ft, keys, "alphapha", true)
}
