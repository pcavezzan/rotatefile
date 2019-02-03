package main

import (
	"flag"
	"github.com/pcavezzan/gorotatefile"
	"log"
	"os"
)

var filesize int
var mode int

func init() {
	flag.IntVar(&filesize, "filesize", 100, "Maximum size of file before creating a rotation file in Mb. Default: 100Mb")
	flag.IntVar(&mode, "mode", 0, "Mode of execution: 0 (None), 1 (Timing), 2 (Verbose), 3 (Timing+Verbose). Default : 0")
}

func main() {
	flag.Parse()

	info, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	if info.Mode()&os.ModeCharDevice != 0 || info.Size() <= 0 {
		log.Println("This tools is intended to work with stding.")
		log.Println("Usage: cat myfile | gorotatefile or gorotatefile < myfile")
		os.Exit(1)
	}

	rf := gorotatefile.NewRotateFile(filesize*1000, gorotatefile.ExecMode(mode))
	rf.Rotate(os.Stdin)
}
