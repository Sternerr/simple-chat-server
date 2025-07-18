package main

import (
	"os"
	"log"
)

func main() {
	logger := getLogger("/tmp/termtalk_log.txt")
	(*logger).Println("I will get the job done!")
}

func getLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic("could not open file")
	}

	return log.New(logfile, "[termtalk]", log.Ldate|log.Ltime|log.Lshortfile)
}
