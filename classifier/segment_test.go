package classifier

import (
	"testing"

	"github.com/yanyiwu/gojieba"
)

func TestJieba(t *testing.T) {
	x := gojieba.NewJieba()
	words := x.Cut("这是一篇Javascript的技巧", true)
	t.Log(words)
	t.Log(filterWord(words))
}
