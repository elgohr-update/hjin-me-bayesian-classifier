package classifier

import (
	"os"
	"strings"
)

// 过滤干扰词 空格，单字词
// 英文单词小写标准化
func filterWord(ws []string) []string {
	result := make([]string, 0)
	for _, w := range ws {
		if strings.Count(strings.TrimSpace(w), "") <= 2 {
			continue
		}
		// 将英文词转化为小写
		w := strings.ToLower(w)
		result = append(result, w)
	}
	return result
}

// 读取一个文件夹返回文件列表
func ReadDir(dirName string) ([]os.FileInfo, error) {
	f, err := os.Open(dirName)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = f.Close()
	}()
	list, err := f.Readdir(-1)
	if err != nil {
		return nil, err
	}
	return list, nil
}
