package sego

import (
	"io"
	"io/ioutil"

	"github.com/hjin-me/bayesian-classifier/segmenter"
)

func New(r io.Reader) segmenter.Segmenter {
	s := Segmenter{}
	if err := s.LoadDictionary(r); err != nil {
		panic(err)
	}
	fn := func(r io.Reader) ([]string, error) {
		b, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}
		ss := s.Segment(b)
		return SegmentsToSlice(ss, false), nil
	}
	return segmenter.SimpleSegmenter(fn)
}
