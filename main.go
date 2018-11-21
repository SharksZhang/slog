package main

import (
	"github.com/SharksZhang/slog/infra/tcp"

	"fmt"
	"os"
	"sync"
	"github.com/SharksZhang/slog/app"
	"github.com/SharksZhang/slog/domain/model"
	"github.com/SharksZhang/slog/infra/impl"
	_ "net/http/pprof"
	"io/ioutil"
	"encoding/json"
	"strconv"
)

func main() {

	bytes, err := ioutil.ReadFile("/opt/config.json")
	if err != nil {
		panic(err)
	}

	config := &config{}
	err = json.Unmarshal(bytes, config)
	if err!= nil{
		panic(err)
	}


	log_dir_path := config.LogDirPath
	err = createLogDir(log_dir_path)
	if err != nil {
		panic(err)
	}

	var stop_control_msg = make(chan int)
	var stopAll = make(chan int)
	logWorkerRepo := &impl.LogWorkerRepoImpl{
	}
	model.SetLogWorkerRepo(logWorkerRepo)

	chTcp := make(chan []byte, config.Cache)
	tcpServer := tcp.NewServer(":" + strconv.Itoa(config.TcpPort), chTcp)
	if tcpServer == nil {
		fmt.Println("tcp.NewServer() error")
		return
	}
	go tcpServer.Serve()


	var logSaverWg = &sync.WaitGroup{}
	var unMarWg = &sync.WaitGroup{}
	for i := 0; i < config.ConcurrentNum; i++ {
		unMarWg.Add(1)
		go func() {
			defer unMarWg.Done()
			app.UnmarshalLog(chTcp, stop_control_msg, logSaverWg, stopAll)
		}()
	}
	<-stop_control_msg
	unMarWg.Wait()
	close(stopAll)
	logSaverWg.Wait()
	fmt.Printf("receive packet:[%v]", tcp.ReceivePact)
	fmt.Printf("main exit\n")
}

func createLogDir(logDirPath string)  error {
	_, err := os.Stat(logDirPath)
	if err != nil && os.IsNotExist(err) {
		err := os.Mkdir(logDirPath, os.ModePerm)
		if err != nil {
			fmt.Printf("os.Mkdir([%v]) err, error is [%v]", logDirPath, err)
			return err
		}
	}
	return nil
}

type config struct {
	LogDirPath string `json:"log_dir_path"`
	ConcurrentNum int `json:"concurrent_num"`
	Cache int `json:"cache"`
	TcpPort int `json:"tcp_port"`

}
