package html2text

import (
	"bytes"
	"strings"

	"golang.org/x/net/html"
)

func Convert(r []byte) ([]byte, error) {
	z := html.NewTokenizer(bytes.NewBuffer(r))
	var b []byte
	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			return b, nil
		}
		if tt == html.TextToken {
			b = append(b, []byte(" "+strings.TrimSpace(z.Token().String()))...)
		}
	}
}
