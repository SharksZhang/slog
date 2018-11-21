package test_log

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)


const (
	DEBUG_LEVEL  = "DEBUG"
	INFO_LEVEL   = "INFO"
	WARN_LEVEL   = "WARN"
	ERROR_LEVEL  = "ERROR"
	FATAL_LEVEL  = "FATAL"
)




func GenerateContentLog() {

	GenerateTestLog("log_debug_error.log", DEBUG_LEVEL, 200, 100,100)
	GenerateTestLog("log_debug.log", DEBUG_LEVEL, 200, 100, 100)
	GenerateTestLog("log_info.log", INFO_LEVEL, 200, 100, 100)
	GenerateTestLog("log_info_error.log", INFO_LEVEL, 200, 100, 100)
	GenerateTestLog("log_warn.log", WARN_LEVEL, 200, 100,100)
	GenerateTestLog("log_warn_error.log", WARN_LEVEL, 200, 100, 100)
	GenerateTestLog("log_error.log", ERROR_LEVEL, 200, 100, 100)
	GenerateTestLog("log_error_error.log", ERROR_LEVEL, 200, 100, 100)
	GenerateTestLog("log_fatal.log", FATAL_LEVEL, 200, 100, 100)
	GenerateTestLog("log_fatal_error.log", FATAL_LEVEL, 200, 100, 100)
}

func GenerateRegisterLog(fileName string ,number int, level string) {
	f, err := os.OpenFile("/testlog/"+fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	defer f.Close()
	if err != nil {
		fmt.Printf("openfile err, err is [%v]", err)
		return
	}

	writer := bufio.NewWriter(f)
	for i := 0; i < number; i++ {
		rl := TestRegisterLog{
			Type:        "register",
			Name:        "test" + strconv.Itoa(i),
			FileSizeKB:  1024,
			FilterLevel: level,
		}

		b, err := json.Marshal(rl)
		b = append(b, []byte("\000")...)
		if err != nil {
			fmt.Printf("Marshal err, error is [%v]", err)
		}

		writer.Write(b)

	}
	defer writer.Flush()

}


func GenerateTestLog(fileName string, level string, numbers int, componetNumber int , sourceNumber int) {
	f, err := os.OpenFile("/testlog/" +fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	defer f.Close()
	if err != nil {
		fmt.Printf("openfile err, err is [%v]", err)
		return
	}

	writer := bufio.NewWriter(f)
	for i := 0; i < numbers; i++ {
		for j := 0; j < componetNumber; j++ {
			for k := 0; k < sourceNumber; k++{

				rl := &TestLog{
					Type:     "content",
					Name:     "test" + strconv.Itoa(j),
					Level:    level,
					Time:     time.Now().String(),
					Source:   "source" + strconv.Itoa(k),
					File:     "aaa",
					Line:     1,
					Function: "test1",
					Content:  "test content " + GetNowTimeStamp(),
				}
				b, err := json.Marshal(rl)
				b = append(b, []byte("\000")...)
				if err != nil {
					fmt.Printf("Marshal err, error is [%v]", err)
				}
				writer.Write(b)

			}

		}
	}
	fmt.Printf("file [%v] last log timestamp :[%v]\n", "/testlog/" +fileName, GetNowTimeStamp())
	defer writer.Flush()

}

func GetNowTimeStamp() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}



type TestRegisterLog struct {
	Type string  `json:"type"`
	Name string  `json:"name"`
	FileSizeKB int `json:"file_size_KB"`
	FilterLevel string `json:"filter_level"`
}

type TestLog struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	Level       string `json:"level"`
	Time        string `json:"time"`
	Source      string `json:"source"`
	File        string `json:"file"`
	Line        int    `json:"line"`
	Function    string `json:"function"`
	Content     string `json:"content"`
}

type TestControlLog struct {
	Type        string `json:"type"`
	Action      string `json:"action"`

}