package impl

import (
	"testing"
	"github.com/SharksZhang/slog/domain/model"
	"fmt"
)

func TestLogWorkerRepoImpl(t *testing.T) {
	testKey := "test"
	testWorker := &model.LogWorker{
		ID          :testKey,
		FileName    :testKey,
		FileSizeKB  :1024,
		FilterLevel :1,
	}

	testLogWorkerRepoImpl := &LogWorkerRepoImpl{}

	testLogWorkerRepoImpl.Add(testKey, testWorker)
	testLogWorkerRepoImpl.Update(testWorker)
	testLogWorkerRepoImpl.Remove(testKey)

	result := testLogWorkerRepoImpl.Get(testKey)
	if result.ID != testKey {
		fmt.Errorf("TestLogWorkerRepoImpl Fail: Get LogWorker [%v] err!", testWorker)
		t.Error("TestLogWorkerRepoImpl Fail")
	}

	result = testLogWorkerRepoImpl.Get("wrong_key")
	if result != nil {
		fmt.Errorf("TestLogWorkerRepoImpl Fail: Get LogWorker [%v] err!", testWorker)
		t.Error("TestLogWorkerRepoImpl Fail")
	}

}
