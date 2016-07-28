package recommend

import (
	"fmt"
	"log"
	"os"
	"time"
)

var Glogger *log.Logger
var GlogFile *os.File

func InitLogger(baseDir string) {
	t := time.Now().Format("2006-01-02")
	GlogFile, err := os.OpenFile("/data/logs/go-recommend-"+baseDir+"."+t+".log", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0)
	if err != nil {
		fmt.Printf("%s\r\n", err.Error())
		os.Exit(-1)
	}
	Glogger = log.New(GlogFile, "", log.Ldate|log.Ltime|log.Llongfile)
	Glogger.Printf(" logger init ! ")
	return
}
