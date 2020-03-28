package segmenter

import "io"

var defaultSegmenter Segmenter

type CutFn = func(r io.Reader) ([]string, error)

type Segmenter interface {
	Cut(r io.Reader) ([]string, error)
}

type s struct {
	fn CutFn
}

func (s s) Cut(r io.Reader) ([]string, error) {
	return s.fn(r)
}

func SimpleSegmenter(fn CutFn) Segmenter {
	return s{fn: fn}
}
func Use(s Segmenter) {
	defaultSegmenter = s
}
func Get() Segmenter {
	if defaultSegmenter == nil {
		panic("segmenter is not init")
	}
	return defaultSegmenter
}
