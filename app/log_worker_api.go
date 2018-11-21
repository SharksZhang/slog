package app

import (
	"fmt"
	"sync"
	"github.com/SharksZhang/slog/domain/model"
	"github.com/SharksZhang/slog/domain/service"
)

var should_close bool

func UnmarshalLog(chTcp chan []byte, stop chan int, group *sync.WaitGroup, stopAll chan int) {

	for {
		if shouldReturn(chTcp) {
			fmt.Printf("unmarshal return\n ")
			return
		}
		select {
		case b := <-chTcp:
			log := &model.Log{}
			err := log.UnmarshalJSON(b[0 : len(b)-1])
			if err != nil {
				fmt.Printf("json.Unmarshal(%v) err, erro is [%v]\n", string(b), err)
				panic(err)
				continue
			}
			switch log.Type {
			case "content":
				service.GetLogWorkerService().AddLog(log)
			case "register":
				fmt.Println("register")
				service.GetLogWorkerService().Create(log, group, stopAll)
			case "control":
				if log.Action == "stop" {
					fmt.Println("stop log:[%v]\n", log)
					should_close = true
					stop <- 1
				}
			}
		default:

		}

	}

}

func shouldReturn(chTcp chan []byte) bool {
	return should_close == true && len(chTcp) == 0
}
