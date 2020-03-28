package server

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/hjin-me/bayesian-classifier/adaptor/segmenter/gojieba"
	"github.com/hjin-me/bayesian-classifier/classifier"
	"github.com/kataras/iris/v12"
	"github.com/spf13/viper"
)

const EnvModelPath = "MODEL_PATH"
const EnvDictDir = "DICT_DIR"

func New(ctx context.Context) *iris.Application {
	modelPath := viper.GetString(EnvModelPath)
	app := iris.New()
	app.Configure(iris.WithoutInterruptHandler, iris.WithoutBodyConsumptionOnUnmarshal)
	// 存活检测
	app.Any("/liveness", func(i iris.Context) {
		i.StatusCode(http.StatusOK)
		_, _ = i.WriteString("live")
	})
	// 可用检测
	app.Any("/readiness", func(i iris.Context) {
		i.StatusCode(http.StatusOK)
		_, _ = i.WriteString("ready")
	})

	v1 := app.Party("/v1", func(i iris.Context) {
		// 所有请求替换 context
		i.ResetRequest(i.Request().WithContext(ctx))
		i.Next()
	})
	v1.Post("/train", func(i iris.Context) {
		// 样本训练接口
	})
	dictDir := viper.GetString(EnvDictDir)
	clsfr := classifier.New()
	err := clsfr.LoadSegmenter(gojieba.New(dictDir))
	if err != nil {
		panic(err)
	}
	b, err := ioutil.ReadFile(modelPath)
	if err != nil {
		panic(err)
	}
	err = clsfr.LoadModel(bytes.NewBuffer(b))
	if err != nil {
		panic(err)
	}
	v1.Post("/categorize", func(i iris.Context) {
		// 文本分类接口
		i.GetHeader("content-type")
		score, err := clsfr.Categorize(i.Request().Body)
		if err != nil {
			i.StatusCode(http.StatusInternalServerError)
			_, _ = i.WriteString(err.Error())
			return
		}
		_, _ = i.JSON(score)
	})
	return app
}

func Run(ctx context.Context) {
	app := New(ctx)
	// 确保能够读取用户的数据库来初始化基础数据
	//e := rbac.Get()
	//rbac.InitPolicy(e)
	// Load the policy from DB.
	//err := e.LoadPolicy()
	//if err != nil {
	//	os.Exit(1)
	//}
	go func() {
		<-ctx.Done()
		timeout := 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		_ = app.Shutdown(ctx)

	}()
	_ = app.Run(iris.Addr(":8080"))
}
