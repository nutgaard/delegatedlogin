package internal

import (
	Trie "github.com/dghubble/trie"
	"github.com/rs/zerolog/log"
)

type myPathTrie[T any] struct {
	*Trie.PathTrie
}

func NewPathTrie[T any](
	config []T,
	keySelector func(T) string,
) *myPathTrie[T] {
	trie := Trie.NewPathTrie()
	for _, t := range config {
		trie.Put(keySelector(t), t)
	}
	return &myPathTrie[T]{trie}
}

func (t *myPathTrie[T]) Search(path string) *T {
	var out *T
	var err error
	err = t.WalkPath(path, func(key string, value interface{}) error {
		casted, ok := value.(T)
		if ok {
			out = &casted
		}
		return nil
	})

	if err != nil {
		log.Error().Msgf("Error while searching for proxyconfig given path: %s", path)
		return nil
	}

	return out
}
