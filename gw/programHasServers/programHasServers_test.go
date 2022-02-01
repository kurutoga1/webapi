package programHasServers_test

import (
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	gp "webapi/gw/programHasServers"
	"webapi/pro/router"
	"webapi/test"
	"webapi/utils/file"
)

var (
	currentDir string
	addrs      []string
)

func init() {
	c, err := file.GetCurrentDir()
	if err != nil {
		log.Fatalln(err.Error())
	}
	currentDir = c

	addrs, _, err = test.GetStartedServers(router.New(), 3)
	if err != nil {
		log.Fatalln(err)
	}
}

func tearDown() {
	os.RemoveAll(filepath.Join(currentDir, "fileserver"))
	os.Remove(filepath.Join(currentDir, "log.txt"))
}

func TestGet(t *testing.T) {
	programHasServersGetter := gp.GetProgramHasServersGetter()

	programHasServers, err := programHasServersGetter.Get(addrs, "/json/program/all", "convertToJson")
	if err != nil {
		t.Errorf("err from Get(): %v \n", err.Error())
	}

	if !reflect.DeepEqual(addrs, programHasServers) {
		t.Errorf("aliveServers(%v) is not equal programHasServers(%v) \n", addrs, programHasServers)
	}

	t.Cleanup(func() {
		tearDown()
	})

}

func TestIsProgramHasServer(t *testing.T) {
	t.Run("success test", func(t *testing.T) {
		testIsProgramHasServer(t, addrs[0], "convertToJson", true)
	})

	t.Run("fail test", func(t *testing.T) {
		testIsProgramHasServer(t, addrs[0], "toJson", false)
	})

}

func testIsProgramHasServer(t *testing.T, url, programName string, wantBool bool) {
	url = url + "/json/program/all"
	ok, err := gp.IsProgramHasServer(url, programName)
	if err != nil {
		t.Errorf(err.Error())
	}

	if ok != wantBool {
		t.Errorf("IsProgramHasServer(): %v, want: %v \n", ok, wantBool)
	}

}
