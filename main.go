// 贝叶斯分类器（Naive Bayesian classifier）支持中文文档解析训练和分类，提供HTTP API访问。
package main

import (
	"flag"
)

func main() {
	var dictionaryPath string
	var sampleDir string
	var cacheDir string
	flag.StringVar(&dictionaryPath, "d", "", "字典路径")
	flag.StringVar(&sampleDir, "s", "", "训练样本目录")
	flag.StringVar(&cacheDir, "c", "", "临时缓存目录")
	flag.Parse()
}
