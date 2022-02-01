package minimumServerSelector_test

import (
	"log"
	"os"
	"path/filepath"
	"testing"
	"webapi/gw/config"
	mg "webapi/gw/memoryGetter"
	"webapi/gw/minimumServerSelector"
	sc "webapi/gw/serverAliveConfirmer"
	"webapi/pro/router"
	"webapi/test"
	"webapi/utils/file"
)

var (
	memoryGetter                mg.Getter
	currentDir                  string
	minimumMemoryServerSelector = minimumServerSelector.NewMinimumMemoryServerSelector()
	serverAliveConfirmer        = sc.NewServerAliveConfirmer()
	// CreateDummyServer関数のcleanUpで削除するファイルたち
	deletes []string
	cfg     = config.NewServerConfig()
	addrs   []string
)

func init() {
	// サーバを立てるとカレントディレクトリにfileserverディレクトリとlog.txt
	// ができるのでそれを削除する。
	c, err := file.GetCurrentDir()
	if err != nil {
		log.Fatalln(err.Error())
	}
	currentDir = c
	memoryGetter = mg.NewMemoryGetter()
	addrs, _, err = test.GetStartedServers(router.New(), 3)
	if err != nil {
		log.Fatalln(err)
	}
}

func tearDown() {
	os.RemoveAll(filepath.Join(currentDir, "fileserver"))
	os.RemoveAll(filepath.Join(currentDir, "log.txt"))
}

func TestSelect(t *testing.T) {

	selectedServer, err := minimumMemoryServerSelector.Select(addrs, serverAliveConfirmer, memoryGetter, cfg.GetMemoryEndPoint, "/health")
	if err != nil {
		t.Errorf("err occur: %v \n", err.Error())
	}

	contain := false
	for _, addr := range addrs {
		if selectedServer == addr {
			contain = true
		}
	}
	if !contain {
		t.Errorf("selected pro: %v doesn't contain servers: %v", selectedServer, addrs)
	}

}

func TestGetMinimumMemoryServer(t *testing.T) {
	serverInfoMap := map[string]uint64{
		"http://127.0.0.1:8083": 3000,
		"http://127.0.0.1:8081": 2597,
		"http://127.0.0.1:8082": 2700,
	}

	wantAddr := "http://127.0.0.1:8081"
	addr := minimumServerSelector.GetMinimumMemoryServer(serverInfoMap)
	if addr != wantAddr {
		t.Errorf("GetMinimumMemoryServer(): %v, want: %v \n", addr, wantAddr)
	}

	t.Cleanup(func() {
		tearDown()
	})
}

func TestGetServerMemoryMap(t *testing.T) {
	serverInfoMap, err := minimumServerSelector.GetServerMemoryMap(addrs, "/json/health/memory", memoryGetter)
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(serverInfoMap) != 3 {
		t.Errorf("len(serverInfoMap) is not 3. got: %v \n", len(serverInfoMap))
	}

	for _, memory := range serverInfoMap {
		if memory < 1 {
			t.Errorf("memory(%v) is not more than 1. \n", memory)
		}
	}
	t.Cleanup(func() {
		tearDown()
	})
}
