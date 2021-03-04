package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

type file struct {
	name    string
	comment string
	modTime time.Time
	data    string
}

func main() {
	var buf = writeGz()
	readGz(buf)
}

func writeGz() *bytes.Buffer {
	var buf bytes.Buffer
	var zw = gzip.NewWriter(&buf)

	var files = []file{
		//Year, month, day, hour, min, sec, nsec, loc
		{"text.txt", "file-header-txt", time.Date(2010, time.July, 7, 7, 47, 35, 0, time.UTC), "Insert text txt here"},
		{"numbers.cvs", "file-header-csv", time.Date(2018, time.April, 6, 12, 18, 9, 8, time.UTC), "1,2,3\n4,5,6\n7,8,9"},
		{"names.txt", "file-header-txt", time.Date(2018, time.March, 1, 4, 5, 6, 1, time.UTC), "Marco - John - Dante"},
	}

	for _, file := range files {
		zw.Name = file.name
		zw.Comment = file.comment
		zw.ModTime = file.modTime

		var _, err = zw.Write([]byte(file.data))
		if err != nil {
			log.Fatal(err)
		}
		err = zw.Close()
		if err != nil {
			log.Fatal(err)
		}
		// Restart zw
		zw.Reset(&buf)
	}
	return &buf
}

func readGz(buf *bytes.Buffer) {
	var zr, err = gzip.NewReader(buf)
	defer zr.Close()
	if err != nil {
		log.Fatal(err)
	}

	for {
		// Change default value (true)
		zr.Multistream(false) // No sequence
		fmt.Printf("Name: %s\nComment: %s\nModTime: %s \nData: \n", zr.Name, zr.Comment, zr.ModTime.UTC())

		var _, err = io.Copy(os.Stdout, zr)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print("\n\n")

		err = zr.Reset(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
		}
	}
}
