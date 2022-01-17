package memoryGetter_test

import (
	"log"
	"testing"
	"webapi/gw/memoryGetter"
	"webapi/tests"
	"webapi/utils"
)

var (
	getter     memoryGetter.Getter
	currentDir string
	// CreateDummyServer関数のcleanUpで削除するファイルたち
	deletes []string
	addrs   []string
)

func init() {
	// サーバを立てるとカレントディレクトリにfileserverディレクトリとlog.txt
	// ができるのでそれを削除する。
	c, err := utils.GetCurrentDir()
	if err != nil {
		log.Fatalln(err.Error())
	}
	currentDir = c
	getter = memoryGetter.NewMemoryGetter()
	addrs, _, err = tests.GetSomeStartedProgramServer(1)
	if err != nil {
		log.Fatalln(err)
	}
}

func TestGet(t *testing.T) {

	url := addrs[0] + "/json/health/memory"

	t.Run("test 1", func(t *testing.T) {
		testGet(t, url)
	})
}

func testGet(t *testing.T, addr string) {
	t.Helper()
	memoryInfo, err := getter.Get(addr)
	if err != nil {
		t.Errorf("err occur: %v \n", err.Error())
	}

	if memoryInfo.Mallocs == 0 {
		t.Errorf("got: 0, want: else 0.")
	}
}
