package bayesianc

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/hjin-me/bayesian-classifier/html2text"
	"github.com/hjin-me/bayesian-classifier/util"
	"github.com/stretchr/testify/assert"
)

var cwd = "当前目录"

func initModel(t *testing.T, cwd string) *SDK {
	// 分类器
	handler := New()
	b, err := ioutil.ReadFile(path.Join(cwd, "/assets/dictionary.txt"))
	assert.Nil(t, err)
	err = handler.LoadDictionary(bytes.NewBuffer(b))
	assert.Nil(t, err)
	b, err = ioutil.ReadFile(path.Join(cwd, "/_temp/storage.json"))
	assert.Nil(t, err)
	err = handler.LoadModel(bytes.NewBuffer(b))
	assert.Nil(t, err)
	return handler
}
func TestSDK_Categorize(t *testing.T) {
	// 分类器
	handler := initModel(t, cwd)
	handler.EnableDebug(true)
	b, err := ioutil.ReadFile(path.Join(cwd, "sample/normal/eula_cn.txt"))
	assert.Nil(t, err)
	t.Log(handler.Categorize(b))
	//t.Log(handler.Categorize([]byte("这是一篇Javascript的技巧")))
}

func TestCategorizeNegative(t *testing.T) {
	// 分类器
	handler := initModel(t, cwd)

	// 分类测试
	sampleDir := path.Join(cwd, "/sample/normal")
	fs, err := util.ReadDir(sampleDir)
	assert.Nil(t, err)
	for _, f := range fs {
		if filepath.Ext(f.Name()) != ".txt" {
			continue
		}
		t.Run("normal text "+f.Name(), func(t *testing.T) {
			doc, err := ioutil.ReadFile(path.Join(sampleDir + "/" + f.Name()))
			assert.Nil(t, err)
			score := handler.Categorize(doc)
			//t.Logf("%3.6f", score[0].Score)
			msg := fmt.Sprintf("%s=%0.10f, %s=%0.10f", score[0].Category, score[0].Score, score[1].Category, score[1].Score)
			assert.Len(t, score, 2, msg)
			assert.Equal(t, score[0].Category, "normal", fmt.Sprintf("wrong %s, %s", f.Name(), msg))
			assert.Greater(t, score[0].Score, 0.9, fmt.Sprintf("score failed %s, %s", f.Name(), msg))
			assert.Less(t, score[1].Score, 0.3, fmt.Sprintf("score failed %s, %s", f.Name(), msg))
		})
	}
}

func TestCategorizePositive(t *testing.T) {
	// 分类器
	handler := initModel(t, cwd)

	sampleDir := path.Join(cwd, "/sample/privacy")
	fs, err := util.ReadDir(sampleDir)
	assert.Nil(t, err)
	for _, f := range fs {
		if filepath.Ext(f.Name()) != ".txt" {
			continue
		}
		t.Run("privacy "+f.Name(), func(t *testing.T) {
			doc, err := ioutil.ReadFile(path.Join(sampleDir + "/" + f.Name()))
			assert.Nil(t, err)
			score := handler.Categorize(doc)
			//t.Logf("%3.6f", score[0].Score)
			msg := fmt.Sprintf("%s=%0.10f, %s=%0.10f", score[0].Category, score[0].Score, score[1].Category, score[1].Score)
			assert.Len(t, score, 2, msg)
			assert.Equal(t, score[0].Category, "privacy", fmt.Sprintf("wrong %s, %s", f.Name(), msg))
			//assert.Less(t, score[0].Score, 1.5, fmt.Sprintf("score failed %s, %s", f.Name(), msg))
			assert.Greater(t, score[0].Score, 0.9, fmt.Sprintf("score failed %s, %s", f.Name(), msg))
			assert.Less(t, score[1].Score, 0.3, fmt.Sprintf("score failed %s, %s", f.Name(), msg))
		})
	}
}

func TestConvertHTML(t *testing.T) {
	t.SkipNow()
	sampleDir := path.Join(cwd, "/sample/normal")
	fs, err := util.ReadDir(sampleDir)
	assert.Nil(t, err)
	for _, f := range fs {
		doc, err := util.ReadFile(path.Join(sampleDir + "/" + f.Name()))
		assert.Nil(t, err)
		textDoc, err := html2text.Convert(doc)
		assert.Nil(t, err)
		err = ioutil.WriteFile(path.Join(sampleDir+"/"+f.Name()+".txt"), textDoc, os.ModePerm)
		assert.Nil(t, err)
	}
}
