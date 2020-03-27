// 贝叶斯分类器（Naive Bayesian classifier）支持中文文档解析训练和分类，提供HTTP API访问。
package main

import (
	"flag"
	"fmt"
	"path"

	"github.com/hjin-me/bayesian-classifier/classifier"
	"github.com/hjin-me/bayesian-classifier/testcase"
)

const (
	DefaultProb      = 0.5            // 默认概率
	DefaultWeight    = 1.0            // 默认概率的权重，假定与一个单词相当
	Debug            = true           // 开启调试
	HTTP             = false          // 开启HTTP服务
	HTTPPort         = ":8812"        // HTTP端口
	Storage          = "file"         // 存储引擎，接受 file,redis，目前只支持file
	StoragePath      = "storage.json" // 文件存储引擎的存储路径
	StorageFrequency = "10"           // 自动存储的频率, 单位: 秒，0 表示不自动存储
)

func main() {
	var dictionaryPath string
	var sampleDir string
	var cacheDir string
	flag.StringVar(&dictionaryPath, "d", "", "字典路径")
	flag.StringVar(&sampleDir, "s", "", "训练样本目录")
	flag.StringVar(&cacheDir, "c", "", "临时缓存目录")
	flag.Parse()
	// 分类器
	handler := classifier.NewClassifier(map[string]interface{}{
		"defaultProb":   DefaultProb,
		"defaultWeight": DefaultWeight,
		"debug":         Debug,
		"http":          HTTP,
		"httpPort":      HTTPPort,
		"storage": map[string]string{
			"adapter":   Storage,
			"path":      path.Join(cacheDir, StoragePath),
			"frequency": StorageFrequency,
		},
	}, dictionaryPath)

	// 训练
	//handler.Training("这是一篇WEB开发的内容", "WEB")
	//handler.Training("这是一篇Javascript的技巧", "WEB")
	//handler.Training("这是一篇养生的内容", "WEB")
	//handler.Training("这是一篇养生的内容2", "健康")
	//handler.Training("这是一篇冬天养生食谱", "健康")
	//handler.Training(testcase.BaiduText, "normal")

	// 从txt文件进行训练
	//classifier.FileTrain("privacy", path.Join(sampleDir, "privcay"), handler)
	//classifier.FileTrain("normal", path.Join(sampleDir, "normal"), handler)

	// 获取训练数据
	testWord(handler, "隐私", "")   // 测试已知分类
	testWord(handler, "个人信息", "") // 测试已知分类
	//testWord(handler, "养生", "XX")  // 测试未知分类
	//testWord(handler, "养生", "")    // 查看所有分类
	//testWord(handler, "不认识", "")   // 测试未知单词
	testWord(handler, "服务器", "") // 测试未知单词

	// 分类测试
	//testDoc(handler, "养生", "养生是什么分类")
	//testDoc(handler, "一段测试文字", "即表示您信赖我们对您的信息的处理方式。我们深知这项责任事关重大，因此一直致力于保护您的隐私信息，并让您拥有控制权。")
	//testDoc(handler, "验证概率", "您可以通过更改您的浏览器设置限制百度公司对Cookie的使用。以百度浏览器为例，您可以在百度浏览器右上方的下拉菜单的“浏览器设置”中，通过“隐私设置——清除浏览数据”，选择清除您的Cookie。")
	//testDoc(handler, "百度文字", testcase.BaiduText)
	//testDoc(handler, "百度隐私", testcase.BaiduPrivacy)
	//testDoc(handler, "bilibili隐私", testcase.BiliBiliPrivacy)
	testDoc(handler, "bilibili普通文本", testcase.BiliBiliText)
	//testDoc(handler, "API Go")
	//testDoc(handler, "服务器")

	// 暂停
	//time.Sleep(time.Second * 15)

	// 开启了HTTP服务，不能结束进程
	//if HTTP {
	//	for {
	//		time.Sleep(time.Second)
	//	}
	//}
}

// 辅助测试：测试单词的频率
func testWord(classifier *classifier.Classifier, word, category string) {
	score := classifier.Score(word, category)
	if category != "" {
		fmt.Printf("单词【%s】在分类【%s】中出现的概率为: \n", word, category)
	} else {
		fmt.Printf("单词【%s】在分类中出现的概率为: \n", word)
	}
	printScore(score)
}

// 辅助测试：测试文档的分类
func testDoc(classifier *classifier.Classifier, name string, doc string) {
	score := classifier.Categorize(doc)
	//fmt.Println("测试文档归类于以下分类的概率为: ")
	//fmt.Println("--------------------------")
	//fmt.Println(name)
	//fmt.Println("--------------------------")
	//printScore(score)
	_ = score
}

// 辅助测试：输出
func printScore(scores []*classifier.ScoreItem) {
	if len(scores) == 0 {
		fmt.Println("未知单词 Orz！")
	}
	for k := range scores {
		fmt.Printf("%s\t%0.6f\n", scores[k].Category, scores[k].Score)
	}
	fmt.Println(".")
}
