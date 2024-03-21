package bits

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSearch(t *testing.T) {
	te := Trie{}
	te.Init()
	insertNotInAlphabeticalOrder(&te)
	teData, _ := te.Encode()
	rd := CreateRankDirectory(teData, te.GetNodeCount()*2+1, L1, L2)

	ft := FrozenTrie{}
	ft.Init(teData, rd.GetData(), te.GetNodeCount())

	assert.Equal(t, []string{"apple", "alphapha"}, ft.GetSuggestedWords("a", 10))
	assert.Equal(t, 0, ft.GetSuggestedWords("b", 10))
	assert.Equal(t, []string{"hello"}, ft.GetSuggestedWords("h", 10))
}
