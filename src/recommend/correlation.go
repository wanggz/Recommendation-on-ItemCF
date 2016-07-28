package recommend

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

type MapSlices map[int64][]int64

func CreateFileForMapSlice(file string, tmps map[int64][]int64) {

	hrFile, err := os.Create(getFilePath(file))
	defer hrFile.Close()
	if err != nil {
		fmt.Println(file, err)
		return
	}
	for k, v := range tmps {
		hrFile.WriteString(strconv.FormatInt(k, 10) + "\t" + JoinForInt64(v, "\t") + "\n")
	}

}

func CreateFileForMapMap(file string, tmps map[int64]map[int64]float64) {

	hrFile, err := os.Create(getFilePath(file))
	defer hrFile.Close()
	if err != nil {
		fmt.Println(file, err)
		return
	}
	for k, v := range tmps {
		var line string = ""
		for resume, score := range v {
			line += "\t"
			line += strconv.FormatInt(resume, 10)
			line += ","
			line += strconv.FormatFloat(score, 'f', 3, 64)
		}
		hrFile.WriteString(strconv.FormatInt(k, 10) + line + "\n")
	}

}

func CreateFileForMapMap2(file string, tmps map[int64]map[int64]float64) {

	hrFile, err := os.Create(file)

	w := bufio.NewWriter(hrFile)
	defer hrFile.Close()
	if err != nil {
		fmt.Println(file, err)
		return
	}
	start := time.Now().Unix()
	for k, v := range tmps {
		var line string = ""
		for resume, score := range v {
			line += "\t"
			line += strconv.FormatInt(resume, 10)
			line += ","
			line += strconv.FormatFloat(score, 'f', 3, 64)
		}
		w.WriteString(strconv.FormatInt(k, 10) + line + "\n")
	}
	w.Flush()
	end := time.Now().Unix()
	Glogger.Println("write file times=", end-start)
}

func BuildStructure(file string, sep string, aIndex int, bIndex int) (MapSlices, MapSlices) {
	//"material/20140822110003_195910.txt"
	aToBMap := make(MapSlices)
	bToAMap := make(MapSlices)

	start := 0
	FileOpen(file, func(position int, line string) {
		tmps := strings.Split(line, sep)
		if len(tmps) < 2 || position < 0 {
			return
		}
		aId, _ := strconv.ParseInt(tmps[aIndex], 10, 64)
		bId, _ := strconv.ParseInt(strings.TrimSpace(tmps[bIndex]), 10, 64)
		if _, ok := aToBMap[aId]; ok {
			if !ContainForInt64(aToBMap[aId], bId) {
				aToBMap[aId] = append(aToBMap[aId], bId)
			}
		} else {
			aToBMap[aId] = []int64{bId}
		}
		if _, ok := bToAMap[bId]; ok {
			if !ContainForInt64(bToAMap[bId], aId) {
				bToAMap[bId] = append(bToAMap[bId], aId)
			}
		} else {
			bToAMap[bId] = []int64{aId}
		}
		start++
	})
	return aToBMap, bToAMap
}

func BuildInviteStructure(file string, sep string, aIndex int, bIndex int) map[int64]int64 {
	bToAMap := make(map[int64]int64)

	start := 0
	FileOpen(file, func(position int, line string) {
		tmps := strings.Split(line, sep)
		if len(tmps) < 2 || position < 0 {
			return
		}
		bId, _ := strconv.ParseInt(strings.TrimSpace(tmps[bIndex]), 10, 64)
		if _, ok := bToAMap[bId]; ok {

		} else {
			bToAMap[bId] = 1
		}
		start++
	})
	return bToAMap
}

func Calculate(filedir string, bToaMap MapSlices, aToBMap MapSlices, inviteBToAMap map[int64]int64) {
	//scoreMap := make(map[int64]map[int64]float64)
	file := strings.Join([]string{filedir, time.Now().Format("2006-01-02-15"), ".txt"}, "")
	hrFile, err := os.Create(file)
	w := bufio.NewWriter(hrFile)
	defer hrFile.Close()
	if err != nil {
		fmt.Println(file, err)
		return
	}

	for k, v := range bToaMap {
		if len(v) < 2 {
			continue
		}

		correlateMap := make(map[int64]int)
		for _, hr := range v {

			for _, resume := range aToBMap[hr] {
				if resume == k {
					continue
				}
				if _, ok := correlateMap[resume]; ok {
					correlateMap[resume] += 1
				} else {
					correlateMap[resume] = 1
				}
			}
		}
		curr := len(bToaMap[k])
		scoreMap := make(map[int64]float64)
		//topCount := 0
		for resume, inter := range correlateMap {
			if inter > 1 {
				count := len(bToaMap[resume])
				if count < 2 {
					//count = inter
					continue
				}

				score := float64(inter) / (math.Pow(float64(curr), 0.5) * math.Pow(float64(count), 1-0.5))
				if _, ok := inviteBToAMap[resume]; ok {
					//if score > 0.2 && topCount < 2 {
					//score = score + 1
					//topCount = topCount + 1
					//} else {
					if score > 0.2 {
						score = math.Pow(score, 0.6)
					} else {
						score = score * 1.9
					}
					//}
				}
				scoreMap[resume] = score
			}
		}

		//var line string = ""
		line := []string{}
		line = append(line, strconv.FormatInt(k, 10))
		for resume, score := range scoreMap {
			line = append(line, "\t")
			line = append(line, strconv.FormatInt(resume, 10))
			line = append(line, ",")
			line = append(line, strconv.FormatFloat(score, 'f', 3, 64))
		}
		line = append(line, "\n")
		w.WriteString(strings.Join(line, ""))
	}
	w.Flush()

	//cp outFile to nfs
	CpFileToNFS(file)
}

func In_slice(val int64, slice []int64) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}

func (mainMap MapSlices) Merge(tmp MapSlices) {
	for k, _ := range tmp {
		if _, ok := mainMap[k]; ok {
			for _, v := range tmp[k] {
				if !In_slice(v, mainMap[k]) {
					mainMap[k] = append(mainMap[k], v)
				}
			}
		} else {
			mainMap[k] = tmp[k]
		}
	}
	return
}
