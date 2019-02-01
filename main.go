package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)


func main() {

	info, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	if info.Mode()&os.ModeCharDevice != 0 || info.Size() <= 0 {
		log.Println("This tools is intended to work with stding.")
		log.Println("Usage: cat myfile | gorotatefile or gorotatefile < myfile")
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)
	now := time.Now()

	maxBytes := 70 * 1024

	buf := bytes.NewBuffer(make([]byte, 1000, maxBytes))
	nFiles := 0

	for {
		l, err := reader.ReadBytes('\n')

		if err != nil {
			if err == io.EOF {
				flushToFile(buf, nFiles)
				break
			} else  {
				log.Fatal(err)
			}
		}


		if len(l) >= (buf.Cap() - buf.Len()) {
			flushToFile(buf, nFiles)
			nFiles++
		}

		_, err = buf.Write(l)
		if err != nil {
			log.Fatalln(err)
		}

	}
	log.Println("took", (time.Now().Nanosecond() - now.Nanosecond()) / 1000000, "ms to split into", nFiles, "files.")

}


func flushToFile(buf *bytes.Buffer, nFiles int) {
	log.Println("Needs to flush into a file .....")
	log.Println()
	fileName := fmt.Sprintf("output-%d.txt", nFiles)

	f, errFile := os.Create(fileName)
	if errFile != nil {
		log.Fatalln(errFile)
		f.Close()
	}

	buf.WriteTo(f)
	f.Close()
}