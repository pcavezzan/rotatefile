package gorotatefile

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

type RotateFile struct {
	MaxSize int
	Mode    ExecMode

	buffer *bytes.Buffer
	nFiles int
}

type ExecMode int

const (
	TIMING = iota + 1
	VERBOSE
)

func NewRotateFile(mSize int, mode ExecMode) RotateFile {
	maxBytes := mSize * 1024
	buf := bytes.NewBuffer(make([]byte, maxBytes))
	return RotateFile{
		MaxSize: mSize,
		buffer:  buf,
		Mode:    mode,
	}
}

func (rf *RotateFile) Rotate(f *os.File) {
	start, stop := rf.timer()
	start()
	rf.rotate(f)
	stop()
}

func (rf *RotateFile) rotate(f *os.File) {
	reader := bufio.NewReader(f)

	for {
		l, err := reader.ReadBytes('\n')

		if err != nil {
			if err == io.EOF {
				rf.flushToFile()
				break
			} else {
				log.Fatal(err)
			}
		}

		if len(l) >= (rf.buffer.Cap() - rf.buffer.Len()) {
			rf.flushToFile()
			rf.nFiles++
		}

		_, err = rf.buffer.Write(l)
		if err != nil {
			log.Fatalln(err)
		}

	}
}

func (rf RotateFile) timer() (func(), func()) {
	start := time.Now()
	return func() {
			if int(rf.Mode)%2 != 0 {
				log.Println("Timer enable - Rotatefile started at", start.Format("2006/01/02 15:04:05"))
			}
		}, func() {
			if int(rf.Mode)%2 != 0 {
				log.Println("Timer enable - Rotatefile ended at", start.Format("2006/01/02 15:04:05"))
				elapsedTime := time.Since(start)
				log.Println("Timer enable - Rotatefile ended took about", elapsedTime, "to split into", rf.nFiles, "files.")
			}
		}
}

func (rf *RotateFile) flushToFile() {
	if int(rf.Mode) == 2 || int(rf.Mode) == 3 {
		// Mode
		log.Println("Needs to flush into a file .....")
		log.Println()
	}

	fileName := fmt.Sprintf("output-%d.txt", rf.nFiles)

	f, errFile := os.Create(fileName)
	if errFile != nil {
		log.Fatalln(errFile)
		f.Close()
	}

	rf.buffer.WriteTo(f)
	f.Close()
}
