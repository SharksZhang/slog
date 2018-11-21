package model

import (
	"bufio"
	"fmt"
	"os"
	"runtime/debug"
	"strconv"
	"time"

)

const (
	LOG_PATH_PREFIX = "/log/"
	LOG_PATH_SUFFIX = ".log"
)

type LogSaver struct {
	Directory         string
	current_file_path string
	ch                chan *Log
}

func NewLogSaver(directory string, ch chan *Log) *LogSaver {
	return &LogSaver{
		Directory:         directory,
		ch:                ch,
		current_file_path: LOG_PATH_PREFIX + directory + "/" + directory + "_" + getNowTimeStamp() + LOG_PATH_SUFFIX,
	}
}

func getNowTimeStamp() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}
func (ls *LogSaver) Run(stop chan int) {
	fmt.Printf("LogSaver.InitLogSaver [%v]\n", ls)
	defer RecoverPanic()
	var f *os.File
	var err error
	var info os.FileInfo

	f, err = os.OpenFile(ls.current_file_path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	defer f.Close()
	if err != nil {
		fmt.Errorf("os.OpenFile(ls.current_file_path,os.O_CREATE|os.O_APPEND, 0666) err, err is [%v]", err)
		panic("os open file " + err.Error())
	}

	var writer *bufio.Writer
	writer = bufio.NewWriter(f)
	defer writer.Flush()

	var timer = time.NewTimer(5 * time.Second)
	//todo
	for {
		select {
		case log := <-ls.ch:
			info, err = f.Stat()
			if err != nil {
				fmt.Printf("f :(%v) err, error is [%v]\n", f, err)
				panic(err)
				continue
			}

			if int(info.Size()/1024) > log.FileSizeKb {
				fmt.Printf("info.size:[%v], info.size/1024:[%v], fileSizeKb:[%v], currentfile:[%v]\n",
					info.Size(), info.Size()/1024, log.FileSizeKb, ls.current_file_path)
				writer, f = reCreateFile(writer, f, ls)
			}

			writer.Write([]byte(log.transferToString()))
			timer.Reset(5 * time.Second)

		default:
			if timer == nil {
				timer = time.NewTimer(5 * time.Second)
			}
			select {
			case c := <-timer.C:
				fmt.Printf("timer.C: c :[%v], directory:[%v ] \n", c, ls.Directory)
				writer.Flush()

			case <-stop:
				writer.Flush()
				f.Close()
				fmt.Printf("ls.ch.len:[%v]", len(ls.ch))
				return
			default:
			}

		}
	}
}

func reCreateFile(writer *bufio.Writer, f *os.File, ls *LogSaver) (*bufio.Writer, *os.File) {
	writer.Flush()
	f.Close()
	ls.current_file_path = LOG_PATH_PREFIX + ls.Directory + "/" + ls.Directory + "_" + getNowTimeStamp() + LOG_PATH_SUFFIX
	var err error
	f, err = os.OpenFile(ls.current_file_path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		//todo check err
		fmt.Errorf("os.OpenFile(ls.current_file_path,os.O_CREATE|os.O_APPEND, 0666) err, err is [%v]\n", err)
		panic("os.openFile err")
	}
	writer = bufio.NewWriter(f)
	return writer, f
}

func RecoverPanic() {
	if p := recover(); p != nil {
		fmt.Printf("panic: [%v] ,Stack:%v", p, string(debug.Stack()))
	}
}
