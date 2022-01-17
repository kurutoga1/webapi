package log_test

import (
	"os"
	"sync"
	"testing"
	"webapi/utils/log"
)

var (
	logFile = "log.txt"
	logger  = log.GetLogger(logFile)
)

func init() {
	file, err := os.Create(logFile)
	if err != nil {
		panic(err.Error())
	}

	lines := []string{
		"aaaaaaaaaaaaaaaaaaaaaaa \n",
		"bbbbbbbbbbbbbbbbbbbbbbb \n",
		"ccccccccccccccccccccccc \n",
	}

	for _, line := range lines {
		b := []byte(line)
		_, err := file.Write(b)
		if err != nil {
			panic(err.Error())
		}
	}

}

func tearDown() {
	os.Remove(logFile)
}

func TestRotate(t *testing.T) {
	f, err := os.Stat(logFile)
	if err != nil {
		panic(err)
	}

	shavingSize := 30
	maxSize := 40
	var mu sync.Mutex

	rotater := log.NewLogRotater(shavingSize, maxSize, &mu, logger, logFile)
	err = rotater.Rotate()
	if err != nil {
		panic(err.Error())
	}

	f, err = os.Stat(logFile)
	if err != nil {
		panic(err)
	}

	if f.Size() > int64(shavingSize) {
		t.Errorf("log file size is more than %v \n", shavingSize)
	}

	t.Cleanup(func() {
		tearDown()
	})
}
