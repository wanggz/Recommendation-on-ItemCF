package recommend

import (
	"io/ioutil"
	"os"
	"strings"
	"time"
)

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

func ExecuteOld(baseDir string) {
	aToBMap, bToAMap := make(MapSlices), make(MapSlices)

	files1 := ListDir(baseDir+"filebyday", ".txt")
	files2 := ListDir(baseDir+"filebyhour", ".txt")
	files := files1
	for _, path := range files2 {
		files = append(files, path)
	}
	Glogger.Println("files:", files)
	for _, path := range files {
		tmpAToBMap, tmpBToAMap := BuildStructure(path, "^", 0, 2)

		aToBMap.Merge(tmpAToBMap)
		bToAMap.Merge(tmpBToAMap)
	}
	Glogger.Println("aToBMap,bToAMap=", len(aToBMap), ",", len(bToAMap))
	inviteBToAMap := make(map[int64]int64)
	start2 := time.Now().Unix()
	Calculate(baseDir+"out/", bToAMap, aToBMap, inviteBToAMap)
	end2 := time.Now().Unix()
	Glogger.Println("write file times=", end2-start2)
}

func Execute(baseDir string) {

	aToBMap, bToAMap := make(MapSlices), make(MapSlices)

	files := ListDir(baseDir+"apply", ".txt")
	Glogger.Println("files:", files)
	//index:=0
	start := time.Now().Unix()
	for _, path := range files {
		tmpAToBMap, tmpBToAMap := BuildStructure(path, "^", 0, 2)

		aToBMap.Merge(tmpAToBMap)
		bToAMap.Merge(tmpBToAMap)
	}
	Glogger.Println("aToBMap,bToAMap=", len(aToBMap), ",", len(bToAMap))

	delAToBMap, delBToAMap := make(MapSlices), make(MapSlices)
	filterfiles := ListDir(baseDir+"filter", ".txt")
	if len(filterfiles) > 0 {
		for _, path := range filterfiles {
			tmpAToBMap, tmpBToAMap := BuildStructure(path, "^", 0, 2)

			delAToBMap.Merge(tmpAToBMap)
			delBToAMap.Merge(tmpBToAMap)
		}
		Glogger.Println("delAToBMap,delBToAMap=", len(delAToBMap), ",", len(delBToAMap))
		for k, v := range delAToBMap {
			if len(aToBMap[k]) > 0 {
				QuickSort(aToBMap[k])
				QuickSort(v)
				aToBMap[k] = removeAll(aToBMap[k], v)
			}
		}
		for k, v := range delBToAMap {
			if len(bToAMap[k]) > 0 {
				QuickSort(bToAMap[k])
				QuickSort(v)
				bToAMap[k] = removeAll(bToAMap[k], v)
			}
		}
		end := time.Now().Unix()
		Glogger.Println("read file & remove filter times=", end-start)
	}
	inviteBToAMap := make(map[int64]int64)
	invitefiles := ListDir(baseDir+"invite", ".txt")
	if len(invitefiles) > 0 {
		for _, path := range invitefiles {
			tmpBToAMap := BuildInviteStructure(path, "^", 0, 2)
			for k, _ := range tmpBToAMap {
				if _, ok := inviteBToAMap[k]; ok {
				} else {
					inviteBToAMap[k] = tmpBToAMap[k]
				}
			}
		}
	}

	start2 := time.Now().Unix()
	Calculate(baseDir+"out/", bToAMap, aToBMap, inviteBToAMap)
	end2 := time.Now().Unix()
	Glogger.Println("write file times=", end2-start2)
	//CreateFileForMapMap2(baseDir+"out/result.txt", result)
}

func removeAll(values []int64, deletes []int64) []int64 {
	j := 0
	newSlice := values[0:0]
	for i := 0; i < len(values); i++ {
		//fmt.Println(newSlice)
		if values[i] < deletes[j] {
			newSlice = append(newSlice, values[i])
		} else if values[i] == deletes[j] {
			j++
			if j == len(deletes) {
				for z := i + 1; z < len(values); z++ {
					newSlice = append(newSlice, values[z])
				}
				break
			}
		} else if values[i] > deletes[j] {
			newSlice = append(newSlice, values[i])
			j++
			if j == len(deletes) {
				for z := i + 1; z < len(values); z++ {
					newSlice = append(newSlice, values[z])
				}
				break
			}
		}
	}
	return newSlice
}
