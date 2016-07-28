package recommend

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func FileOpen(file string, process func(int, string)) {
	//"material/download.txt"
	index := 0

	f, err := os.Open(file)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	defer f.Close()
	br := bufio.NewReader(f)
	for {
		line, err := br.ReadString(byte('\n'))
		if err == io.EOF {
			if num := len(line); num > 0 {
				process(index, line)
			}
			break
		} else {
			process(index, line)
		}
		index++
	}
}

func getFilePath(fileName string) string {
	basePath, err := os.Getwd()
	if err != nil {
		return fileName
	}
	filePath := basePath + "/" + fileName
	return filePath
}

func CpFileToNFS(outFile string) {
	cpmd := exec.Command(
		"cp",
		outFile,
		"/data/nfs_share/society/out/position/")
	cpmd.Stderr = os.Stdout
	err := cpmd.Run()
	if err != nil {
		Glogger.Fatalln("copy error:"+outFile, err)
	}
}
