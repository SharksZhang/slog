package model

type LogWorkerRepo interface {
	Add(key string ,worker *LogWorker)(actual *LogWorker, loaded bool)
	Get(workerID string) *LogWorker
	Update(worker *LogWorker)
	Remove(WorkerID string)
}

var logWorkerRepo LogWorkerRepo

func SetLogWorkerRepo(repo LogWorkerRepo) {
	logWorkerRepo = repo
}

func GetLogWorkerRepo() LogWorkerRepo {
	return logWorkerRepo
}
