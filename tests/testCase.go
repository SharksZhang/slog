package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
	"time"
	"github.com/SharksZhang/slog/tests/test_log"
	"github.com/SharksZhang/slog/ui/tcp"
)

func main() {
	TestCase1()
	TestCase2()
	TestCase3()
}

var server = "127.0.0.1:8080"

//一个日志对象，一个source, 一个连接.运行两次，覆盖重复注册。
func TestCase1() {
	start_time := time.Now().UnixNano()
	fmt.Printf("start time stamp :[%v]", start_time)
	test_log.GenerateRegisterLog("register_testcase1.log", 1, test_log.INFO_LEVEL)
	SendFileData("register_testcase1.log")
	test_log.GenerateTestLog("INFO_testcase1.log", test_log.INFO_LEVEL, 10000000, 1, 1)
	SendFileData("INFO_testcase1.log")
	//SendFileData("INFO_testcase1.log")
	//SendFileData("INFO_testcase1.log")

	SendStopLog()
	end_time := time.Now().UnixNano()
	fmt.Printf("end time stamp :[%v]\n", end_time)

	fmt.Printf("testcase1 duration is [%v]", strconv.FormatInt((end_time-start_time), 10))

	//time.Sleep(2*time.Minute)
}

//生成50个日志对象
func TestCase2() {

	start_time := time.Now().Unix()
	fmt.Printf("start time stamp :[%v]", start_time)

	const TEST_CASE2_STR = "testcase2"
	const tc2RegLogFile = "register_testcase2.log"
	test_log.GenerateRegisterLog(tc2RegLogFile, 100, test_log.INFO_LEVEL)

	SendFileData(tc2RegLogFile)
	fs := []string{"log_debug.log", "log_info.log", "log_warn.log", "log_error.log", "log_fatal.log"}

	//test_log.GenerateTestLog(TEST_CASE2_STR+fs[0], test_log.DEBUG_LEVEL, 2000, 100,1)
	//test_log.GenerateTestLog(TEST_CASE2_STR+fs[1], test_log.INFO_LEVEL, 2000, 100,1)
	//test_log.GenerateTestLog(TEST_CASE2_STR+fs[2], test_log.WARN_LEVEL, 2000, 100,1)
	//test_log.GenerateTestLog(TEST_CASE2_STR+fs[3], test_log.ERROR_LEVEL, 2000, 100,1)
	//test_log.GenerateTestLog(TEST_CASE2_STR+fs[4], test_log.FATAL_LEVEL, 2000, 100,1)

	waitGroup1 := sync.WaitGroup{}
	for _, f := range fs {
		waitGroup1.Add(1)
		go func(file string) {
			defer waitGroup1.Done()
			SendFileData(file)
		}(TEST_CASE2_STR + f)
	}
	waitGroup1.Wait()
	SendStopLog()

	end_time := time.Now().Unix()
	fmt.Printf("end time stamp :[%v]", end_time)
	fmt.Printf("testcase2 duration is [%v]", strconv.FormatInt((end_time-start_time), 10))

}

func TestCase3() {
	errorFiles := []string{"log_debug_error.log", "log_error_error.log", "log_fatal_error.log", "log_info_error.log", "log_warn_error.log"}
	waitGroup := sync.WaitGroup{}
	for _, errorFile := range errorFiles {
		waitGroup.Add(1)
		go func(file string) {
			fmt.Printf("wg.add\n")
			defer waitGroup.Done()
			SendFileData(file)
			fmt.Printf("wg.done\n")
		}(errorFile)
	}
	waitGroup.Wait()

	registerFile := "register.log"
	SendFileData(registerFile)

	fs := []string{"log_debug.log", "log_error.log", "log_fatal.log", "log_info.log", "log_warn.log"}
	waitGroup1 := &sync.WaitGroup{}
	for _, File := range fs {
		waitGroup1.Add(1)
		go func(file string) {
			defer waitGroup1.Done()
			SendFileData(file)
		}(File)
	}
	waitGroup1.Wait()

}

func SendFileData(file string) {
	f, err := os.OpenFile("/testlog/"+file, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	defer f.Close()
	if err != nil {
		fmt.Printf("openfile err, err is [%v]", err)
		return
	}
	fmt.Printf("send [%v] data start \n", "/testlog/"+file)
	reader := bufio.NewReader(f)
	client, err := tcp.NewTcpClient(server)
	defer client.Close()
	if err != nil {
		fmt.Printf("New tcp client err, error os [%v]", err)
		return
	}
	for {
		bytes, err := reader.ReadBytes('\000')
		if err == io.EOF {
			fmt.Printf("file [%v] end\n", file)
			break
		}
		if err != nil {
			fmt.Printf("readBytes err, error is [%v]", err)
		}
		//fmt.Printf("file:[%v]\n",file)
		client.SendData(bytes)
	}
}

func SendStopLog() {
	client, err := tcp.NewTcpClient(server)
	if err != nil {
		fmt.Printf("NewTcpClient err, error is [%v]", err)
	}
	defer client.Close()
	stopLog := test_log.TestControlLog{
		Type:   "control",
		Action: "stop",
	}

	bytes, err := json.Marshal(stopLog)
	bytes = append(bytes, []byte("\000")...)
	client.SendData(bytes)
	fmt.Printf("send stop end \n")

}
