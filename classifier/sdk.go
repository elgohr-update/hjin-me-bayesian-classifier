package classifier

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"sync"

	"github.com/hjin-me/go-utils/logex"
	"github.com/yanyiwu/gojieba"
)

var ErrDuplicateDocs = errors.New("已经学习过的样本")
var ErrCategoryEmpty = errors.New("分类为空")

type data struct {
	Category map[string]float64            `json:"category"` // 分类数据
	Words    map[string]map[string]float64 `json:"words"`    // 单词数据
	Docs     map[string]bool               `json:"docs"`     // 文档数据
}
type SDK struct {
	mutex     sync.Mutex // 训练字典读写锁
	data      data
	debug     bool
	segmenter *gojieba.Jieba
}

func (s *SDK) EnableDebug(b bool) {
	s.debug = b
}

func (s *SDK) LoadDictionary() error {
	s.segmenter = gojieba.NewJieba()
	return nil
}

func (s *SDK) LoadModel(r io.Reader) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	err = json.Unmarshal(b, &s.data)
	if err != nil {
		return err
	}
	return nil
}
func (s *SDK) SaveModel(w io.Writer) error {
	s.mutex.Lock()
	b, err := json.Marshal(s.data)
	s.mutex.Unlock()
	if err != nil {
		return err
	}

	_, err = io.Copy(w, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	return nil
}

// 训练样本
func (s *SDK) Train(r io.Reader, category string) error {
	category = strings.TrimSpace(category)
	if category == "" {
		return ErrCategoryEmpty
	}

	var buf bytes.Buffer
	mr := io.TeeReader(r, &buf)
	hash := sha256.New()
	_, err := io.Copy(hash, mr)
	if err != nil {
		return err
	}
	// 样本去重
	docHash := hex.EncodeToString(hash.Sum(nil))
	if _, ok := s.data.Docs[docHash]; ok {
		return ErrDuplicateDocs
	}
	defer func() {
		s.mutex.Lock()
		s.data.Docs[docHash] = true
		s.mutex.Unlock()
	}()

	b, err := ioutil.ReadAll(mr)
	if err != nil {
		return err
	}
	doc := string(b)

	// 更新单词数据
	// 同一个文档中单词出现多次，仅记录一次
	fwords := make(map[string]bool)
	x := gojieba.NewJieba()
	defer x.Free()
	words := filterWord(x.Cut(doc, true))
	//words := t.segmenter.Segment(doc)
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, word := range words {
		if _, ok := fwords[word]; ok {
			continue
		}
		fwords[word] = true
		if _, ok := s.data.Words[word]; !ok {
			s.data.Words[word] = make(map[string]float64)
		}
		s.data.Words[word][category]++
		//log.Println("单词训练：", word)
	}
	// 更新分类统计
	s.data.Category[category]++

	return nil
}

// 计算样本总数
func (s *SDK) calcSampleCounts() float64 {
	total := 0.0
	for _, n := range s.data.Category {
		total += n
	}
	return total
}
func (s *SDK) calcSampleCountsByCategory(category string) float64 {
	return s.data.Category[category]
}
func (s *SDK) calcWordCountsByCategory(word, category string) float64 {
	if set, ok := s.data.Words[word]; ok {
		return set[category]
	}
	return 0
}
func (s *SDK) calcWordCounts(word string) float64 {
	var total float64 = 0
	for _, s := range s.data.Words[word] {
		total += s
	}
	return total
}
func (s *SDK) factor(word, category string) float64 {
	totalCategoryCounts := s.calcSampleCounts()
	targetCategoryCounts := s.calcSampleCountsByCategory(category)
	wordCountsTotal := s.calcWordCounts(word)
	wordCountsInCategory := s.calcWordCountsByCategory(word, category)
	// 频数较小的作为废弃规则
	//if wordCountsTotal < 5 {
	//	//prob *= (1 / targetCategoryCounts) / (1 / totalCategoryCounts)
	//	continue
	//}

	//log.Printf("%s = %0.6f / %0.6f", word, wordCountsInCategory/targetCategoryCounts, wordCountsTotal/totalCategoryCounts)
	// 拉普拉斯平滑
	wordCountsInCategory += 1
	targetCategoryCounts += 2
	wordCountsTotal += 1
	totalCategoryCounts += 2

	invalidWord := false
	if wordCountsInCategory == 1 && wordCountsTotal == 1 {
		invalidWord = true
	}
	//log.Printf("[%s], %s = %0.6f / %0.6f, laplace, = %0.6f", category, word, wordCountsInCategory/targetCategoryCounts, wordCountsTotal/totalCategoryCounts, (wordCountsInCategory/targetCategoryCounts)/(wordCountsTotal/totalCategoryCounts))
	if s.debug {
		if invalidWord {
			logex.Debugf("factor %s[%s] = %0.6f = ( %0.6f / %0.6f ) / ( %0.6f / %0.6f )",
				"*"+word, category, 1.0,
				wordCountsInCategory, targetCategoryCounts, wordCountsTotal, totalCategoryCounts)
		} else {
			logex.Debugf("factor %s[%s] = %0.6f = ( %0.6f / %0.6f ) / ( %0.6f / %0.6f )",
				word, category, wordCountsInCategory/targetCategoryCounts/(wordCountsTotal/totalCategoryCounts),
				wordCountsInCategory, targetCategoryCounts, wordCountsTotal, totalCategoryCounts)
		}
	}
	if invalidWord {
		return 1.0
	}

	return wordCountsInCategory / targetCategoryCounts / (wordCountsTotal / totalCategoryCounts)
}
func (s *SDK) Categorize(doc string) []*ScoreItem {
	words := s.segmenter.Cut(doc, true)
	words = filterWord(words)
	scores := NewScores()
	for category, categoryCounts := range s.data.Category {
		prob := 1.0
		for _, word := range words {
			prob *= s.factor(word, category)
		}
		scores.Append(
			category,
			prob*categoryCounts/s.calcSampleCounts(),
		)
		if s.debug {
			logex.Debugf("P(%s) = %0.6f = %0.0f / %0.0f", category, categoryCounts/s.calcSampleCounts(), categoryCounts, s.calcSampleCounts())
		}
	}
	return scores.Top(10)
}

func New() *SDK {
	s := SDK{}
	return &s
}

func PrettyScore(s []*ScoreItem) string {
	result := ""
	for _, item := range s {
		result += fmt.Sprintf("\n%s\t%0.10f", item.Category, item.Score)
	}
	return result
}
