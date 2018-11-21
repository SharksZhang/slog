package service

import (
	"sync"
	"github.com/SharksZhang/slog/domain/model"
	"fmt"
)

type LogWorkerService struct {
	repo model.LogWorkerRepo
}

var lw = &LogWorkerService{}
var lwOnce sync.Once

func GetLogWorkerService() *LogWorkerService {
	lwOnce.Do(func() {
		lw.repo = model.GetLogWorkerRepo()
	})
	return lw
}

func (ls *LogWorkerService) Create(l *model.Log, group *sync.WaitGroup, stopAll chan int) {
	worker := model.NewLogWorker(l)
	actual, loaded := ls.repo.Add(l.Name, worker)
	if !loaded{
		fmt.Print("register log name:[%v]", l.Name)
		actual.InitLogSaver(group, stopAll)
	}else{
		if actual != nil{
			actual.FileSizeKB = worker.FileSizeKB
			actual.FilterLevel = worker.FilterLevel
		}
	}


}

func (ls *LogWorkerService) AddLog(l *model.Log) {
	worker := ls.repo.Get(l.Name)
	if worker != nil {
		worker.SaveLog(l)
	}

}

