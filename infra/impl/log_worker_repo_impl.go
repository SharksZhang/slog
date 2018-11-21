package impl

import (
	"fmt"

	"sync"
	"github.com/SharksZhang/slog/domain/model"
)

type LogWorkerRepoImpl struct {

	sync.Map
}

func (l *LogWorkerRepoImpl) Add(key string, worker *model.LogWorker) (actual *model.LogWorker, loaded bool)  {
	actual2, loaded := l.LoadOrStore(key, worker)

	logWorker, ok := actual2.(*model.LogWorker)
	if ok{
		return logWorker, loaded
	}
	return nil, false

}

func (l *LogWorkerRepoImpl) Get(workerID string) *model.LogWorker {
	value, ok := l.Load(workerID)
	if ok {
		worker, ok1 := value.(*model.LogWorker)
		if ok1 {
			return worker
		}
	}
	return nil
}

func (l *LogWorkerRepoImpl) Update(worker *model.LogWorker) {
	fmt.Println("implement me")
}

func (l *LogWorkerRepoImpl) Remove(WorkerID string) {
	fmt.Println("implement me")
}
