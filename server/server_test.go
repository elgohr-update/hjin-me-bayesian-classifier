package server

import (
	"context"
	"log"
	"net/http"
	"path"
	"runtime"
	"testing"

	"github.com/kataras/iris/v12/httptest"
	"github.com/spf13/viper"
)

var cwd string

func init() {
	_, filePath, _, _ := runtime.Caller(0)
	log.Println("current path is ", filePath)
	cwd = path.Join(filePath, "../..")
}

func TestNew(t *testing.T) {
	viper.Set(EnvModelPath, path.Join(cwd, "/data/jieba_storage.json"))
	viper.Set(EnvDictDir, path.Join(cwd, "/assets/dict/"))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	app := New(ctx)

	expect := httptest.New(t, app)
	resp := expect.POST("/v1/categorize").WithBytes([]byte("这是一篇Javascript的技巧")).Expect()
	t.Log(resp.Status(http.StatusOK).Body().Raw())
}
