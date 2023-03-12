package set

import (
	"testing"

	"github.com/shoenig/test/must"
)

func TestBox(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		b := Box(42)
		must.Eq(t, 42, b.item)
	})

	t.Run("string", func(t *testing.T) {
		b := Box("hello")
		must.Eq(t, "hello", b.item)
	})
}
