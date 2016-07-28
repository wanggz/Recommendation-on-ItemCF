package main

import (
	"os"
	"os/exec"
	"recommend"
	"strings"
	"time"
)

func main() {
	dir := "/data/go-recommend/society/"

	//log file init
	s := strings.Split(dir, "/")
	recommend.InitLogger(s[len(s)-2])

	//init file list
	recommend.InitFileList(dir)

	t1 := time.Now().Unix()
	recommend.Execute(dir)
	t2 := time.Now().Unix()
	recommend.Glogger.Println(" +++++++ end ++++++ time-taken:%d", (t2 - t1))

	cmd := exec.Command("/bin/sh", "/data/services/society/startsociety.sh")
	cmd.Stderr = os.Stdout
	err := cmd.Run()
	if err != nil {
		recommend.Glogger.Fatalln("cmd.Output:", err)
	}
	defer recommend.GlogFile.Close()
}
