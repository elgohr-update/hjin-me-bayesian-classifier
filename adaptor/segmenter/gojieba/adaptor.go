package gojieba

import (
	"io"
	"io/ioutil"
	"path"

	"github.com/hjin-me/bayesian-classifier/segmenter"
	"github.com/yanyiwu/gojieba"
)

func New(args ...string) segmenter.Segmenter {
	var s *gojieba.Jieba
	if len(args) == 1 {
		dictDir := args[0]
		DictPath := path.Join(dictDir, "jieba.dict.utf8")
		HmmPath := path.Join(dictDir, "hmm_model.utf8")
		UserDictPath := path.Join(dictDir, "user.dict.utf8")
		IDFPath := path.Join(dictDir, "idf.utf8")
		StopWordsPath := path.Join(dictDir, "stop_words.utf8")

		dicts := [5]string{
			DictPath,
			HmmPath,
			UserDictPath,
			IDFPath,
			StopWordsPath,
		}
		s = gojieba.NewJieba(dicts[:]...)
	} else {
		s = gojieba.NewJieba()
	}

	fn := func(r io.Reader) ([]string, error) {
		b, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}
		return s.Cut(string(b), true), nil
	}
	return segmenter.SimpleSegmenter(fn)
}
