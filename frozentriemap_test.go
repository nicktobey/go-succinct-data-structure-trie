package bits

import "testing"

func tlookupMap(t *testing.T, ftm *FrozenTrieMap, word string, expected bool) {
	index, found := ftm.LookupIndex(word)
	if found != expected {
		t.Error(word)
	}
	t.Log(word, index, found)
	if found {
		t.Log(ftm.keys.Rank(1, index))
	}
}

func TestMapLookup(t *testing.T) {
	te := Trie{}
	te.Init()
	insertNotInAlphabeticalOrder(&te)

	ftm := FrozenTrieMap{}
	ftm.Create(te)

	tlookupMap(t, &ftm, "apple", true)
	tlookupMap(t, &ftm, "appl", false)
	tlookupMap(t, &ftm, "applea", false)
	tlookupMap(t, &ftm, "orange", true)
	tlookupMap(t, &ftm, "lamp", true)
	tlookupMap(t, &ftm, "hello", true)
	tlookupMap(t, &ftm, "jello", true)
	tlookupMap(t, &ftm, "quiz", true)
	tlookupMap(t, &ftm, "quize", false)
	tlookupMap(t, &ftm, "alphaph", false)
	tlookupMap(t, &ftm, "alphapha", true)
}

func TestMapReverseLookup(t *testing.T) {
	te := Trie{}
	te.Init()

	words := []string{}

	for _, word := range words {
		te.Insert(word)
	}
	insertNotInAlphabeticalOrder(&te)

	ftm := FrozenTrieMap{}
	ftm.Create(te)

	// for i := range words {
	for i := 0; i < 7; i++ {
		key := ftm.ReverseLookup(uint(i + 1))
		t.Log(key)
	}
}
