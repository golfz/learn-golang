package randtest

import (
	"fmt"
	word "github.com/golfz/learn-golang/test/word2"
	"math/rand"
	"testing"
	"time"
)

func TestRandomPalindromes(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	t.Logf("Random seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))

	for i := 0; i < 1000; i++ {
		p := RandomPalindrome(rng)
		fmt.Println(p)
		if !word.IsPalindrome(p) {
			t.Errorf("IsPanlindrome(%q) = false", p)
		}
	}
}
