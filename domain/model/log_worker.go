package model

import (
	"fmt"
	"os"
	"sync"
)

const (
	DEBUG = iota
	INFO
	WARN
	ERROR
	FATAL
)

type LogWorker struct {
	ID          string
	FileName    string
	FileSizeKB  int
	FilterLevel int
	logChan     chan *Log
}

func NewLogWorker(log *Log) *LogWorker {
	lw := &LogWorker{
		ID:         log.Name,
		FileName:   log.Name,
		FileSizeKB: log.FileSizeKb,
	}
	lw.FilterLevel = transferIntLogLevel(log.FilterLevel)


	return lw
}

func (lw *LogWorker) InitLogSaver(wg *sync.WaitGroup, stopAll chan int)  {
	dir_path := "/log/" + lw.FileName
	_, err := os.Stat(dir_path)
	if err != nil && os.IsNotExist(err) {
		err := os.Mkdir(dir_path, os.ModePerm)
		if err != nil && !os.IsExist(err) {
			panic(err)
			fmt.Printf("os.Mkdir([%v]) err, error is [%v]", dir_path, err)
		}
	}

	wg.Add(1)
	lw.logChan = make(chan *Log, 64)
	go func() {
		defer wg.Done()
		NewLogSaver(lw.FileName, lw.logChan).Run(stopAll)
		fmt.Printf("wg.Donw, directory :[%v]\n", lw.FileName)
	}()
}
func (lw *LogWorker) SaveLog(log *Log) {
	level := transferIntLogLevel(log.Level)
	if level >= lw.FilterLevel {
		log.FileSizeKb = lw.FileSizeKB
		lw.logChan <- log

	}
}

func transferIntLogLevel(level_str string) int {
	var level int
	switch level_str {
	case "DEBUG":
		level = DEBUG
	case "INFO":
		level = INFO
	case "WARN":
		level = WARN
	case "ERROR":
		level = ERROR
	case "FATAL":
		level = FATAL
	}
	return level
}

