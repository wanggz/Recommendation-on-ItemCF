package recommend

import (
	"os"
	"os/exec"
	"strings"
	"time"
)

func InitFileList(dir string) {
	timeNow := time.Now().Format("2006-01-02")
	cpmd := exec.Command(
		"cp",
		strings.Join([]string{"/data/nfs_share/society/apply/", timeNow, ".txt"}, ""),
		strings.Join([]string{dir, "apply/"}, ""))
	cpmd.Stderr = os.Stdout
	err := cpmd.Run()
	if err != nil {
		Glogger.Fatalln("copy error:", err)
	}
	cpmd = exec.Command(
		"cp",
		strings.Join([]string{"/data/nfs_share/society/filter/", "4-", timeNow, ".txt"}, ""),
		strings.Join([]string{dir, "filter/"}, ""))
	cpmd.Stderr = os.Stdout
	err = cpmd.Run()
	if err != nil {
		Glogger.Fatalln("copy error:", err)
	}
	cpmd = exec.Command(
		"cp",
		strings.Join([]string{"/data/nfs_share/society/invite/", "3-", timeNow, ".txt"}, ""),
		strings.Join([]string{dir, "invite/"}, ""))
	cpmd.Stderr = os.Stdout
	err = cpmd.Run()
	if err != nil {
		Glogger.Fatalln("copy error:", err)
	}
}
