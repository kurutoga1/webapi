package getAllPrograms_test

import (
	"log"
	"os"
	"path/filepath"
	"testing"
	gg "webapi/gw/getAllPrograms"
	gs "webapi/gw/serverAliveConfirmer"
	"webapi/tests"
	"webapi/utils"
)

var currentDir string
var addrs []string
var ports []string

func init() {
	c, err := utils.GetCurrentDir()
	if err != nil {
		log.Fatalln(err.Error())
	}
	currentDir = c

	addrs, ports, err = tests.GetSomeStartedProgramServer(3)
	if err != nil {
		panic(err)
	}
}

func tearDown() {
	os.RemoveAll(filepath.Join(currentDir, "fileserver"))
	os.Remove(filepath.Join(currentDir, "log.txt"))
}

func TestGetAllPrograms(t *testing.T) {
	serverAliveConfirmer := gs.NewServerAliveConfirmer()
	aliveServers, err := gs.GetAliveServers(addrs, "/health", serverAliveConfirmer)
	if err != nil {
		t.Errorf("err from GetAliveServers(): %v \n", err.Error())

	}
	allProgramGetter := gg.NewAllProgramGetter()

	endPoint := "/json/program/all"
	allProgramMap, err := allProgramGetter.Get(aliveServers, endPoint)
	if err != nil {
		t.Errorf("err from Get(): %v \n", err.Error())
	}

	if _, ok := allProgramMap["convertToJson"]; !ok {
		t.Errorf("convertToJson is not found. of %v. \n", allProgramMap)
	}

	t.Cleanup(func() {
		tearDown()
	})

}
