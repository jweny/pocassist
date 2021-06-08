package file

import (
	"bufio"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// Reading file and return content as []string
func ReadingLines(filename string) []string {
	var result []string
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		return result
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		val := scanner.Text()
		if val == "" {
			continue
		}
		result = append(result, val)
	}

	if err := scanner.Err(); err != nil {
		return result
	}
	return result
}

func UploadTargetsPath(extstring string) string {
	// 生成文件名
	fileNameInt := time.Now().Unix()
	fileNameStr := strconv.FormatInt(fileNameInt,10)
	fileName := fileNameStr + extstring
	// 格式化当前时间
	folderName := time.Now().Format("2006/01/02")
	folderPath := filepath.Join("upload",folderName)
	// 创建多层级目录
	os.MkdirAll(folderPath, os.ModePerm)

	filePath := filepath.Join(folderPath, fileName)
	return  filePath
}

// 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path)    //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

