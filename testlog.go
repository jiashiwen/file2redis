package main

import (
	"fmt"
	// "io"
	"errors"
	"log"
	"os"
	"time"
)

type ErrLog struct {
	logfile *os.File
	prefix  string
	silent  bool
	content interface{}
}

func main() {
	var st string
	st = time.Now().Format("20060102150405")
	fmt.Println(st)
	var log ErrLog
	file, err := os.OpenFile("/home/develop/goproject/src/file2redis/testlog", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0660)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	log.logfile = file
	er := errors.New("error:cccccccccccccccccc!")
	log.content = er
	log.DealLog()
	// loger := log.New(io.Writer(logfile), "Log:", log.Ldate|log.Ltime|log.Llongfile)

}

func (errlog *ErrLog) DealLog() {
	log.SetPrefix(errlog.prefix)
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	if !errlog.silent {
		log.SetOutput(os.Stdout)
		log.Println(errlog.content)
		log.SetOutput(errlog.logfile)
		log.Println(errlog.content)
	} else {
		log.SetOutput(errlog.logfile)
		log.Println(errlog.content)
	}
}
