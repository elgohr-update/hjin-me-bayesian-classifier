package gojieba

import (
	"io"
	"io/ioutil"

	"github.com/hjin-me/bayesian-classifier/segmenter"
	"github.com/yanyiwu/gojieba"
)

func New() segmenter.Segmenter {
	s := gojieba.NewJieba()
	fn := func(r io.Reader) ([]string, error) {
		b, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}
		return s.Cut(string(b), true), nil
	}
	return segmenter.SimpleSegmenter(fn)
}
