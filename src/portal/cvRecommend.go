package main

import (
	"fmt"
	"os"
	"recommend"
	"strings"
	"time"
)

func main() {
	args := os.Args
	if args == nil || len(args) < 2 {
		fmt.Println(" no args ! ")
		return
	}
	s := strings.Split(args[1], "/")
	recommend.InitLogger(s[len(s)-2])
	t1 := time.Now().Unix()
	recommend.ExecuteOld(args[1])
	t2 := time.Now().Unix()
	recommend.Glogger.Println(" +++++++ end ++++++ time-taken:%d", (t2 - t1))
	defer recommend.GlogFile.Close()
}
