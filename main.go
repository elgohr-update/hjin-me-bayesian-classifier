// 贝叶斯分类器（Naive Bayesian classifier）支持中文文档解析训练和分类，提供HTTP API访问。
package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hjin-me/bayesian-classifier/server"
	"github.com/spf13/viper"
)

func main() {
	//var dictionaryPath string
	//var sampleDir string
	//var cacheDir string
	//flag.StringVar(&dictionaryPath, "d", "", "字典路径")
	//flag.StringVar(&sampleDir, "s", "", "训练样本目录")
	//flag.StringVar(&cacheDir, "c", "", "临时缓存目录")
	flag.Parse()
	viper.AutomaticEnv()

	// 注册一个 cancel，当进程被退出时触发
	closeContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch,
			// kill -SIGINT XXXX 或 Ctrl+c
			os.Interrupt,
			syscall.SIGINT, // register that too, it should be ok
			// kill -SIGTERM XXXX
			syscall.SIGTERM,
		)
		<-ch
		cancel()
	}()
	go server.Run(closeContext)

	<-closeContext.Done()
	<-time.After(5 * time.Second)

}
