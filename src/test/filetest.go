package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	files1 := ListDir("E:\\result", ".txt")
	fmt.Println(files1)
	for _, path := range files1 {
		fmt.Println(path)
	}
	fmt.Println("--------------------")
	for i := 0; i < len(files1); i++ {
		fmt.Println(files1[i])
	}
	fmt.Println("--------------------")
	for i := (len(files1) - 1); i > 0; i-- {
		fmt.Println(files1[i])
	}
}

//获取指定目录下的所有文件，不进入下一级目录搜索，可以匹配后缀过滤。
func ListDir(dirPth string, suffix string) (files []string) {
	files = make([]string, 0, 10)
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil
	}
	PthSep := string(os.PathSeparator)
	suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写
	for _, fi := range dir {
		if fi.IsDir() { // 忽略目录
			continue
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) { //匹配文件
			files = append(files, dirPth+PthSep+fi.Name())
		}
	}
	return files
}
