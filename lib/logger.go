package lib

import (
	"log"
	"os"
)

func GetLogger(logfile string) *log.Logger {
	file, err := os.OpenFile(logfile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalln("fail to create " + logfile + " file!")
	}
	logger := log.New(file, "", log.LstdFlags|log.Llongfile)
	return logger
}
